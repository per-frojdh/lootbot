// @APIVersion 0.0.1
// @APITitle Lootplanner
// @APIDescription Create and manipulate lootlists for your character
// @Contact N/A
// @TermsOfServiceUrl N/A
// @License TBD
// @LicenseUrl 

package main

import (
	"fmt"
    "os"
	"github.com/gin-gonic/gin"
    "github.com/jinzhu/gorm"
    _ "github.com/lib/pq"
    "github.com/gin-gonic/contrib/jwt"
    //
    items "github.com/per-frojdh/lootbot/api/v1/items"
    users "github.com/per-frojdh/lootbot/api/v1/users"
    lootlists "github.com/per-frojdh/lootbot/api/v1/lootlists"
    characters "github.com/per-frojdh/lootbot/api/v1/characters"
    misc "github.com/per-frojdh/lootbot/api/v1/public"
    config "github.com/per-frojdh/lootbot/config"
    lib "github.com/per-frojdh/lootbot/lib"
    // v2
    v2 "github.com/per-frojdh/lootbot/api/v2"
)

func main() {
	fmt.Println("Go App starting..")
	fmt.Println("Loading config")
   
    cfg := config.LoadConfig()
    db, err := gorm.Open("postgres", cfg.ConnectionString)
    
    if err != nil {
        fmt.Println("Failed to connect to database")
        panic(err)
    }
    
    db.DB()
    db.DB().Ping()
    db.DB().SetMaxIdleConns(10)
    db.DB().SetMaxOpenConns(100)
    db.LogMode(true)
    
    router := gin.Default()
    router.Use(gin.Recovery())
    router.Use(lib.AddDBContext(*db))
    router.Use(lib.AddConfigContext(cfg))
    
    // TODO: Figure out if I can break out these to a separate file (in a nice way)
    api := router.Group("/api")
    {
        version1 := api.Group("/v1")
        
        version1.Use(lib.AuthorizeSource(), lib.Authorization())
        {
            itemEndpoint := version1.Group("/items")
            itemEndpoint.GET("/:id", items.GetItem)            
            itemEndpoint.GET("/", items.SearchItems)                
            
            userEndpoint := version1.Group("/users")
            userEndpoint.GET("/", users.GetUsers)
            userEndpoint.GET("/:name", users.GetUser)
            userEndpoint.POST("/delete", users.DeleteUser)
            userEndpoint.POST("/register", users.Register)
            
            lootEndpoint := version1.Group("/lootlist")
            lootEndpoint.GET("/:name", lootlists.GetLootLists)
            lootEndpoint.POST("/add/:id", lootlists.AddItem)
            lootEndpoint.POST("/delete/:id", lootlists.RemoveItem)
                       
            charEndpoint := version1.Group("/characters")
            charEndpoint.GET("/", characters.GetCharacters)
            charEndpoint.POST("/import", characters.CreateCharacter)
            charEndpoint.POST("/delete", characters.DeleteCharacter)
        }
        
        version2 := api.Group("/v2")
        
        version2.Use(jwt.Auth(os.Getenv("SUPER_SECRET_TOKEN")))
        {
            authEndpoint := version2.Group("/auth")
            authEndpoint.GET("/", v2.ResetPassword)
            
            userEndpoint := version2.Group("/users")
            userEndpoint.GET("/", users.GetUsers)
                       
        }
        
        // These should be all of the public endpoints (in the future)
        // TODO: Add authenticate here
        public := api.Group("/public") 
        {
            public.GET("/health", misc.HealthCheck)
            public.POST("/register", users.RegisterUser)
            public.POST("/login", v2.Authenticate)    
        }

    }
    
	router.Run(":1234")
}


