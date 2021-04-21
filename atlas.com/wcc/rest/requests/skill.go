package requests

import (
	"atlas-wcc/rest/attributes"
	"fmt"
	"github.com/sirupsen/logrus"
)

const (
	skillsByCharacter = charactersResource + "%d/skills"
)

var Skill = func() *skill {
	return &skill{}
}

type skill struct {
	l logrus.FieldLogger
}

func (s *skill) GetForCharacter(characterId uint32) (*attributes.SkillDataContainer, error) {
	ar := &attributes.SkillDataContainer{}
	err := get(fmt.Sprintf(skillsByCharacter, characterId), ar)
	if err != nil {
		return nil, err
	}
	return ar, nil
}
