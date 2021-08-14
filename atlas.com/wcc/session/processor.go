package session

import (
	"atlas-wcc/character"
	"github.com/sirupsen/logrus"
)

// Operator function which performs an operation on a single session
type Operator func(*Model)

// SliceOperator function which performs a operation on a slice of sessions
type SliceOperator func([]*Model)

// Filter function which dictates whether a session should be considered or not
type Filter func(*Model) bool

// Getter function which retrieves a slice of sessions
type Getter func() []*Model

// ExecuteForEach executes a Operator over a slice of sessions
func ExecuteForEach(f Operator) SliceOperator {
	return func(sessions []*Model) {
		for _, session := range sessions {
			f(session)
		}
	}
}

// GetByCharacterId gets a session (if one exists) for the given characterId
func GetByCharacterId(characterId uint32) *Model {
	sessions := getAllFiltered(CharacterIdFilter(characterId))
	if len(sessions) >= 1 {
		return sessions[0]
	}
	return nil
}

// ForSessionByCharacterId executes a Operator if a session exists for the characterId
func ForSessionByCharacterId(characterId uint32, f Operator) {
	s := GetByCharacterId(characterId)
	if s == nil {
		return
	}
	f(s)
	return
}

// GetOtherInMap a Getter which will retrieve all sessions for characters in the given map, not identified by the supplied characterId
func GetOtherInMap(l logrus.FieldLogger) func(worldId byte, channelId byte, mapId uint32, characterId uint32) Getter {
	return func(worldId byte, channelId byte, mapId uint32, characterId uint32) Getter {
		return func() []*Model {
			cs, err := character.GetCharacterIdsInMap(l)(worldId, channelId, mapId)
			if err != nil {
				return nil
			}
			if len(cs) <= 0 {
				return nil
			}
			return getAllFiltered(CharacterIdInFilter(cs), CharacterIdNotFilter(characterId))
		}
	}
}

// CharacterIdFilter a filter which yields true when the characterId matches the one in the session
func CharacterIdFilter(characterId uint32) Filter {
	return func(session *Model) bool {
		return session.CharacterId() == characterId
	}
}

// CharacterIdNotFilter a filter which yields true when the characterId does not match the one in the session
func CharacterIdNotFilter(characterId uint32) Filter {
	return func(session *Model) bool {
		return session.CharacterId() != characterId
	}
}

// CharacterIdInFilter a filter which yields true when the characterId of the session is in the slice of provided characterIds
func CharacterIdInFilter(validIds []uint32) Filter {
	return func(session *Model) bool {
		return contains(validIds, session.CharacterId())
	}
}

// GetInMap a Getter which retrieve all sessions which reside in the identified map
func GetInMap(l logrus.FieldLogger) func(worldId byte, channelId byte, mapId uint32) Getter {
	return func(worldId byte, channelId byte, mapId uint32) Getter {
		return func() []*Model {
			cs, err := character.GetCharacterIdsInMap(l)(worldId, channelId, mapId)
			if err != nil {
				return nil
			}
			if len(cs) <= 0 {
				return nil
			}
			return getAllFiltered(CharacterIdInFilter(cs))
		}
	}
}

// ForEachOtherInMap executes a Operator for all sessions in the identified map, aside from the session of the provided characterId
func ForEachOtherInMap(l logrus.FieldLogger) func(worldId byte, channelId byte, characterId uint32, f Operator) {
	return func(worldId byte, channelId byte, characterId uint32, f Operator) {
	ForOtherInMap(l)(worldId, channelId, characterId, ExecuteForEach(f))
	}
}

// ForOtherInMap executes a SliceOperator for all sessions in the identified map, aside from the session of the provided characterId
func ForOtherInMap(l logrus.FieldLogger) func(worldId byte, channelId byte, characterId uint32, f SliceOperator) {
	return func(worldId byte, channelId byte, characterId uint32, f SliceOperator) {
		c, err := character.GetCharacterAttributesById(l)(characterId)
		if err != nil {
			return
		}
		forSessions(GetOtherInMap(l)(worldId, channelId, c.MapId(), characterId), f)
	}
}

// ForEachInMap executes a Operator for all sessions in the identified map
func ForEachInMap(l logrus.FieldLogger) func(worldId byte, channelId byte, mapId uint32, f Operator) {
	return func(worldId byte, channelId byte, mapId uint32, f Operator) {
		ForSessionsInMap(l)(worldId, channelId, mapId, ExecuteForEach(f))
	}
}

// ForSessionsInMap executes a SliceOperator for all sessions in the identified map
func ForSessionsInMap(l logrus.FieldLogger) func(worldId byte, channelId byte, mapId uint32, f SliceOperator) {
	return func(worldId byte, channelId byte, mapId uint32, f SliceOperator) {
		forSessions(GetInMap(l)(worldId, channelId, mapId), f)
	}
}

// ForEachGM executes a Operator for all sessions which correspond to GMs
func ForEachGM(f Operator) {
	ForGMs(ExecuteForEach(f))
}

// ForGMs executes a SliceOperator for all sessions which correspond to GMs
func ForGMs(f SliceOperator) {
	forSessions(GetGMs, f)
}

// executes a Operator for all sessions retrieved by a Getter
func forSessions(getter Getter, f SliceOperator) {
	f(getter())
}

// GetGMs retrieves all sessions which correspond to GMs
func GetGMs() []*Model {
	return getAllFiltered(OnlyGMs())
}

// OnlyGMs a Filter which yields true when the session is a GM
func OnlyGMs() Filter {
	return func(session *Model) bool {
		return session.GM()
	}
}

// retrieves all sessions which pass the provided Filter slice.
func getAllFiltered(filters ...Filter) []*Model {
	sessions := GetRegistry().GetAll()

	var results []*Model
	for _, session := range sessions {
		if len(filters) == 0 {
			results = append(results, session)
		} else {
			good := true
			for _, filter := range filters {
				if !filter(session) {
					good = false
					break
				}
			}
			if good {
				results = append(results, session)
			}
		}
	}
	return results
}

func contains(set []uint32, id uint32) bool {
	for _, s := range set {
		if s == id {
			return true
		}
	}
	return false
}
