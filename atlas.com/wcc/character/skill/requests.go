package skill

import (
	"atlas-wcc/rest/requests"
	"fmt"
)

const (
	charactersServicePrefix string = "/ms/cos/"
	charactersService              = requests.BaseRequest + charactersServicePrefix
	charactersResource             = charactersService + "characters/"
	skillsByCharacter              = charactersResource + "%d/skills"
)

func requestForCharacter(characterId uint32) requests.Request[attributes] {
	return requests.MakeGetRequest[attributes](fmt.Sprintf(skillsByCharacter, characterId))
}
