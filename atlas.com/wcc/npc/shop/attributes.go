package shop

type attributes struct {
	NPC   uint32           `json:"npc"`
	Items []itemAttributes `json:"items"`
}

type itemAttributes struct {
	ItemId   uint32 `json:"itemId"`
	Price    uint32 `json:"price"`
	Pitch    uint32 `json:"pitch"`
	Position uint32 `json:"position"`
}
