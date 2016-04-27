package models

// User is a thing..
type User struct {
    DatabaseModel
	Login          string       `gorm:"column:login"           json:"login"`
	Name           string       `gorm:"column:name"            json:"name"`
	Password       string       `gorm:"column:passphrase"      json:"-"` 
	Email          string       `gorm:"column:email"           json:"email"`
	GuildID        uint         `gorm:"column:guildid"         json:"guild"`
    SecretQuestion string       `gorm:"column:secretquestion"  json:"-"`
    SecretAnswer   string       `gorm:"column:secretanswer"    json:"-"`
    Token          string       `gorm:"column:token"           json:"-"`
    Role           uint         `gorm:"column:role"            json:"-"`
    Characters     []Character  `gorm:"ForeignKey:userid"`
}