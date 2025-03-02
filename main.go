package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"gorm.io/driver/postgres"
)

type Book struct {
	ID     string `json:"id" binding:"required"`
	Title  string `json:"title" binding:"required"`
	Author string `json:"author" binding:"required"`
}

type Handler struct {
	db *gorm.DB
}

func newHandler(db *gorm.DB) *Handler {
	return &Handler{db}
}

var books = []Book{
	{ID: "1", Title: "Harry Potter", Author: "J. K. Rowling"},
	{ID: "2", Title: "The Lord of the Rings", Author: "J. R. R. Tolkien"},
	{ID: "3", Title: "The Wizard of Oz", Author: "L. Frank Baum"},
}

func (h *Handler) listBooksHandler(c *gin.Context) {
	var books []Book

	if result := h.db.Find(&books); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": result.Error.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, &books)
}

func pingHandler(c *gin.Context) {
	c.String(http.StatusOK, "pong :)")
}

func (h *Handler) createBookHandler(c *gin.Context) {
	var book Book
	if err := c.ShouldBindJSON(&book); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	
	if result := h.db.Create(&book); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": result.Error.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, &book)
}

func (h *Handler) deleteBookHandler(c *gin.Context) {
	id := c.Param("id")

	if result := h.db.Delete(&Book{}, id); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": result.Error.Error(),
		})
		return
	}
	c.Status(http.StatusNoContent)
}

var db *gorm.DB

func main() {
	var err error
	
	db, err = gorm.Open(postgres.New(postgres.Config{
		DSN: "host=localhost user=testuser password=testpass dbname=testdb port=5432 sslmode=disable",
		PreferSimpleProtocol: true,
	}), &gorm.Config{})

	if err != nil {
		panic("failed to connect to database")
	}

	db.AutoMigrate(&Book{})

	handler := newHandler(db)

	r := gin.New()

	r.GET("/books", handler.listBooksHandler)
	r.GET("/ping", pingHandler)
	r.POST("/books", handler.createBookHandler)
	r.DELETE("/books/:id", handler.deleteBookHandler)
	
	r.Run()
}
