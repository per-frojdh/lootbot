// @SubApi User  [/users]
// @SubApi Allows you access to different features of the users , login , get status etc [/users]
package users

import (
    "log"
	"github.com/gin-gonic/gin"
    "github.com/jinzhu/gorm"
    "net/http"
    models "github.com/per-frojdh/lootbot/models"
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
    login := c.PostForm("Username")
    // name := c.PostForm("username")
    // email := c.PostForm("email")
    // password := c.PostForm("password")
    token := c.PostForm("ID")
    
    log.Println("Token is: ", token)
    log.Println("Username is: ", login)
    
    db, ok := c.MustGet("databaseConnection").(gorm.DB)
    if !ok {
        log.Fatal("DatabaseConnection failed")
        // Do something
    }
    
    if len(login) == 0 || len(token) == 0 {
        c.JSON(http.StatusBadRequest, gin.H{ "message" : models.ErrorMessages["BAD_INPUT_PARAMETERS"]})
        return
    }
    
    // Don't need this right now
    // hash := make([]byte, 32)
    // rand.Read(hash)
    // secureString := base64.URLEncoding.EncodeToString(hash)
    
    user := models.User{
        Login: login,
        Name: login,
        // Email: email,
        GuildID: 3,
        SecretQuestion: "Test",
        SecretAnswer: "Test",
        Token: token,
    }
    
    // Don't need this right now
    // This should give them a hashed password    
    // hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost) 
    // user.Password = string(hashedPassword)
    
    if db.Where(&models.User{
        Token: token,
    }).First(&user).RecordNotFound() {
        db.NewRecord(user)
        db.Create(&user)
        c.JSON(http.StatusOK, user)
        return    
    }
    
    c.JSON(http.StatusBadRequest, gin.H{ "message" : models.ErrorMessages["FAILED_CREATING_USER"] })
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