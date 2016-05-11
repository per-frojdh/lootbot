// @SubApi User  [/users]
// @SubApi Allows you access to different features of the users , login , get status etc [/users]
package users

import (
    "log"
	"github.com/gin-gonic/gin"
    "github.com/jinzhu/gorm"
    "net/http"
    bcrypt "golang.org/x/crypto/bcrypt"
    models "github.com/per-frojdh/lootbot/models"
    util "github.com/per-frojdh/lootbot/lib"
)

// GetUsers ...
func GetUsers(c *gin.Context) {
    // Get the DB context
    db, ok := c.MustGet("databaseConnection").(gorm.DB)
    if !ok {
        c.Error(util.CreatePanicResponse("DATABASE_ERROR")).
            SetMeta(util.CreateErrorResponse(http.StatusInternalServerError, "DATABASE_ERROR"))
        c.Abort()
        return
    }
    
    var returnedUser[] models.User
    db.Find(&returnedUser)
    
    if db.Error != nil || len(returnedUser) == 0 {
        c.Error(util.CreatePanicResponse("RESOURCE_NOT_FOUND")).
            SetMeta(util.CreateErrorResponse(http.StatusNotFound, "RESOURCE_NOT_FOUND"))
        c.Abort()
        return
    }
    c.JSON(http.StatusOK, returnedUser)    
    
}

// GetUser ...
func GetUser(c *gin.Context) {
    login := c.Param("name");
    
    if len(login) == 0 {
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
    
    var user models.User
    
    if db.Where(&models.User{
        Login: login,
    }).First(&user).RecordNotFound() {
        c.Error(util.CreatePanicResponse("RESOURCE_NOT_FOUND")).
            SetMeta(util.CreateErrorResponse(http.StatusNotFound, "RESOURCE_NOT_FOUND"))
        c.Abort()
        return
    }
    
    db.Model(&user).Association("Characters").Find(&user.Characters)
    c.JSON(http.StatusOK, user)
}

// RegisterUser ...
func RegisterUser(c *gin.Context) {
    login := c.PostForm("Username")
    token := c.PostForm("ID")
    
    log.Println("Token is: ", token)
    log.Println("Username is: ", login)
    
    if len(login) == 0 || len(token) == 0 {
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
    
    user := models.User{
        Login: login,
        Name: login,
        GuildID: 3, // 3 happens to be Aeon on live
        SecretQuestion: "Test",
        SecretAnswer: "Test",
        Token: token,
    } 
    
    if db.Where(&models.User{
        Token: token,
    }).First(&user).RecordNotFound() {
        db.NewRecord(user)
        db.Create(&user)
        c.JSON(http.StatusOK, user)
        return    
    }
    
    c.Error(util.CreatePanicResponse("FAILED_CREATING_USER")).
            SetMeta(util.CreateErrorResponse(http.StatusBadRequest, "FAILED_CREATING_USER"))
    c.Abort()
}

// DeleteUser ...
func DeleteUser(c *gin.Context) {
    authUser, ok := c.MustGet("authUser").(models.User)
    if !ok {
        c.Error(util.CreatePanicResponse("AUTHENTICATION_FAILED")).
            SetMeta(util.CreateErrorResponse(http.StatusForbidden, "AUTHENTICATION_FAILED"))
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
    
    var user models.User
    if db.First(&user, authUser.ID).RecordNotFound() {
        c.Error(util.CreatePanicResponse("RESOURCE_NOT_FOUND")).
            SetMeta(util.CreateErrorResponse(http.StatusNotFound, "RESOURCE_NOT_FOUND"))
        c.Abort()
        return
    }
    
    db.Unscoped().Delete(&user)
    c.JSON(http.StatusOK, gin.H{"message": "User deleted"})
}

// ResetPassword ...
func ResetPassword(c *gin.Context) {
    c.JSON(http.StatusOK, gin.H{"message": "ResetPassword"})
}

// Register ...
func Register(c *gin.Context) {
    passphrase := c.PostForm("password")
    email := c.PostForm("email")
    
    if len(passphrase) == 0 || len(email) == 0 {
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
    
    user, ok := c.MustGet("authUser").(models.User)
    if !ok {
        c.Error(util.CreatePanicResponse("AUTHENTICATION_FAILED")).
            SetMeta(util.CreateErrorResponse(http.StatusForbidden, "AUTHENTICATION_FAILED"))
        c.Abort()
        return
    }
    
    hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(passphrase), bcrypt.MinCost) 
    user.Password = string(hashedPassword)
    
    user.Email = email
    user.Password = string(hashedPassword[:])
    
    db.Save(user)
    if db.Error != nil {
        c.Error(util.CreatePanicResponse("FAILED_CREATING_USER")).
            SetMeta(util.CreateErrorResponse(http.StatusInternalServerError, "FAILED_CREATING_USER"))
        c.Abort()
        return
    }
    c.JSON(http.StatusOK, gin.H{"message": "Successfully registered user to web interface"})
}