package session

import (
	"atlas-wcc/model"
	"errors"
)

func AllModelProvider() ([]Model, error) {
	return getRegistry().GetAll(), nil
}

func ByIdModelProvider(sessionId uint32) model.Provider[Model] {
	return func() (Model, error) {
		s, found := getRegistry().Get(sessionId)
		if found {
			return s, nil
		}
		return Model{}, errors.New("not found")
	}
}

func ByCharacterIdModelProvider(characterId uint32) model.Provider[Model] {
	return model.SliceProviderToProviderAdapter[Model](AllModelProvider, CharacterIdPreciselyOneFilter(characterId))
}

func GetById(sessionId uint32) (Model, error) {
	return ByIdModelProvider(sessionId)()
}

// GetByCharacterId gets a session (if one exists) for the given characterId
func GetByCharacterId(characterId uint32) (Model, error) {
	return ByCharacterIdModelProvider(characterId)()
}

// IfPresentByCharacterId executes an Operator if a session exists for the characterId
func IfPresentByCharacterId(characterId uint32, f model.Operator[Model]) {
	model.IfPresent(ByCharacterIdModelProvider(characterId), f)
}

func ForEachByCharacterId(provider model.SliceProvider[uint32], f model.Operator[Model]) {
	model.ForEach(model.SliceMap[uint32, Model](provider, GetByCharacterId), f)
}

func IfPresentById(sessionId uint32, f model.Operator[Model]) {
	model.IfPresent[Model](ByIdModelProvider(sessionId), f)
}

func ForAll(f model.Operator[Model]) {
	model.ForEach(AllModelProvider, f)
}

func CharacterIdFilter(referenceId uint32) model.Filter[Model] {
	return func(model Model) bool {
		return model.CharacterId() == referenceId
	}
}

// CharacterIdPreciselyOneFilter a filter which yields true when the characterId matches the one in the session
func CharacterIdPreciselyOneFilter(characterId uint32) model.PreciselyOneFilter[Model] {
	return func(models []Model) (Model, error) {
		return model.First(model.FixedSliceProvider(models), CharacterIdFilter(characterId))
	}
}

// ForEachGM executes an Operator for all sessions which correspond to GMs
func ForEachGM(f model.Operator[Model]) {
	model.ForEach(OnlyGMModelProvider(), f)
}

func OnlyGMModelProvider() model.SliceProvider[Model] {
	return model.FilteredProvider(AllModelProvider, OnlyGMFilter)
}

// OnlyGMFilter a Filter which yields true when the session is a GM
func OnlyGMFilter(session Model) bool {
	return session.GM()
}

func Announce(s Model, bytes ...[]byte) error {
	return AnnounceOperator(bytes...)(s)
}

func AnnounceOperator(bytes ...[]byte) func(s Model) error {
	return func(s Model) error {
		if l, ok := getRegistry().GetLock(s.SessionId()); ok {
			l.Lock()
			for _, b := range bytes {
				err := s.announceEncrypted(b)
				if err != nil {
					l.Unlock()
					return err
				}
			}
			l.Unlock()
			return nil
		}
		return errors.New("invalid session")
	}
}

func SetAccountId(accountId uint32) func(id uint32) Model {
	return func(id uint32) Model {
		s := Model{}
		var ok bool
		if s, ok = getRegistry().Get(id); ok {
			s = s.setAccountId(accountId)
			getRegistry().Update(s)
			return s
		}
		return s
	}
}

func SetCharacterId(characterId uint32) func(id uint32) Model {
	return func(id uint32) Model {
		s := Model{}
		var ok bool
		if s, ok = getRegistry().Get(id); ok {
			s = s.setCharacterId(characterId)
			getRegistry().Update(s)
			return s
		}
		return s
	}
}

func SetGm(gm bool) func(id uint32) Model {
	return func(id uint32) Model {
		s := Model{}
		var ok bool
		if s, ok = getRegistry().Get(id); ok {
			s = s.setGm(gm)
			getRegistry().Update(s)
			return s
		}
		return s
	}
}

func UpdateLastRequest() func(id uint32) Model {
	return func(id uint32) Model {
		s := Model{}
		var ok bool
		if s, ok = getRegistry().Get(id); ok {
			s = s.updateLastRequest()
			getRegistry().Update(s)
			return s
		}
		return s
	}
}
