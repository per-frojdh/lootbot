package v1

import (
	"github.com/gin-gonic/gin"
    "github.com/jinzhu/gorm"
    "net/http"
    "fmt"
    "log"
    models "github.com/per-frojdh/lootbot/models"
    config "github.com/per-frojdh/lootbot/config"
    util "github.com/per-frojdh/lootbot/lib"
)

// GetCharacters ...
func GetCharacters(c *gin.Context) {
    // Get the DB context
    db, ok := c.MustGet("databaseConnection").(gorm.DB)
    if !ok {
        // Do something
    }
    
    authUser, ok := c.MustGet("authUser").(models.User)
    if !ok {
        c.AbortWithStatus(http.StatusBadRequest)
        // Do something
    }
    
    var characters []models.Character
    if db.Where(&models.Character{
        UserID: authUser.ID,
    }).Find(&characters).RecordNotFound() {
        c.JSON(http.StatusNotFound, gin.H{ "message" : models.ErrorMessages["RESOURCE_NOT_FOUND"] })
        return
    }
    
    c.JSON(http.StatusOK, characters)
}

func CreateCharacter(c *gin.Context) {
    realm := c.PostForm("realm")
    character := c.PostForm("character")
    
    cfg, ok := c.MustGet("config").(config.Configuration)
    if !ok {
        c.JSON(http.StatusInternalServerError, gin.H{ "message" : models.ErrorMessages["DATABASE_ERROR"]})
        return
    }
    
    log.Println("Attempting import of character:", realm, character)
    
    if len(realm) == 0 || len(character) == 0 {
        c.JSON(http.StatusBadRequest, gin.H{ "message" : models.ErrorMessages["BAD_INPUT_PARAMETER"]})
        return
    }

    importedCharacter, errs := util.FetchCharacter(realm, character, cfg.ApiKey)
    
    if errs != nil && importedCharacter == nil {
        c.JSON(http.StatusBadRequest, gin.H{ "message" : models.ErrorMessages["FAILED_BNET"]})
        return
    }
    
    if importedCharacter != nil {
        authUser, ok := c.MustGet("authUser").(models.User)
        if !ok {
            c.AbortWithStatus(http.StatusBadRequest)
            // Do something
        }
        importedCharacter.UserID = authUser.ID
        
        db, ok := c.MustGet("databaseConnection").(gorm.DB)
        if !ok {
            c.JSON(http.StatusInternalServerError, gin.H{ "message" : models.ErrorMessages["DATABASE_ERROR"]})
            return
        }
        
        var character models.Character
        if db.Where(&models.Character{
            Realm: importedCharacter.Realm,
            Name: importedCharacter.Name,
        }).First(&character).RecordNotFound() {
            db.NewRecord(importedCharacter)
            db.Create(&importedCharacter)
            
            if db.Error != nil {
                c.JSON(http.StatusInternalServerError, gin.H{ "message" : models.ErrorMessages["DATABASE_ERROR"]})
                return
            }    
        } else {
            c.JSON(http.StatusBadRequest, gin.H{ "message" : models.ErrorMessages["CHARACTER_EXIST"]})
            return
        }

        c.JSON(http.StatusOK, importedCharacter)
    }
}

func DeleteCharacter(c *gin.Context) {
    name := c.PostForm("name")
    realm := c.PostForm("realm")

    if len(name) == 0 || len(realm) == 0 {
        c.JSON(http.StatusBadRequest, gin.H{ "message" : models.ErrorMessages["BAD_INPUT_PARAMETER"] })
        return
    }
    
    
    name = util.CapitalizeString(name)    
    
    db, ok := c.MustGet("databaseConnection").(gorm.DB)
    if !ok {
        c.AbortWithStatus(http.StatusBadRequest)
    }
    
    authUser, ok := c.MustGet("authUser").(models.User)
    if !ok {
        c.AbortWithStatus(http.StatusForbidden)
    }
    
    var character models.Character

    if db.Where(&models.Character{
        Name: name,
        Realm: realm,
        UserID: authUser.ID,
    }).First(&character).RecordNotFound() {
        c.JSON(http.StatusNotFound, gin.H{ "message" : models.ErrorMessages["RESOURCE_NOT_FOUND"] })
        return
    }
    
    db.Delete(&character)
    msg := fmt.Sprintf("Character %[1]s successfully deleted", name)
    c.JSON(http.StatusOK, gin.H{ "message" : msg })
    
}
