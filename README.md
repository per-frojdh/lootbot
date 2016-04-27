## Aeon Lootbot

### Currently working on

### Maybe sort of finished
* Adding items to the lootlist
* Deleting items off the lootlist
* View all items in the lootlist
* Register a new user (double check)
* Import a character from battle.net
* Get all registered users
* Search
* Remove a character
* Show an item
* Show characters
* Show specific characters lootlist
* Delete a user (will delete everything)

### Unknown stuff
* Should also try to look into future data (not sure how, since I'm reliant on b.net API)

### Commands

| Command | Method | URI | Params | Authentication |
| --- | --- | --- | --- | --- |
|`!register` | POST | /api/public/register | - | No |
|`!lootlist item <id>` | GET | /api/items/:id | - | Yes |
|`!lootlist search <query>` | GET | /api/items/ | `?search=text` | Yes |
|`!lootlist users` | GET | /api/users/ | - | Yes |
|`!lootlist users <name>` | GET | /api/users/:name | - | Yes |
|`!lootlist users delete` | POST | /api/users/delete | - | Yes |
|`!lootlist <name>` | GET | /api/lootlist/:name | - | Yes |
|`!lootlist add <charname> <id>` | POST | /api/lootlist/add/:id | `name=<charname>` | Yes |
|`!lootlist remove <charname> <id>` | POST | /api/lootlist/delete/:id | `name=<charname>` | Yes |
|`!lootlist chars` | GET | /api/characters/ | - | Yes |
|`!lootlist import <realm> <name>` | POST | /api/characters/import | `name=<name> realm=<realm>` | Yes |
|`!lootlist chars remove <realm> <name>` | POST | /api/characters/delete | `name=<name> realm=<realm>` | Yes |



 
Register a user (we need to see what info we can use from Discord, we need a login, name, email (maybe password))
> `!register`

Adding items to a characters lootlist
> `!lootlist add <charname> <itemid>`

Deleting item off a given characters lootlist
> `!lootlist remove <name> <itemid>`

Search for available items
> `!lootlist search <itemname>`

View all items in a characters lootlist
> `!lootlist <name>`

Import character from battle.net
> `!lootlist char import <realm> <name>`

List all users (maybe need higher role?)
> `!lootlist users`

Remove a character from your user
> `!lootlist char remove <name>`

Show an item
> `!lootlist item <itemid>`

Show your characters
> `!lootlist char list`

Delete a user (will delete everything)
> `!lootlist users delete <name>`
