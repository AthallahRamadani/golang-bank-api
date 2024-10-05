package main

import (
	"fmt"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
)

// Customer represents a registered user
type Customer struct {
	Username string
	Balance  float64
}

// Global variables
var (
	customers     = make(map[string]*Customer)
	loggedInUsers = make(map[string]bool)
	mu            sync.Mutex
)

func main() {
	// Initialize some customers
	customers["alice"] = &Customer{Username: "alice", Balance: 1000}
	customers["bob"] = &Customer{Username: "bob", Balance: 500}

	r := gin.Default()

	// Login endpoint
	r.POST("/login", login)

	// Payment endpoint
	r.POST("/payment", payment)

	// Logout endpoint
	r.POST("/logout", logout)

	r.Run(":8080")
}

func login(c *gin.Context) {
	username := c.PostForm("username")

	mu.Lock()
	defer mu.Unlock()

	if customer, exists := customers[username]; exists {
		loggedInUsers[username] = true
		c.JSON(http.StatusOK, gin.H{"message": "Login successful", "customer": customer.Username})
	} else {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Customer does not exist"})
	}
}

func payment(c *gin.Context) {
	from := c.PostForm("from")
	to := c.PostForm("to")
	amount := c.PostForm("amount")

	mu.Lock()
	defer mu.Unlock()

	// Check if user is logged in
	if !loggedInUsers[from] {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "User not logged in"})
		return
	}

	// Check if both customers exist
	fromCustomer, fromExists := customers[from]
	toCustomer, toExists := customers[to]
	if !fromExists || !toExists {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid customer(s)"})
		return
	}

	// Parse amount
	transferAmount := 0.0
	_, err := fmt.Sscanf(amount, "%f", &transferAmount)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid amount"})
		return
	}

	// Check if sender has enough balance
	if fromCustomer.Balance < transferAmount {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Insufficient balance"})
		return
	}

	// Perform transfer
	fromCustomer.Balance -= transferAmount
	toCustomer.Balance += transferAmount

	c.JSON(http.StatusOK, gin.H{"message": "Payment successful"})
}

func logout(c *gin.Context) {
	username := c.PostForm("username")

	mu.Lock()
	defer mu.Unlock()

	if loggedInUsers[username] {
		delete(loggedInUsers, username)
		c.JSON(http.StatusOK, gin.H{"message": "Logout successful"})
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"message": "User not logged in"})
	}
}
