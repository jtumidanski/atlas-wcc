package requests

import (
	"atlas-wcc/rest/attributes"
	"fmt"
	"github.com/sirupsen/logrus"
)

const (
	skillsByCharacter = charactersResource + "%d/skills"
)

func GetForCharacter(l logrus.FieldLogger) func(characterId uint32) (*attributes.SkillDataContainer, error) {
	return func(characterId uint32) (*attributes.SkillDataContainer, error) {
		ar := &attributes.SkillDataContainer{}
		err := Get(l)(fmt.Sprintf(skillsByCharacter, characterId), ar)
		if err != nil {
			return nil, err
		}
		return ar, nil
	}
}
