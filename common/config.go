package common

import (
	"flag"
	"github.com/koding/multiconfig"
	"os"
)

type FlagConfig struct {
	ConfigFile string `default:"config.json"`
}
type LoggerConfig struct {
	Enabled    bool `default:true`
	LogFile    string
	TraceLevel int `default:3`
}
type ServerConfig struct {
	Listen      string `default:":5000"`
	RuntimePath string `default:"runtime"`
	Redis       struct {
		Network     string
		Address     string
		Database    int `default:"0"`
		Password    string
		ShowCommand bool `default:"false"`
	}
	Db struct {
		Host        string
		Port        string `default:"3306"`
		Name        string
		User        string
		Password    string
		SlaveConfig struct {
			User     string
			Password string
		}
		Slaves []struct {
			Host string
			Port string
			Name string
		}
		MaxOpenConns int  `default:"0"`
		ShowSQL      bool `default:"false"`
	}
	Log struct {
		LogFile    string `default:""`
		SaveType   string `default:"d"`
		TraceLevel int    `default:"3"`
		Logger     struct {
			Trace LoggerConfig
			Info  LoggerConfig
			Warn  LoggerConfig
			Error LoggerConfig
		}
	}
	RpcServer struct{
		//项目RPC设置
		Project struct{
			User     string  `default:""`
			Vendor   string  `default:""`
			Cart     string  `default:""`
			Products string  `default:""`
			Coupon   string  `default:""`
			Push     string  `default:""`
		}
	}
	System struct{
		Debug bool `default:"false"`
	}
}

type AppConfig struct {
	Project struct{
		//定时取消-待支付下单时间大于5*60 (5分钟)
		SCHEDULE_CANCEL_TIME 					int 	`default:"300"`

		//定时完成-订单超过8*60*60自动变为已完成(8小时)
		SCHEDULE_COMPLETE_TIME					int 	`default:"28800"`

		//定时取消-10*60内商家未接单，订单自动变为已取消(10分钟)
		SCHEDULE_CANCEL_MERCHANTS_TIME			int 	`default:"600"`

		//定时取消-2*60内商家未接单，邮件提醒(2分钟)
		SCHEDULE_CANCEL_EMAIL_MERCHANTS_TIME	int  	`default:"120"`

		//定时取消-在用户取餐时间前5*60，推送消息提醒取餐(5分钟)
		BEFORE_CANCEL_MERCHANTS_TIME			int 	`default:"300"`

		//十五分钟没有确认收货，推送消息提醒取餐(15分钟)
		WITHOUT_CONFIRM_TAKE_NOTIFY_TIME       int 	`default:"900"`

		//商家接收售后订单，超过24*60*60未处理，售后订单自动确认退款(24小时)
		AUTO_AFTER_SALE_REFUND_TIME			int 	`default:"86400"`
	}
}

func (c *FlagConfig) load() error {
	t := &multiconfig.TagLoader{}
	f := &multiconfig.FlagLoader{}
	m := multiconfig.MultiLoader(t, f)
	if err := m.Load(c);err == flag.ErrHelp {
		os.Exit(0)
	} else if err != nil {
		return err
	}
	return nil
}
func (c *ServerConfig) load() error {
	f := &FlagConfig{}
	err := f.load()
	if err == flag.ErrHelp {
		os.Exit(0)
	} else if err != nil {
		return err
	}
	t   := &multiconfig.TagLoader{}
	j   := &multiconfig.JSONLoader{Path:f.ConfigFile}
	m   := multiconfig.MultiLoader(t, j)
	err =m.Load(c)
	return err
}
