package models

// Guild ...
type Guild struct {
    DatabaseModel
    Name        string      `gorm:"column:guildname"    json:"guildname"`
    Realm       string      `gorm:"column:realm"        json:"realm"`
    Faction     string      `gorm:"column:faction"      json:"faction"`
    Battlegroup string      `gorm:"column:battlegroup"  json:"battlegroup"`
    Status      string      `gorm:"column:status"       json:"status"`
    Users       []User                                `json:"users"`
}

// RequestAccess ...
type RequestAccess struct {
    DatabaseModel
    Login       string      `gorm:"column:login"        json:"login"`
    Guild       string      `gorm:"column:guildname"    json:"guildname"`
    Password    string      `gorm:"column:password"     json:"password"`
    Realm       string      `gorm:"column:realm"        json:"realm"`
    Email       string      `gorm:"column:email"        json:"email"`
    Status      string      `gorm:"column:status"       json:"status"`
    Faction     string      `gorm:"column:faction"      json:"faction"`
    Battlegroup string      `gorm:"column:battlegroup"  json:"battlegroup"`
    Hash        string      `gorm:"column:hash"         json:"hash"`
}

func (r RequestAccess) TableName() string {
    return "guild_request"
}