package config

import (
	"fmt"
    "log"
    "io/ioutil"
    "encoding/json"
    
)

// Configuration ...
type Configuration struct {
    Foo     string
    ConnectionString    string
    ApiKey string
}

// LoadConfig ...
func LoadConfig() Configuration{
    file, err := ioutil.ReadFile("config.json")
    
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