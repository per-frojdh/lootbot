package lootlists

import (
	"github.com/gin-gonic/gin"
    "github.com/jinzhu/gorm"
    "net/http"
    "strconv"
    models "lootbot/models"
    util "lootbot/lib"
	"fmt"
)

// GetLootLists ...
func GetLootLists(c *gin.Context) {
    name := c.Param("name");
    
    // Get the DB context
    db, ok := c.MustGet("databaseConnection").(gorm.DB)
    if !ok {
        c.AbortWithStatus(http.StatusBadRequest)
    }
    
    var char models.Character
    
    if db.Where(&models.Character{
        Name: name,
    }).First(&char).RecordNotFound() {
        c.JSON(http.StatusNotFound, gin.H{ "message" : models.ErrorMessages["RESOURCE_NOT_FOUND"] })
        return
    }
    db.Model(&char).Association("Lootlist").Find(&char.Lootlist)
    if len(char.Lootlist) == 0 {
        c.JSON(http.StatusNotFound, gin.H{ "message": "No items added to lootlist"})
        return
    }
    someData, _ := util.ParseItems(char.Lootlist)
    char.Lootlist = someData
    c.JSON(http.StatusOK, char)
}

// AddItem ...
func AddItem(c *gin.Context) {

    db, ok := c.MustGet("databaseConnection").(gorm.DB)
    if !ok {
        // Do something
    }
    
    character := c.PostForm("character")
    id := c.Param("id")
    
    if len(character) == 0 || len(id) == 0 {
        c.JSON(http.StatusBadRequest, gin.H{ "message" : models.ErrorMessages["BAD_INPUT_PARAMETERS"]})
        return
    }
    
    itemID, err := strconv.Atoi(id)
    
    // Convert Parameter to int, for db query
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{ "message" : models.ErrorMessages["BAD_INPUT_PARAMETERS"]})
        return
    }
    
    authUser, ok := c.MustGet("authUser").(models.User)
    if !ok {
        c.AbortWithStatus(http.StatusBadRequest)
    }
    
    var char models.Character
    var item models.Item

    if db.Where(&models.Item{
        ItemID: itemID,
        Context: "raid-mythic",
    }).First(&item).RecordNotFound() {
        c.JSON(http.StatusNotFound, gin.H{ "message" : models.ErrorMessages["RESOURCE_NOT_FOUND"] })
        return
    }
    
    if db.Where(&models.Character{
        Name: character,
        UserID: authUser.ID,
    }).First(&char).RecordNotFound() {
        c.JSON(http.StatusNotFound, gin.H{ "message" : models.ErrorMessages["RESOURCE_NOT_FOUND"] })
        return
    }
    
    lootlist := db.Model(&char).Association("Lootlist").Find(&char.Lootlist)
    found := false    
    for _, obj := range char.Lootlist {
        fmt.Println("Comparing: ", obj.ItemID, itemID)
        if (obj.ItemID == itemID) {
            found = true
        }
    }
    
    if (!found) {
        lootlist.Append(&item)
        c.JSON(http.StatusOK, char)     
        return   
    }
    c.String(http.StatusNotModified, "You can't add an item you already have on your itemlist") 
}

// RemoveItem ... 
func RemoveItem(c *gin.Context) {
    db, ok := c.MustGet("databaseConnection").(gorm.DB)
    if !ok {
        // Do something
    }
    
    character := c.PostForm("character")
    id := c.Param("id")
    
    if len(character) == 0 || len(id) == 0 {
        c.JSON(http.StatusBadRequest, gin.H{ "message" : models.ErrorMessages["BAD_INPUT_PARAMETERS"]})
        return
    }
    
    itemID, err := strconv.Atoi(id)
    
    // Convert Parameter to int, for db query
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{ "message" : models.ErrorMessages["BAD_INPUT_PARAMETERS"]})
        return
    }
    
    authUser, ok := c.MustGet("authUser").(models.User)
    if !ok {
        c.AbortWithStatus(http.StatusBadRequest)
    }
    
    var char models.Character
    var item models.Item
    
    if db.Where(&models.Item{
        ItemID: itemID,
        Context: "raid-mythic",
    }).First(&item).RecordNotFound() {
        c.JSON(http.StatusNotFound, gin.H{ "message" : models.ErrorMessages["RESOURCE_NOT_FOUND"] })
        return
    }
    
    if db.Where(&models.Character{
        Name: character,
        UserID: authUser.ID,
    }).First(&char).RecordNotFound() {
        c.JSON(http.StatusNotFound, gin.H{ "message" : models.ErrorMessages["RESOURCE_NOT_FOUND"] })
        return
    }
    
    lootlist := db.Model(&char).Association("Lootlist").Find(&char.Lootlist)
    found := false    
    for _, obj := range char.Lootlist {
        fmt.Println("Comparing: ", obj.ItemID, itemID)
        if (obj.ItemID == itemID) {
            found = true
        }
    }
    
    if (found) {
        lootlist.Delete(&item)
        c.JSON(http.StatusOK, char)
        return   
    }
    
    c.String(http.StatusNotModified, "You don't have the item you're trying to delete")
       
}
