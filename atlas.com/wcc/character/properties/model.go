package properties

import (
	"strconv"
	"strings"
)

type Properties struct {
	id                 uint32
	accountId          uint32
	worldId            byte
	name               string
	gender             byte
	skinColor          byte
	face               uint32
	hair               uint32
	level              byte
	jobId              uint16
	strength           uint16
	dexterity          uint16
	intelligence       uint16
	luck               uint16
	hp                 uint16
	maxHp              uint16
	mp                 uint16
	maxMp              uint16
	ap                 uint16
	sp                 string
	experience         uint32
	fame               int16
	gachaponExperience uint32
	mapId              uint32
	spawnPoint         byte
	gm                 bool
	gmJob              bool
	rank               int
	rankMove           int
	jobRank            int
	jobRankMove        int
	meso               uint32
	x                  int16
	y                  int16
	stance             byte
}

func (a Properties) Gm() bool {
	return a.gm
}

func (a Properties) GmJob() bool {
	return a.gmJob
}

func (a Properties) Rank() int {
	return a.rank
}

func (a Properties) RankMove() int {
	return a.rankMove
}

func (a Properties) JobRank() int {
	return a.jobRank
}

func (a Properties) JobRankMove() int {
	return a.jobRankMove
}

func (a Properties) Id() uint32 {
	return a.id
}

func (a Properties) Name() string {
	return a.name
}

func (a Properties) Gender() byte {
	return a.gender
}

func (a Properties) SkinColor() byte {
	return a.skinColor
}

func (a Properties) Face() uint32 {
	return a.face
}

func (a Properties) Hair() uint32 {
	return a.hair
}

func (a Properties) Level() byte {
	return a.level
}

func (a Properties) JobId() uint16 {
	return a.jobId
}

func (a Properties) Strength() uint16 {
	return a.strength
}

func (a Properties) Dexterity() uint16 {
	return a.dexterity
}

func (a Properties) Intelligence() uint16 {
	return a.intelligence
}

func (a Properties) Luck() uint16 {
	return a.luck
}

func (a Properties) Hp() uint16 {
	return a.hp
}

func (a Properties) MaxHp() uint16 {
	return a.maxHp
}

func (a Properties) Mp() uint16 {
	return a.mp
}

func (a Properties) MaxMp() uint16 {
	return a.maxMp
}

func (a Properties) Ap() uint16 {
	return a.ap
}

func (a Properties) HasSPTable() bool {
	switch a.jobId {
	case 2001:
		return true
	case 2200:
		return true
	case 2210:
		return true
	case 2211:
		return true
	case 2212:
		return true
	case 2213:
		return true
	case 2214:
		return true
	case 2215:
		return true
	case 2216:
		return true
	case 2217:
		return true
	case 2218:
		return true
	default:
		return false
	}
}

func (a Properties) Sp() []uint16 {
	s := strings.Split(a.sp, ",")
	var sps = make([]uint16, 0)
	for _, x := range s {
		sp, err := strconv.ParseUint(x, 10, 16)
		if err == nil {
			sps = append(sps, uint16(sp))
		}
	}
	return sps
}

func (a Properties) RemainingSp() uint16 {
	return a.Sp()[a.skillBook()]

}

func (a Properties) skillBook() uint16 {
	if a.jobId >= 2210 && a.jobId <= 2218 {
		return a.jobId - 2209
	}
	return 0
}

func (a Properties) Experience() uint32 {
	return a.experience
}

func (a Properties) Fame() int16 {
	return a.fame
}

func (a Properties) GachaponExperience() uint32 {
	return a.gachaponExperience
}

func (a Properties) SpawnPoint() byte {
	return a.spawnPoint
}

func (a Properties) WorldId() byte {
	return a.worldId
}

func (a Properties) MapId() uint32 {
	return a.mapId
}

func (a Properties) AccountId() uint32 {
	return a.accountId
}

func (a Properties) Meso() uint32 {
	return a.meso
}

func (a Properties) X() int16 {
	return a.x
}

func (a Properties) Y() int16 {
	return a.y
}

func (a Properties) Stance() byte {
	return a.stance
}


type propertiesBuilder struct {
	id                 uint32
	accountId          uint32
	worldId            byte
	name               string
	gender             byte
	skinColor          byte
	face               uint32
	hair               uint32
	level              byte
	jobId              uint16
	strength           uint16
	dexterity          uint16
	intelligence       uint16
	luck               uint16
	hp                 uint16
	maxHp              uint16
	mp                 uint16
	maxMp              uint16
	ap                 uint16
	sp                 string
	experience         uint32
	fame               int16
	gachaponExperience uint32
	mapId              uint32
	spawnPoint         byte
	gm                 bool
	gmJob              bool
	rank               int
	rankMove           int
	jobRank            int
	jobRankMove        int
	meso               uint32
	x                  int16
	y                  int16
	stance             byte
}

func NewPropertiesBuilder() *propertiesBuilder {
	return &propertiesBuilder{}
}

func (c *propertiesBuilder) SetId(id uint32) *propertiesBuilder {
	c.id = id
	return c
}

func (c *propertiesBuilder) SetAccountId(accountId uint32) *propertiesBuilder {
	c.accountId = accountId
	return c
}

func (c *propertiesBuilder) SetWorldId(worldId byte) *propertiesBuilder {
	c.worldId = worldId
	return c
}

