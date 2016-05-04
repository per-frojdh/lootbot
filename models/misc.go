package models

import (
    "time"
)

// Easy way to decode classes from API
var Class = map[int]string{
	0:  "None",
	1:  "Warrior",
	2:  "Paladin",
	3:  "Hunter",
	4:  "Rogue",
	5:  "Priest",
	6:  "Death Knight",
	7:  "Shaman",
	8:  "Mage",
	9:  "Warlock",
	10: "Monk",
	11: "Druid",
}

var Faction = map[int]string{
    0: "Alliance",
    1: "Horde",
}

var Gender = map[int]string{
    0: "Male",
    1: "Female",
}

var Race = map[int]string{   
    1: "Human",
    2: "Orc",
    3: "Dwarf",
    4: "Night Elf",
    5: "Undead",
    6: "Tauren",
    7: "Gnome",
    8: "Troll",
    9: "Goblin",
    10: "Blood Elf",
    11: "Draenei",
    22: "Worgen",
    24: "Pandaren",
    25: "Pandaren",
}

// ErrorMessages ...
var ErrorMessages = map[string]string {
    "RESOURCE_NOT_FOUND": "The requested resource was not found",
    "FAILED_CREATING_USER": "We could not register a new user for you",
    "AUTHENTICATION_FAILED": "Authentication failed",
    "FORBIDDEN": "Not authorized",
    "BAD_INPUT_PARAMETERS": "Bad input parameters, try again",
    "DATABASE_ERROR": "Can't connect to database, tell someone",
    "CHARACTER_EXIST": "That character already exists",
}

// Basic model all other models have referenced
type DatabaseModel struct {
    ID              uint      `gorm:"primary_key"             json:"-"`
    CreatedAt       time.Time `gorm:"column:created_at"       json:"-"`
    UpdatedAt       time.Time `gorm:"column:updated_at"       json:"-"`
    DeletedAt       *time.Time `gorm:"column:deleted_at"      json:"-"`
}

type AccessToken struct {
    ID              uint         `gorm:"primary_key"             json:"-"`
    Key             string       `gorm:"column:key"              json:"-"`
	Token           string       `gorm:"column:token"            json:"-"`
}