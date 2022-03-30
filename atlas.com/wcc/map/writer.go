package _map

import (
	"atlas-wcc/character"
	"atlas-wcc/character/inventory"
	"atlas-wcc/pet"
	"atlas-wcc/socket/response"
	"github.com/sirupsen/logrus"
	"math/rand"
	"time"
)

const OpCodeSetField uint16 = 0x7D
const OpCodeSpawnCharacter uint16 = 0xA0
const OpCodeRemoveCharacterFromMap uint16 = 0xA1
const OpCodeChatText uint16 = 0xA2
const OpCodeShowForeignEffect uint16 = 0xC6
const OpCodeShowItemGainInChat uint16 = 0xCE
const OpCodeMoveCharacter uint16 = 0xB9
const OpCodeCharacterDamage uint16 = 0xC0
const OpCodeCharacterExpression uint16 = 0xC1
const OpCodeUpdateCharacterLook uint16 = 0xC5

func WriteGetCharacterInfo(l logrus.FieldLogger) func(channelId byte, character character.Model) []byte {
	return func(channelId byte, character character.Model) []byte {
		w := response.NewWriter(l)
		w.WriteShort(OpCodeSetField)
		w.WriteInt(uint32(channelId - 1))
		w.WriteByte(1)
		w.WriteByte(1)
		w.WriteShort(0)
		for i := 0; i < 3; i++ {
			w.WriteInt(rand.Uint32())
		}
		addCharacterInfo(w, character)
		w.WriteInt64(getTime(timeNow()))
		return w.Bytes()
	}
}

func addCharacterInfo(w *response.Writer, character character.Model) {
	w.WriteInt64(-1)
	w.WriteByte(0)
	addCharacterStats(w, character)
	//buddy list capacity
	w.WriteByte(0)
	//      if (character.getLinkedName() == null) {
	w.WriteByte(0)
	//      } else {
	//         writer.write(1);
	//         writer.writeMapleAsciiString(character.getLinkedName());
	//      }
	w.WriteInt(character.Attributes().Meso())
	inventory.AddInventoryInfo(w, character.Equipment(), character.Inventory())
	addSkillInfo(w, character)
	addQuestInfo(w, character)
	addMiniGameInfo(w, character)
	addRingInfo(w, character)
	addTeleportInfo(w, character)
	addMonsterBookInfo(w, character)
	addNewYearInfo(w, character)
	addAreaInfo(w, character)
	w.WriteShort(0)
}

func addAreaInfo(w *response.Writer, _ character.Model) {
	w.WriteShort(0)
}

func addNewYearInfo(w *response.Writer, _ character.Model) {
	w.WriteShort(0)
}

func addMonsterBookInfo(w *response.Writer, _ character.Model) {
	w.WriteInt(0)
	w.WriteByte(0)
	w.WriteShort(0)
}

func addTeleportInfo(w *response.Writer, _ character.Model) {
	for i := 0; i < 5; i++ {
		w.WriteInt(999999999)
	}
	for j := 0; j < 10; j++ {
		w.WriteInt(999999999)
	}
}

func addRingInfo(w *response.Writer, _ character.Model) {
	w.WriteShort(0)
	w.WriteShort(0)
	w.WriteShort(0)
}

func addMiniGameInfo(w *response.Writer, _ character.Model) {
	w.WriteShort(0)
}

func addQuestInfo(w *response.Writer, _ character.Model) {
	w.WriteShort(0)
	w.WriteShort(0)
}

func addSkillInfo(w *response.Writer, character character.Model) {
	w.WriteByte(0)

	sc := uint16(0)
	for _, s := range character.Skills() {
		if !s.Hidden() {
			sc += 1
		}
	}
	w.WriteShort(sc)

	for _, s := range character.Skills() {
		if !s.Hidden() {
			w.WriteInt(s.Id())
			w.WriteInt(s.Level())
			w.WriteLong(uint64(getTime(s.Expiration())))
			if s.FourthJob() {
				w.WriteInt(s.MasterLevel())
			}
		}
	}

	//      writer.writeShort(character.getAllCoolDowns().size());
	w.WriteShort(0)
	//      for (PlayerCoolDownValueHolder cooling : character.getAllCoolDowns()) {
	//         writer.writeInt(cooling.skillId);
	//         int timeLeft = (int) (cooling.length + cooling.startTime - System.currentTimeMillis());
	//         writer.writeShort(timeLeft / 1000);
	//      }
}

