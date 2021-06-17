package session

import (
	"atlas-wcc/character"
)

// function which performs an operation on a single session
type SessionOperator func(Model)

// function which performs a operation on a slice of sessions
type SessionsOperator func([]Model)

// function which dictates whether a session should be considered or not
type SessionFilter func(session Model) bool

// function which retrieves a slice of sessions
type SessionGetter func() []Model

// executes a SessionOperator over a slice of sessions
func ExecuteForEachSession(f SessionOperator) SessionsOperator {
	return func(sessions []Model) {
		for _, session := range sessions {
			f(session)
		}
	}
}

// gets a session (if one exists) for the given characterId
func GetSessionByCharacterId(characterId uint32) *Model {
	sessions := getFilteredSessions(CharacterIdFilter(characterId))
	if len(sessions) >= 1 {
		return &sessions[0]
	}
	return nil
}

// executes a SessionOperator if a session exists for the characterId
func ForSessionByCharacterId(characterId uint32, f SessionOperator) {
	s := GetSessionByCharacterId(characterId)
	if s == nil {
		return
	}
	f(*s)
	return
}

// a SessionGetter which will retrieve all sessions for characters in the given map, not identified by the supplied characterId
func GetOtherSessionsInMap(worldId byte, channelId byte, mapId uint32, characterId uint32) SessionGetter {
	return func() []Model {
		cs, err := character.GetCharacterIdsInMap(worldId, channelId, mapId)
		if err != nil {
			return nil
		}
		if len(cs) <= 0 {
			return nil
		}
		return getFilteredSessions(CharacterIdInFilter(cs), CharacterIdNotFilter(characterId))
	}
}

// a filter which yields true when the characterId matches the one in the session
func CharacterIdFilter(characterId uint32) SessionFilter {
	return func(session Model) bool {
		return session.CharacterId() == characterId
	}
}

// a filter which yields true when the characterId does not match the one in the session
func CharacterIdNotFilter(characterId uint32) SessionFilter {
	return func(session Model) bool {
		return session.CharacterId() != characterId
	}
}

// a filter which yields true when the characterId of the session is in the slice of provided characterIds
func CharacterIdInFilter(validIds []uint32) SessionFilter {
	return func(session Model) bool {
		return contains(validIds, session.CharacterId())
	}
}

// a SessionGetter which which retrieve all sessions which reside in the identified map
func GetSessionsInMap(worldId byte, channelId byte, mapId uint32) SessionGetter {
	return func() []Model {
		cs, err := character.GetCharacterIdsInMap(worldId, channelId, mapId)
		if err != nil {
			return nil
		}
		if len(cs) <= 0 {
			return nil
		}
		return getFilteredSessions(CharacterIdInFilter(cs))
	}
}

// executes a SessionOperator for all sessions in the identified map, aside from the session of the provided characterId
func ForEachOtherSessionInMap(worldId byte, channelId byte, characterId uint32, f SessionOperator) {
	ForOtherSessionsInMap(worldId, channelId, characterId, ExecuteForEachSession(f))
}

// executes a SessionsOperator for all sessions in the identified map, aside from the session of the provided characterId
func ForOtherSessionsInMap(worldId byte, channelId byte, characterId uint32, f SessionsOperator) {
	c, err := character.GetCharacterAttributesById(characterId)
	if err != nil {
		return
	}
	forSessions(GetOtherSessionsInMap(worldId, channelId, c.MapId(), characterId), f)
}

// executes a SessionOperator for all sessions in the identified map
func ForEachSessionInMap(worldId byte, channelId byte, mapId uint32, f SessionOperator) {
	ForSessionsInMap(worldId, channelId, mapId, ExecuteForEachSession(f))
}

// executes a SessionsOperator for all sessions in the identified map
func ForSessionsInMap(worldId byte, channelId byte, mapId uint32, f SessionsOperator) {
	forSessions(GetSessionsInMap(worldId, channelId, mapId), f)
}

// executes a SessionOperator for all sessions which correspond to GMs
func ForEachGMSession(f SessionOperator) {
	ForGMSessions(ExecuteForEachSession(f))
}

// executes a SessionsOperator for all sessions which correspond to GMs
func ForGMSessions(f SessionsOperator) {
	forSessions(GetGMSessions, f)
}

// executes a SessionOperator for all sessions retrieved by a SessionGetter
func forSessions(getter SessionGetter, f SessionsOperator) {
	f(getter())
}

// retrieves all sessions which correspond to GMs
func GetGMSessions() []Model {
	return getFilteredSessions(OnlyGMs())
}

// a SessionFilter which yields true when the session is a GM
func OnlyGMs() SessionFilter {
	return func(session Model) bool {
		return session.GM()
	}
}

func ForEachSession(f SessionOperator) {
	ForAllSessions(ExecuteForEachSession(f))
}

func ForAllSessions(f SessionsOperator) {
	getAll := func() []Model {
		return getFilteredSessions()
	}
	forSessions(getAll, f)
}

// retrieves all sessions which pass the provided SessionFilter slice.
func getFilteredSessions(filters ...SessionFilter) []Model {
	sessions := GetRegistry().GetAll()

	var results []Model
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
