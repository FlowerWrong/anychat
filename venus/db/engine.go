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
	dbEngine, err = xorm.NewPostgreSQL(viper.Get("db_url").(string))
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
	dbOnce.Do(func() {
		err := initEngine()
		if err != nil {
			panic(err)
		}
	})
	return dbEngine
}
