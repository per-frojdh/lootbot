package lib

import (
    models "github.com/per-frojdh/lootbot/models"
    "encoding/json"
)

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
