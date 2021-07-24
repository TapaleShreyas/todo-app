package router

import (
	"go-server/middleware"

	"github.com/gorilla/mux"
)

func Router() *mux.Router {
	router := mux.NewRouter()
	router.HandleFunc("/api/task/{id}", middleware.GetTask).Methods("GET", "OPTIONS")
	router.HandleFunc("/api/task", middleware.GetAllTask).Methods("GET", "OPTIONS")
	router.HandleFunc("/api/task", middleware.CreateTask).Methods("POST", "OPTIONS")
	router.HandleFunc("/api/task/delete/{id}", middleware.DeleteTask).Methods("DELETE", "OPTIONS")
	router.HandleFunc("/api/task/delete/all/", middleware.DeleteAllTask).Methods("DELETE", "OPTIONS")
	router.HandleFunc("/api/task/complete/{id}", middleware.CompleteTask).Methods("PUT", "OPTIONS")
	router.HandleFunc("/api/task/undo/{id}", middleware.UndoTask).Methods("PUT", "OPTIONS")
	return router
}
