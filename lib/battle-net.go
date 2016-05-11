package lib

import (
    "github.com/parnurzeal/gorequest"
    models "github.com/per-frojdh/lootbot/models"
    "encoding/json"
    "fmt"
    "errors"
    "unicode"
)

type ActivityFeedContainer struct {
    Feed        []ActivityFeed  `json:"feed"`
}

type ActivityFeed struct{
    Type        string      `json:"type"`
    Timestamp        int      `json:"timestamp"`
    ItemID        int      `json:"itemid"`
    Context        string      `json:"context"`
}

func FetchCharacter(realm string, character string, apiKey string) (*models.Character, []error){
    request := gorequest.New()
    url := BuildURL(realm, character, apiKey)
    fmt.Println(url)
    resp, body, errs := request.Get(url).End()
    fmt.Println(resp.StatusCode)
    if resp.StatusCode != 200 {
        return nil, []error{errors.New("Failed to get character")};
    }
    if errs != nil {
        // Do something
        return nil, errs;
    }
    var returnedCharacter models.APICharacter
    err := json.Unmarshal([]byte(body), &returnedCharacter)
    if err != nil {
        return nil, []error{err};
    }    
    // TODO: Pandaren isn't set as a race
    return ParseCharacterValues(returnedCharacter), nil;
}

func ParseCharacterValues(char models.APICharacter) *models.Character {
    var dbCharacter models.Character
    dbCharacter.Faction = models.Faction[char.Faction]
    dbCharacter.Class = models.Class[char.Class]
    dbCharacter.Gender = models.Gender[char.Gender]
    dbCharacter.Race = models.Race[char.Race]
    dbCharacter.Level = char.Level
    dbCharacter.Name = char.Name
    dbCharacter.Realm = char.Realm
    dbCharacter.Thumbnail = char.Thumbnail
    dbCharacter.Battlegroup = char.Battlegroup
    return &dbCharacter
}

func FetchActivityFeed(realm string, character string, apiKey string) *[]ActivityFeed{
    request := gorequest.New()
    _, body, errs := request.Get(BuildURLWithFields(realm, character, apiKey, "feed")).End()
    if errs != nil {
        // Do something
    }
    
    var feed ActivityFeedContainer
    err := json.Unmarshal([]byte(body), &feed)
    if err != nil {
        // Do something
    }
    return ParseFeed(&feed)
}

func ParseFeed(feed *ActivityFeedContainer) *[]ActivityFeed{
    returnData := []ActivityFeed{}
    for _, event := range feed.Feed {
        if (event.Type == "LOOT") {
            returnData = append(returnData, event)        
        }
    }
    return &returnData
}

func CapitalizeString(s string) string {
    word := []rune(s)
    word[0] = unicode.ToUpper(word[0])
    return string(word)
}