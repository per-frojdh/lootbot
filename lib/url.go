package lib

import (
    "fmt"
)

func BuildURL(realm string, character string, apiKey string) string {
    baseURL := "https://eu.api.battle.net/wow/character/"
    locale := "locale=en_GB"
    return fmt.Sprintf("%[1]s%[2]s/%[3]s?%[4]s&apikey=%[5]s", baseURL, realm, character, locale, apiKey)
}

func BuildURLWithFields(realm string, character string, apiKey string, fields string) string {
    baseURL := "https://eu.api.battle.net/wow/character/"
    locale := "locale=en_GB"
    field := fmt.Sprintf("fields=%[1]s", fields)
    return fmt.Sprintf("%[1]s%[2]s/%[3]s?%[4]s&%[5]s&apikey=%[6]s", baseURL, realm, character, field, locale, apiKey)
}