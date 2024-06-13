package db

import (
	"fmt"
	"library-service/config"
	"log/slog"
	"os"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func connect(dbConfig config.DBConfig) *sqlx.DB {
	dbSource := fmt.Sprintf(
		"user=%s password=%s host=%s port=%d dbname=%s sslmode=",
		dbConfig.User,
		dbConfig.Pass,
		dbConfig.Host,
		dbConfig.Port,
		dbConfig.Name,
	)

	if !dbConfig.EnableSSLMode {
		dbSource += "disable"
	} else {
		dbSource += "enable"
	}

	dbCon, err := sqlx.Connect("postgres", dbSource)
	if err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}

	dbCon.SetConnMaxIdleTime(
		time.Duration(dbConfig.MaxIdleTimeInMinute * int(time.Minute)),
	)

	return dbCon
}

func ConnectDB() {
	conf := config.GetConfig()

	readDb = connect(conf.DB.Read)
	slog.Info("Connected to read database")

	writeDb = connect(conf.DB.Write)
	slog.Info("Connected to write database")
}

func CloseDB() {
	if err := readDb.Close(); err != nil {
		slog.Error(err.Error())
		return
	}
	slog.Info("Disconnected from read database")

	if err := writeDb.Close(); err != nil {
		slog.Error(err.Error())
		return
	}
	slog.Info("Disconnected from write database")
}
