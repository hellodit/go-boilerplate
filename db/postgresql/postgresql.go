package postgresql

import (
	"fmt"
	"github.com/go-pg/pg/v10"
	"github.com/spf13/viper"
)

func Connect() *pg.DB {
	dbHost := viper.GetString("DB_HOTS")
	dbPort := viper.GetString("DB_PORT")
	dbUser := viper.GetString("DB_USER")
	dbPass := viper.GetString("DB_PASSWORD")
	dbName := viper.GetString("DB_NAME")
	dbSslMode := viper.GetString("DB_SSL_MODE")

	parse := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s", dbUser, dbPass, dbHost, dbPort, dbName, dbSslMode)
	opt, err := pg.ParseURL(parse)

	if err != nil {
		panic(err)
	}
	db := pg.Connect(opt)

	if db == nil {
		panic(fmt.Errorf("failed to connect database: %s", db))
	}
	return db
}
