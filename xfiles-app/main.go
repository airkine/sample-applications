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

func main() {
	router := gin.Default()

	// ✅ Serve static files
	router.Static("/static", "./static")

	// ✅ Load HTML templates from ./templates directory
	router.LoadHTMLGlob("templates/*")

	// ✅ Serve the HTML page at "/"
	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", nil)
	})

	// ✅ API returns JSON for cases
	router.GET("/api/cases", func(c *gin.Context) {
		dbURL := os.Getenv("DATABASE_URL")
		pool, err := pgxpool.Connect(context.Background(), dbURL)
		if err != nil {
			log.Fatal(err)
		}
		defer pool.Close()

		cases, err := getCases(pool)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch cases"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"cases": cases})
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	router.Run(":" + port)
}