func addCharacterStats(w *response.Writer, character character.Model) {
	w.WriteInt(character.Attributes().Id())
	addPaddedCharacterName(w, character)
	w.WriteByte(character.Attributes().Gender())
	w.WriteByte(character.Attributes().SkinColor())
	w.WriteInt(character.Attributes().Face())
	w.WriteInt(character.Attributes().Hair())
	pet.WriteForEachPet(w, character.Pets(), writePetId, writeEmptyPetId)
	w.WriteByte(character.Attributes().Level())
	w.WriteShort(character.Attributes().JobId())
	w.WriteShort(character.Attributes().Strength())
	w.WriteShort(character.Attributes().Dexterity())
	w.WriteShort(character.Attributes().Intelligence())
	w.WriteShort(character.Attributes().Luck())
	w.WriteShort(character.Attributes().Hp())
	w.WriteShort(character.Attributes().MaxHp())
	w.WriteShort(character.Attributes().Mp())
	w.WriteShort(character.Attributes().MaxMp())
	w.WriteShort(character.Attributes().Ap())

	if character.Attributes().HasSPTable() {
		addRemainingSkillInfo(w, character)
	} else {
		w.WriteShort(character.Attributes().RemainingSp())
	}

	w.WriteInt(character.Attributes().Experience())
	w.WriteShort(uint16(character.Attributes().Fame()))
	w.WriteInt(character.Attributes().GachaponExperience())
	w.WriteInt(character.Attributes().MapId())
	w.WriteByte(character.Attributes().SpawnPoint())
	w.WriteInt(0)
}

func addRemainingSkillInfo(_ *response.Writer, _ character.Model) {

}

func addPaddedCharacterName(w *response.Writer, character character.Model) {
	name := character.Attributes().Name()
	if len(name) > 13 {
		name = name[:13]
	}
	padSize := 13 - len(name)
	w.WriteByteArray([]byte(name))
	for i := 0; i < padSize; i++ {
		w.WriteByte(0x0)
	}
}

func writePetId(w *response.Writer, pet pet.Model) {
	w.WriteLong(pet.Id())
}

func writeEmptyPetId(w *response.Writer) {
	w.WriteLong(0)
}

func getTime(utcTimestamp int64) int64 {
	if utcTimestamp < 0 && utcTimestamp >= -3 {
		if utcTimestamp == -1 {
			return DefaultTime //high number ll
		} else if utcTimestamp == -2 {
			return ZeroTime
		} else {
			return Permanent
		}
	}

	ftUtOffset := 116444736010800000 + (10000 * timeNow())
	return utcTimestamp*10000 + ftUtOffset
}

func timeNow() int64 {
	return time.Now().UnixNano() / int64(time.Millisecond)
}

const (
	DefaultTime int64 = 150842304000000000
	ZeroTime    int64 = 94354848000000000
	Permanent   int64 = 150841440000000000
)

func WriteWarpToMap(l logrus.FieldLogger) func(channelId byte, mapId uint32, portalId uint32, hp uint16) []byte {
	return func(channelId byte, mapId uint32, portalId uint32, hp uint16) []byte {
		w := response.NewWriter(l)
		w.WriteShort(OpCodeSetField)
		w.WriteInt(uint32(channelId) - 1)
		w.WriteInt(0)
		w.WriteByte(0)
		w.WriteInt(mapId)
		w.WriteByte(byte(portalId))
		w.WriteShort(hp)
		w.WriteBool(false)
		w.WriteLong(uint64(getTime(timeNow())))
		return w.Bytes()
	}
}

