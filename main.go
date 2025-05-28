package main

import (
	"database/sql"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
)

type Todo struct {
	ID        int    `json:"id"`
	Title     string `json:"title"`
	Completed bool   `json:"completed"`
}

// DB connection
var db *sql.DB

func main() {
	initDB()
	defer db.Close()

	router := gin.Default()

	// Routes
	api := router.Group("/api/todos")
	{
		api.GET("", getTodos)
		api.GET("/:id", getTodo)
		api.POST("", createTodo)
		api.PUT("/:id", updateTodo)
		api.DELETE("/:id", deleteTodo)
	}

	log.Println("http://localhost:8080")

	router.Run(":8080")
}

func initDB() {
	var err error

	db, err = sql.Open("sqlite3", "./todos.db")

	if err != nil {
		log.Fatal("Failed to connect to db:", err)
	}

	// Create Table
	createTableSQL := `
	CREATE TABLE IF NOT EXISTS todos (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		title TEXT NOT NULL,
		completed BOOLEAN NOT NULL DEFAULT 0
	);
	`

	_, err = db.Exec(createTableSQL)

	if err != nil {
		log.Fatal("Failed to create table:", err)
	}

	// Sample data
	// TODO: created every time
	_, err = db.Exec("INSERT INTO todos (title, completed) VALUES (?, ?)", "Hello World", false)

	if err != nil {
		log.Fatal("Failed to insert data:", err)
	}

	log.Println("Database created")
}

/*
 * APIs
 */
func getTodos(context *gin.Context) {
	// var parameter
	rows, err := db.Query("SELECT id, title, completed FROM todos")

	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch todos"})
		return
	}
	defer rows.Close()

	var todos []Todo

	for rows.Next() {
		var todo Todo

		err := rows.Scan(&todo.ID, &todo.Title, &todo.Completed)
		if err != nil {
			context.JSON(http.StatusInternalServerError, gin.H{"error": "Error reading todo data"})
			return
		}

		todos = append(todos, todo)
	}

	if err := rows.Err(); err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": "Error iterating todos"})
		return
	}

	context.JSON(http.StatusOK, todos)
}

func getTodo(context *gin.Context) {
	idStr := context.Param("id")
	// casts string to int
	id, err := strconv.Atoi(idStr)

	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	var todo Todo
	err = db.QueryRow("SELECT id, title, completed FROM todos WHERE id = ?", id).Scan(&todo.ID, &todo.Title, &todo.Completed)

	if err != nil {
		if err == sql.ErrNoRows {
			context.JSON(http.StatusNotFound, gin.H{"error": "Todo not found"})
		} else {
			context.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch todo"})
		}

		return
	}

	context.JSON(http.StatusOK, todo)
}

func createTodo(context *gin.Context) {
	var newTodo Todo
	if err := context.ShouldBindJSON(&newTodo); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})

		return
	}

	result, err := db.Exec("INSERT INTO todos (title, completed) VALUES (?, ?)", newTodo.Title, newTodo.Completed)

	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create todo"})

		return
	}

	id, err := result.LastInsertId()
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get inserted ID"})

		return
	}

	newTodo.ID = int(id)
	context.JSON(http.StatusCreated, newTodo)
}

func updateTodo(context *gin.Context) {
	idStr := context.Param("id")
	id, err := strconv.Atoi(idStr)

	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": "Invalid id"})
		return
	}

	var updatedTodo Todo
	if err := context.ShouldBindJSON(&updatedTodo); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var existingTodo Todo
	err = db.QueryRow("SELECT id FROM todos WHERE id = ?", id).Scan(&existingTodo.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			context.JSON(http.StatusNotFound, gin.H{"error": "Not found"})
		} else {
			context.JSON(http.StatusInternalServerError, gin.H{"error": ""})
		}
		return
	}

	_, err = db.Exec("UPDATE todos SET title = ?, completed = ? WHERE id = ?",
		updatedTodo.Title, updatedTodo.Completed, id)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update"})
		return
	}

	updatedTodo.ID = id
	context.JSON(http.StatusOK, updatedTodo)
}

// TODO
func deleteTodo(context *gin.Context) {}
