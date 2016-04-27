package models

// Item is a thing..
type Item struct {
    DatabaseModel
    ItemID          int       `gorm:"column:itemid"          json:"itemid"`
    Quality         int       `gorm:"column:quality"         json:"quality"`
    ItemLevel       int       `gorm:"column:itemlevel"       json:"itemLevel"`
    Armor           int       `gorm:"column:armor"           json:"armor,omitempty"`
    Context         string    `gorm:"column:context"         json:"context"`
    DisplayInfoID   string    `gorm:"column:displayinfoid"   json:"displayinfoid,omitempty"`
    ItemClass       string    `gorm:"column:itemclass"       json:"itemClass"`
    Description     string    `gorm:"column:description"     json:"description,omitempty"`
    ItemSubClass    string    `gorm:"column:itemsubclass"    json:"itemSubClass"`
    Name            string    `gorm:"column:name"            json:"name"`
    Icon            string    `gorm:"column:icon"            json:"icon"`
    RequiredLevel   int       `gorm:"column:requiredlevel"   json:"requiredlevel"`
    InventoryType   string    `gorm:"column:inventorytype"   json:"inventorytype"`
    Stats           []byte    `gorm:"column:stats"           json:",omitempty"`    
    WeaponStats     []byte    `gorm:"column:weaponinfo"      json:",omitempty"`
    ItemSpell       []byte    `gorm:"column:itemspells"      json:",omitempty"`
    Classes*        []byte    `gorm:"column:allowableclasses"   json:"allowableclasses,omitempty"`
    StatsList       []ItemInfo                              `json:"stats"`    
    WeaponInfo*     WeaponInfo                              `json:"weaponstats,omitempty"`
    SpellInfo       []SpellInfo                             `json:"itemspell"`
    AllowClasses    []string                                `json:"classes,omitempty"`
}

// ItemInfo ...
type ItemInfo struct {
    Stat        string      `json: "stat"`
    Amount      int         `json: "amount"`
}

// WeaponInfo ...
type WeaponInfo struct {
    Damage          WeaponDamage    `json:"damage"`
    WeaponSpeed     float64         `json:"weaponspeed"`
    DPS             float64         `json:"dps"`
}

// WeaponDamage ...
type WeaponDamage struct {
    Min             int             `json:"min"`
    Max             int             `json:"max"`
    ExactMin        float64         `json:"exactmin"`
    ExactMax        float64         `json:"exactmax"`
}

// SpellInfo ...
type SpellInfo struct {
    SpellID     int     `json:"spellId"`
    Spell       Spell   `json:"spell"`
    Charges     int     `json:"nCharges"`
    Consumable  bool    `json:"consumable"`
    CategoryID  int     `json:"categoryId"`
    Trigger     string  `json:"trigger"`
}

// Spell ...
type Spell struct {
    ID              int        `json:"id"`
    Name            string     `json:"name"`
    Icon            string     `json:"icon"`
    Description     string     `json:"description"`
    CastTime        string     `json:"casttime"`  
}