package requests

import (
   "atlas-wcc/rest/attributes"
   "fmt"
)

const (
   SkillsByCharacter = CharactersResource + "%d/skills"
)

func GetForCharacter(characterId uint32) (*attributes.SkillDataContainer, error) {
   ar := &attributes.SkillDataContainer{}
   err := Get(fmt.Sprintf(SkillsByCharacter, characterId), ar)
   if err != nil {
      return nil, err
   }
   return ar, nil
}
