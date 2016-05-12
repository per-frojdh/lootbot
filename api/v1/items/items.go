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
    context := c.Param("context")
    
    // Convert Parameter to int, for db query
    if len(c.Param("id")) == 0 || err != nil || id < 1 || len(context) == 0{
        c.Error(util.CreatePanicResponse("BAD_INPUT_PARAMETERS")).
            SetMeta(util.CreateErrorResponse(http.StatusBadRequest, "BAD_INPUT_PARAMETERS"))
        c.Abort()
        return
    }   
     
    if !util.CheckValidContext(context) {
        c.Error(util.CreatePanicResponse("BAD_INPUT_PARAMETERS")).
            SetMeta(util.CreateErrorResponse(http.StatusBadRequest, "BAD_INPUT_PARAMETERS"))
        c.Abort()
        return
    }
    
    // Get the DB context
    db, ok := c.MustGet("databaseConnection").(gorm.DB)
    if !ok {
        c.Error(util.CreatePanicResponse("DATABASE_ERROR")).
            SetMeta(util.CreateErrorResponse(http.StatusInternalServerError, "DATABASE_ERROR"))
        c.Abort()
        return
    }
    
    // Hold the structified item here.
    var returnedItem models.Item
    
    // Get the db row    
    if db.Where(&models.Item{
        ItemID: id,
        Context: context,
    }).First(&returnedItem).RecordNotFound() {
        c.Error(util.CreatePanicResponse("RESOURCE_NOT_FOUND")).
            SetMeta(util.CreateErrorResponse(http.StatusNotFound, "RESOURCE_NOT_FOUND"))
        c.Abort()
        return
    }
   
    // Loop through the items from DB
    returnData, err := util.ParseItem(returnedItem)
    
    if err != nil {
        c.Error(util.CreatePanicResponse("DATABASE_ERROR")).
            SetMeta(util.CreateErrorResponse(http.StatusInternalServerError, "DATABASE_ERROR"))
        c.Abort()
        return
    }
    
    // Respond with the struct as json
    c.JSON(http.StatusOK, returnData)
}

// SearchItems ...
func SearchItems(c *gin.Context) {
    query := c.Query("search")
    
    if len(query) == 0 {
        c.Error(util.CreatePanicResponse("BAD_INPUT_PARAMETERS")).
            SetMeta(util.CreateErrorResponse(http.StatusBadRequest, "BAD_INPUT_PARAMETERS"))
        c.Abort()
        return
    }
    
    // Get the DB context
    db, ok := c.MustGet("databaseConnection").(gorm.DB)
    if !ok {
        c.Error(util.CreatePanicResponse("DATABASE_ERROR")).
            SetMeta(util.CreateErrorResponse(http.StatusInternalServerError, "DATABASE_ERROR"))
        c.Abort()
        return
    }
    arr := []string{"%", query, "%"}
    searchQuery := strings.Join(arr, "")
    var items []models.Item
    
    if db.Where("name ILIKE ?", searchQuery).
        Limit(10).
        Find(&items).
        RecordNotFound() {
            c.Error(util.CreatePanicResponse("RESOURCE_NOT_FOUND")).
                SetMeta(util.CreateErrorResponse(http.StatusNotFound, "RESOURCE_NOT_FOUND"))
            c.Abort()
            return
        }
    
    returnData, _ := util.ParseItems(items)
    
    c.JSON(http.StatusOK, returnData)
}
