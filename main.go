package main

import (
	"math/big"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type shortLink struct {
	ID  int64  `json:"-" gorm:"primaryKey"`
	URL string `json:"url" binding:"required,url"`
}

func main() {
	dsn := "user:pass@tcp(127.0.0.1:3306)/dbname?charset=utf8mb4&parseTime=True&loc=Local"
	db, _ := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	db.AutoMigrate(&shortLink{})
	db.Raw("alter table short_links AUTO_INCREMENT=2001").Scan(nil)

	r := gin.Default()
	r.POST("/links", func(c *gin.Context) {
		var links shortLink
		if err := c.ShouldBindJSON(&links); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		db.Create(&links)
		c.JSON(http.StatusOK, gin.H{"short": big.NewInt(links.ID).Text(62)})
	})
	r.GET("/:id", func(c *gin.Context) {
		id, ok := new(big.Int).SetString(c.Param("id"), 62)
		if !ok {
			c.JSON(http.StatusBadRequest, gin.H{"error": "wrong id format"})
			return
		}
		var links shortLink
		if err := db.First(&links, id.Int64()).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.Redirect(http.StatusPermanentRedirect, links.URL)
	})
	r.Run()
}
