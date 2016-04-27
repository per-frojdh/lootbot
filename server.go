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
	"github.com/gin-gonic/gin"
    "github.com/jinzhu/gorm"
    _ "github.com/lib/pq"
    //
    items "github.com/per-frojdh/lootbot/api/items"
    users "github.com/per-frojdh/lootbot/api/users"
    lootlists "github.com/per-frojdh/lootbot/api/lootlists"
    characters "github.com/per-frojdh/lootbot/api/characters"
    misc "github.com/per-frojdh/lootbot/api/public"
    config "github.com/per-frojdh/lootbot/config"
    lib "github.com/per-frojdh/lootbot/lib"
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
        v1 := api.Group("/v1")
        
        v1.Use(lib.Authorization())
        {
            itemEndpoint := v1.Group("/items")
            itemEndpoint.GET("/:id", items.GetItem)            
            itemEndpoint.GET("/", items.SearchItems)                
            
            userEndpoint := v1.Group("/users")
            userEndpoint.GET("/", users.GetUsers)
            userEndpoint.GET("/:name", users.GetUser)
            userEndpoint.POST("/delete", users.DeleteUser)
            
            lootEndpoint := v1.Group("/lootlist")
            lootEndpoint.GET("/:name", lootlists.GetLootLists)
            lootEndpoint.POST("/add/:id", lootlists.AddItem)
            lootEndpoint.POST("/delete/:id", lootlists.RemoveItem)
                       
            charEndpoint := v1.Group("/characters")
            charEndpoint.GET("/", characters.GetCharacters)
            charEndpoint.POST("/import", characters.CreateCharacter)
            charEndpoint.POST("/delete", characters.DeleteCharacter)
        }
        
        // These should be all of the public endpoints (in the future)
        // TODO: Add authenticate here
        public := api.Group("/public") 
        {
            public.GET("/health", misc.HealthCheck)
            public.POST("/register", users.RegisterUser)
        }

    }
    
	router.Run(":3000")
}


