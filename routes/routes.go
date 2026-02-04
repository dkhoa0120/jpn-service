package routes

import (
	"jpn-service/handlers"

	"github.com/gorilla/mux"
)

func SetupRoutes() *mux.Router {
	router := mux.NewRouter()

	// User routes
	router.HandleFunc("/api/users", handlers.GetAllUsers).Methods("GET")
	router.HandleFunc("/api/users/{id}", handlers.GetUser).Methods("GET")
	router.HandleFunc("/api/users", handlers.CreateUser).Methods("POST")
	router.HandleFunc("/api/users/{id}", handlers.UpdateUser).Methods("PUT")
	router.HandleFunc("/api/users/{id}", handlers.DeleteUser).Methods("DELETE")

	// vocabulary

	router.HandleFunc("/api/vocabularies", handlers.GetVocabularies).Methods("GET")
	router.HandleFunc("/api/vocabularies", handlers.CreateVocabulary).Methods("POST")

	return router
}
