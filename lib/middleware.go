package lib

import (
	"github.com/jinzhu/gorm"
	"github.com/gin-gonic/gin"
	"net/http"
	"lootbot/models"
	"lootbot/config"
    "log"
)

// AddDBContext ...
func AddDBContext(db gorm.DB) gin.HandlerFunc {
    return func(c *gin.Context) {
        c.Set("databaseConnection", db)
        c.Next()
    }
}

func Authorization() gin.HandlerFunc {
    return func(c *gin.Context) {
        Token := c.Request.Header.Get("X-Auth-Token")
        if len(Token) == 0 {
            log.Println("No authorization header")
            c.AbortWithStatus(http.StatusForbidden)
        }

        db, ok := c.MustGet("databaseConnection").(gorm.DB)
        if !ok {
            c.AbortWithStatus(http.StatusInternalServerError)
        }
        var user models.User
        
        if db.Where(&models.User{
                Token: Token,
        }).First(&user).RecordNotFound() {
            log.Println("User not found with token", Token)    
            c.AbortWithStatus(http.StatusForbidden)
        }

        c.Set("authUser", user)
        c.Next()
    }
    
}

func AddConfigContext(cfg config.Configuration) gin.HandlerFunc {
    return func(c *gin.Context) {
        c.Set("config", cfg)
        c.Next()
    }
}