package main

import (
	"net/http"
	"github.com/gin-gonic/gin"
	"errors"
)

//1. CARA MELIHAT OUTPUT 
// -> go run main.go
// -> curl localhost:8080/books

//2. CARA RUN FILE JSON
// -> curl localhost:8080/books --include --header "Content-Type: application/json" -d @body.json --request "POST"

//3. CARA RUN CHECKOUT BUKU
// curl localhost:8080/checkout?id=2 --request "PATCH"

//4. CARA RUN MISSING ID/ID TANPA QUERY
// curl localhost:8080/checkout --request "PATCH"

//5. CARA RUN RETURN
// curl localhost:8080/return?id=2 --request "PATCH"

type book struct {
	ID       string `json:"id"`
	Title    string `json:"title"`
	Author   string `json:"author"`
	Quantity int    `json:"quantity"`
}

var books = []book{
	{ID: "1", Title: "In Search of Lost time", Author: "Marcel Proust", Quantity: 2},
	{ID: "2", Title: "The Great Gatsby", Author: "F. Scott Fitzgerald", Quantity: 5},
	{ID: "3", Title: "War and Peace", Author: "Leo Tolstoy", Quantity: 6},
}
// mengambil data 
func getBooks(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, books)
}
// fungsi mendapatkan buku berdasarkan id
func bookById(c *gin.Context) {
	id := c.Param("id") // akan mengakses semua parameter yang ada id nya
	book, err := getBookById(id)

	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "Book not found."})
		return
	}

	c.IndentedJSON(http.StatusOK, book)
}
// checkout books & query parameters
func checkoutBook(c *gin.Context) {
	id, ok := c.GetQuery("id")

	if !ok {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Missing id query parameter."})
		return
	}

	book, err := getBookById(id)

	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "Book not found."})
		return
	}

	if book.Quantity <= 0 {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "Book not available."})
		return
	}

	book.Quantity -= 1
	c.IndentedJSON(http.StatusOK, book)
}
// returning books
func returnBook(c *gin.Context) {
	id, ok := c.GetQuery("id")

	if !ok {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Missing id query parameter."})
		return
	}

	book, err := getBookById(id)

	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "Book not found."})
		return
	}

	book.Quantity += 1
	c.IndentedJSON(http.StatusOK, book)
}
// mengecek data berdasarkan id
func getBookById(id string) (*book, error) {
	for i, b := range books {
		if b.ID == id {
			return &books[i], nil
		}
	}

	return nil, errors.New("book not found")
}
// membuat data
func createBook(c *gin.Context) {
	var newBook book

	if err := c.BindJSON(&newBook); err != nil {
		return
	}

	books = append(books, newBook)
	c.IndentedJSON(http.StatusCreated, newBook)
}

func main() {
	router := gin.Default()
	router.GET("/books", getBooks) 
	router.GET("/books/:id", bookById)
	router.POST("/books", createBook)
	router.PATCH("/checkout", checkoutBook)
	router.PATCH("/return", returnBook)
	router.Run("localhost:8080")
}

// GET   -> Mendapatkan informasi / data
// POST  -> Menambahkan informasi / data 
// PATCH -> Update/memperbarui informasi/data/jumlah produk