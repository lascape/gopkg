package seleniumx

import (
	"github.com/sirupsen/logrus"
	"github.com/tebeka/selenium"
	"github.com/tebeka/selenium/chrome"
	"os"
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

type Config struct {
	DriverPath string `json:"driver_path" yaml:"driver_path"`
	Port       int    `json:"port" yaml:"port"`
}

func Must(opts ...Option) selenium.Capabilities {
	client, err := must(opts...)
	if err != nil {
		logrus.Error(err)
		return nil
	}

	return client
}

func must(opts ...Option) (selenium.Capabilities, error) {
	var (
		o            Options
		seleniumOpts = []selenium.ServiceOption{
			selenium.Output(os.Stderr), // Output debug information to STDERR.
		}
		err error
	)

	for _, opt := range opts {
		opt(&o)
	}

	//SetDebug 设置调试模式
	selenium.SetDebug(true)
	//在后台启动一个ChromeDriver实例
	_, err = selenium.NewChromeDriverService("/usr/bin/chromedriver", 8080, seleniumOpts...)
	//_, err = selenium.NewChromeDriverService("/Users/jianfei/Downloads/chromedriver", 8080, seleniumOpts...)
	if err != nil {
		return nil, err // panic is used only as an example and is not otherwise recommended.
	}

	seleniumCaps := selenium.Capabilities{
		"browserName": "chrome",
	}
	seleniumCaps.AddChrome(chrome.Capabilities{
		Args: []string{
			"--no-sandbox",
			"--headless",
		},
	})

	return seleniumCaps, nil
}
