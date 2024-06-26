package routes

import (
	"github.com/gorilla/mux"
	"github.com/temuka-api-service/handlers"
)

func PostRoutes() *mux.Router {
	r := mux.NewRouter()

	r.HandleFunc("/", handlers.CreatePost).Methods("POST")
	r.HandleFunc("/timeline", handlers.GetTimelinePosts).Methods("GET")
	r.HandleFunc("/like/{id}", handlers.LikePost).Methods("PUT")
	r.HandleFunc("/{id}", handlers.DeletePost).Methods("DELETE")

	return r
}
