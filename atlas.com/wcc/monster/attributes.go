package monster

import "atlas-wcc/rest/response"

type MonsterDataContainer struct {
   data response.DataSegment
}

type MonsterData struct {
   Id         string            `json:"id"`
   Type       string            `json:"type"`
   Attributes MonsterAttributes `json:"attributes"`
}

type DamageEntry struct {
   CharacterId uint32 `json:"characterId"`
   Damage      int64  `json:"damage"`
}

type MonsterAttributes struct {
   WorldId            byte          `json:"worldId"`
   ChannelId          byte          `json:"channelId"`
   MapId              uint32        `json:"mapId"`
   MonsterId          uint32        `json:"monsterId"`
   ControlCharacterId uint32        `json:"controlCharacterId"`
   X                  int16         `json:"x"`
   Y                  int16         `json:"y"`
   FH                 int16        `json:"fh"`
   Stance             byte          `json:"stance"`
   Team               int8          `json:"team"`
   HP                 uint32        `json:"hp"`
   DamageEntries      []DamageEntry `json:"damageEntries"`
}

func (c *MonsterDataContainer) UnmarshalJSON(data []byte) error {
   d, _, err := response.UnmarshalRoot(data, response.MapperFunc(EmptyMonsterData))
   if err != nil {
      return err
   }
   c.data = d
   return nil
}

func (c *MonsterDataContainer) Data() *MonsterData {
   if len(c.data) >= 1 {
      return c.data[0].(*MonsterData)
   }
   return nil
}

func (c *MonsterDataContainer) DataList() []MonsterData {
   var r = make([]MonsterData, 0)
   for _, x := range c.data {
      r = append(r, *x.(*MonsterData))
   }
   return r
}

func EmptyMonsterData() interface{} {
   return &MonsterData{}
}
