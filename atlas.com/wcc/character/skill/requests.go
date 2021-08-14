package skill

import (
	"atlas-wcc/rest/requests"
	"fmt"
	"github.com/sirupsen/logrus"
)

const (
	charactersServicePrefix string = "/ms/cos/"
	charactersService              = requests.BaseRequest + charactersServicePrefix
	charactersResource             = charactersService + "characters/"
	skillsByCharacter              = charactersResource + "%d/skills"
)

func requestForCharacter(l logrus.FieldLogger) func(characterId uint32) (*DataContainer, error) {
	return func(characterId uint32) (*DataContainer, error) {
		ar := &DataContainer{}
		err := requests.Get(l)(fmt.Sprintf(skillsByCharacter, characterId), ar)
		if err != nil {
			return nil, err
		}
		return ar, nil
	}
}
