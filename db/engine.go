package db

import (
	"sync"

	// justifying it
	_ "github.com/lib/pq"
	"github.com/spf13/viper"
	"github.com/xormplus/xorm"
)

var (
	dbEngine *xorm.Engine
	dbOnce   sync.Once
)

func initEngine() error {
	var err error
	dbEngine, err = xorm.NewPostgreSQL(viper.GetString("db_url"))
	if err != nil {
		return err
	}
	err = dbEngine.RegisterSqlTemplate(xorm.Default("./models/sql", ".tpl"))
	if err != nil {
		return err
	}
	err = dbEngine.StartFSWatcher()
	if err != nil {
		return err
	}
	err = dbEngine.Ping()
	if err != nil {
		return err
	}
	dbEngine.ShowSQL(true)
	return nil
}

// Engine return the uniq db engine
func Engine() *xorm.Engine {
	if dbEngine == nil {
		dbOnce.Do(func() {
			err := initEngine()
			if err != nil {
				panic(err)
			}
		})
	}
	return dbEngine
}
