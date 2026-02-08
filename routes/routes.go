package routes

import (
	"jpn-service/handlers"
	"jpn-service/middleware"
	"net/http"

	"github.com/gorilla/mux"
)

func SetupRoutes() *mux.Router {
    router := mux.NewRouter()

    // Apply CORS middleware globally
    router.Use(middleware.CORS)

    // Health check
    router.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Content-Type", "application/json")
        w.WriteHeader(http.StatusOK)
        w.Write([]byte(`{"status": "ok"}`))
    }).Methods("GET", "OPTIONS")

    // Daily routes - THÊM OPTIONS VÀO MỖI ROUTE
    router.HandleFunc("/api/dailies", handlers.GetAllDailies).Methods("GET", "OPTIONS")
    router.HandleFunc("/api/dailies/{id}", handlers.GetDaily).Methods("GET", "OPTIONS")
    router.HandleFunc("/api/dailies", handlers.CreateDaily).Methods("POST", "OPTIONS")
    router.HandleFunc("/api/dailies/{id}", handlers.UpdateDaily).Methods("PUT", "OPTIONS")
    router.HandleFunc("/api/dailies/{id}", handlers.DeleteDaily).Methods("DELETE", "OPTIONS")

    // User routes
    router.HandleFunc("/api/users", handlers.GetAllUsers).Methods("GET", "OPTIONS")
    router.HandleFunc("/api/users/{id}", handlers.GetUser).Methods("GET", "OPTIONS")
    router.HandleFunc("/api/users", handlers.CreateUser).Methods("POST", "OPTIONS")
    router.HandleFunc("/api/users/{id}", handlers.UpdateUser).Methods("PUT", "OPTIONS")
    router.HandleFunc("/api/users/{id}", handlers.DeleteUser).Methods("DELETE", "OPTIONS")

    // Vocabulary routes
    router.HandleFunc("/api/vocabularies", handlers.GetVocabularies).Methods("GET", "OPTIONS")
    router.HandleFunc("/api/vocabularies", handlers.CreateVocabulary).Methods("POST", "OPTIONS")
    router.HandleFunc("/api/vocabularies/{id}", handlers.UpdateVocabulary).Methods("PUT", "OPTIONS")

    // Grammar routes
    router.HandleFunc("/api/grammars", handlers.CreateGrammar).Methods("POST", "OPTIONS")
    router.HandleFunc("/api/grammars", handlers.GetGrammars).Methods("GET", "OPTIONS")
    router.HandleFunc("/api/grammars/search", handlers.SearchGrammars).Methods("GET", "OPTIONS")
    router.HandleFunc("/api/grammars/level/{level}", handlers.GetGrammarsByLevel).Methods("GET", "OPTIONS")
    router.HandleFunc("/api/grammars/{id}", handlers.GetGrammarByID).Methods("GET", "OPTIONS")
    router.HandleFunc("/api/grammars/{id}", handlers.UpdateGrammar).Methods("PUT", "OPTIONS")
    router.HandleFunc("/api/grammars/{id}", handlers.DeleteGrammar).Methods("DELETE", "OPTIONS")

    return router
}