func (c *propertiesBuilder) SetName(name string) *propertiesBuilder {
	c.name = name
	return c
}

func (c *propertiesBuilder) SetGender(gender byte) *propertiesBuilder {
	c.gender = gender
	return c
}

func (c *propertiesBuilder) SetSkinColor(skinColor byte) *propertiesBuilder {
	c.skinColor = skinColor
	return c
}

func (c *propertiesBuilder) SetFace(face uint32) *propertiesBuilder {
	c.face = face
	return c
}

func (c *propertiesBuilder) SetHair(hair uint32) *propertiesBuilder {
	c.hair = hair
	return c
}

func (c *propertiesBuilder) SetLevel(level byte) *propertiesBuilder {
	c.level = level
	return c
}

func (c *propertiesBuilder) SetJobId(jobId uint16) *propertiesBuilder {
	c.jobId = jobId
	return c
}

func (c *propertiesBuilder) SetStrength(strength uint16) *propertiesBuilder {
	c.strength = strength
	return c
}

func (c *propertiesBuilder) SetDexterity(dexterity uint16) *propertiesBuilder {
	c.dexterity = dexterity
	return c
}

func (c *propertiesBuilder) SetIntelligence(intelligence uint16) *propertiesBuilder {
	c.intelligence = intelligence
	return c
}

func (c *propertiesBuilder) SetLuck(luck uint16) *propertiesBuilder {
	c.luck = luck
	return c
}

func (c *propertiesBuilder) SetHp(hp uint16) *propertiesBuilder {
	c.hp = hp
	return c
}

func (c *propertiesBuilder) SetMaxHp(maxHp uint16) *propertiesBuilder {
	c.maxHp = maxHp
	return c
}

func (c *propertiesBuilder) SetMp(mp uint16) *propertiesBuilder {
	c.mp = mp
	return c
}

func (c *propertiesBuilder) SetMaxMp(maxMp uint16) *propertiesBuilder {
	c.maxMp = maxMp
	return c
}

func (c *propertiesBuilder) SetAp(ap uint16) *propertiesBuilder {
	c.ap = ap
	return c
}

func (c *propertiesBuilder) SetSp(sp string) *propertiesBuilder {
	c.sp = sp
	return c
}

func (c *propertiesBuilder) SetExperience(experience uint32) *propertiesBuilder {
	c.experience = experience
	return c
}

func (c *propertiesBuilder) SetFame(fame int16) *propertiesBuilder {
	c.fame = fame
	return c
}

func (c *propertiesBuilder) SetGachaponExperience(gachaponExperience uint32) *propertiesBuilder {
	c.gachaponExperience = gachaponExperience
	return c
}

func (c *propertiesBuilder) SetMapId(mapId uint32) *propertiesBuilder {
	c.mapId = mapId
	return c
}

func (c *propertiesBuilder) SetSpawnPoint(spawnPoint byte) *propertiesBuilder {
	c.spawnPoint = spawnPoint
	return c
}

func (c *propertiesBuilder) SetGm(gm bool) *propertiesBuilder {
	c.gm = gm
	return c
}

func (c *propertiesBuilder) SetGmJob(gmJob bool) *propertiesBuilder {
	c.gmJob = gmJob
	return c
}

func (c *propertiesBuilder) SetRank(rank int) *propertiesBuilder {
	c.rank = rank
	return c
}

func (c *propertiesBuilder) SetRankMove(rankMove int) *propertiesBuilder {
	c.rankMove = rankMove
	return c
}

func (c *propertiesBuilder) SetJobRank(jobRank int) *propertiesBuilder {
	c.jobRank = jobRank
	return c
}

func (c *propertiesBuilder) SetJobRankMove(jobRankMove int) *propertiesBuilder {
	c.jobRankMove = jobRankMove
	return c
}

func (c *propertiesBuilder) SetMeso(meso uint32) *propertiesBuilder {
	c.meso = meso
	return c
}

func (c *propertiesBuilder) SetX(x int16) *propertiesBuilder {
	c.x = x
	return c
}

func (c *propertiesBuilder) SetY(y int16) *propertiesBuilder {
	c.y = y
	return c
}

func (c *propertiesBuilder) SetStance(stance byte) *propertiesBuilder {
	c.stance = stance
	return c
}

func (c *propertiesBuilder) Build() Properties {
	return Properties{
		id:                 c.id,
		accountId:          c.accountId,
		worldId:            c.worldId,
		name:               c.name,
		gender:             c.gender,
		skinColor:          c.skinColor,
		face:               c.face,
		hair:               c.hair,
		level:              c.level,
		jobId:              c.jobId,
		strength:           c.strength,
		dexterity:          c.dexterity,
		intelligence:       c.intelligence,
		luck:               c.luck,
		hp:                 c.hp,
		maxHp:              c.maxHp,
		mp:                 c.mp,
		maxMp:              c.maxMp,
		ap:                 c.ap,
		sp:                 c.sp,
		experience:         c.experience,
		fame:               c.fame,
		gachaponExperience: c.gachaponExperience,
		mapId:              c.mapId,
		spawnPoint:         c.spawnPoint,
		gm:                 c.gm,
		gmJob:              c.gmJob,
		rank:               c.rank,
		rankMove:           c.rankMove,
		jobRank:            c.jobRank,
		jobRankMove:        c.jobRankMove,
		meso:               c.meso,
		x:                  c.x,
		y:                  c.y,
		stance:             c.stance,
	}
}