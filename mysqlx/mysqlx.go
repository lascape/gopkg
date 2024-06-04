package mysqlx

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"reflect"
	"time"

	"github.com/pkg/errors"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"gorm.io/plugin/dbresolver"
)

type Options struct {
	conf        Config
	hooks       []func(*gorm.DB)
	prepareStmt bool
}

type Option func(*Options)

func WithConfig(conf Config) Option {
	return func(o *Options) {
		o.conf = conf
	}
}

func WithHooks(hks ...func(*gorm.DB)) Option {
	return func(o *Options) {
		for _, hk := range hks {
			if hk == nil {
				continue
			}
			o.hooks = append(o.hooks, hk)
		}
	}
}

// WithPrepareStmt 我们现在支持用配置项`prepare_stmt`来设置在不同环境的prepare_stmt效果，
// 同时你可以使用这个方法来强制修改db实例的prepare_stmt开关
func WithPrepareStmt() Option {
	return func(o *Options) {
		o.conf.PrepareStmt = true
	}
}

// Config mysql的配置字段
type Config struct {
	Source             string          `yaml:"source" json:"source"` // 如果该字段不为空，则直接适用打开该链接
	Addr               string          `yaml:"addr" json:"addr"`
	Password           string          `yaml:"password" json:"password" kms:"decode"`
	User               string          `yaml:"user" json:"user"`
	DbName             string          `yaml:"db_name" json:"db_name"`
	Timeout            int             `yaml:"timeout" json:"timeout"`                           // 单位：秒钟
	WriteTimeout       int             `yaml:"write_timeout" json:"write_timeout"`               // 单位：秒钟
	ReadTimeout        int             `yaml:"read_timeout" json:"read_timeout"`                 // 单位：秒钟
	MaxOpenConnections int             `yaml:"max_open_connections" json:"max_open_connections"` // 设置打开数据库连接的最大数量。
	MaxIdleConnections int             `yaml:"max_idle_connections" json:"max_idle_connections"` // 用于设置连接池中空闲连接的最大数量
	MaxLifetime        int             `yaml:"max_lifetime" json:"max_lifetime"`                 // 单位：秒钟 设置了连接可复用的最大时间。
	LogLevel           logger.LogLevel `yaml:"log_level" json:"log_level"`
	ReadConfig         []Config        `yaml:"read_config" json:"read_config"`   // 读配置的数据库配置，内部为空的字段将用写的配置替代
	PrepareStmt        bool            `yaml:"prepare_stmt" json:"prepare_stmt"` //预编译语句开关
}

func Must(opts ...Option) *gorm.DB {
	db, err := mustDb(opts...)
	if err != nil {
		fmt.Println(err.Error())
		log.Fatal(err)
	}
	return db
}

func mustDb(opts ...Option) (*gorm.DB, error) {
	o := &Options{}
	for _, opt := range opts {
		opt(o)
	}
	if !o.prepareStmt {
		o.prepareStmt = o.conf.PrepareStmt
	}
	o.conf.configInitialize()
	db, err := gorm.Open(mysql.Open(o.conf.DSN()), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true, // 使用单数表名，启用该选项，此时，`User` 的表名应该是 `t_user`
		},
		Logger: func() logger.Interface {
			if o.conf.LogLevel > 0 {
				return logger.New(log.New(os.Stdout, "\r\n", log.LstdFlags), logger.Config{
					SlowThreshold: 0 * time.Millisecond,
					LogLevel:      o.conf.LogLevel,
					Colorful:      true,
				})
			}
			return logger.Discard
		}(),
		DisableForeignKeyConstraintWhenMigrating: true,          // 禁用自动创建外键关联
		PrepareStmt:                              o.prepareStmt, // 预编译语句开关
	})
	if err != nil {
		return nil, errors.Errorf("init mysqlx err: %v", err)
	}

	for _, m := range o.hooks {
		m(db)
	}

	originSqlDBConnPool, err := getOriginSqlDBConnPool(db, o.prepareStmt)
	if err != nil {
		return nil, errors.Errorf("init mysqlx getOriginSqlDBConnPool err: %v", err)
	}
	sourceDialector := mysql.New(mysql.Config{
		Conn: originSqlDBConnPool,
	})

	// 注册读写分离
	_ = db.Use(dbresolver.Register(
		dbresolver.Config{
			Sources:  []gorm.Dialector{sourceDialector},
			Replicas: dialector(o.conf.ReadConfig...),
		}).
		SetMaxIdleConns(o.conf.MaxIdleConnections).
		SetMaxOpenConns(o.conf.MaxOpenConnections).
		SetConnMaxLifetime(time.Duration(o.conf.MaxLifetime) * time.Second))
	return db, nil
}

func getOriginSqlDBConnPool(db *gorm.DB, prepareStmt bool) (gorm.ConnPool, error) {
	originSqlDBConnPool := db.ConnPool

	if prepareStmt {
		preparedStmtDB, ok := db.ConnPool.(*gorm.PreparedStmtDB)
		if !ok {
			return nil, errors.New("db.ConnPool.(*gorm.PreparedStmtDB) fail")
		}
		originSqlDBConnPool = preparedStmtDB.ConnPool
	}

	_, ok := originSqlDBConnPool.(*sql.DB)
	if !ok {
		return nil, errors.New("originConnPool.(*sql.DB) fail")
	}

	return originSqlDBConnPool, nil
}

func dialector(c ...Config) []gorm.Dialector {
	var dialectors []gorm.Dialector
	for _, config := range c {
		dialectors = append(dialectors,
			mysql.Open(config.DSN()),
		)
	}
	return dialectors
}

func (c *Config) DSN() string {
	if c.Source != "" {
		return c.Source
	}
	c.defaultTimeout()
	dsn := fmt.Sprintf(`%s:%s@tcp(%s)/%s?timeout=%ds&readTimeout=%ds&writeTimeout=%ds&charset=utf8mb4,utf8&parseTime=true`,
		c.User, c.Password, c.Addr, c.DbName, c.Timeout, c.ReadTimeout, c.WriteTimeout) + "&loc=Asia%2FShanghai"
	return dsn
}

// defaultTimeout 设置默认的超时时间
func (c *Config) defaultTimeout() {
	if c.Timeout == 0 {
		c.Timeout = 3
	}
	if c.WriteTimeout == 0 {
		c.WriteTimeout = 3
	}
	if c.ReadTimeout == 0 {
		c.ReadTimeout = 3
	}
}

// configInitialize 初始化配置文件，根据写数据库动态变更读数据库配置
func (c *Config) configInitialize() {
	reflectRe := func(cv, v reflect.Value) {
		for i := 0; i < v.NumField(); i++ {
			if v.Field(i).IsZero() && v.Field(i).Kind() != reflect.Slice {
				cvName := v.Type().Field(i).Name
				cvValue := cv.FieldByName(cvName)
				v.Field(i).Set(cvValue)
			}
		}
	}
	value := reflect.ValueOf(c).Elem()
	readConfig := value.FieldByName("ReadConfig")
	for i := 0; i < readConfig.Len(); i++ {
		v := readConfig.Index(i)
		reflectRe(value, v)
	}
}
