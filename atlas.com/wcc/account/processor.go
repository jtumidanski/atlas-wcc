package account

import (
	"atlas-wcc/rest/requests"
	"github.com/opentracing/opentracing-go"
	"github.com/sirupsen/logrus"
	"strconv"
)

type ModelProvider func() (*Model, error)

func requestModelProvider(l logrus.FieldLogger, span opentracing.Span) func(r requests.Request[attributes]) ModelProvider {
	return func(r requests.Request[attributes]) ModelProvider {
		return func() (*Model, error) {
			resp, err := r(l, span)
			if err != nil {
				return nil, err
			}

			p, err := makeModel(resp.Data())
			if err != nil {
				return nil, err
			}
			return p, nil
		}
	}
}

func ByIdModelProvider(l logrus.FieldLogger, span opentracing.Span) func(id uint32) ModelProvider {
	return func(id uint32) ModelProvider {
		return requestModelProvider(l, span)(requestAccountById(id))
	}
}

func GetById(l logrus.FieldLogger, span opentracing.Span) func(id uint32) (*Model, error) {
	return func(id uint32) (*Model, error) {
		return ByIdModelProvider(l, span)(id)()
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

func makeModel(body requests.DataBody[attributes]) (*Model, error) {
	id, err := strconv.ParseUint(body.Id, 10, 32)
	if err != nil {
		return nil, err
	}
	att := body.Attributes
	m := NewBuilder().
		SetId(uint32(id)).
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
	return &m, nil
}