func WriteSpawnCharacter(l logrus.FieldLogger) func(target character.Model, character character.Model, enteringField bool) []byte {
	return func(target character.Model, c character.Model, enteringField bool) []byte {
		w := response.NewWriter(l)
		w.WriteShort(OpCodeSpawnCharacter)
		w.WriteInt(c.Attributes().Id())
		w.WriteByte(c.Attributes().Level())
		w.WriteAsciiString(c.Attributes().Name())
		//      if (chr.getGuildId() < 1) {
		w.WriteAsciiString("")
		w.WriteByteArray([]byte{0, 0, 0, 0, 0, 0})
		//      } else {
		//         MapleGuildSummary gs = chr.getClient().getWorldServer().getGuildSummary(chr.getGuildId(), chr.getWorld());
		//         if (gs != null) {
		//            writer.writeMapleAsciiString(gs.getName());
		//            writer.writeShort(gs.getLogoBG());
		//            writer.write(gs.getLogoBGColor());
		//            writer.writeShort(gs.getLogo());
		//            writer.write(gs.getLogoColor());
		//         } else {
		//            writer.writeMapleAsciiString("");
		//            writer.write(new byte[6]);
		//         }
		//      }

		writeForeignBuffs(w, c)
		w.WriteShort(c.Attributes().JobId())
		character.AddCharacterLook(w, c, false)
		//writer.writeInt(chr.getInventory(MapleInventoryType.CASH).countById(5110000));
		w.WriteInt(0)
		//writer.writeInt(chr.getItemEffect());
		w.WriteInt(0)
		//writer.writeInt(ItemConstants.getInventoryType(chr.getChair()) == MapleInventoryType.SETUP ? chr.getChair() : 0);
		w.WriteInt(0)

		if enteringField {
			w.WriteInt16(c.Attributes().X())
			w.WriteInt16(c.Attributes().Y() - int16(42))
			w.WriteByte(6)
		} else {
			w.WriteInt16(c.Attributes().X())
			w.WriteInt16(c.Attributes().Y())
			w.WriteByte(c.Attributes().Stance())
		}

		w.WriteShort(0) //chr.getFh()
		w.WriteByte(0)

		pet.WriteForEachPet(w, c.Pets(), addPetInfoButDoNotShow, noOpWrite)

		//end of pets
		w.WriteByte(0)

		//      if (chr.getMount() == null) {
		w.WriteInt(1)  // mob level
		w.WriteLong(0) // mob exp + tiredness
		//      } else {
		//         writer.writeInt(chr.getMount().level());
		//         writer.writeInt(chr.getMount().exp());
		//         writer.writeInt(chr.getMount().tiredness());
		//      }

		//      MaplePlayerShop mps = chr.getPlayerShop();
		//      if (mps != null && mps.isOwner(chr)) {
		//         if (mps.hasFreeSlot()) {
		//            addAnnounceBox(writer, mps, mps.getVisitors().length);
		//         } else {
		//            addAnnounceBox(writer, mps, 1);
		//         }
		//      } else {
		//         MapleMiniGame miniGame = chr.getMiniGame();
		//         if (miniGame != null && miniGame.isOwner(chr)) {
		//            if (miniGame.hasFreeSlot()) {
		//               addAnnounceBox(writer, miniGame, 1, 0);
		//            } else {
		//               addAnnounceBox(writer, miniGame, 2, miniGame.isMatchInProgress() ? 1 : 0);
		//            }
		//         } else {
		w.WriteByte(0)
		//         }
		//      }

		//      if (chr.getChalkboard() != null) {
		//         writer.write(1);
		//         writer.writeMapleAsciiString(chr.getChalkboard());
		//      } else {
		w.WriteByte(0)
		//      }
		character.AddRingLook(w, c, true)  // crush
		character.AddRingLook(w, c, false) // friendship
		character.AddMarriageRingLook(w, target, c)
		encodeNewYearCardInfo(w, c) // new year seems to crash sometimes...
		w.WriteByte(0)
		w.WriteByte(0)
		//writer.write(chr.getTeam());//only needed in specific fields
		w.WriteByte(0)

		return w.Bytes()
	}
}

func noOpWrite(_ *response.Writer) {
}

func addPetInfoButDoNotShow(w *response.Writer, p pet.Model) {
	addPetInfo(w, p, false)
}

func addPetInfo(w *response.Writer, p pet.Model, showPet bool) {
	w.WriteByte(1)
	if showPet {
		w.WriteByte(0)
	}
	w.WriteInt(uint32(p.Id()))
	w.WriteAsciiString(p.Name())
	w.WriteLong(p.Id())
	w.WriteInt16(p.X())
	w.WriteInt16(p.Y())
	w.WriteByte(p.Stance())
	w.WriteInt(p.Fh())
}

