package data

import (
	"fmt"
	"geek-user-service/internal/conf"
	"github.com/go-kratos/kratos/v2/log"
	_ "github.com/go-sql-driver/mysql"
	"github.com/google/wire"
	"github.com/jinzhu/gorm"
	"time"
)

// ProviderSet is data providers.
var ProviderSet = wire.NewSet(NewData, NewGreeterRepo, NewUserRepo)
var AutoMigrates = make([]interface{}, 0)


// Data .
type Data struct {
	// TODO warpped database client
	db *gorm.DB
	log *log.Helper
}

// NewData .
func NewData(c *conf.Data, logger log.Logger) (*Data, func(), error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True", c.Database.User, c.Database.Password, c.Database.Host, c.Database.Schema)
	logHelper := log.NewHelper(log.With(logger, "module", "data"))
	db, err := gorm.Open(c.Database.Driver, dsn)
	if err != nil {
		logHelper.Errorf("NewData error: %v", err)
		panic(err)
	}

	db.DB().SetMaxIdleConns(int(c.Database.MaxIdleConns))
	db.DB().SetMaxOpenConns(int(c.Database.MaxOpenConns))
	db.DB().SetConnMaxLifetime(time.Duration(c.Database.ConnMaxLifeTime) * time.Second)
	db.LogMode(c.Database.Debug)

	db.Set("gorm:table_options", "ENGINE=InnoDB CHARSET=utf8mb4").AutoMigrate(AutoMigrates...)
	//db.AutoMigrate(AutoMigrates...)
	cleanup := func() {
		logHelper.Info("NewData: closing the data resources")
		db.DB().Close()
	}
	return &Data{
		log: logHelper,
		db: db,
	}, cleanup, nil
}

func RegisterAutoMigrates(migrates interface{}) {
	AutoMigrates = append(AutoMigrates, migrates)
}
