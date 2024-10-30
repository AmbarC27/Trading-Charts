package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
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

var jwtSecret = []byte("sample-secret-key")

// Helper function to generate JWT
func GenerateJWT(username string) (string, error) {
	// Create the token with claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": username,
		"exp":      time.Now().Add(time.Hour * 2).Unix(), // Token expires in 2 hours
	})

	// Sign the token with our secret key
	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func login(c *gin.Context) {
	var user struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	// Parse the JSON body into the user struct
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	if user.Username == "admin" && user.Password == "password" {
		// Generate JWT token
		token, err := GenerateJWT(user.Username)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not generate token"})
			return
		}

		// Return the token
		c.JSON(http.StatusOK, gin.H{"token": token})
	} else {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
	}
}

// Middleware to protect routes and validate JWT
func JWTAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get the Authorization header (format: "Bearer <token>")
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization token required"})
			c.Abort() // Stop further processing
			return
		}

		// Token usually comes as "Bearer <token>"
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid Authorization header format"})
			c.Abort()
			return
		}

		tokenString := parts[1]

		// Parse and validate the token
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// Ensure the token's signing method is HMAC (HS256)
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return jwtSecret, nil // Use the secret to validate the token
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			c.Abort()
			return
		}

		// Token is valid, extract claims (like username)
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
			c.Abort()
			return
		}

		// Save the username from the token into the context for use in other routes
		c.Set("username", claims["username"])

		// Proceed to the next middleware/handler
		c.Next()
	}
}

func main() {
	initDB()
	defer db.Close()

	r := gin.Default()
	r.GET("/stocks", getAllStocks)
	r.GET("/latest", getLatestTickerData)
	r.GET("/stocks/:ticker", getStockDataByTicker)
	r.POST("/login", login)

	// Protected route - only accessible with a valid JWT
	r.GET("/dashboard", JWTAuthMiddleware(), func(c *gin.Context) {
		// Access the username from the token (stored in context by the middleware)
		username := c.MustGet("username").(string)

		// Return a response with the username
		c.JSON(http.StatusOK, gin.H{
			"message":  "Welcome to your dashboard!",
			"username": username,
		})
	})
	r.Run(":8080")
}
