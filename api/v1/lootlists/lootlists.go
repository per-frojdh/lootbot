package lootlists

import (
	"github.com/gin-gonic/gin"
    "github.com/jinzhu/gorm"
    "net/http"
    "strconv"
    models "github.com/per-frojdh/lootbot/models"
    util "github.com/per-frojdh/lootbot/lib"
)

// GetLootLists ...
func GetLootLists(c *gin.Context) {
    name := c.Param("name");
    
    if len(name) == 0 {
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
    
    var char models.Character
    
    if db.Where(&models.Character{
        Name: name,
    }).First(&char).RecordNotFound() {
        c.Error(util.CreatePanicResponse("RESOURCE_NOT_FOUND")).
            SetMeta(util.CreateErrorResponse(http.StatusNotFound, "RESOURCE_NOT_FOUND"))
        c.Abort()
        return
    }
    
    db.Model(&char).Association("Lootlist").Find(&char.Lootlist)
    
    if len(char.Lootlist) == 0 {
        c.Error(util.CreatePanicResponse("NO_LOOTLIST_ITEMS")).
            SetMeta(util.CreateErrorResponse(http.StatusNotFound, "NO_LOOTLIST_ITEMS"))
        c.Abort()
        return
    }
    
    someData, _ := util.ParseItems(char.Lootlist)
    char.Lootlist = someData
    c.JSON(http.StatusOK, char)
}

// AddItem ...
func AddItem(c *gin.Context) {
    character := c.PostForm("character")
    context := c.PostForm("context")
    id := c.Param("id")
    
    itemID, err := strconv.Atoi(id)
    
    if len(character) == 0 || (len(id) == 0 || itemID < 1) || len(context) == 0 || err != nil {
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
    
    db, ok := c.MustGet("databaseConnection").(gorm.DB)
    if !ok {
        c.Error(util.CreatePanicResponse("DATABASE_ERROR")).
            SetMeta(util.CreateErrorResponse(http.StatusInternalServerError, "DATABASE_ERROR"))
        c.Abort()
        return
    }
    
    authUser, ok := c.MustGet("authUser").(models.User)
    if !ok {
        c.Error(util.CreatePanicResponse("AUTHENTICATION_FAILED")).
            SetMeta(util.CreateErrorResponse(http.StatusForbidden, "AUTHENTICATION_FAILED"))
        c.Abort()
        return
    }
    
    var char models.Character
    var item models.Item

    if db.Where(&models.Item{
        ItemID: itemID,
        Context: context,
    }).First(&item).RecordNotFound() {
        c.Error(util.CreatePanicResponse("RESOURCE_NOT_FOUND")).
            SetMeta(util.CreateErrorResponse(http.StatusNotFound, "RESOURCE_NOT_FOUND"))
        c.Abort()
        return
    }
    
    if db.Where(&models.Character{
        Name: character,
        UserID: authUser.ID,
    }).First(&char).RecordNotFound() {
        c.Error(util.CreatePanicResponse("RESOURCE_NOT_FOUND")).
            SetMeta(util.CreateErrorResponse(http.StatusNotFound, "RESOURCE_NOT_FOUND"))
        c.Abort()
        return
    }
    
    lootlist := db.Model(&char).Association("Lootlist").Find(&char.Lootlist)
    found := false    
    for _, obj := range char.Lootlist {
        if (obj.ItemID == itemID) {
            found = true
        }
    }
    
    char.Lootlist, _ = util.ParseItems(char.Lootlist)    
        
    if (found) {
        c.Error(util.CreatePanicResponse("ITEM_ALREADY_ADDED")).
            SetMeta(util.CreateErrorResponse(http.StatusBadRequest, "ITEM_ALREADY_ADDED"))
        c.Abort()
        return   
    }
    
    item, _ = util.ParseItem(item)
    lootlist.Append(&item)
    
    c.JSON(http.StatusOK, char)    
}

// RemoveItem ... 
func RemoveItem(c *gin.Context) {
    character := c.PostForm("character")
    id := c.Param("id")
    itemID, err := strconv.Atoi(id)
    
    if len(character) == 0 || len(id) == 0 || itemID < 1 || err != nil {
        c.Error(util.CreatePanicResponse("BAD_INPUT_PARAMETERS")).
            SetMeta(util.CreateErrorResponse(http.StatusBadRequest, "BAD_INPUT_PARAMETERS"))
        c.Abort()
        return
    }
    
    db, ok := c.MustGet("databaseConnection").(gorm.DB)
    if !ok {
        c.Error(util.CreatePanicResponse("DATABASE_ERROR")).
            SetMeta(util.CreateErrorResponse(http.StatusInternalServerError, "DATABASE_ERROR"))
        c.Abort()
        return
    }
    
    authUser, ok := c.MustGet("authUser").(models.User)
    if !ok {
        c.Error(util.CreatePanicResponse("AUTHENTICATION_FAILED")).
            SetMeta(util.CreateErrorResponse(http.StatusForbidden, "AUTHENTICATION_FAILED"))
        c.Abort()
        return
    }
    
    var char models.Character
    var item models.Item
    
    if db.Where(&models.Item{
        ItemID: itemID,
        Context: "raid-mythic",
    }).First(&item).RecordNotFound() {
        c.Error(util.CreatePanicResponse("RESOURCE_NOT_FOUND")).
            SetMeta(util.CreateErrorResponse(http.StatusNotFound, "RESOURCE_NOT_FOUND"))
        c.Abort()
        return
    }
    
    if db.Where(&models.Character{
        Name: character,
        UserID: authUser.ID,
    }).First(&char).RecordNotFound() {
        c.Error(util.CreatePanicResponse("RESOURCE_NOT_FOUND")).
            SetMeta(util.CreateErrorResponse(http.StatusNotFound, "RESOURCE_NOT_FOUND"))
        c.Abort()
        return
    }
    
    lootlist := db.Model(&char).Association("Lootlist").Find(&char.Lootlist)
    
    found := false    
    for _, obj := range char.Lootlist {
        if (obj.ItemID == itemID) {
            found = true
        }
    }
    
    if (!found) {
        c.Error(util.CreatePanicResponse("ITEM_NOT_ADDED")).
            SetMeta(util.CreateErrorResponse(http.StatusBadRequest, "ITEM_NOT_ADDED"))
        c.Abort()
        return
    }
    
    lootlist.Delete(&item)
    c.JSON(http.StatusOK, char)
}
