// @SubApi User  [/users]
// @SubApi Allows you access to different features of the users , login , get status etc [/users]
package users

import (
    "log"
    "crypto/rand"
    "encoding/base64"
	"github.com/gin-gonic/gin"
    "github.com/jinzhu/gorm"
    bcrypt "golang.org/x/crypto/bcrypt"
    "net/http"
    models "lootbot/models"
)

// GetUsers ...
func GetUsers(c *gin.Context) {
    // Get the DB context
    db, ok := c.MustGet("databaseConnection").(gorm.DB)
    if !ok {
        // Do something
    }
    
    var returnedUser[] models.User
    db.Find(&returnedUser)
    
    if db.Error != nil {
        c.JSON(http.StatusNotFound, gin.H{ "message" : models.ErrorMessages["RESOURCE_NOT_FOUND"] })
    } else {
        c.JSON(http.StatusOK, returnedUser)    
    }
}

// GetUser ...
func GetUser(c *gin.Context) {
    login := c.Param("name");
    
    // Get the DB context
    db, ok := c.MustGet("databaseConnection").(gorm.DB)
    if !ok {
        // Do something
    }
    
    var user models.User
    
    if db.Where(&models.User{
        Login: login,
    }).First(&user).RecordNotFound() {
        c.JSON(http.StatusNotFound, gin.H{ "message" : models.ErrorMessages["RESOURCE_NOT_FOUND"] })
        return
    }
    
    db.Model(&user).Association("Characters").Find(&user.Characters)
    c.JSON(http.StatusOK, user)
}

// RegisterUser ...
func RegisterUser(c *gin.Context) {
    login := c.PostForm("login")
    name := c.PostForm("name")
    email := c.PostForm("email")
    password := c.PostForm("password")
    
    db, ok := c.MustGet("databaseConnection").(gorm.DB)
    if !ok {
        log.Fatal("DatabaseConnection failed")
        // Do something
    }
    
    hash := make([]byte, 32)
    rand.Read(hash)
    secureString := base64.URLEncoding.EncodeToString(hash)
    
    user := models.User{
        Login: login,
        Name: name,
        Email: email,
        GuildID: 1,
        SecretQuestion: "Test",
        SecretAnswer: "Test",
        Token: secureString,
    }
    
    // This should give them a hashed password    
    hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost) 
    user.Password = string(hashedPassword)
    
    db.NewRecord(user)
    db.Create(&user)
    
    if db.Error != nil {
        c.JSON(http.StatusBadRequest, gin.H{ "message" : models.ErrorMessages["FAILED_CREATING_USER"] })
    } else {
        c.JSON(http.StatusOK, user)    
    }    
}

// DeleteUser ...
func DeleteUser(c *gin.Context) {
    authUser, ok := c.MustGet("authUser").(models.User)
    if !ok {
        c.AbortWithStatus(http.StatusBadRequest)
    }
    
    db, ok := c.MustGet("databaseConnection").(gorm.DB)
    if !ok {
        c.AbortWithStatus(http.StatusInternalServerError)
    }
    
    var user models.User
    if db.First(&user, authUser.ID).RecordNotFound() {
        c.AbortWithStatus(http.StatusBadRequest)
    }
    
    db.Unscoped().Delete(&user)
    c.JSON(http.StatusOK, gin.H{"message": "User deleted"})
}