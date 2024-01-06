package main

import (
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
)

type book struct {
	ID       string `json:"id"`
	Title    string `json:"title"`
	Author   string `json:"author"`
	Quantity int    `json:"quantity"`
}

var books = []book{
	{ID: "1", Title: "In Search of Lost Time", Author: "Marcel Proust", Quantity: 2},
	{ID: "2", Title: "The Great Gatsby", Author: "F. Scott Fitzgerald", Quantity: 5},
	{ID: "3", Title: "War and Peace", Author: "Leo Tolstoy", Quantity: 6},
}

func getBooks(context *gin.Context) {
	context.JSON(http.StatusOK, books)
}

func createBook(context *gin.Context) {
	var newBook book
	if err := context.BindJSON(&newBook); err != nil {
		return
	}
	books = append(books, newBook)
	context.IndentedJSON(http.StatusCreated, newBook)
}

func bookById(context *gin.Context) {
	id := context.Param("id")
	book, err := getBookById(id)

	if err != nil {
		context.IndentedJSON(http.StatusNotFound, gin.H{"message": "book not found"})
		return
	}
	context.IndentedJSON(http.StatusOK, book)
}

func getBookById(id string) (*book, error) {
	for i, b := range books {
		if b.ID == id {
			return &books[i], nil
		}
	}
	return nil, errors.New("book not found")
}

func checkoutBook(context *gin.Context) {
	id, ok := context.GetQuery("id")

	if ok == false {
		context.IndentedJSON(http.StatusBadRequest, gin.H{"message": "id is required"})
		return
	}

	book, err := getBookById(id)

	if err != nil {
		context.IndentedJSON(http.StatusNotFound, gin.H{"message": "book not found"})
		return
	}

	if book.Quantity <= 0 {
		context.IndentedJSON(http.StatusNotFound, gin.H{"message": "book not available"})
		return
	}

	book.Quantity = book.Quantity - 1
	context.IndentedJSON(http.StatusOK, book)
}

func returnBook(context *gin.Context) {
	id, ok := context.GetQuery("id")

	if ok == false {
		context.IndentedJSON(http.StatusBadRequest, gin.H{"message": "id is required"})
		return
	}

	book, err := getBookById(id)

	if err != nil {
		context.IndentedJSON(http.StatusNotFound, gin.H{"message": "book not found"})
		return
	}

	book.Quantity = book.Quantity + 1
	context.IndentedJSON(http.StatusOK, book)

}

func main() {
	router := gin.Default()
	router.GET("/books", getBooks)
	router.POST("/books", createBook)
	router.GET("/books/:id", bookById)
	router.PATCH("/books/checkout", checkoutBook)
	router.PATCH("/books/return", returnBook)
	router.Run("localhost:8080")
}