func encodeNewYearCardInfo(w *response.Writer, _ character.Model) {
	w.WriteByte(0)
}

func writeForeignBuffs(w *response.Writer, _ character.Model) {
	w.WriteInt(0)
	w.WriteShort(0)
	w.WriteByte(0xFC)
	w.WriteByte(1)
	//      if (chr.getBuffedValue(MapleBuffStat.MORPH) != null) {
	//         writer.writeInt(2);
	//      } else {
	w.WriteInt(0)
	//      }
	bm := uint64(0)
	//      Integer buffValue = null;
	//      if (chr.getBuffedValue(MapleBuffStat.DARK_SIGHT) != null && !chr.isHidden()) {
	//         buffMask |= MapleBuffStat.DARK_SIGHT.getValue();
	//      }
	//      if (chr.getBuffedValue(MapleBuffStat.COMBO) != null) {
	//         buffMask |= MapleBuffStat.COMBO.getValue();
	//         buffValue = chr.getBuffedValue(MapleBuffStat.COMBO);
	//      }
	//      if (chr.getBuffedValue(MapleBuffStat.SHADOW_PARTNER) != null) {
	//         buffMask |= MapleBuffStat.SHADOW_PARTNER.getValue();
	//      }
	//      if (chr.getBuffedValue(MapleBuffStat.SOUL_ARROW) != null) {
	//         buffMask |= MapleBuffStat.SOUL_ARROW.getValue();
	//      }
	//      if (chr.getBuffedValue(MapleBuffStat.MORPH) != null) {
	//         buffValue = chr.getBuffedValue(MapleBuffStat.MORPH);
	//      }
	//      if (chr.getBuffedValue(MapleBuffStat.ENERGY_CHARGE) != null) {
	//         buffMask |= MapleBuffStat.ENERGY_CHARGE.getValue();
	//         buffValue = chr.getBuffedValue(MapleBuffStat.ENERGY_CHARGE);
	//      }//AREN'T THESE
	w.WriteInt(uint32((bm >> 32) & 0xffffffff))
	//      if (buffValue != null) {
	//         if (chr.getBuffedValue(MapleBuffStat.MORPH) != null) { //TEST
	//            writer.writeShort(buffValue);
	//         } else {
	//            writer.write(buffValue.byteValue());
	//         }
	//      }
	w.WriteInt(uint32(bm & 0xffffffff))
	cms := rand.Uint32()
	w.Skip(6)
	w.WriteInt(cms)
	w.Skip(11)
	w.WriteInt(cms)
	w.Skip(11)
	w.WriteInt(cms)
	w.WriteShort(0)
	w.WriteByte(0)

	//      Integer bv = chr.getBuffedValue(MapleBuffStat.MONSTER_RIDING);
	//      if (bv != null) {
	//         MapleMount mount = chr.getMount();
	//         if (mount != null) {
	//            writer.writeInt(mount.itemId());
	//            writer.writeInt(mount.skillId());
	//         } else {
	//            writer.writeLong(0);
	//         }
	//      } else {
	w.WriteLong(0)
	//      }

	w.WriteInt(cms)
	w.Skip(9)
	w.WriteInt(cms)
	w.WriteShort(0)
	w.WriteInt(0)
	w.Skip(10)
	w.WriteInt(cms)
	w.Skip(13)
	w.WriteInt(cms)
	w.WriteShort(0)
	w.WriteByte(0)
}

func WriteRemoveCharacterFromMap(l logrus.FieldLogger) func(characterId uint32) []byte {
	return func(characterId uint32) []byte {
		w := response.NewWriter(l)
		w.WriteShort(OpCodeRemoveCharacterFromMap)
		w.WriteInt(characterId)
		return w.Bytes()
	}
}

func WriteChatText(l logrus.FieldLogger) func(characterId uint32, message string, gm bool, show bool) []byte {
	return func(characterId uint32, message string, gm bool, show bool) []byte {
		w := response.NewWriter(l)
		w.WriteShort(OpCodeChatText)
		w.WriteInt(characterId)
		w.WriteBool(gm)
		w.WriteAsciiString(message)
		w.WriteBool(show)
		return w.Bytes()
	}
}

