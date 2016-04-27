package models

// Equipment is a thing..
// Probably deprecated, see docs/todo.md
type Equipment struct {
    DatabaseModel
	AverageItemLevel         int       `json:"averageItemLevel"`
    AverageItemLevelEquipped int       `json:"averageItemLevelEquipped"`
    Head                     Item      `json:"head"`
	Neck                     Item      `json:"neck"`
	Shoulder                 Item      `json:"shoulder"`
	Back                     Item      `json:"back"`
	Chest                    Item      `json:"chest"`
	Shirt                    Item      `json:"shirt"`
	Wrist                    Item      `json:"wrist"`
	Hands                    Item      `json:"hands"`
	Waist                    Item      `json:"waist"`
	Legs                     Item      `json:"legs"`
	Feet                     Item      `json:"feet"`
	Finger1                  Item      `json:"finger1"`
	Finger2                  Item      `json:"finger2"`
	Trinket1                 Item      `json:"trinket1"`
	Trinket2                 Item      `json:"trinket2"`
	MainHand                 Item      `json:"mainHand"`
	OffHand                  Item      `json:"offHand"`
}