package handler

import (
   "atlas-wcc/kafka/producers"
   "atlas-wcc/mapleSession"
   "atlas-wcc/processors"
   "context"
   "github.com/jtumidanski/atlas-socket/request"
   "log"
)

const OpMoveCharacter uint16 = 0x29

type MoveCharacterRequest struct {
   hasMovement  bool
   movementData []interface{}
   movementList []byte
}

func ReadMoveCharacterRequest(reader *request.RequestReader) *MoveCharacterRequest {
   reader.Skip(9)

   movementDataStart := reader.Position()
   movementDataList := updatePosition(reader, 0)
   if len(movementDataList) == 0 {
      return nil
   }

   movementDataLength := reader.Position() - movementDataStart
   hasMovement := movementDataLength > 0

   movementList := make([]byte, 0)
   if hasMovement {
      reader.Seek(movementDataStart)
      for i := 0; i < movementDataLength; i++ {
         movementList = append(movementList, reader.ReadByte())
      }
   }
   return &MoveCharacterRequest{hasMovement, movementDataList, movementList}
}

type AbsoluteMovement struct {
   X     int16
   Y     int16
   State byte
}

type RelativeMovement struct {
   State byte
}

func updatePosition(reader *request.RequestReader, offset int16) []interface{} {
   commands := reader.ReadByte()

   mdl := make([]interface{}, 0)
   if commands < 1 {
      return mdl
   }

   for i := byte(0); i < commands; i++ {
      command := reader.ReadByte()
      switch command {
      case 0:
      case 5:
      case 17:
         x := reader.ReadInt16()
         y := reader.ReadInt16()
         reader.Skip(6)
         ns := reader.ReadByte()
         reader.ReadUint16()
         md := &AbsoluteMovement{
            X:     x,
            Y:     y + offset,
            State: ns,
         }
         mdl = append(mdl, md)
         break
      case 1:
      case 2:
      case 6:
      case 12:
      case 13:
      case 16:
      case 18:
      case 19:
      case 20:
      case 22:
         reader.Skip(4)
         ns := reader.ReadByte()
         reader.ReadUint16()
         md := &RelativeMovement{State: ns}
         mdl = append(mdl, md)
         break
      case 3:
      case 4:
      case 7:
      case 8:
      case 9:
      case 11:
         reader.Skip(8)
         ns := reader.ReadByte()
         md := &RelativeMovement{State: ns}
         mdl = append(mdl, md)
         break
      case 14:
         reader.Skip(9)
         break
      case 10:
         reader.ReadByte()
         break
      case 15:
         reader.Skip(12)
         ns := reader.ReadByte()
         reader.ReadUint16()
         md := &RelativeMovement{State: ns}
         mdl = append(mdl, md)
         break
      case 21:
         reader.Skip(3)
         break
      }
   }
   return mdl
}

type MoveCharacterHandler struct {
}

func (h *MoveCharacterHandler) IsValid(l *log.Logger, ms *mapleSession.MapleSession) bool {
   v := processors.IsLoggedIn((*ms).AccountId())
   if !v {
      l.Printf("[ERROR] attempting to process a [MoveCharacterRequest] when the account %d is not logged in.", (*ms).SessionId())
   }
   return v
}

func (h *MoveCharacterHandler) HandleRequest(l *log.Logger, ms *mapleSession.MapleSession, r *request.RequestReader) {
   p := ReadMoveCharacterRequest(r)
   if p == nil {
      return
   }

   summary := processMovementList(p.movementData)
   producers.NewCharacterMovement(l, context.Background()).EmitMovement((*ms).WorldId(), (*ms).ChannelId(), (*ms).CharacterId(), summary.X, summary.Y, summary.State, p.movementList)
}

func processMovementList(movementList []interface{}) MovementSummary {
   ms := &MovementSummary{0, 0, 0}
   for _, m := range movementList {
      ms = appendMovement(ms, m)
   }
   return *ms
}

func appendMovement(ms *MovementSummary, m interface{}) *MovementSummary {
   switch m.(type) {
   case AbsoluteMovement:
      am := m.(AbsoluteMovement)
      return &MovementSummary{am.X, am.Y, am.State}
   case RelativeMovement:
      rm := m.(RelativeMovement)
      return ms.SetState(rm.State)
   }
   return ms
}

type MovementSummary struct {
   X     int16
   Y     int16
   State byte
}

func (m *MovementSummary) SetState(state byte) *MovementSummary {
   return &MovementSummary{m.X, m.Y, state}
}
