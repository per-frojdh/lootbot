package models

// User is a thing..
type Character struct {
    DatabaseModel
	Name              string 	`gorm:"column:name"				json:"name"`
	Realm             string 	`gorm:"column:realm"			json:"realm"`
	Class             string 	`gorm:"column:class"			json:"class"`
	Race              string 	`gorm:"column:race"		    	json:"race"`
	UserID			  uint    	`gorm:"column:user_id"          	json:"-"`
	Gender            string 	`gorm:"column:gender"			json:"gender"`
	Level             int    	`gorm:"column:level"			json:"level"`
	Thumbnail         string 	`gorm:"column:thumbnail"		json:"thumbnail"`
	Faction           string 	`gorm:"column:faction"			json:"faction"`
	Battlegroup       string 	`gorm:"column:battlegroup"		json:"battlegroup"`
	Lootlist		  []Item    `gorm:"many2many:lootlist"`
}



// To make things easier we're gonna create a db character and a api character

type APICharacter struct {
    Name              string `json:"name"`
	Realm             string `json:"realm"`
	Class             int    `json:"class"`
	Race              int    `json:"race"`
	Gender            int    `json:"gender"`
	Level             int    `json:"level"`
	Thumbnail         string `json:"thumbnail"`
	Faction           int    `json:"faction"`
	Battlegroup       string `json:"battlegroup"`
    Owner             int    `json:"owner"`
}