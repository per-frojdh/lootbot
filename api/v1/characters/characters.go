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
    
    var characters []models.Character
    if db.Where(&models.Character{
        UserID: authUser.ID,
    }).Find(&characters).RecordNotFound() {
        c.Error(util.CreatePanicResponse("RESOURCE_NOT_FOUND")).
            SetMeta(util.CreateErrorResponse(http.StatusNotFound, "RESOURCE_NOT_FOUND"))
        c.Abort()
        return
    }
    
    c.JSON(http.StatusOK, characters)
}

func CreateCharacter(c *gin.Context) {
    realm := c.PostForm("realm")
    character := c.PostForm("character")
    
    if len(realm) == 0 || len(character) == 0 {
        c.Error(util.CreatePanicResponse("BAD_INPUT_PARAMETERS")).
            SetMeta(util.CreateErrorResponse(http.StatusBadRequest, "BAD_INPUT_PARAMETERS"))
        c.Abort()
        return
    }
    
    cfg, ok := c.MustGet("config").(config.Configuration)
    if !ok {
        // Can't really happen
    }

    log.Println("Attempting import of character:", realm, character)
    importedCharacter, errs := util.FetchCharacter(realm, character, cfg.ApiKey)
    
    if errs != nil && importedCharacter == nil {
        c.Error(util.CreatePanicResponse("FAILED_BNET")).
            SetMeta(util.CreateErrorResponse(http.StatusInternalServerError, "FAILED_BNET"))
        c.Abort()
        return
    }
    
    if importedCharacter != nil {
        
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
        
        importedCharacter.UserID = authUser.ID        
        
        var character models.Character
        if db.Where(&models.Character{
            Realm: importedCharacter.Realm,
            Name: importedCharacter.Name,
        }).First(&character).RecordNotFound() {
            db.NewRecord(importedCharacter)
            db.Create(&importedCharacter)
            
            if db.Error != nil {
                c.Error(util.CreatePanicResponse("FAILED_DATABASE_CREATION")).
                    SetMeta(util.CreateErrorResponse(http.StatusInternalServerError, "FAILED_DATABASE_CREATION"))
                c.Abort()
                return
            }    
        } else {
            c.Error(util.CreatePanicResponse("CHARACTER_EXIST")).
                    SetMeta(util.CreateErrorResponse(http.StatusBadRequest, "CHARACTER_EXIST"))
            c.Abort()
            return
        }
        c.JSON(http.StatusOK, importedCharacter)
    }
}

func DeleteCharacter(c *gin.Context) {
    name := c.PostForm("name")
    realm := c.PostForm("realm")

    if len(name) == 0 || len(realm) == 0 {
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
    
    var character models.Character

    if db.Where(&models.Character{
        Name: name,
        Realm: realm,
        UserID: authUser.ID,
    }).First(&character).RecordNotFound() {
        c.Error(util.CreatePanicResponse("RESOURCE_NOT_FOUND")).
            SetMeta(util.CreateErrorResponse(http.StatusNotFound, "RESOURCE_NOT_FOUND"))
        c.Abort()
    } else {
        db.Delete(&character)
        msg := fmt.Sprintf("Character %[1]s successfully deleted", name)
        c.JSON(http.StatusOK, gin.H{ "message" : msg })   
    }   
}
