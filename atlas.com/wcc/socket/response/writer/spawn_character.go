package writer

import (
   "atlas-wcc/character"
   "atlas-wcc/inventory"
   "atlas-wcc/pet"
   "atlas-wcc/socket/response"
   "math/rand"
)

const OpCodeSpawnCharacter uint16 = 0xA0

func WriteSpawnCharacter(target character.Model, character character.Model, enteringField bool) []byte {
   w := response.NewWriter()
   w.WriteShort(OpCodeSpawnCharacter)
   w.WriteInt(character.Attributes().Id())
   w.WriteByte(character.Attributes().Level())
   w.WriteAsciiString(character.Attributes().Name())
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

   writeForeignBuffs(w, character)
   w.WriteShort(character.Attributes().JobId())
   addCharacterLook(w, character, false)
   //writer.writeInt(chr.getInventory(MapleInventoryType.CASH).countById(5110000));
   w.WriteInt(0)
   //writer.writeInt(chr.getItemEffect());
   w.WriteInt(0)
   //writer.writeInt(ItemConstants.getInventoryType(chr.getChair()) == MapleInventoryType.SETUP ? chr.getChair() : 0);
   w.WriteInt(0)

   if enteringField {
      w.WriteInt16(character.Attributes().X())
      w.WriteInt16(character.Attributes().Y() - int16(42))
      w.WriteByte(6)
   } else {
      w.WriteInt16(character.Attributes().X())
      w.WriteInt16(character.Attributes().Y())
      w.WriteByte(character.Attributes().Stance())
   }

   w.WriteShort(0) //chr.getFh()
   w.WriteByte(0)

   writeForEachPet(w, character.Pets(), addPetInfoButDoNotShow, noOpWrite)

   //end of pets
   w.WriteByte(0)

   //      if (chr.getMount() == null) {
   w.WriteInt(1) // mob level
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
   addRingLook(w, character, true)  // crush
   addRingLook(w, character, false) // friendship
   addMarriageRingLook(w, target, character)
   encodeNewYearCardInfo(w, character)  // new year seems to crash sometimes...
   w.WriteByte(0)
   w.WriteByte(0)
   //writer.write(chr.getTeam());//only needed in specific fields
   w.WriteByte(0)


   return w.Bytes()
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

func encodeNewYearCardInfo(w *response.Writer, character character.Model) {
   w.WriteByte(0)
}

func addMarriageRingLook(w *response.Writer, target character.Model, character character.Model) {
   w.WriteByte(0)
}

func addRingLook(w *response.Writer, character character.Model, crush bool) {
   w.WriteByte(0)
}

func addCharacterLook(w *response.Writer, character character.Model, mega bool) {
   w.WriteByte(character.Attributes().Gender())
   w.WriteByte(character.Attributes().SkinColor())
   w.WriteInt(character.Attributes().Face())
   if mega {
      w.WriteByte(0)
   } else {
      w.WriteByte(1)
   }
   w.WriteInt(character.Attributes().Hair())
   addCharacterEquips(w, character)
}

func addCharacterEquips(w *response.Writer, character character.Model) {
   var equips = getEquippedItemSlotMap(character.Equipment())
   var maskedEquips = make(map[int16]uint32)
   writeEquips(w, equips, maskedEquips)

   var weapon *inventory.EquippedItem
   for _, x := range character.Equipment() {
      if x.InWeaponSlot() {
         weapon = &x
         break
      }
   }
   if weapon != nil {
      w.WriteInt(weapon.ItemId())
   } else {
      w.WriteInt(0)
   }

   writeForEachPet(w, character.Pets(), writePetItemId, writeEmptyPetItemId)
}

func writeEquips(w *response.Writer, equips map[int16]uint32, maskedEquips map[int16]uint32) {
   for k, v := range equips {
      w.WriteKeyValue(byte(k), v)
   }
   w.WriteByte(0xFF)
   for k, v := range maskedEquips {
      w.WriteKeyValue(byte(k), v)
   }
   w.WriteByte(0xFF)
}

func getEquippedItemSlotMap(e []inventory.EquippedItem) map[int16]uint32 {
   var equips = make(map[int16]uint32, 0)
   for _, x := range e {
      if x.NotInWeaponSlot() {
         y := x.InvertSlot()
         equips[y.Slot()] = y.ItemId()
      }
   }
   return equips
}

func writePetItemId(w *response.Writer, p pet.Model) {
   w.WriteInt(p.ItemId())
}

func writeEmptyPetItemId(w *response.Writer) {
   w.WriteInt(0)
}

func writeForeignBuffs(w *response.Writer, character character.Model) {
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
