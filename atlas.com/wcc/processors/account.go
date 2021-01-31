package processors

import (
   "atlas-wcc/domain"
   "atlas-wcc/rest/attributes"
   "atlas-wcc/rest/requests"
   "strconv"
)

func GetAccountByName(name string) (*domain.Account, error) {
   resp, err := requests.GetAccountByName(name)
   if err != nil {
      return nil, err
   }

   d := resp.Data()
   aid, err := strconv.ParseUint(d.Id, 10, 32)
   if err != nil {
      return nil, err
   }

   a := makeAccount(uint32(aid), d.Attributes)
   return &a, nil
}

func GetAccountById(id uint32) (*domain.Account, error) {
   resp, err := requests.GetAccountById(id)
   if err != nil {
      return nil, err
   }

   d := resp.Data()
   aid, err := strconv.ParseUint(d.Id, 10, 32)
   if err != nil {
      return nil, err
   }

   a := makeAccount(uint32(aid), d.Attributes)
   return &a, nil
}

func IsLoggedIn(id uint32) bool {
   a, err := GetAccountById(id)
   if err != nil {
      return false
   } else if a.LoggedIn() != 0 {
      return true
   } else {
      return false
   }
}

func makeAccount(id uint32, att attributes.AccountAttributes) domain.Account {
   return domain.NewAccountBuilder().
      SetId(id).
      SetPassword(att.Password).
      SetPin(att.Pin).
      SetPic(att.Pic).
      SetLoggedIn(att.LoggedIn).
      SetLastLogin(att.LastLogin).
      SetGender(att.Gender).
      SetBanned(att.Banned).
      SetTos(att.TOS).
      SetLanguage(att.Language).
      SetCountry(att.Country).
      SetCharacterSlots(att.CharacterSlots).
      Build()
}
