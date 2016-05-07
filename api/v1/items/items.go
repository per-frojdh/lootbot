package v1

import (
	"github.com/gin-gonic/gin"
    "github.com/jinzhu/gorm"
    "net/http"
    "strconv"
    util "github.com/per-frojdh/lootbot/lib"
    models "github.com/per-frojdh/lootbot/models"
	"strings"
)

// GetItem ...
func GetItem(c *gin.Context) {
    id, err := strconv.Atoi(c.Param("id"))
    
    // Convert Parameter to int, for db query
    if err != nil || id < 1 {
        c.AbortWithStatus(http.StatusBadRequest)
    }
    
    // Get the DB context
    db, ok := c.MustGet("databaseConnection").(gorm.DB)
    if !ok {
        c.AbortWithStatus(http.StatusInternalServerError)
    }
    
    // Hold the structified item here.
    var returnedItems []models.Item
    
    // Get the db row    
    if db.Where(&models.Item{
        ItemID: id,
    }).Find(&returnedItems).RecordNotFound() {
        c.JSON(http.StatusNotFound, gin.H{ "message" : models.ErrorMessages["RESOURCE_NOT_FOUND"] })
        return
    }
    
    // Loop through the items from DB
    returnData, _ := util.ParseItems(returnedItems)
    
    // Respond with the struct as json
    c.JSON(http.StatusOK, returnData)
}

// SearchItems ...
func SearchItems(c *gin.Context) {
    query := c.Query("search")
    
    if len(query) == 0 {
        c.AbortWithStatus(http.StatusBadRequest)
    }
    // Get the DB context
    db, ok := c.MustGet("databaseConnection").(gorm.DB)
    if !ok {
        c.AbortWithStatus(http.StatusInternalServerError)
    }
    arr := []string{"%", query, "%"}
    searchQuery := strings.Join(arr, "")
    var items []models.Item
    if db.Where("name ILIKE ?", searchQuery).Find(&items).RecordNotFound() {
        c.JSON(http.StatusBadRequest, gin.H{"message" : "Could not find any items with that name"})
    }
    
    returnData, _ := util.ParseItems(items)
    
    c.JSON(http.StatusOK, returnData)
}
