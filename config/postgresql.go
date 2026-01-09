package config

import (
	"database/sql"
	"fmt"
	"github.com/spf13/viper"
)

func ConnectToDB() (*sql.DB, error) {
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./config")
	viper.SetConfigName("env")
	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s "+
		"password=%s dbname=%s sslmode=disable",
		viper.GetString("database.host"), viper.GetString("database.port"),
		viper.GetString("database.username"), viper.GetString("database.password"),
		viper.GetString("database.name"))
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		return nil, err
	}
	fmt.Println("Successfully connected to DB!")
	return db, nil
}
