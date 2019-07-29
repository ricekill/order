package common

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/core"
	"github.com/go-xorm/xorm"
)

func OpenDb() error {
	var err error
	master, err := xorm.NewEngine("mysql",
		fmt.Sprintf("%s:%s@tcp(%s:%s)/%s",
			Config.Db.User,
			Config.Db.Password,
			Config.Db.Host,
			Config.Db.Port,
			Config.Db.Name))
	if err != nil {
		return err
	}
	slaves := make([]*xorm.Engine, len(Config.Db.Slaves))
	for i, slave := range Config.Db.Slaves {
		slaves[i], err=xorm.NewEngine("mysql",
			fmt.Sprintf("%s:%s@tcp(%s:%s)/%s",
				Config.Db.SlaveConfig.User,
				Config.Db.SlaveConfig.Password,
				slave.Host,
				slave.Port,
				slave.Name))
		if err != nil {
			return err
		}
	}
	if len(slaves)<=0 {
		slaves =append(slaves,master)
	}
	DB, err = xorm.NewEngineGroup(master, slaves)
	DB.SetMaxOpenConns(Config.Db.MaxOpenConns)
	DB.ShowSQL(Config.Db.ShowSQL)
	DB.ShowExecTime(Config.Db.ShowSQL)
	DB.SetLogger(&XOrmLogger{logger:Log})
	return err
}
type XOrmLogger struct {
	logger *Logger
	level core.LogLevel
	showSQL bool
}

// Error implement core.ILogger
func (x *XOrmLogger) Error(v ...interface{}) {
	if x.level <= core.LOG_ERR {
		Log.Errorln(v...)
	}
	return
}

// Errorf implement core.ILogger
func (x *XOrmLogger) Errorf(format string, v ...interface{}) {
	if x.level <= core.LOG_ERR {
		Log.Errorf(format, v...)
	}
	return
}

// Debug implement core.ILogger
func (x *XOrmLogger) Debug(v ...interface{}) {
	if x.level <= core.LOG_DEBUG {
		Log.Traceln(v...)
	}
	return
}

// Debugf implement core.ILogger
func (x *XOrmLogger) Debugf(format string, v ...interface{}) {
	if x.level <= core.LOG_DEBUG {
		Log.Tracef(format, v...)
	}
	return
}

// Info implement core.ILogger
func (x *XOrmLogger) Info(v ...interface{}) {
	if x.level <= core.LOG_INFO {
		Log.Infoln(v...)
	}
	return
}

// Infof implement core.ILogger
func (x *XOrmLogger) Infof(format string, v ...interface{}) {
	if x.level <= core.LOG_INFO {
		Log.Infof(format, v...)
	}
	return
}

// Warn implement core.ILogger
func (x *XOrmLogger) Warn(v ...interface{}) {
	if x.level <= core.LOG_WARNING {
		Log.Warnln(v...)
	}
	return
}

// Warnf implement core.ILogger
func (x *XOrmLogger) Warnf(format string, v ...interface{}) {
	if x.level <= core.LOG_WARNING {
		Log.Warnf(format, v...)
	}
	return
}

// Level implement core.ILogger
func (x *XOrmLogger) Level() core.LogLevel {
	return x.level
}

// SetLevel implement core.ILogger
func (x *XOrmLogger) SetLevel(l core.LogLevel) {
	x.level = l
	return
}

// ShowSQL implement core.ILogger
func (x *XOrmLogger) ShowSQL(show ...bool) {
	if len(show) == 0 {
		x.showSQL = true
		return
	}
	x.showSQL = show[0]
}

// IsShowSQL implement core.ILogger
func (x *XOrmLogger) IsShowSQL() bool {
	return x.showSQL
}
