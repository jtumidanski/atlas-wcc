package session

import (
	"atlas-wcc/model"
	"errors"
)

// GetByCharacterId gets a session (if one exists) for the given characterId
func GetByCharacterId(characterId uint32) (Model, error) {
	sessions, err := getAllFiltered(CharacterIdFilter(characterId))
	if err != nil {
		return Model{}, err
	}
	if len(sessions) >= 1 {
		return sessions[0], nil
	}
	return Model{}, errors.New("not found")
}

// ForSessionByCharacterId executes a Operator if a session exists for the characterId
func ForSessionByCharacterId(characterId uint32, f model.Operator[Model]) {
	s, err := GetByCharacterId(characterId)
	if err != nil {
		return
	}
	f(s)
}

func ForEachByCharacterId(characterIds []uint32, f model.Operator[Model]) {
	for _, id := range characterIds {
		s, err := GetByCharacterId(id)
		if err != nil {
			return
		}
		f(s)
	}
}

// CharacterIdFilter a filter which yields true when the characterId matches the one in the session
func CharacterIdFilter(characterId uint32) model.Filter[Model] {
	return func(session Model) bool {
		return session.CharacterId() == characterId
	}
}

// ForEachGM executes a Operator for all sessions which correspond to GMs
func ForEachGM(f model.Operator[Model]) {
	ForGMs(model.ExecuteForEach(f))
}

// ForGMs executes a SliceOperator for all sessions which correspond to GMs
func ForGMs(f model.SliceOperator[Model]) {
	forSessions(GetGMs, f)
}

// executes a Operator for all sessions retrieved by a Getter
func forSessions(p model.SliceProvider[Model], f model.SliceOperator[Model]) {
	s, err := p()
	if err == nil {
		f(s)
	}
}

// GetGMs retrieves all sessions which correspond to GMs
func GetGMs() ([]Model, error) {
	return getAllFiltered(OnlyGMs())
}

// OnlyGMs a Filter which yields true when the session is a GM
func OnlyGMs() model.Filter[Model] {
	return func(session Model) bool {
		return session.GM()
	}
}

// retrieves all sessions which pass the provided Filter slice.
func getAllFiltered(filters ...model.Filter[Model]) ([]Model, error) {
	sessions := Registry().GetAll()

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
	return results, nil
}

func Announce(b []byte) func(s Model) error {
	return func(s Model) error {
		if l, ok := Registry().GetLock(s.SessionId()); ok {
			l.Lock()
			err := s.announceEncrypted(b)
			l.Unlock()
			return err
		}
		return errors.New("invalid session")
	}
}

func SetAccountId(accountId uint32) func(id uint32) Model {
	return func(id uint32) Model {
		s := Model{}
		var ok bool
		if s, ok = Registry().Get(id); ok {
			s = s.setAccountId(accountId)
			Registry().Update(s)
			return s
		}
		return s
	}
}

func SetWorldId(worldId byte) func(id uint32) Model {
	return func(id uint32) Model {
		s := Model{}
		var ok bool
		if s, ok = Registry().Get(id); ok {
			s = s.setWorldId(worldId)
			Registry().Update(s)
			return s
		}
		return s
	}
}

func SetChannelId(channelId byte) func(id uint32) Model {
	return func(id uint32) Model {
		s := Model{}
		var ok bool
		if s, ok = Registry().Get(id); ok {
			s = s.setChannelId(channelId)
			Registry().Update(s)
			return s
		}
		return s
	}
}

func SetCharacterId(characterId uint32) func(id uint32) Model {
	return func(id uint32) Model {
		s := Model{}
		var ok bool
		if s, ok = Registry().Get(id); ok {
			s = s.setCharacterId(characterId)
			Registry().Update(s)
			return s
		}
		return s
	}
}

func SetGm(gm bool) func(id uint32) Model {
	return func(id uint32) Model {
		s := Model{}
		var ok bool
		if s, ok = Registry().Get(id); ok {
			s = s.setGm(gm)
			Registry().Update(s)
			return s
		}
		return s
	}
}

func UpdateLastRequest() func(id uint32) Model {
	return func(id uint32) Model {
		s := Model{}
		var ok bool
		if s, ok = Registry().Get(id); ok {
			s = s.updateLastRequest()
			Registry().Update(s)
			return s
		}
		return s
	}
}
