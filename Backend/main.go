package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type Stock struct {
	Datetime    time.Time `json:"datetime"`
	Open        float64   `json:"open"`
	High        float64   `json:"high"`
	Low         float64   `json:"low"`
	Close       float64   `json:"close"`
	Volume      int64     `json:"volume"`
	Dividends   float64   `json:"dividends"`
	StockSplits float64   `json:"stock_splits"`
	Ticker      string    `json:"ticker"`
}

var db *sql.DB

func initDB() {
	// Load the environment variables
	err := godotenv.Load("../Data_Store/.env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")

	connStr := fmt.Sprintf("user=%s password=%s dbname=%s host=%s port=%s sslmode=disable",
		user, password, dbname, host, port)

	// Open a connection to the database
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
}

func getAllStocks(c *gin.Context) {
	rows, err := db.Query("SELECT * FROM stocks")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error retrieving data"})
		return
	}
	defer rows.Close()

	var stocks []Stock

	for rows.Next() {
		var stock Stock
		err := rows.Scan(
			&stock.Datetime,
			&stock.Open,
			&stock.High,
			&stock.Low,
			&stock.Close,
			&stock.Volume,
			&stock.Dividends,
			&stock.StockSplits,
			&stock.Ticker,
		)

		// Check if there was an error during the scan
		if err != nil {
			// Return an error response if scanning failed
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error scanning data"})
			return
		}
		stocks = append(stocks, stock)
	}

	if err = rows.Err(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error during row iteration"})
		return
	}

	// Respond with the retrieved stocks as JSON
	c.JSON(http.StatusOK, stocks)
}

func getLatestTickerData(c *gin.Context) {
	query := `
        SELECT DISTINCT ON (ticker) *
        FROM stocks
        ORDER BY ticker, Datetime DESC;
    `

	rows, err := db.Query(query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error retrieving data"})
		return
	}
	defer rows.Close()

	var stockData []Stock

	for rows.Next() {
		var stock Stock
		err := rows.Scan(
			&stock.Datetime,
			&stock.Open,
			&stock.High,
			&stock.Low,
			&stock.Close,
			&stock.Volume,
			&stock.Dividends,
			&stock.StockSplits,
			&stock.Ticker,
		)

		// Check if there was an error during the scan
		if err != nil {
			// Return an error response if scanning failed
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error scanning data"})
			return
		}
		stockData = append(stockData, stock)
	}

	if err = rows.Err(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error during row iteration"})
		return
	}

	// Return the result as a JSON array
	c.JSON(http.StatusOK, stockData)
}

func getStockDataByTicker(c *gin.Context) {
	// Get the ticker from the URL parameter
	ticker := c.Param("ticker")

	query := `
        SELECT *
        FROM stocks
        WHERE ticker = $1;
    `

	rows, err := db.Query(query, ticker)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error retrieving data"})
		return
	}
	defer rows.Close()

	var stockData []Stock

	// Iterate over the rows and scan them into the Stock struct
	for rows.Next() {
		var stock Stock
		err := rows.Scan(
			&stock.Datetime,
			&stock.Open,
			&stock.High,
			&stock.Low,
			&stock.Close,
			&stock.Volume,
			&stock.Dividends,
			&stock.StockSplits,
			&stock.Ticker,
		)

		// Check if there was an error during the scan
		if err != nil {
			// Return an error response if scanning failed
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error scanning data"})
			return
		}
		stockData = append(stockData, stock)
	}

	if err = rows.Err(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error during row iteration"})
		return
	}

	// Return the result as a JSON array
	c.JSON(http.StatusOK, stockData)
}

func main() {
	initDB()
	defer db.Close()

	r := gin.Default()
	r.GET("/stocks", getAllStocks)
	r.GET("/latest", getLatestTickerData)
	r.GET("/stocks/:ticker", getStockDataByTicker)
	r.Run(":8080")
}
