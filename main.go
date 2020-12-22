package main

import (
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/eltropycal/controllers"
	"github.com/eltropycal/dataservices"
	"github.com/eltropycal/prometheus"
	"github.com/eltropycal/router"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

var serviceName = "Eltropycal"

func main() {
	fmt.Printf("Starting %v\n", serviceName)
	initializeLogrus()
	initializeViper()
	initializeDatabase()
	prometheus.RegisterPrometheusMetrics()
	startServer("9090")
}

func initializeDatabase() {
	controllers.DataService = &dataservices.PostgresClient{}
	controllers.DataService.Connect()
}

func initializeViper() {
	viper.SetConfigName("appConfig")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatal("Can't open config file.")
	}
}

func initializeLogrus() {
	Formatter := new(log.JSONFormatter)
	Formatter.TimestampFormat = "02-01-2006 15:04:05"
	log.SetFormatter(Formatter)

	var filename = "/tmp/eltropycal/log/logfile_rolling.log"
	f, err := os.OpenFile(filename, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)
	if err != nil {
		log.Error("Error occurred while opening log file, make sure you have the directory created, logging to stdout : " + err.Error())
	} else {
		multiWriter := io.MultiWriter(os.Stdout, f)
		log.SetOutput(multiWriter)
	}
}

func startServer(port string) {
	r := router.NewRouter()
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./static/"))))
	http.Handle("/", promhttp.InstrumentHandlerCounter(prometheus.AllMetrics.HttpRequestsTotal, r))
	http.Handle("/metrics", promhttp.Handler())
	log.Info("Starting HTTP service at " + port)
	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		log.Error("An error occurred starting HTTP listener at port " + port + " error: " + err.Error())
	}
}
