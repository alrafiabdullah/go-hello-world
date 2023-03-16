package main

import (
	"fmt"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Todo struct {
	ID		uint	`gorm:"primary_key" json:"id"`
	Name 	string	`gorm:"type:varchar(512);not null" json:"name"`
	IsDone	bool	`gorm:"default:false" json:"is_done"`
	CreatedAt 	time.Time 	`gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt 	time.Time 	`gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
}

// get todos using GORM
func getTodos(c *gin.Context, db *gorm.DB) {
	var todos []Todo
	db.Find(&todos)
	c.JSON(200, gin.H{"data": todos})
}

func getTodo(c *gin.Context, db *gorm.DB) {
	var todo Todo
	id := c.Param("id")
	if err := db.Where("id = ?", id).First(&todo).Error; err != nil {
		c.JSON(400, gin.H{"error": "Could not retrieve todo with id " + id})
		return
	}
	c.JSON(200, gin.H{"data": todo})
}

func createTodo(c *gin.Context, db *gorm.DB) {
	var todo Todo
	err := c.BindJSON(&todo)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	err = db.Create(&todo).Error
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	c.JSON(201, gin.H{"data": todo})
}

func updateTodo(c *gin.Context, db *gorm.DB) {
	var todo Todo
	id := c.Param("id")
	if err := db.Where("id = ?", id).First(&todo).Error; err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	err := c.BindJSON(&todo)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	err = db.Save(&todo).Error
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"data": todo})
}

func deleteTodo(c *gin.Context, db *gorm.DB) {
	var todo Todo
	id := c.Param("id")
	if err := db.Where("id = ?", id).First(&todo).Error; err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	err := db.Delete(&todo).Error
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"data": true})
}


func main() {
	godotenv.Load()
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=UTC",
        os.Getenv("DB_HOST"),
        os.Getenv("DB_USER"),
        os.Getenv("DB_PASSWORD"),
        os.Getenv("DB_NAME"),
        os.Getenv("DB_PORT"))
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	// ping the db
	postgresDB, err := db.DB()
	err = postgresDB.Ping()
	if err != nil {
		panic(err)
	}

	db.AutoMigrate(&Todo{})


	router := gin.Default()

	router.GET("/todos", func(c *gin.Context) {
		getTodos(c, db)
	})
	router.GET("/todos/:id", func(c *gin.Context) {
		getTodo(c, db)
	})
	router.POST("/todos", func(c *gin.Context) {
		createTodo(c, db)
	})
	router.PATCH("/todos/:id", func(c *gin.Context) {
		updateTodo(c, db)
	})
	router.DELETE("/todos/:id", func(c *gin.Context) {
		deleteTodo(c, db)
	})


	router.Run(":8000")
}
