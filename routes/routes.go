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

	// Daily routes
	router.HandleFunc("/api/dailies", handlers.GetAllDailies).Methods("GET")
	router.HandleFunc("/api/dailies/{id}", handlers.GetDaily).Methods("GET")
	router.HandleFunc("/api/dailies", handlers.CreateDaily).Methods("POST")
	router.HandleFunc("/api/dailies/{id}", handlers.UpdateDaily).Methods("PUT")
	router.HandleFunc("/api/dailies/{id}", handlers.DeleteDaily).Methods("DELETE")

	// vocabulary
	router.HandleFunc("/api/vocabularies", handlers.GetVocabularies).Methods("GET")
	router.HandleFunc("/api/vocabularies", handlers.CreateVocabulary).Methods("POST")
	router.HandleFunc("/api/vocabularies/{id}", handlers.UpdateVocabulary).Methods("PUT")

	// Grammar routes
	router.HandleFunc("/api/grammars", handlers.CreateGrammar).Methods("POST")
	router.HandleFunc("/api/grammars", handlers.GetGrammars).Methods("GET")
	router.HandleFunc("/api/grammars/search", handlers.SearchGrammars).Methods("GET")
	router.HandleFunc("/api/grammars/level/{level}", handlers.GetGrammarsByLevel).Methods("GET")
	router.HandleFunc("/api/grammars/{id}", handlers.GetGrammarByID).Methods("GET")
	router.HandleFunc("/api/grammars/{id}", handlers.UpdateGrammar).Methods("PUT")
	router.HandleFunc("/api/grammars/{id}", handlers.DeleteGrammar).Methods("DELETE")

	return router
}
