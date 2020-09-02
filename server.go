package main

import (
	"os"
	"log"
	"time"
	"fmt"
	"net/http"
	// "database/sql"

	"github.com/gorilla/mux"
	"github.com/urfave/negroni"
)

var (
	tokenInfoEndpoint string
	tokenInfoUser     string
	tokenInfoPassword string

	mysqlHost     string
	mysqlDatabase string
	mysqlUsername string
	mysqlPassword string
	mysqlDSN      string // Data Source Name

	disableAuth = false // disable token introspection for testing purposes

	mainRouter *mux.Router
)
func init() {

	// token info
	tokenInfoEndpoint = os.Getenv("tokenInfoEndpoint")
	tokenInfoUser = os.Getenv("tokenInfoUser")
	tokenInfoPassword = os.Getenv("tokenInfoPassword")


	if os.Getenv("TESTING_NOAUTH") == "1" {
		disableAuth = true
		log.Printf("WARNING: token validation is disabled, use only for testing/development")
		time.Sleep(time.Second * 2)
	}

}

func createRouter() {
	mainRouter = mux.NewRouter()
	r := mainRouter

	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, `{"id": "SAGE Edge Scheduler","available_resources":["api/v1/","metrics/"]}`)
	})

	log.Println("Sage Edge Scheduler API")
	api := r.PathPrefix("/api/v1").Subrouter()
	api.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, `{"id": "SAGE Edge Scheduler","available_resources":["goals"]}`)
	})

	api.Handle("/metrics", negroni.New(
		negroni.HandlerFunc(authMW),
		negroni.Wrap(http.HandlerFunc(getSESstats)),
	)).Methods(http.MethodGet)

	api.Handle("/goals", negroni.New(
		negroni.HandlerFunc(authMW),
		negroni.Wrap(http.HandlerFunc(listGoals)),
	)).Methods(http.MethodGet)

	api.Handle("/goals", negroni.New(
		negroni.HandlerFunc(authMW),
		negroni.Wrap(http.HandlerFunc(createGoal)),
	)).Methods(http.MethodPost)

	api.NewRoute().PathPrefix("/goals/{goal}/status").Handler(negroni.New(
		negroni.HandlerFunc(authMW),
		negroni.Wrap(http.HandlerFunc(getGoalStatus)),
	)).Methods(http.MethodGet)

	api.NewRoute().PathPrefix("/goals/{goal}/status").Handler(negroni.New(
		negroni.HandlerFunc(authMW),
		negroni.Wrap(http.HandlerFunc(addGoalStatus)),
	)).Methods(http.MethodPost)

	api.NewRoute().PathPrefix("/goals/metrics").Handler(negroni.New(
		negroni.HandlerFunc(authMW),
		negroni.Wrap(http.HandlerFunc(getGoalmetrics)),
	)).Methods(http.MethodGet)
	
	// match everything else...
	api.NewRoute().PathPrefix("/").HandlerFunc(defaultHandler)

	log.Fatalln(http.ListenAndServe(":8080", r))

}

func main() {

	createRouter()
}