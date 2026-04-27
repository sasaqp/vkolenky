package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB

func main() {
	initDB()

	r := gin.Default()

	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	api := r.Group("/api")
	{
		api.POST("/auth", handleAuth)
		api.GET("/users/me", handleGetMe)
		api.PUT("/users/me", handleUpdateMe)

		api.GET("/motorcycles", handleGetMotorcycles)
		api.POST("/motorcycles", handleCreateMotorcycle)
		api.PUT("/motorcycles/:id", handleUpdateMotorcycle)
		api.DELETE("/motorcycles/:id", handleDeleteMotorcycle)

		api.GET("/races", handleGetRaces)
		api.POST("/races", handleCreateRace)
		api.GET("/races/:id/laps", handleGetLaps)
		api.POST("/laps", handleCreateLap)
		api.POST("/laps/manual", handleManualLap)
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("Server starting on port %s", port)
	r.Run(":" + port)
}

func initDB() {
	db, err := sql.Open("sqlite3", "./data/laptiming.db")
	if err != nil {
		log.Fatal(err)
	}
	DB = db

	DB.Exec(`
		CREATE TABLE IF NOT EXISTS users (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			telegram_id TEXT UNIQUE NOT NULL,
			name TEXT NOT NULL,
			role TEXT DEFAULT 'rider',
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP
		);
		CREATE TABLE IF NOT EXISTS motorcycles (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			user_id INTEGER NOT NULL,
			brand TEXT NOT NULL,
			model TEXT NOT NULL,
			number TEXT NOT NULL,
			cubature INTEGER,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (user_id) REFERENCES users(id)
		);
		CREATE TABLE IF NOT EXISTS races (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL,
			date DATETIME,
			status TEXT DEFAULT 'pending',
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP
		);
		CREATE TABLE IF NOT EXISTS laps (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			race_id INTEGER NOT NULL,
			motorcycle_id INTEGER NOT NULL,
			lap_number INTEGER DEFAULT 1,
			time REAL NOT NULL,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (race_id) REFERENCES races(id),
			FOREIGN KEY (motorcycle_id) REFERENCES motorcycles(id)
		);
	`)
}

func handleAuth(c *gin.Context) {
	var req struct {
		InitData string `json:"initData"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user struct {
		ID          int
		TelegramID string
		Name        string
		Role        string
	}

	err := DB.QueryRow(`
		INSERT INTO users (telegram_id, name) VALUES (?, 'User')
		ON CONFLICT(telegram_id) DO UPDATE SET name=excluded.name
		RETURNING id, telegram_id, name, role
	`, req.InitData).Scan(&user.ID, &user.TelegramID, &user.Name, &user.Role)

	if err != nil {
		log.Printf("Auth error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "auth failed"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"user": user})
}

func handleGetMe(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"id": 1, "name": "User", "role": "rider"})
}

func handleUpdateMe(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "updated"})
}

func handleGetMotorcycles(c *gin.Context) {
	rows, err := DB.Query("SELECT id, user_id, brand, model, number, cubature FROM motorcycles")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	var motorcycles []map[string]interface{}
	for rows.Next() {
		var m map[string]interface{}
		rows.Scan(&m["id"], &m["user_id"], &m["brand"], &m["model"], &m["number"], &m["cubature"])
		motorcycles = append(motorcycles, m)
	}
	c.JSON(http.StatusOK, motorcycles)
}

func handleCreateMotorcycle(c *gin.Context) {
	var req struct {
		Brand     string `json:"brand"`
		Model    string `json:"model"`
		Number   string `json:"number"`
		Cubature int    `json:"cubature"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	_, err := DB.Exec(`
		INSERT INTO motorcycles (user_id, brand, model, number, cubature)
		VALUES (1, ?, ?, ?, ?)
	`, req.Brand, req.Model, req.Number, req.Cubature)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"status": "created"})
}

func handleUpdateMotorcycle(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "updated"})
}

func handleDeleteMotorcycle(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "deleted"})
}

func handleGetRaces(c *gin.Context) {
	c.JSON(http.StatusOK, []map[string]interface{}{})
}

func handleCreateRace(c *gin.Context) {
	var req struct {
		Name string `json:"name"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	_, err := DB.Exec("INSERT INTO races (name) VALUES (?)", req.Name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"status": "created"})
}

func handleGetLaps(c *gin.Context) {
	c.JSON(http.StatusOK, []map[string]interface{}{})
}

func handleCreateLap(c *gin.Context) {
	c.JSON(http.StatusCreated, gin.H{"status": "created"})
}

func handleManualLap(c *gin.Context) {
	var req struct {
		RaceID        int     `json:"raceId"`
		MotorcycleID int     `json:"motorcycleId"`
		LapNumber    int     `json:"lapNumber"`
		Time         float64 `json:"time"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	_, err := DB.Exec(`
		INSERT INTO laps (race_id, motorcycle_id, lap_number, time)
		VALUES (?, ?, ?, ?)
	`, req.RaceID, req.MotorcycleID, req.LapNumber, req.Time)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"status": "created"})
}