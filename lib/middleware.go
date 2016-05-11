package lib

import (
	"github.com/jinzhu/gorm"
	"github.com/gin-gonic/gin"
	"net/http"
	"github.com/per-frojdh/lootbot/models"
	"github.com/per-frojdh/lootbot/config"
    "log"
    "os"
)

// AddDBContext ...
func AddDBContext(db gorm.DB) gin.HandlerFunc {
    return func(c *gin.Context) {
        c.Set("databaseConnection", db)
        c.Next()
    }
}

func AuthorizeSource() gin.HandlerFunc {
    return func(c *gin.Context) {
        Source := c.Request.Header.Get("X-Source")
        if len(Source) == 0 {
            log.Println("No verifiable source")
            c.AbortWithStatus(http.StatusForbidden)
        }
        
        db, ok := c.MustGet("databaseConnection").(gorm.DB)
        if !ok {
            c.AbortWithStatus(http.StatusInternalServerError)
        }
        
        var token models.AccessToken
        
        if db.Where(&models.AccessToken{
                Token: Source,
        }).First(&token).RecordNotFound() {
            log.Println("Source not found with token", Source)    
            c.AbortWithStatus(http.StatusForbidden)
        }
        
        c.Next()
    }
}

func Authorization() gin.HandlerFunc {
    return func(c *gin.Context) {
        Token := c.Request.Header.Get("X-Auth-Token")
        if len(Token) == 0 {
            log.Println("No authorization header")
            c.AbortWithStatus(http.StatusForbidden)
        }

        db, ok := c.MustGet("databaseConnection").(gorm.DB)
        if !ok {
            c.AbortWithStatus(http.StatusInternalServerError)
        }
        var user models.User
        
        if db.Where(&models.User{
                Token: Token,
        }).First(&user).RecordNotFound() {
            log.Println("User not found with token", Token)    
            c.AbortWithStatus(http.StatusForbidden)
        }

        c.Set("authUser", user)
        c.Next()
    }
    
}

func AddConfigContext(cfg config.Configuration) gin.HandlerFunc {
    return func(c *gin.Context) {
        c.Set("config", cfg)
        c.Next()
    }
}

func DebugPostForm() gin.HandlerFunc {
    return func(c *gin.Context) {
        log.Println("Parsing form values")
        log.Println(c.Request.Header.Get("Content-Type"))
        c.Request.ParseForm()
        values := c.Request.Form
        log.Println(values)
        for k, v := range values {
            log.Println(k, v[0])   
        }
    }
}

func ErrorHandler() gin.HandlerFunc {
    return func(c *gin.Context) {
        c.Next()
        
        lastError := c.Errors.ByType(gin.ErrorTypeAny).Last()
        if lastError != nil && lastError.Meta != nil {
            errorList := []ResponseMessage{}
            for _, err := range c.Errors.ByType(gin.ErrorTypeAny) {
                errorList = append(errorList, err.Meta.(ResponseMessage))
            }
            
            apiError := APIError{}
            c.Request.ParseForm()
            list := []Parameter{}
            values := c.Request.Form
            for k, v := range values {
                list = append(list, Parameter{ Key: k, Value: v[0]})
            }
            
            statusCode := lastError.Meta.(ResponseMessage).StatusCode
            if len(os.Getenv("dev")) == 0 && c.Request.Method == "POST" {
                apiError.Request.ContentType = c.Request.Header.Get("Content-Type")
                apiError.Request.Params = list    
            }

            apiError.Errors = errorList
            c.JSON(statusCode, apiError)    
        }
        
        
    }
}