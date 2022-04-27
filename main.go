package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type healthCheckResponse struct {
	Status string `json:"status"`
}

type CassandraNode struct {
	Name string `json:"name"`
	IPA  string `json:"ipa"`
}

type CassandraCluster []CassandraNode

var counter float64 = 1
var (
	onlineUsers = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "oneline_user_count",
		Help: "Current number of online users",
	})
	error404 = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "notfound_errors_total",
		Help: "Number of 404 errors.",
	})
)

func init() {
	prometheus.MustRegister(onlineUsers)
	prometheus.MustRegister(error404)
}

func main() {
	onlineUsers.Set(65.3)
	error404.Inc()

	router := mux.NewRouter()
	router.Handle("/metrics", promhttp.Handler())

	router.HandleFunc("/health", health).Methods("GET")
	router.HandleFunc("/shownodes", showNodes).Methods("GET")
	//router.LoggingHandlers(os.Stdout, r)
	loggedRouter := handlers.LoggingHandler(os.Stdout, router)
	startupMsg := `Starting router
	Listening on: -
	http://localhost:8080/metrics
	http://localhost:8080/health
	http://localhost:8080/shownodes`
	log.Println(startupMsg)
	http.ListenAndServe(":8080", loggedRouter)
}

func health(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "application/json")

	onlineUsers.Set(counter)
	counter = counter + 1
	error404.Inc() // increment the hd failures
	outgoingJSON, error := json.Marshal(healthCheckResponse{Status: "UP"})
	if error != nil {
		log.Println(error.Error())
		http.Error(res, error.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Fprint(res, string(outgoingJSON))
}

func showNodes(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "application/json")
	log.Println("Someone requested the cassandra nodelist")
	outgoingJSON, error := json.Marshal(CassandraNode{Name: "test", IPA: "1.1.1.1"})
	if error != nil {
		log.Println(error.Error())
		http.Error(res, error.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Fprint(res, string(outgoingJSON))
}
