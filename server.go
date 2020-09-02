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

	// mysqlHost = os.Getenv("MYSQL_HOST")
	// mysqlDatabase = os.Getenv("MYSQL_DATABASE")
	// mysqlUsername = os.Getenv("MYSQL_USER")
	// mysqlPassword = os.Getenv("MYSQL_PASSWORD")

	// // example: "root:password1@tcp(127.0.0.1:3306)/test"
	// mysqlDSN = fmt.Sprintf("%s:%s@tcp(%s)/%s?parseTime=true", mysqlUsername, mysqlPassword, mysqlHost, mysqlDatabase)

	// log.Printf("mysqlHost: %s", mysqlHost)
	// log.Printf("mysqlDatabase: %s", mysqlDatabase)
	// log.Printf("mysqlUsername: %s", mysqlUsername)
	// log.Printf("mysqlDSN: %s", mysqlDSN)
	// count := 0
	// for {
	// 	count++
	// 	db, err := sql.Open("mysql", mysqlDSN)
	// 	if err != nil {
	// 		if count > 1000 {
	// 			log.Fatalf("(sql.Open) Unable to connect to database: %v", err)
	// 			return
	// 		}
	// 		log.Printf("(sql.Open) Unable to connect to database: %v, retrying...", err)
	// 		time.Sleep(time.Second * 3)
	// 		continue
	// 	}
	// 	//err = db.Ping()
	// 	for {
	// 		_, err = db.Exec("DO 1")
	// 		if err != nil {
	// 			if count > 1000 {
	// 				log.Fatalf("(db.Ping) Unable to connect to database: %v", err)
	// 				return
	// 			}
	// 			log.Printf("(db.Ping) Unable to connect to database: %v, retrying...", err)
	// 			time.Sleep(time.Second * 3)
	// 			continue
	// 		}
	// 		break
	// 	}
	// 	break
	// }

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