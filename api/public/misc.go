package misc

import (
	"github.com/gin-gonic/gin"
    "github.com/jinzhu/gorm"
    "net/http"
    //
    models "lootbot/models"
)

// HealthCheck ...
func HealthCheck(c *gin.Context) {
    // Get the DB context
    db, ok := c.MustGet("databaseConnection").(gorm.DB)
    if !ok {
        // Do something
    }
    
    var returnedUser models.User
    db.First(&returnedUser)
    
    c.JSON(http.StatusOK, gin.H{"HealthCheck": "OK"})
   
}