package consumers

import (
   "atlas-wcc/domain"
   "atlas-wcc/mapleSession"
   "atlas-wcc/processors"
   "atlas-wcc/registries"
   "atlas-wcc/rest/requests"
   "atlas-wcc/socket/response/writer"
   "context"
   "encoding/json"
   "fmt"
   "github.com/segmentio/kafka-go"
   "log"
   "os"
   "time"
)

type MapCharacter struct {
   l   *log.Logger
   ctx context.Context
}

func NewMapCharacter(l *log.Logger, ctx context.Context) *MapCharacter {
   return &MapCharacter{l, ctx}
}

func (mc *MapCharacter) Init(worldId byte, channelId byte) {
   t := requests.NewTopic(mc.l)
   td, err := t.GetTopic("TOPIC_MAP_CHARACTER_EVENT")
   if err != nil {
      mc.l.Fatal("[ERROR] Unable to retrieve topic for consumer.")
   }

   r := kafka.NewReader(kafka.ReaderConfig{
      Brokers: []string{os.Getenv("BOOTSTRAP_SERVERS")},
      Topic:   td.Attributes.Name,
      GroupID: fmt.Sprintf("World Channel Coordinator %d %d", worldId, channelId),
      MaxWait: 50 * time.Millisecond,
   })
   for {
      msg, err := r.ReadMessage(mc.ctx)
      if err != nil {
         panic("Could not successfully read message " + err.Error())
      }

      var event MapCharacterEvent
      err = json.Unmarshal(msg.Value, &event)
      if err != nil {
         mc.l.Println("Could not unmarshal event into event class ", msg.Value)
      } else {
         mc.processEvent(event)
      }
   }
}

func (mc *MapCharacter) processEvent(event MapCharacterEvent) {
   as := getSessionForCharacterId(event.CharacterId)
   if as == nil {
      return
   }

   if event.Type == "ENTER" {
      mc.enter(*as, event)
   } else if event.Type == "EXIT" {
      mc.exit(*as, event)
   }
}

func (mc *MapCharacter) enter(as mapleSession.MapleSession, event MapCharacterEvent) {
   cIds, err := processors.GetCharacterIdsInMap(event.WorldId, event.ChannelId, event.MapId)
   if err != nil {
      return
   }

   cm := make(map[uint32]*domain.Character)
   for _, cId := range cIds {
      c, err := processors.GetCharacterById(cId)
      if err != nil {
         //log something
      } else {
         cm[c.Attributes().Id()] = c
      }
   }

   // Spawn new character for other character.
   for k, v := range cm {
      if k != event.CharacterId {
         s := *getSessionForCharacterId(k)
         s.Announce(writer.WriteSpawnCharacter(*v, *cm[event.CharacterId], true))
      }
   }

   // Spawn other characters for incoming character.
   for k, v := range cm {
      if k != event.CharacterId {
         as.Announce(writer.WriteSpawnCharacter(*cm[event.CharacterId], *v, false))
      }
   }

   // Spawn NPCs for incoming character.
   ns, err := processors.GetNPCsInMap(event.MapId)
   if err != nil {
      return
   }
   for _, n := range ns {
      spawnNPCForSession(as, n)
   }

   // Spawn monsters for incoming character.
   ms, err := processors.GetMonstersInMap(event.WorldId, event.ChannelId, event.MapId)
   if err != nil {
      return
   }
   for _, m := range ms {
      spawnMonsterForSession(as, m)
   }

   // Spawn drops for incoming character.
   ds, err := processors.GetDropsInMap(event.WorldId, event.ChannelId, event.MapId)
   if err != nil {
      return
   }
   for _, d := range ds {
      spawnDropForSession(as, d)
   }
}

func spawnDropForSession(s mapleSession.MapleSession, d domain.Drop) {
   var a = uint32(0)
   if d.ItemId() != 0 {
      a = 0
   } else {
      a = d.Meso()
   }
   s.Announce(writer.WriteDropItemFromMapObject(d.UniqueId(), d.ItemId(), d.Meso(), a, d.DropperUniqueId(), d.DropType(), d.OwnerId(), d.OwnerPartyId(), s.CharacterId(), 0, d.DropTime(), d.DropX(), d.DropY(), d.DropperX(), d.DropperY(), d.CharacterDrop(), d.Mod()))
}

func spawnMonsterForSession(s mapleSession.MapleSession, m domain.Monster) {
   s.Announce(writer.WriteSpawnMonster(m, false))
}

func spawnNPCForSession(s mapleSession.MapleSession, n domain.NPC) {
   s.Announce(writer.WriteSpawnNPC(n))
   s.Announce(writer.WriteSpawnNPCController(n, true))
}

func getSessionsForThoseInMap(worldId byte, channelId byte, mapId uint32) ([]mapleSession.MapleSession, error) {
   cs, err := processors.GetCharacterIdsInMap(worldId, channelId, mapId)
   if err != nil {
      return nil, err
   }

   sl := getSessionsForCharacterIds(cs)
   return sl, nil
}

func (mc *MapCharacter) exit(as mapleSession.MapleSession, event MapCharacterEvent) {
   sl, err := getSessionsForThoseInMap(event.WorldId, event.ChannelId, event.MapId)
   if err != nil {
      return
   }
   for _, s := range sl {
      removeCharacterForSession(s, event.CharacterId)
   }

   ns, err := processors.GetNPCsInMap(event.MapId)
   if err != nil {
      return
   }
   for _, n := range ns {
      removeNpcForSession(as, n)
   }
}

func removeNpcForSession(as mapleSession.MapleSession, n domain.NPC) {
   as.Announce(writer.WriteRemoveNPCController(n.ObjectId()))
   as.Announce(writer.WriteRemoveNPC(n.ObjectId()))
}

func removeCharacterForSession(s mapleSession.MapleSession, characterId uint32) {
   s.Announce(writer.WriteRemoveCharacterFromMap(characterId))
}

func getSessionsForCharacterIds(cids []uint32) []mapleSession.MapleSession {
   sl := make([]mapleSession.MapleSession, 0)
   for _, s := range registries.GetSessionRegistry().GetAll() {
      if contains(cids, s.CharacterId()) {
         sl = append(sl, s)
      }
   }
   return sl
}

func getSessionForCharacterId(cid uint32) *mapleSession.MapleSession {
   for _, s := range registries.GetSessionRegistry().GetAll() {
      if cid == s.CharacterId() {
         return &s
      }
   }
   return nil
}

func contains(set []uint32, id uint32) bool {
   for _, s := range set {
      if s == id {
         return true
      }
   }
   return false
}

type MapCharacterEvent struct {
   WorldId     byte   `json:"worldId"`
   ChannelId   byte   `json:"channelId"`
   MapId       uint32 `json:"mapId"`
   CharacterId uint32 `json:"characterId"`
   Type        string `json:"type"`
}
