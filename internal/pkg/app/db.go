package app

import (
	"os"

	"github.com/leonideliseev/songLibraryCrud/pkg/postgresql"
	"github.com/leonideliseev/songLibraryCrud/schema"
	"github.com/spf13/viper"
)

func (a *App) initDBConn() {
	config := postgresql.Config{
		Host:     viper.GetString("db.host"),
		Port:     viper.GetString("db.port"),
		Username: viper.GetString("db.username"),
		Password: os.Getenv("DB_PASSWORD"),
		DBName:   viper.GetString("db.dbname"),
		SSLMode:  viper.GetString("db.sslmode"),
	}

	configTest := config
	configTest.DBName = "postgres"
	connTest, err := postgresql.ConnWithPgxPool(configTest)
	if err != nil {
		a.logger.WithError(err).Fatal("error connect to database")
	}
	defer connTest.Close()

	postgresql.CreateDatabaseIfNotExists(connTest, viper.GetString("db.dbname"), a.logger)

	conn, err := postgresql.ConnWithPgxPool(config)
	if err != nil {
		a.logger.WithError(err).Fatal("error connect to database")
	}

	postgresql.Migrate(a.logger, &schema.DB, &config)
	a.logger.Info("database ready")

	a.conn = conn
}
