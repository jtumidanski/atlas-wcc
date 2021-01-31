package domain

import (
   "strconv"
   "strings"
)

type CharacterAttributes struct {
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

func (a CharacterAttributes) Gm() bool {
   return a.gm
}

func (a CharacterAttributes) GmJob() bool {
   return a.gmJob
}

func (a CharacterAttributes) Rank() int {
   return a.rank
}

func (a CharacterAttributes) RankMove() int {
   return a.rankMove
}

func (a CharacterAttributes) JobRank() int {
   return a.jobRank
}

func (a CharacterAttributes) JobRankMove() int {
   return a.jobRankMove
}

func (a CharacterAttributes) Id() uint32 {
   return a.id
}

func (a CharacterAttributes) Name() string {
   return a.name
}

func (a CharacterAttributes) Gender() byte {
   return a.gender
}

func (a CharacterAttributes) SkinColor() byte {
   return a.skinColor
}

func (a CharacterAttributes) Face() uint32 {
   return a.face
}

func (a CharacterAttributes) Hair() uint32 {
   return a.hair
}

func (a CharacterAttributes) Level() byte {
   return a.level
}

func (a CharacterAttributes) JobId() uint16 {
   return a.jobId
}

func (a CharacterAttributes) Strength() uint16 {
   return a.strength
}

func (a CharacterAttributes) Dexterity() uint16 {
   return a.dexterity
}

func (a CharacterAttributes) Intelligence() uint16 {
   return a.intelligence
}

func (a CharacterAttributes) Luck() uint16 {
   return a.luck
}

func (a CharacterAttributes) Hp() uint16 {
   return a.hp
}

func (a CharacterAttributes) MaxHp() uint16 {
   return a.maxHp
}

func (a CharacterAttributes) Mp() uint16 {
   return a.mp
}

func (a CharacterAttributes) MaxMp() uint16 {
   return a.maxMp
}

func (a CharacterAttributes) Ap() uint16 {
   return a.ap
}

func (a CharacterAttributes) HasSPTable() bool {
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

func (a CharacterAttributes) Sp() []uint16 {
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

func (a CharacterAttributes) RemainingSp() uint16 {
   return a.Sp()[a.skillBook()]

}

func (a CharacterAttributes) skillBook() uint16 {
   if a.jobId >= 2210 && a.jobId <= 2218 {
      return a.jobId - 2209
   }
   return 0
}

func (a CharacterAttributes) Experience() uint32 {
   return a.experience
}

func (a CharacterAttributes) Fame() int16 {
   return a.fame
}

func (a CharacterAttributes) GachaponExperience() uint32 {
   return a.gachaponExperience
}

func (a CharacterAttributes) SpawnPoint() byte {
   return a.spawnPoint
}

func (a CharacterAttributes) WorldId() byte {
   return a.worldId
}

func (a CharacterAttributes) MapId() uint32 {
   return a.mapId
}

func (a CharacterAttributes) AccountId() uint32 {
   return a.accountId
}

func (a CharacterAttributes) Meso() uint32 {
   return a.meso
}

func (a CharacterAttributes) X() int16 {
   return a.x
}

func (a CharacterAttributes) Y() int16 {
   return a.y
}

func (a CharacterAttributes) Stance() byte {
   return a.stance
}

type characterAttributeBuilder struct {
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

func NewCharacterAttributeBuilder() *characterAttributeBuilder {
   return &characterAttributeBuilder{}
}

func (c *characterAttributeBuilder) SetId(id uint32) *characterAttributeBuilder {
   c.id = id
   return c
}

func (c *characterAttributeBuilder) SetAccountId(accountId uint32) *characterAttributeBuilder {
   c.accountId = accountId
   return c
}

func (c *characterAttributeBuilder) SetWorldId(worldId byte) *characterAttributeBuilder {
   c.worldId = worldId
   return c
}

func (c *characterAttributeBuilder) SetName(name string) *characterAttributeBuilder {
   c.name = name
   return c
}

func (c *characterAttributeBuilder) SetGender(gender byte) *characterAttributeBuilder {
   c.gender = gender
   return c
}

func (c *characterAttributeBuilder) SetSkinColor(skinColor byte) *characterAttributeBuilder {
   c.skinColor = skinColor
   return c
}

func (c *characterAttributeBuilder) SetFace(face uint32) *characterAttributeBuilder {
   c.face = face
   return c
}

func (c *characterAttributeBuilder) SetHair(hair uint32) *characterAttributeBuilder {
   c.hair = hair
   return c
}

func (c *characterAttributeBuilder) SetLevel(level byte) *characterAttributeBuilder {
   c.level = level
   return c
}

func (c *characterAttributeBuilder) SetJobId(jobId uint16) *characterAttributeBuilder {
   c.jobId = jobId
   return c
}

func (c *characterAttributeBuilder) SetStrength(strength uint16) *characterAttributeBuilder {
   c.strength = strength
   return c
}

func (c *characterAttributeBuilder) SetDexterity(dexterity uint16) *characterAttributeBuilder {
   c.dexterity = dexterity
   return c
}

func (c *characterAttributeBuilder) SetIntelligence(intelligence uint16) *characterAttributeBuilder {
   c.intelligence = intelligence
   return c
}

func (c *characterAttributeBuilder) SetLuck(luck uint16) *characterAttributeBuilder {
   c.luck = luck
   return c
}

func (c *characterAttributeBuilder) SetHp(hp uint16) *characterAttributeBuilder {
   c.hp = hp
   return c
}

func (c *characterAttributeBuilder) SetMaxHp(maxHp uint16) *characterAttributeBuilder {
   c.maxHp = maxHp
   return c
}

func (c *characterAttributeBuilder) SetMp(mp uint16) *characterAttributeBuilder {
   c.mp = mp
   return c
}

func (c *characterAttributeBuilder) SetMaxMp(maxMp uint16) *characterAttributeBuilder {
   c.maxMp = maxMp
   return c
}

func (c *characterAttributeBuilder) SetAp(ap uint16) *characterAttributeBuilder {
   c.ap = ap
   return c
}

func (c *characterAttributeBuilder) SetSp(sp string) *characterAttributeBuilder {
   c.sp = sp
   return c
}

func (c *characterAttributeBuilder) SetExperience(experience uint32) *characterAttributeBuilder {
   c.experience = experience
   return c
}

func (c *characterAttributeBuilder) SetFame(fame int16) *characterAttributeBuilder {
   c.fame = fame
   return c
}

func (c *characterAttributeBuilder) SetGachaponExperience(gachaponExperience uint32) *characterAttributeBuilder {
   c.gachaponExperience = gachaponExperience
   return c
}

func (c *characterAttributeBuilder) SetMapId(mapId uint32) *characterAttributeBuilder {
   c.mapId = mapId
   return c
}

func (c *characterAttributeBuilder) SetSpawnPoint(spawnPoint byte) *characterAttributeBuilder {
   c.spawnPoint = spawnPoint
   return c
}

func (c *characterAttributeBuilder) SetGm(gm bool) *characterAttributeBuilder {
   c.gm = gm
   return c
}

func (c *characterAttributeBuilder) SetGmJob(gmJob bool) *characterAttributeBuilder {
   c.gmJob = gmJob
   return c
}

func (c *characterAttributeBuilder) SetRank(rank int) *characterAttributeBuilder {
   c.rank = rank
   return c
}

func (c *characterAttributeBuilder) SetRankMove(rankMove int) *characterAttributeBuilder {
   c.rankMove = rankMove
   return c
}

func (c *characterAttributeBuilder) SetJobRank(jobRank int) *characterAttributeBuilder {
   c.jobRank = jobRank
   return c
}

func (c *characterAttributeBuilder) SetJobRankMove(jobRankMove int) *characterAttributeBuilder {
   c.jobRankMove = jobRankMove
   return c
}

func (c *characterAttributeBuilder) SetMeso(meso uint32) *characterAttributeBuilder {
   c.meso = meso
   return c
}

func (c *characterAttributeBuilder) SetX(x int16) *characterAttributeBuilder {
   c.x = x
   return c
}

func (c *characterAttributeBuilder) SetY(y int16) *characterAttributeBuilder {
   c.y = y
   return c
}

func (c *characterAttributeBuilder) SetStance(stance byte) *characterAttributeBuilder {
   c.stance = stance
   return c
}

func (c *characterAttributeBuilder) Build() CharacterAttributes {
   return CharacterAttributes{
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
