package processors

import (
	"atlas-wcc/mapleSession"
	"atlas-wcc/registries"
	"log"
)

type SessionOperator func(*log.Logger, mapleSession.MapleSession)
type SessionsOperator func(*log.Logger, []mapleSession.MapleSession)

func ExecuteForEachSession(f SessionOperator) SessionsOperator {
	return func(l *log.Logger, sessions []mapleSession.MapleSession) {
		for _, session := range sessions {
			f(l, session)
		}
	}
}

func GetSessionByCharacterId(characterId uint32) *mapleSession.MapleSession {
	for _, s := range registries.GetSessionRegistry().GetAll() {
		if characterId == s.CharacterId() {
			return &s
		}
	}
	return nil
}

func GetSessionsByCharacterIds(characterIds []uint32) []mapleSession.MapleSession {
	sl := make([]mapleSession.MapleSession, 0)
	for _, s := range registries.GetSessionRegistry().GetAll() {
		if contains(characterIds, s.CharacterId()) {
			sl = append(sl, s)
		}
	}
	return sl
}

func GetOtherSessionsInMap(worldId byte, channelId byte, mapId uint32, characterId uint32) ([]mapleSession.MapleSession, error) {
	all, err := GetCharacterIdsInMap(worldId, channelId, mapId)
	if err != nil {
		return nil, err
	}

	var cs []uint32
	for _, id := range all {
		if id != characterId {
			cs = append(cs, id)
		}
	}

	sl := GetSessionsByCharacterIds(cs)
	return sl, nil
}

func GetSessionsInMap(worldId byte, channelId byte, mapId uint32) ([]mapleSession.MapleSession, error) {
	cs, err := GetCharacterIdsInMap(worldId, channelId, mapId)
	if err != nil {
		return nil, err
	}

	sl := GetSessionsByCharacterIds(cs)
	return sl, nil
}

func ForSessionByCharacterId(l *log.Logger, characterId uint32, f SessionOperator) {
	s := GetSessionByCharacterId(characterId)
	if s == nil {
		return
	}
	f(l, *s)
	return
}

func ForEachOtherSessionInMap(l *log.Logger, worldId byte, channelId byte, characterId uint32, f SessionOperator) {
	ForOtherSessionsInMap(l, worldId, channelId, characterId, ExecuteForEachSession(f))
}

func ForOtherSessionsInMap(l *log.Logger, worldId byte, channelId byte, characterId uint32, f SessionsOperator) {
	c, err := GetCharacterAttributesById(characterId)
	if err != nil {
		return
	}
	sessions, err := GetOtherSessionsInMap(worldId, channelId, c.MapId(), characterId)
	if err != nil {
		return
	}
	f(l, sessions)
	return
}

func ForEachSessionInMap(l *log.Logger, worldId byte, channelId byte, mapId uint32, f SessionOperator) {
	ForSessionsInMap(l, worldId, channelId, mapId, ExecuteForEachSession(f))
}

func ForSessionsInMap(l *log.Logger, worldId byte, channelId byte, mapId uint32, f SessionsOperator) {
	sessions, err := GetSessionsInMap(worldId, channelId, mapId)
	if err != nil {
		return
	}
	f(l, sessions)
	return
}

func contains(set []uint32, id uint32) bool {
	for _, s := range set {
		if s == id {
			return true
		}
	}
	return false
}
