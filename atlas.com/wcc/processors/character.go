package processors

import (
   "atlas-wcc/domain"
   "atlas-wcc/rest/attributes"
   "atlas-wcc/rest/requests"
   "errors"
   "strconv"
)

func makeCharacterAttributes(ca *attributes.CharacterAttributesData) *domain.CharacterAttributes {
   cid, err := strconv.ParseUint(ca.Id, 10, 32)
   if err != nil {
      return nil
   }
   att := ca.Attributes
   r := domain.NewCharacterAttributeBuilder().
      SetId(uint32(cid)).
      SetAccountId(att.AccountId).
      SetWorldId(att.WorldId).
      SetName(att.Name).
      SetGender(att.Gender).
      SetSkinColor(att.SkinColor).
      SetFace(att.Face).
      SetHair(att.Hair).
      SetLevel(att.Level).
      SetJobId(att.JobId).
      SetStrength(att.Strength).
      SetDexterity(att.Dexterity).
      SetIntelligence(att.Intelligence).
      SetLuck(att.Luck).
      SetHp(att.Hp).
      SetMaxHp(att.MaxHp).
      SetMp(att.Mp).
      SetMaxMp(att.MaxMp).
      SetAp(att.Ap).
      SetSp(att.Sp).
      SetExperience(att.Experience).
      SetFame(att.Fame).
      SetGachaponExperience(att.GachaponExperience).
      SetMapId(att.MapId).
      SetSpawnPoint(att.SpawnPoint).
      SetMeso(att.Meso).
      Build()
   return &r
}

func GetCharacterById(characterId uint32) (*domain.Character, error) {
   cs, err := requests.GetCharacterAttributesById(characterId)
   if err != nil {
      return nil, err
   }

   c, err := getCharacterForAttributes(cs.Data())
   if err != nil {
      return nil, err
   }
   return c, nil
}

func GetCharacterAttributesById(characterId uint32) (*domain.CharacterAttributes, error) {
   cs, err := requests.GetCharacterAttributesById(characterId)
   if err != nil {
      return nil, err
   }
   ca := makeCharacterAttributes(cs.Data())
   if ca == nil {
      return nil, errors.New("unable to make character attributes")
   }
   return ca, nil
}

func getCharacterForAttributes(data *attributes.CharacterAttributesData) (*domain.Character, error) {
   ca := makeCharacterAttributes(data)
   if ca == nil {
      return nil, errors.New("unable to make character attributes")
   }

   eq, err := getEquippedItemsForCharacter(ca.Id())
   if err != nil {
      return nil, err
   }

   ps, err := getPetsForCharacter()
   if err != nil {
      return nil, err
   }

   ss, err := getSkillsForCharacter(ca.Id())
   if err != nil {
      return nil, err
   }

   c := domain.NewCharacter(*ca, eq, ss, ps)
   return &c, nil
}

func getSkillsForCharacter(characterId uint32) ([]domain.Skill, error) {
   r, err := requests.GetForCharacter(characterId)
   if err != nil {
      return nil, err
   }

   ss := make([]domain.Skill, 0)
   for _, s := range r.DataList() {
      sid, err := strconv.ParseUint(s.Id, 10, 32)
      if err != nil {
         break
      }
      sr := domain.NewSkill(uint32(sid), s.Attributes.Level, s.Attributes.MasterLevel, s.Attributes.Expiration, false, false)
      ss = append(ss, sr)
   }
   return ss, nil
}

func getPetsForCharacter() ([]domain.Pet, error) {
   return make([]domain.Pet, 0), nil
}

func getEquippedItemsForCharacter(characterId uint32) ([]domain.EquippedItem, error) {
   r, err := requests.GetEquippedItemsForCharacter(characterId)
   if err != nil {
      return nil, err
   }

   eis := make([]domain.EquippedItem, 0)
   for _, e := range r.GetIncludedEquippedItems() {
      ea := r.GetEquipmentStatistics(e.Attributes.EquipmentId)
      if ea != nil {
         ei := domain.NewEquippedItemBuilder().
            SetItemId(ea.ItemId).
            SetSlot(e.Attributes.Slot).
            SetStrength(ea.Strength).
            SetDexterity(ea.Dexterity).
            SetIntelligence(ea.Intelligence).
            SetLuck(ea.Luck).
            SetHp(ea.Hp).
            SetMp(ea.Mp).
            SetWeaponAttack(ea.WeaponAttack).
            SetMagicAttack(ea.MagicAttack).
            SetWeaponDefense(ea.WeaponDefense).
            SetMagicDefense(ea.MagicDefense).
            SetAccuracy(ea.Accuracy).
            SetAvoidability(ea.Avoidability).
            SetHands(ea.Hands).
            SetSpeed(ea.Speed).
            SetJump(ea.Jump).
            SetSlots(ea.Slots).
            Build()
         eis = append(eis, ei)
      }
   }

   return eis, nil
}
