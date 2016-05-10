package v2

import (
    "log"
    "strings"
    "time"
    "os"
    "encoding/base64"
	"github.com/gin-gonic/gin"
    "github.com/jinzhu/gorm"
    bcrypt "golang.org/x/crypto/bcrypt"
    "net/http"
    models "github.com/per-frojdh/lootbot/models"
    jwt_lib "github.com/dgrijalva/jwt-go"
    // "github.com/joho/godotenv"
)

// Authenticate ...
func Authenticate(c *gin.Context) {
    basicHeader := c.Request.Header.Get("Authorization")
    if len(basicHeader) == 0 {
        c.AbortWithStatus(http.StatusForbidden)
        return
    }
    auth := strings.SplitN(basicHeader, " ", 2)
    if auth[0] != "Basic" || len(auth) != 2 {
        log.Println("Basic authentication doesn't look right")
        c.AbortWithStatus(http.StatusForbidden)
        return
    }
    
    payload, _ := base64.StdEncoding.DecodeString(auth[1])
    pair := strings.SplitN(string(payload), ":", 2)
    
    if len(pair) != 2 {
        log.Println("Couldn't decode authorization header")
        c.AbortWithStatus(http.StatusForbidden)
        return
    }
    
    db, ok := c.MustGet("databaseConnection").(gorm.DB)
    if !ok {
        log.Fatal("DatabaseConnection failed")
        // Do something
    }
    
    password := pair[1]
    username := pair[0]
    
    var user models.User
    
    db.Where(&models.User{
        Login: username,
    }).First(&user)
    
    err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
    if err != nil {
        log.Println(err)
        c.JSON(http.StatusBadRequest, gin.H{ "message": models.ErrorMessages["AUTHENTICATION_FAILED"]})
    } else {
        // Give token and stuff
        token := jwt_lib.New(jwt_lib.GetSigningMethod("HS256"))
        token.Claims["token"] = user.Token
        token.Claims["login"] = user.Login
        token.Claims["exp"] = time.Now().Add(time.Hour * 1).Unix()
        
        log.Println("Supersecretstring is: ", os.Getenv("SUPER_SECRET_TOKEN"))
        tokenString, err := token.SignedString([]byte(os.Getenv("SUPER_SECRET_TOKEN")))
        if err != nil {
            c.JSON(http.StatusBadRequest, gin.H{ "message": models.ErrorMessages["AUTHENTICATION_FAILED"]})
        }
        
        c.JSON(http.StatusOK, gin.H{"message": tokenString})
    }    
}

func ResetPassword(c *gin.Context) {
    c.JSON(http.StatusOK, gin.H{"message": "Successfully authenticated"})
}
