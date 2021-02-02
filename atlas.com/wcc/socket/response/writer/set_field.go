package writer

import (
   "atlas-wcc/domain"
   "atlas-wcc/socket/response"
   "math/rand"
   "time"
)

const OpCodeSetField uint16 = 0x7D

func WriteGetCharacterInfo(channelId byte, character domain.Character) []byte {
   w := response.NewWriter()
   w.WriteShort(OpCodeSetField)
   w.WriteInt(uint32(channelId - 1))
   w.WriteByte(1)
   w.WriteByte(1)
   w.WriteShort(0)
   for i := 0; i < 3; i++ {
      w.WriteInt(rand.Uint32())
   }
   addCharacterInfo(w, character)
   w.WriteLong(uint64(getTime(timeNow())))
   return w.Bytes()
}

func addCharacterInfo(w *response.Writer, character domain.Character) {
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
   addInventoryInfo(w, character)
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

func addAreaInfo(w *response.Writer, _ domain.Character) {
   w.WriteShort(0)
}

func addNewYearInfo(w *response.Writer, _ domain.Character) {
   w.WriteShort(0)
}

func addMonsterBookInfo(w *response.Writer, _ domain.Character) {
   w.WriteInt(0)
   w.WriteByte(0)
   w.WriteShort(0)
}

func addTeleportInfo(w *response.Writer, _ domain.Character) {
   for i := 0; i < 5; i++ {
      w.WriteInt(999999999)
   }
   for j := 0; j < 10; j++ {
      w.WriteInt(999999999)
   }
}

func addRingInfo(w *response.Writer, _ domain.Character) {
   w.WriteShort(0)
   w.WriteShort(0)
   w.WriteShort(0)
}

func addMiniGameInfo(w *response.Writer, _ domain.Character) {
   w.WriteShort(0)
}

func addQuestInfo(w *response.Writer, _ domain.Character) {
   w.WriteShort(0)
   w.WriteShort(0)
}

func addSkillInfo(w *response.Writer, character domain.Character) {
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

func addInventoryInfo(w *response.Writer, character domain.Character) {
   w.WriteByte(character.Inventory().EquipInventory().Capacity())
   w.WriteByte(character.Inventory().UseInventory().Capacity())
   w.WriteByte(character.Inventory().SetupInventory().Capacity())
   w.WriteByte(character.Inventory().EtcInventory().Capacity())
   w.WriteByte(character.Inventory().CashInventory().Capacity())

   w.WriteLong(uint64(getTime(-2)))

   for _, e := range character.Equipment() {
      if e.IsRegularEquipment() {
         addEquipmentInfo(w, e)
      }
   }

   w.WriteShort(0)
   for _, e := range character.Equipment() {
      if e.IsEquippedCashItem() {
         addEquipmentInfo(w, e)
      }
   }

   w.WriteShort(0)
   for _, e := range character.Inventory().EquipInventory().Items() {
      addEquipmentInfo(w, e)
   }

   w.WriteInt(0)
   for _, i := range character.Inventory().UseInventory().Items() {
      addItemInfo(w, i)
   }

   w.WriteByte(0)
   for _, i := range character.Inventory().SetupInventory().Items() {
      addItemInfo(w, i)
   }

   w.WriteByte(0)
   for _, i := range character.Inventory().EtcInventory().Items() {
      addItemInfo(w, i)
   }

   w.WriteByte(0)
   for _, i := range character.Inventory().CashInventory().Items() {
      addItemInfo(w, i)
   }
}

func addItemInfo(w *response.Writer, i domain.Item) {
   addItemInfoZero(w, i, false)
}

func addItemInfoZero(w *response.Writer, i domain.Item, zeroPosition bool) {
   if !zeroPosition {
      w.WriteInt8(i.Slot() + 1)
   }
   w.WriteByte(2)
   w.WriteInt(i.ItemId())
   w.WriteBool(false)
   w.WriteLong(uint64(getTime(i.Expiration())))
   w.WriteShort(i.Quantity())
   w.WriteAsciiString(i.Owner())
   w.WriteShort(i.Flag())
}

func addEquipmentInfo(w *response.Writer, e domain.EquippedItem) {
   addEquipmentInfoZero(w, e, false)
}

func addEquipmentInfoZero(w *response.Writer, e domain.EquippedItem, zeroPosition bool) {
   slot := e.Slot()
   if !zeroPosition {
      if slot < 0 {
         slot *= -1
      }
      if slot > 100 {
         w.WriteShort(uint16(slot - 100))
      } else {
         w.WriteShort(uint16(slot))
      }
   }

   w.WriteByte(1)
   w.WriteInt(e.ItemId())
   w.WriteBool(false)
   w.WriteLong(uint64(getTime(e.Expiration())))
   w.WriteByte(e.Slots())
   w.WriteByte(e.Level())
   w.WriteShort(e.Strength())
   w.WriteShort(e.Dexterity())
   w.WriteShort(e.Intelligence())
   w.WriteShort(e.Luck())
   w.WriteShort(e.Hp())
   w.WriteShort(e.Mp())
   w.WriteShort(e.WeaponAttack())
   w.WriteShort(e.MagicAttack())
   w.WriteShort(e.WeaponDefense())
   w.WriteShort(e.MagicDefense())
   w.WriteShort(e.Accuracy())
   w.WriteShort(e.Avoidability())
   w.WriteShort(e.Hands())
   w.WriteShort(e.Speed())
   w.WriteShort(e.Jump())
   w.WriteAsciiString(e.OwnerName())
   w.WriteShort(e.Flags())

   w.WriteByte(0)
   w.WriteByte(0)
   w.WriteInt(0)
   w.WriteInt(0)
   w.WriteLong(0)
   w.WriteLong(uint64(getTime(-2)))
   w.WriteInt32(-1)
}

func addCharacterStats(w *response.Writer, character domain.Character) {
   w.WriteInt(character.Attributes().Id())
   addPaddedCharacterName(w, character)
   w.WriteByte(character.Attributes().Gender())
   w.WriteByte(character.Attributes().SkinColor())
   w.WriteInt(character.Attributes().Face())
   w.WriteInt(character.Attributes().Hair())
   writeForEachPet(w, character.Pets(), writePetId, writeEmptyPetId)
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

func addRemainingSkillInfo(w *response.Writer, character domain.Character) {

}

func addPaddedCharacterName(w *response.Writer, character domain.Character) {
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

func writeForEachPet(w *response.Writer, ps []domain.Pet, pe func(w *response.Writer, p domain.Pet), pne func(w *response.Writer)) {
   for i := 0; i < 3; i++ {
      if ps != nil && len(ps) > i {
         pe(w, ps[i])
      } else {
         pne(w)
      }
   }
}

func writePetId(w *response.Writer, pet domain.Pet) {
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

   ftUtOffset := 116444736010800000 + (10000 * time.Now().UnixNano())
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

func WriteWarpToMap(channelId byte, mapId uint32, portalId uint32, hp uint16) []byte {
   w := response.NewWriter()
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