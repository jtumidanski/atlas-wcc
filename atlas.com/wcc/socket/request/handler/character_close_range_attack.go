package handler

import (
	"atlas-wcc/kafka/producers"
	"atlas-wcc/mapleSession"
	"atlas-wcc/processors"
	request2 "atlas-wcc/socket/request"
	"atlas-wcc/socket/response/writer"
	"context"
	"github.com/jtumidanski/atlas-socket/request"
	"log"
	"math"
)

const OpCharacterCloseRangeAttack uint16 = 0x2C

type attackPacket struct {
	numberAttacked           byte
	numberDamaged            byte
	numberAttackedAndDamaged byte
	skill                    uint32
	skillLevel               byte
	stance                   byte
	direction                byte
	rangedDirection          byte
	charge                   uint32
	display                  byte
	ranged                   bool
	magic                    bool
	speed                    byte
	allDamage                map[uint32][]uint32
	x                        int16
	y                        int16
}

func (p attackPacket) Skill() uint32 {
	return p.skill
}

func (p attackPacket) SkillLevel() byte {
	return p.skillLevel
}

func (p attackPacket) Stance() byte {
	return p.stance
}

func (p attackPacket) NumberAttackedAndDamaged() byte {
	return p.numberAttackedAndDamaged
}

func (p attackPacket) AllDamage() map[uint32][]uint32 {
	return p.allDamage
}

func (p attackPacket) Speed() byte {
	return p.speed
}

func (p attackPacket) Direction() byte {
	return p.direction
}

func (p attackPacket) Display() byte {
	return p.display
}

func readAttackPacket(reader *request.RequestReader, characterId uint32, ranged bool, magic bool) attackPacket {
	reader.ReadByte()
	numberAttackedAndDamaged := reader.ReadByte()
	numberAttacked := numberAttackedAndDamaged >> 4 & 0xF
	numberDamaged := numberAttackedAndDamaged & 0xF
	skillId := reader.ReadUint32()
	skillLevel := byte(0)
	charge := uint32(0)
	reader.Skip(8)
	display := reader.ReadByte()
	direction := reader.ReadByte()
	stance := reader.ReadByte()
	var speed byte
	var rangedDirection byte
	if ranged {
		reader.ReadByte()
		speed = reader.ReadByte()
		reader.ReadByte()
		rangedDirection = reader.ReadByte()
		reader.Skip(7)
	} else {
		reader.ReadByte()
		speed = reader.ReadByte()
		reader.Skip(4)
	}
	calculatedMaximumDamage := processors.GetCharacterWeaponDamage(characterId)
	bonusDamageBuff := uint32(100)
	if bonusDamageBuff != 100 {
		damageBuff := bonusDamageBuff / 100
		calculatedMaximumDamage = uint32(math.Float64bits(math.Ceil(float64(calculatedMaximumDamage * damageBuff))))
	}
	canCritical := false
	shadowPartner := false

	damage := make(map[uint32][]uint32, 0)
	for i := byte(0); i < numberAttacked; i++ {
		oid := reader.ReadUint32()
		reader.Skip(14)

		var allDamageNumbers []uint32
		for j := byte(0); j < numberDamaged; j++ {
			damage := reader.ReadUint32()
			hitDamageMax := calculatedMaximumDamage
			if shadowPartner {
				if j >= numberDamaged/2 {
					hitDamageMax = uint32(math.Ceil(float64(hitDamageMax) * 0.5))
				}
			}
			maxWithCritical := hitDamageMax
			if canCritical {
				maxWithCritical *= 2
			}
			allDamageNumbers = append(allDamageNumbers, damage)
		}
		reader.Skip(4)
		damage[oid] = allDamageNumbers
	}
	return attackPacket{
		numberAttacked:           numberAttacked,
		numberDamaged:            numberDamaged,
		numberAttackedAndDamaged: numberAttackedAndDamaged,
		skill:                    skillId,
		skillLevel:               skillLevel,
		stance:                   stance,
		direction:                direction,
		rangedDirection:          rangedDirection,
		charge:                   charge,
		display:                  display,
		ranged:                   ranged,
		magic:                    magic,
		speed:                    speed,
		allDamage:                damage,
		x:                        0,
		y:                        0,
	}
}

func CharacterCloseRangeAttackHandler() request2.SessionRequestHandler {
	return func(l *log.Logger, s *mapleSession.MapleSession, r *request.RequestReader) {
		p := readAttackPacket(r, (*s).CharacterId(), false, false)

		catt, err := processors.GetCharacterAttributesById((*s).CharacterId())
		if err != nil {
			l.Printf("[ERROR] unable to retrieve character attributes for character %d.", (*s).CharacterId())
			return
		}

		processors.ForEachSessionInMap((*s).WorldId(), (*s).ChannelId(), catt.MapId(), writeCloseRangeAttack(p))

		attackCount := uint32(1)
		applyAttack(l, p, (*s).WorldId(), (*s).ChannelId(), catt.MapId(), (*s).CharacterId(), attackCount)
	}
}

func writeCloseRangeAttack(p attackPacket) processors.SessionOperator {
	return func(session mapleSession.MapleSession) {
		session.Announce(writer.WriteCloseRangeAttack(session.CharacterId(), p.Skill(), p.SkillLevel(), p.Stance(), p.NumberAttackedAndDamaged(), p.AllDamage(), p.Speed(), p.Direction(), p.Display()))
	}
}

func applyAttack(l *log.Logger, p attackPacket, worldId byte, channelId byte, mapId uint32, characterId uint32, attackCount uint32) {
	for k, v := range p.allDamage {
		m, err := processors.GetMonster(k)
		if err != nil {
			l.Printf("[ERROR] cannot locate monster %d which the attack from %d hit.", k, characterId)
		} else {
			totalDamage := uint32(0)
			for _, e := range v {
				totalDamage += e
			}
			producers.MonsterDamage(l, context.Background()).Emit(worldId, channelId, mapId, m.UniqueId(), characterId, totalDamage)
		}
	}
}
