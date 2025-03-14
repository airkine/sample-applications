package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v4/pgxpool"
)

type Case struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}

func getCases(pool *pgxpool.Pool) ([]Case, error) {
	rows, err := pool.Query(context.Background(), "SELECT title, description FROM cases ORDER BY RANDOM() LIMIT 3")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var cases []Case
	for rows.Next() {
		var c Case
		if err := rows.Scan(&c.Title, &c.Description); err != nil {
			return nil, err
		}
		cases = append(cases, c)
	}

	return cases, nil
}

func connectDB() (*pgxpool.Pool, error) {
	connString := os.Getenv("DATABASE_URL")
	pool, err := pgxpool.Connect(context.Background(), connString)
	if err != nil {
		return nil, err
	}
	return pool, nil
}

func main() {
	pool, err := connectDB()
	if err != nil {
		log.Fatalf("Database connection error: %v", err)
	}
	defer pool.Close()

	r := gin.Default()

	// Define Static and Templates ONCE only
	r.Static("/static", "./static")
	r.LoadHTMLGlob("templates/*.html")

	r.GET("/api/truth", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "The truth is out there..."})
	})

	r.GET("/api/cases", func(c *gin.Context) {
		cases, err := getCases(pool)
		if err != nil {
			log.Printf("Error retrieving cases: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch cases"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"cases": cases})
	})

	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", gin.H{"title": "X-Files Database"})
	})

	r.Run(":8080")
}
