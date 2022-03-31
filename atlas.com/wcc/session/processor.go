package session

import (
	"atlas-wcc/model"
	"errors"
)

func AllModelProvider() ([]Model, error) {
	return Registry().GetAll(), nil
}

func ByIdModelProvider(id uint32) model.Provider[Model] {
	return model.SliceProviderToProviderAdapter[Model](AllModelProvider, CharacterIdPreciselyOneFilter(id))
}

// GetByCharacterId gets a session (if one exists) for the given characterId
func GetByCharacterId(characterId uint32) (Model, error) {
	return ByIdModelProvider(characterId)()
}

// ForSessionByCharacterId executes an Operator if a session exists for the characterId
func ForSessionByCharacterId(characterId uint32, f model.Operator[Model]) {
	model.IfPresent(ByIdModelProvider(characterId), f)
}

func ForEachByCharacterId(provider model.SliceProvider[uint32], f model.Operator[Model]) {
	model.ForEach(model.SliceMap[uint32, Model](provider, GetByCharacterId), f)
}

func CharacterIdFilter(referenceId uint32) model.Filter[Model] {
	return func(model Model) bool {
		return model.CharacterId() == referenceId
	}
}

// CharacterIdPreciselyOneFilter a filter which yields true when the characterId matches the one in the session
func CharacterIdPreciselyOneFilter(characterId uint32) model.PreciselyOneFilter[Model] {
	return func(models []Model) (Model, error) {
		return model.First(models, CharacterIdFilter(characterId))
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
