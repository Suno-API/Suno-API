package po

import (
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"os"
	"strings"
	"suno-api/common"
	"suno-api/lib/gormlogger"
	"time"
)

var DB *gorm.DB

func InitDB() (err error) {
	db, err := chooseDB(os.Getenv("SQL_DSN"))
	if err == nil {
		if common.DebugEnabled {
			db = db.Debug()
		}
		DB = db
		sqlDB, err := DB.DB()
		if err != nil {
			return err
		}
		sqlDB.SetMaxIdleConns(common.GetOrDefault("SQL_MAX_IDLE_CONNS", 100))
		sqlDB.SetMaxOpenConns(common.GetOrDefault("SQL_MAX_OPEN_CONNS", 1000))
		sqlDB.SetConnMaxLifetime(time.Second * time.Duration(common.GetOrDefault("SQL_MAX_LIFETIME", 60)))

		common.Logger.Info("database migration started")
		err = db.AutoMigrate(&Task{})
		if err != nil {
			return err
		}
		common.Logger.Info("database migrated")
		return err
	} else {
		common.Logger.Fatal(err.Error())
	}
	return err
}

func chooseDB(sqlDns string) (*gorm.DB, error) {
	gormLog, err := gormlogger.NewGormV2Logger(common.Logger)
	if err != nil {
		return nil, err
	}
	gormConfig := &gorm.Config{
		PrepareStmt: true, // precompile SQL
		Logger:      gormLog,
	}

	if sqlDns != "" {
		dsn := sqlDns
		if strings.HasPrefix(dsn, "postgres://") {
			// Use PostgreSQL
			common.Logger.Info("using PostgreSQL as database")
			//common.UsingPostgreSQL = true
			return gorm.Open(postgres.New(postgres.Config{
				DSN:                  dsn,
				PreferSimpleProtocol: true, // disables implicit prepared statement usage
			}), gormConfig)
		}
		// Use MySQL
		common.Logger.Info("using MySQL as database")
		// check parseTime
		if !strings.Contains(dsn, "parseTime") {
			if strings.Contains(dsn, "?") {
				dsn += "&parseTime=true"
			} else {
				dsn += "?parseTime=true"
			}
		}
		return gorm.Open(mysql.Open(dsn), gormConfig)
	}
	// Use SQLite
	common.Logger.Info("SQL_DSN not set, using SQLite as database")
	//common.UsingSQLite = true
	return gorm.Open(sqlite.Open(common.SQLitePath), gormConfig)
}

func CloseDB() error {
	sqlDB, err := DB.DB()
	if err != nil {
		return err
	}
	err = sqlDB.Close()
	return err
}
