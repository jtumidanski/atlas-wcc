package _map

import "atlas-wcc/rest/response"

type CharacterDataContainer struct {
   data response.DataSegment
}

type CharacterDataBody struct {
   Id         string                 `json:"id"`
   Type       string              `json:"type"`
   Attributes CharacterAttributes `json:"attributes"`
}

type CharacterAttributes struct {
}

func (c *CharacterDataContainer) UnmarshalJSON(data []byte) error {
   d, _, err := response.UnmarshalRoot(data, response.MapperFunc(EmptyMapCharacterData))
   if err != nil {
      return err
   }
   c.data = d
   return nil
}

func (c *CharacterDataContainer) Data() *CharacterDataBody {
   if len(c.data) >= 1 {
      return c.data[0].(*CharacterDataBody)
   }
   return nil
}

func (c *CharacterDataContainer) DataList() []CharacterDataBody {
   var r = make([]CharacterDataBody, 0)
   for _, x := range c.data {
      r = append(r, *x.(*CharacterDataBody))
   }
   return r
}

func EmptyMapCharacterData() interface{} {
   return &CharacterDataBody{}
}
