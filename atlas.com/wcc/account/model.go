package account

type Model struct {
	id             uint32
	name           string
	password       string
	pin            string
	pic            string
	loggedIn       int
	lastLogin      uint64
	gender         byte
	banned         bool
	tos            bool
	language       string
	country        string
	characterSlots int16
}

func (a Model) Id() uint32 {
	return a.id
}

func (a Model) Name() string {
	return a.name
}

func (a Model) Gender() byte {
	return a.gender
}

func (a Model) PIC() string {
	return a.pic
}

func (a Model) CharacterSlots() int16 {
	return a.characterSlots
}

func (a Model) LoggedIn() int {
	return a.loggedIn
}

type accountBuilder struct {
	id             uint32
	name           string
	password       string
	pin            string
	pic            string
	loggedIn       int
	lastLogin      uint64
	gender         byte
	banned         bool
	tos            bool
	language       string
	country        string
	characterSlots int16
}

func NewAccountBuilder() *accountBuilder {
	return &accountBuilder{}
}

func (a *accountBuilder) SetId(id uint32) *accountBuilder {
	a.id = id
	return a
}

func (a *accountBuilder) SetName(name string) *accountBuilder {
	a.name = name
	return a
}

func (a *accountBuilder) SetPassword(password string) *accountBuilder {
	a.password = password
	return a
}

func (a *accountBuilder) SetPin(pin string) *accountBuilder {
	a.pic = pin
	return a
}

func (a *accountBuilder) SetPic(pic string) *accountBuilder {
	a.pin = pic
	return a
}

func (a *accountBuilder) SetLoggedIn(loggedIn int) *accountBuilder {
	a.loggedIn = loggedIn
	return a
}

func (a *accountBuilder) SetLastLogin(lastLogin uint64) *accountBuilder {
	a.lastLogin = lastLogin
	return a
}

func (a *accountBuilder) SetGender(gender byte) *accountBuilder {
	a.gender = gender
	return a
}

func (a *accountBuilder) SetBanned(banned bool) *accountBuilder {
	a.banned = banned
	return a
}

func (a *accountBuilder) SetTos(tos bool) *accountBuilder {
	a.tos = tos
	return a
}

func (a *accountBuilder) SetLanguage(language string) *accountBuilder {
	a.language = language
	return a
}

func (a *accountBuilder) SetCountry(country string) *accountBuilder {
	a.country = country
	return a
}

func (a *accountBuilder) SetCharacterSlots(characterSlots int16) *accountBuilder {
	a.characterSlots = characterSlots
	return a
}

func (a *accountBuilder) Build() Model {
	return Model{
		id:             a.id,
		name:           a.name,
		password:       a.password,
		pin:            a.pin,
		pic:            a.pic,
		loggedIn:       a.loggedIn,
		lastLogin:      a.lastLogin,
		gender:         a.gender,
		banned:         a.banned,
		tos:            a.tos,
		language:       a.language,
		country:        a.country,
		characterSlots: a.characterSlots,
	}
}
