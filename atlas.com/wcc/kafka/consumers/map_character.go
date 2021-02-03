package consumers

import (
	"atlas-wcc/domain"
	"atlas-wcc/mapleSession"
	"atlas-wcc/processors"
	"atlas-wcc/socket/response/writer"
	"log"
)

type mapCharacterEvent struct {
	WorldId     byte   `json:"worldId"`
	ChannelId   byte   `json:"channelId"`
	MapId       uint32 `json:"mapId"`
	CharacterId uint32 `json:"characterId"`
	Type        string `json:"type"`
}

func MapCharacterEventCreator() EmptyEventCreator {
	return func() interface{} {
		return &mapCharacterEvent{}
	}
}

type MapCharacterHandler struct {
}

func (h MapCharacterHandler) topicToken() string {
	return "TOPIC_MAP_CHARACTER_EVENT"
}

func HandleMapCharacterEvent() ChannelEventProcessor {
	return func(l *log.Logger, wid byte, cid byte, e interface{}) {
		if event, ok := e.(*mapCharacterEvent); ok {
			as := processors.GetSessionByCharacterId(event.CharacterId)
			if as == nil {
				l.Printf("[ERROR] unable to locate session for character %d.", event.CharacterId)
				return
			}

			l.Printf("[INFO] processing MapCharacterEvent of type %s", event.Type)
			if event.Type == "ENTER" {
				enterMap(l, *as, *event)
			} else if event.Type == "EXIT" {
				exitMap(l, *as, *event)
			}

		} else {
			l.Printf("[ERROR] unable to cast event provided to handler [HandleMapCharacterEvent]")
		}
	}
}

func enterMap(l *log.Logger, as mapleSession.MapleSession, event mapCharacterEvent) {
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
			s := *processors.GetSessionByCharacterId(k)
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
		l.Printf("[INFO] spawning monster %d type %d for character %d", m.UniqueId(), m.MonsterId(), as.CharacterId())
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

func exitMap(_ *log.Logger, as mapleSession.MapleSession, event mapCharacterEvent) {
	sl, err := processors.GetSessionsInMap(event.WorldId, event.ChannelId, event.MapId)
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



