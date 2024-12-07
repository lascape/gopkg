package gormx

import (
	"fmt"
	"log"
	"os"
	"reflect"
	"time"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"gorm.io/plugin/dbresolver"
)

type Options struct {
	conf  Config
	debug bool
}

type Option func(*Options)

func WithConfig(conf Config) Option {
	return func(o *Options) {
		o.conf = conf
	}
}

// WithDebug 我们现在支持用配置项`debug`来设置在不同环境的debug效果，
// 同时你可以使用这个方法来强制修改db实例的debug开关
func WithDebug() Option {
	return func(o *Options) {
		o.debug = true
	}
}

// Config mysql的配置字段
type Config struct {
	Source             string   `yaml:"source"` //如果该字段不为空，则直接适用打开该链接
	Addr               string   `yaml:"addr"`
	Password           string   `yaml:"password" kms:"encode"`
	User               string   `yaml:"user"`
	DbName             string   `yaml:"db_name"`
	Timeout            int      `yaml:"timeout"`              //单位：秒钟
	WriteTimeout       int      `yaml:"write_timeout"`        //单位：秒钟
	ReadTimeout        int      `yaml:"read_timeout"`         //单位：秒钟
	MaxOpenConnections int      `yaml:"max_open_connections"` //设置打开数据库连接的最大数量。
	MaxIdleConnections int      `yaml:"max_idle_connections"` //用于设置连接池中空闲连接的最大数量
	MaxLifetime        int      `yaml:"max_lifetime"`         //单位：秒钟 设置了连接可复用的最大时间。
	LogLevel           string   `yaml:"log_level"`
	SslMode            string   `yaml:"ssl_mode"`
	ReadConfig         []Config `yaml:"read_config"`
}

func Must(opts ...Option) *gorm.DB {
	db, err := mustDb(opts...)
	if err != nil {
		logrus.Error(err)
	}
	return db
}

func mustDb(opts ...Option) (*gorm.DB, error) {
	k := &Options{}
	for _, o := range opts {
		o(k)
	}
	k.conf.configInitialize()
	db, err := gorm.Open(mysql.Open(k.conf.DSN()), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true, // 使用单数表名，启用该选项，此时，`User` 的表名应该是 `t_user`
		},
		Logger: func() logger.Interface {
			if k.debug {
				return logger.New(log.New(os.Stdout, "\r\n", log.LstdFlags), logger.Config{
					SlowThreshold: 0 * time.Millisecond,
					LogLevel:      logger.Info,
					Colorful:      true,
				})
			}
			return logger.Discard
		}(),
		DisableForeignKeyConstraintWhenMigrating: true,  //禁用自动创建外键关联
		PrepareStmt:                              false, //执行任何 SQL 时都创建并缓存预编译语句，可以提高后续的调用速度
	})
	if err != nil {
		return nil, errors.Errorf("init mysqlx err: %v", err)
	}

	//注册读写分离
	err = db.Use(dbresolver.Register(
		dbresolver.Config{
			Sources:  dialector(k.conf),
			Replicas: dialector(k.conf.ReadConfig...),
		}).
		SetMaxIdleConns(k.conf.MaxIdleConnections).
		SetMaxOpenConns(k.conf.MaxOpenConnections).
		SetConnMaxLifetime(time.Duration(k.conf.MaxLifetime) * time.Second))
	if err != nil {
		return nil, errors.Errorf("init mysqlx.dbresolver err: %v", err)
	}
	return db, nil
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
		c.Timeout = 20
	}
	if c.WriteTimeout == 0 {
		c.WriteTimeout = 20
	}
	if c.ReadTimeout == 0 {
		c.ReadTimeout = 20
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
