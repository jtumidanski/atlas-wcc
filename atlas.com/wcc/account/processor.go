package account

import (
	"github.com/opentracing/opentracing-go"
	"github.com/sirupsen/logrus"
	"strconv"
)

func GetById(l logrus.FieldLogger, span opentracing.Span) func(id uint32) (*Model, error) {
	return func(id uint32) (*Model, error) {
		resp, err := requestById(l, span)(id)
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
}

func IsLoggedIn(l logrus.FieldLogger, span opentracing.Span) func(id uint32) bool {
	return func(id uint32) bool {
		a, err := GetById(l, span)(id)
		if err != nil {
			return false
		} else if a.LoggedIn() != 0 {
			return true
		} else {
			return false
		}
	}
}

func makeAccount(id uint32, att attributes) Model {
	return NewAccountBuilder().
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
