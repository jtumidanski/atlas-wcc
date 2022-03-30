package session

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
}

func ForEachByCharacterId(characterIds []uint32, f Operator) {
	for _, id := range characterIds {
		s := GetByCharacterId(id)
		if s != nil {
			f(s)
		}
	}
}

// CharacterIdFilter a filter which yields true when the characterId matches the one in the session
func CharacterIdFilter(characterId uint32) Filter {
	return func(session *Model) bool {
		return session.CharacterId() == characterId
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
	sessions := Registry().GetAll()

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
