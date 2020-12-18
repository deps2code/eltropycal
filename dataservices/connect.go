package dataservices

import (
	"database/sql"
	"fmt"

	"github.com/eltropycal/constants"
	_ "github.com/lib/pq"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// Connect opens a connection to the DB
func (pc *PostgresClient) Connect() {
	config := dbConfig()
	var err error
	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s",
		config[constants.DBHOST], config[constants.DBPORT], config[constants.DBUSER], config[constants.DBPASS], config[constants.DBNAME])

	pc.DB, err = sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Fatal("Fatal error encountered: ", err.Error())
	}
	err = pc.DB.Ping()
	if err != nil {
		log.Fatal("Unable to ping database. Shutting down server.")
	}
	pc.DB.SetMaxOpenConns(10)
	pc.DB.SetMaxIdleConns(2)
	log.Info("Successfully connected to database.")
}

func dbConfig() map[string]string {
	var host, port, user, password, name string
	conf := make(map[string]string)

	host = viper.GetString(constants.BACKENDDB + "." + constants.DBHOST)
	if host == "" {
		log.Fatal("DBHOST variable required but not set")
	}
	port = viper.GetString(constants.BACKENDDB + "." + constants.DBPORT)
	if port == "" {
		log.Fatal("DBPORT variable required but not set")
	}
	user = viper.GetString(constants.BACKENDDB + "." + constants.DBUSER)
	if user == "" {
		log.Fatal("DBUSER variable required but not set")
	}
	password = viper.GetString(constants.BACKENDDB + "." + constants.DBPASS)
	if password == "" {
		log.Fatal("DBPASS variable required but not set")
	}
	name = viper.GetString(constants.BACKENDDB + "." + constants.DBNAME)
	if name == "" {
		log.Fatal("DBNAME variable required but not set")
	}
	conf[constants.DBHOST] = host
	conf[constants.DBPORT] = port
	conf[constants.DBUSER] = user
	conf[constants.DBPASS] = password
	conf[constants.DBNAME] = name
	return conf
}
