package lib

import (
    models "github.com/per-frojdh/lootbot/models"
    "encoding/json"
)

func CheckValidContext(s string) bool {
    for _, context := range models.RaidContexts {
        if s == context {
            return true
        }
    }
    return false
}

func ParseItem (item models.Item) (models.Item, error) {
    // Fix stats
    err := json.Unmarshal(item.Stats, &item.StatsList)    

    // Fix weapon stats
    err = json.Unmarshal(item.WeaponStats, &item.WeaponInfo)    
    
    // Fix spells    
    err = json.Unmarshal(item.ItemSpell, &item.SpellInfo)
        
    if err != nil {
        return item, err
    }
    
    item.Stats = nil
    item.WeaponStats = nil
    item.ItemSpell = nil
    
    return item, nil
}

func ParseItems (items []models.Item) ([]models.Item, []error) {
    returnData := []models.Item{}
    for _, item := range items {

        // Fix stats
        err := json.Unmarshal(item.Stats, &item.StatsList)    

        // Fix weapon stats
        err = json.Unmarshal(item.WeaponStats, &item.WeaponInfo)    
        
        // Fix spells    
        err = json.Unmarshal(item.ItemSpell, &item.SpellInfo)
            
        if err != nil {
            panic(err)
        }
        
        item.Stats = nil
        item.WeaponStats = nil
        item.ItemSpell = nil
        
        returnData = append(returnData, item)
    }
    return returnData, nil
}
