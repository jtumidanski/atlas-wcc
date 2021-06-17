package character

import "atlas-wcc/rest/response"

type MapCharacterDataContainer struct {
   data response.DataSegment
}

type MapCharacterData struct {
   Id         string                 `json:"id"`
   Type       string                 `json:"type"`
   Attributes MapCharacterAttributes `json:"attributes"`
}

type MapCharacterAttributes struct {
}

func (c *MapCharacterDataContainer) UnmarshalJSON(data []byte) error {
   d, _, err := response.UnmarshalRoot(data, response.MapperFunc(EmptyMapCharacterData))
   if err != nil {
      return err
   }
   c.data = d
   return nil
}

func (c *MapCharacterDataContainer) Data() *MapCharacterData {
   if len(c.data) >= 1 {
      return c.data[0].(*MapCharacterData)
   }
   return nil
}

func (c *MapCharacterDataContainer) DataList() []MapCharacterData {
   var r = make([]MapCharacterData, 0)
   for _, x := range c.data {
      r = append(r, *x.(*MapCharacterData))
   }
   return r
}

func EmptyMapCharacterData() interface{} {
   return &MapCharacterData{}
}
