package data

import "time"

type Task struct {
	ID          string    `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	DueDate     time.Time `json:"due_date"`
	IsCompleted bool      `json:"is_completed"`
}

type CreateTaskRequest struct {
	Title       string `json:"title" validate:"required"`
	Description string `json:"description" validate:"max=200"`
}

type User struct {
	ID           string `json:"id"`
	Username     string `json:"username"`
	PasswordHash string `json:"password_hash"`
	Role         string `json:"role"`
}

type UserRegister struct {
	Username        string `json:"username"`
	Password        string `json:"password"`
	ConfirmPassword string `json:"confirm_password"`
}

type UserLogin struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// curl -L https://github.com/golang-migrate/migrate/releases/download/v4.0.0/migrate.linux-amd64.tar.gz | tar xvz | mv migrate.linux-amd64 ~/go/bin/migrate
