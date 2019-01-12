package handlers

import (
	"encoding/json"
	"net/http"
	"src/internal/data"

	"github.com/gorilla/mux"
)

func (h MyHandler) CreateTask(w http.ResponseWriter, r *http.Request) {
	// Parse the request body to extract task data
	var newTask data.Task
	err := json.NewDecoder(r.Body).Decode(&newTask)
	if err != nil {
		// Handle parsing errors and return a response
		http.Error(w, "Invalid request data", http.StatusBadRequest)
		return
	}

	// Implement logic to create the task in your data source (PostgreSQL)
	// Assuming you have a database connection h.DB

	// Prepare the SQL statement to insert a new task
	insertQuery := `
        INSERT INTO tasks (title, description, due_date, is_completed)
        VALUES ($1, $2, $3, $4)
        RETURNING id, title, description, due_date, is_completed
    `

	// Execute the SQL query and retrieve the ID of the newly created task
	var task data.Task
	err = h.DB.QueryRow(insertQuery, newTask.Title, newTask.Description, newTask.DueDate, newTask.IsCompleted).Scan(
		&task.ID, &task.Title, &task.Description, &task.DueDate, &task.IsCompleted)
	if err != nil {
		http.Error(w, "Failed to create task", http.StatusInternalServerError)
		return
	}

	// Set the generated task ID

	// Return a success response with the created task
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(task)
}

func (h MyHandler) GetAllTasks(w http.ResponseWriter, r *http.Request) {
	// Get a database connection
	db := h.DB

	// Implement logic to fetch all tasks from your data source (PostgreSQL database)
	// Replace "tasks" with the logic to retrieve tasks from your database
	tasks := []data.Task{} // Create an empty slice to hold the tasks

	// Query the database to retrieve tasks
	rows, err := db.Query("SELECT id, title, description, due_date, is_completed FROM tasks")
	if err != nil {
		http.Error(w, "Failed to fetch tasks from the database", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var task data.Task
		if err := rows.Scan(&task.ID, &task.Title, &task.Description, &task.DueDate, &task.IsCompleted); err != nil {
			http.Error(w, "Failed to scan task row", http.StatusInternalServerError)
			return
		}
		tasks = append(tasks, task)
	}
	if err := rows.Err(); err != nil {
		http.Error(w, "Failed to iterate over task rows", http.StatusInternalServerError)
		return
	}

	// Encode tasks as JSON and write the response
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(tasks); err != nil {
		http.Error(w, "Failed to encode tasks", http.StatusInternalServerError)
		return
	}
}

func (h MyHandler) UpdateTask(w http.ResponseWriter, r *http.Request) {
	// Extract the task ID from the URL
	vars := mux.Vars(r)
	taskID := vars["id"]
	// Parse the request body to extract task data for updating
	var updatedTask data.Task
	err := json.NewDecoder(r.Body).Decode(&updatedTask)
	if err != nil {
		// Handle parsing errors and return a response
		http.Error(w, "Invalid request data", http.StatusBadRequest)
		return
	}

	// Implement logic to update the task in your data source (PostgreSQL)
	// Assuming you have a database connection h.DB

	updateQuery := `
        UPDATE tasks
        SET title = $1, description = $2, due_date = $3, is_completed = $4
        WHERE id = $5
        RETURNING id, title, description, due_date, is_completed
    `

	var updated data.Task
	err = h.DB.QueryRow(updateQuery, updatedTask.Title, updatedTask.Description, updatedTask.DueDate, updatedTask.IsCompleted, taskID).Scan(
		&updated.ID, &updated.Title, &updated.Description, &updated.DueDate, &updated.IsCompleted,
	)
	if err != nil {
		http.Error(w, "Failed to update task", http.StatusInternalServerError)
		return
	}

	// Return the updated task as JSON in the response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updated)
}

func (h MyHandler) DeleteTask(w http.ResponseWriter, r *http.Request) {
	// Extract the task ID from the URL
	vars := mux.Vars(r)
	taskID := vars["id"]

	// Implement logic to delete the task with the given ID from your data source (PostgreSQL)
	// Assuming you have a database connection h.DB

	deleteQuery := `
        DELETE FROM tasks
        WHERE id = $1
    `

	_, err := h.DB.Exec(deleteQuery, taskID)
	if err != nil {
		http.Error(w, "Failed to delete task", http.StatusInternalServerError)
		return
	}

	// Return a success response with a status message
	w.WriteHeader(http.StatusNoContent)
}
