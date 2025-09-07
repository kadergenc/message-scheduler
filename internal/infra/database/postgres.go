package database

import (
	"fmt"
	"message-scheduler/config"
	"message-scheduler/log"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewPostgresDB(conf config.PostgresConfig) *gorm.DB {
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		conf.WriteHost, conf.WritePort, conf.User, conf.Password, conf.DbName)

	var err error
	var db *gorm.DB

	for i := 0; i < 5; i++ { // Retry logic for database readiness
		db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
		if err == nil {
			break
		}
		log.Logger.Info().Msgf("Database connection failed. Retrying... (%d/5)\n", i+1)
		time.Sleep(2 * time.Second)
	}

	if err != nil {
		log.Logger.Fatal().Err(err).Msg("Could not connect to the database:")
	}

	log.Logger.Info().Msg("Connected to the database successfully!")
	return db
}