func WriteShowForeignEffect(l logrus.FieldLogger) func(characterId uint32, effect byte) []byte {
	return func(characterId uint32, effect byte) []byte {
		w := response.NewWriter(l)
		w.WriteShort(OpCodeShowForeignEffect)
		w.WriteInt(characterId)
		w.WriteByte(effect)
		return w.Bytes()
	}
}

func WriteShowBuffEffect(l logrus.FieldLogger) func(characterId uint32, effect byte, skillId uint32, direction byte) []byte {
	return func(characterId uint32, effect byte, skillId uint32, direction byte) []byte {
		w := response.NewWriter(l)
		w.WriteShort(OpCodeShowForeignEffect)
		w.WriteInt(characterId)
		w.WriteByte(effect)
		w.WriteInt(skillId)
		w.WriteByte(direction)
		return w.Bytes()
	}
}

func WriteShowOwnBuff(l logrus.FieldLogger) func(effect byte, skillId uint32) []byte {
	return func(effect byte, skillId uint32) []byte {
		w := response.NewWriter(l)
		w.WriteShort(OpCodeShowItemGainInChat)
		w.WriteByte(effect)
		w.WriteInt(skillId)
		w.WriteByte(0xA9)
		w.WriteByte(1)
		return w.Bytes()
	}
}

func WriteMoveCharacter(l logrus.FieldLogger) func(characterId uint32, movementList []byte) []byte {
	return func(characterId uint32, movementList []byte) []byte {
		w := response.NewWriter(l)
		w.WriteShort(OpCodeMoveCharacter)
		w.WriteInt(characterId)
		w.WriteInt(0)
		for _, b := range movementList {
			w.WriteByte(b)
		}
		return w.Bytes()
	}
}

func WriteCharacterDamaged(l logrus.FieldLogger) func(skillId int8, monsterId uint32, characterId uint32, damage int32, fake uint32, direction int8, pgmr bool, pgmr1 byte, pg bool, monsterUniqueId uint32, x int16, y int16) []byte {
	return func(skillId int8, monsterId uint32, characterId uint32, damage int32, fake uint32, direction int8, pgmr bool, pgmr1 byte, pg bool, monsterUniqueId uint32, x int16, y int16) []byte {
		w := response.NewWriter(l)
		w.WriteShort(OpCodeCharacterDamage)
		w.WriteInt(characterId)
		w.WriteInt8(skillId)
		w.WriteInt32(damage)
		if skillId != -4 {
			w.WriteInt(monsterId)
			w.WriteInt8(direction)
			if pgmr {
				w.WriteByte(pgmr1)
				if pg {
					w.WriteByte(1)
				} else {
					w.WriteByte(0)
				}
				w.WriteInt(monsterUniqueId)
				w.WriteByte(6)
				w.WriteInt16(x)
				w.WriteInt16(y)
				w.WriteByte(0)
			} else {
				w.WriteShort(0)
			}
			w.WriteInt32(damage)
			if fake > 0 {
				w.WriteInt(fake)
			}
		} else {
			w.WriteInt32(damage)
		}
		return w.Bytes()
	}
}

func WriteCharacterExpression(l logrus.FieldLogger) func(characterId uint32, expression uint32) []byte {
	return func(characterId uint32, expression uint32) []byte {
		w := response.NewWriter(l)
		w.WriteShort(OpCodeCharacterExpression)
		w.WriteInt(characterId)
		w.WriteInt(expression)
		return w.Bytes()
	}
}

func WriteCharacterLookUpdated(l logrus.FieldLogger) func(r character.Model, c character.Model) []byte {
	return func(r character.Model, c character.Model) []byte {
		w := response.NewWriter(l)
		w.WriteShort(OpCodeUpdateCharacterLook)
		w.WriteInt(c.Attributes().Id())
		w.WriteByte(1)
		character.AddCharacterLook(w, c, false)
		character.AddRingLook(w, c, true)
		character.AddRingLook(w, c, false)
		character.AddMarriageRingLook(w, r, c)
		w.WriteInt(0)
		return w.Bytes()
	}
}
