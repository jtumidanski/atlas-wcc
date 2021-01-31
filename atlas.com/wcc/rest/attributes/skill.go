package attributes

type SkillDataContainer struct {
   data     dataSegment
   included dataSegment
}

type SkillData struct {
   Id         string          `json:"id"`
   Type       string          `json:"type"`
   Attributes SkillAttributes `json:"attributes"`
}

type SkillAttributes struct {
   Level       uint32 `json:"level"`
   MasterLevel uint32 `json:"masterLevel"`
   Expiration  int64  `json:"expiration"`
}

func (a *SkillDataContainer) UnmarshalJSON(data []byte) error {
   d, i, err := unmarshalRoot(data, mapperFunc(EmptySkillData))
   if err != nil {
      return err
   }

   a.data = d
   a.included = i
   return nil
}

func (a *SkillDataContainer) Data() *SkillData {
   if len(a.data) >= 1 {
      return a.data[0].(*SkillData)
   }
   return nil
}

func (a *SkillDataContainer) DataList() []SkillData {
   var r = make([]SkillData, 0)
   for _, x := range a.data {
      r = append(r, *x.(*SkillData))
   }
   return r
}

func EmptySkillData() interface{} {
   return &SkillData{}
}
