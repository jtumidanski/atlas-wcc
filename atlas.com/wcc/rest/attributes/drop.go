package attributes

type DropDataContainer struct {
   data dataSegment
}

type DropData struct {
   Id         string         `json:"id"`
   Type       string         `json:"type"`
   Attributes DropAttributes `json:"attributes"`
}

type DropAttributes struct {
   WorldId         byte   `json:"worldId"`
   ChannelId       byte   `json:"channelId"`
   MapId           uint32 `json:"mapId"`
   ItemId          uint32 `json:"itemId"`
   Quantity        uint32 `json:"quantity"`
   Meso            uint32 `json:"meso"`
   DropType        byte   `json:"dropType"`
   DropX           int16  `json:"dropX"`
   DropY           int16  `json:"dropY"`
   OwnerId         uint32 `json:"ownerId"`
   OwnerPartyId    uint32 `json:"ownerPartyId"`
   DropTime        uint64 `json:"dropTime"`
   DropperUniqueId uint32 `json:"dropperUniqueId"`
   DropperX        int16  `json:"dropperX"`
   DropperY        int16  `json:"dropperY"`
   CharacterDrop   bool   `json:"playerDrop"`
   Mod             byte   `json:"mod"`
}

func (c *DropDataContainer) UnmarshalJSON(data []byte) error {
   d, _, err := unmarshalRoot(data, mapperFunc(EmptyDropData))
   if err != nil {
      return err
   }
   c.data = d
   return nil
}

func (c *DropDataContainer) Data() *DropData {
   if len(c.data) >= 1 {
      return c.data[0].(*DropData)
   }
   return nil
}

func (c *DropDataContainer) DataList() []DropData {
   var r = make([]DropData, 0)
   for _, x := range c.data {
      r = append(r, *x.(*DropData))
   }
   return r
}

func EmptyDropData() interface{} {
   return &DropData{}
}
