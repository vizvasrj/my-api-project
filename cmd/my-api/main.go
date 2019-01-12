package main

import (
	"net/http"
	"src/config"
	"src/internal/api/handlers"
	"src/internal/api/middelware"
	"src/internal/data"

	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter()
	db := data.GetDatabase()
	config := config.Conf{
		DB: db,
	}
	h := handlers.MyHandler{
		DB: config.DB,
	}
	r.HandleFunc("/register", h.RegisterUser).Methods("POST")
	r.HandleFunc("/login", h.LoginUser).Methods("POST")

	r.HandleFunc("/tasks", middelware.AuthMiddleware(h.GetAllTasks)).Methods("GET")

	// Create a new task
	r.HandleFunc("/tasks", middelware.AuthMiddleware(h.CreateTask)).Methods("POST")

	// Update an existing task by ID
	r.HandleFunc("/tasks/{id}", middelware.AuthMiddleware(h.UpdateTask)).Methods("PUT")

	// Delete an existing task by ID
	r.HandleFunc("/tasks/{id}", middelware.AuthMiddleware(h.DeleteTask)).Methods("DELETE")

	http.ListenAndServe(":8080", r)
}

// 	// Encode the task as JSON and write the response
// 	w.Header().Set("Content-Type", "application/json")
// 	json.NewEncoder(w).Encode(task)
// }

// 	// Implement logic to update the task with the given ID in your data source
// 	// ...

// 	// Return a success response with the updated task
// 	w.WriteHeader(http.StatusOK)
// 	json.NewEncoder(w).Encode(updatedTask)
// }
