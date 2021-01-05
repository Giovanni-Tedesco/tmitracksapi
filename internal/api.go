package api

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/Giovanni-Tedesco/tmitracksapi/internal/auth"
	"github.com/Giovanni-Tedesco/tmitracksapi/internal/mechanic"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/mongo"
)

type App struct {
	router *mux.Router
	// Implement DB
	DB *mongo.Database
}

type ErrorResponse struct {
	StatusCode   int    `json:"status"`
	ErrorMessage string `json:"message"`
}

func ConnectDB(client *mongo.Client) *mongo.Database {
	database := client.Database("TMI")

	return database
}

type RequestHandlerFunc func(db *mongo.Database, w http.ResponseWriter, r *http.Request)

func SubRouterRequests(handler RequestHandlerFunc, db *mongo.Database) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		handler(db, w, r)
	}
}

func (a *App) Initialize(client *mongo.Client) {

	a.DB = ConnectDB(client)

	a.router = mux.NewRouter()
	a.setRouters()

	authRoute := a.router.PathPrefix("/v1/auth").Subrouter()
	authRoute.Use(auth.AuthMiddleWare)
	authRoute.HandleFunc("/test", SubRouterRequests(auth.TestSomething, a.DB)).Methods("GET")
	authRoute.HandleFunc("/signup", SubRouterRequests(auth.SignUp, a.DB)).Methods("POST")
	authRoute.HandleFunc("/signin", SubRouterRequests(auth.SignIn, a.DB)).Methods("POST")
}

func (a *App) setRouters() {
	// a.Get("/get_hello", a.handleRequest(bins.Hello))
	a.Post("/create_report", a.handleRequest(mechanic.CreateReport))
	a.Get("/get_report_by_date", a.handleRequest(mechanic.GetReportByDate))
	a.Get("/get_reports", a.handleRequest(mechanic.GetAllReports))
	a.Get("/get_reports_range", a.handleRequest(mechanic.GetReportByDateRange))
	a.Get("/get_report", a.handleRequest(mechanic.GetReportById))
	a.Delete("/delete_report", a.handleRequest(mechanic.DeleteReport))
}

// Wrappers for GET, POST, PUT, and DELETE
func (a *App) Get(path string, f func(w http.ResponseWriter, r *http.Request)) {
	a.router.HandleFunc(path, f).Methods("GET")
}

func (a *App) Post(path string, f func(w http.ResponseWriter, r *http.Request)) {
	a.router.HandleFunc(path, f).Methods("POST")
}

func (a *App) Put(path string, f func(w http.ResponseWriter, r *http.Request)) {
	a.router.HandleFunc(path, f).Methods("PUT")
}

func (a *App) Delete(path string, f func(w http.ResponseWriter, r *http.Request)) {
	a.router.HandleFunc(path, f).Methods("DELETE")
}

func (a *App) Run(host string) {
	log.Fatal(http.ListenAndServe(host, a.router))
}

func (a *App) handleRequest(handler RequestHandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		handler(a.DB, w, r)
	}
}

func GetError(err error, w http.ResponseWriter) {
	log.Fatal(err.Error())

	var response = ErrorResponse{
		ErrorMessage: err.Error(),
		StatusCode:   http.StatusInternalServerError,
	}

	message, _ := json.Marshal(response)

	w.WriteHeader(response.StatusCode)
	w.Write(message)
}
