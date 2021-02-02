package requests

import (
	"atlas-wcc/rest/attributes"
	"fmt"
	"log"
)

const (
	skillsByCharacter = charactersResource + "%d/skills"
)

var Skill = func() *skill {
	return &skill{}
}

type skill struct {
	l *log.Logger
}

func (s *skill) GetForCharacter(characterId uint32) (*attributes.SkillDataContainer, error) {
	ar := &attributes.SkillDataContainer{}
	err := get(fmt.Sprintf(skillsByCharacter, characterId), ar)
	if err != nil {
		return nil, err
	}
	return ar, nil
}
