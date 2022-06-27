package properties

import (
	"atlas-wcc/rest/requests"
	"fmt"
)

const (
	charactersServicePrefix string = "/ms/cos/"
	charactersService              = requests.BaseRequest + charactersServicePrefix
	charactersResource             = charactersService + "characters/"
	charactersById                 = charactersResource + "%d"
	charactersByName               = charactersService + "characters?name=%s"
)

func requestById(id uint32) requests.Request[attributes] {
	return requests.MakeGetRequest[attributes](fmt.Sprintf(charactersById, id))
}

func requestByName(name string) requests.Request[attributes] {
	return requests.MakeGetRequest[attributes](fmt.Sprintf(charactersByName, name))
}
