package config

import (
	"fmt"
    "log"
    "io/ioutil"
    "encoding/json"
    "github.com/joho/godotenv"
    "os"
    
)

// Configuration ...
type Configuration struct {
    Foo     string
    ConnectionString    string
    ApiKey string
}

// LoadConfig ...
func LoadConfig() Configuration{
    err := godotenv.Load("src/github.com/per-frojdh/lootbot/.env")
    if err != nil {
        log.Fatal("Error loading .env file", err)
    }
    
    var file []byte
    environment := os.Getenv("SETUP_ENVIRONMENT")
    
    if (environment == "dev") {
        file, err = ioutil.ReadFile("config.json")
    } else {
        file, err = ioutil.ReadFile("/home/wowutils/go/config.json")    
    }

    if err != nil {
		log.Fatal("Config File Missing. ", err)
	}
    
    var configuration Configuration
    
    err = json.Unmarshal(file, &configuration)
    if err != nil {
        log.Fatal("Config Parse Error: ", err)
    }
    
    if err != nil {
        fmt.Println("Error loading configuration: ", err)    
    }
    return configuration
}