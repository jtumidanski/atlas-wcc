package buff

import (
	"atlas-wcc/socket/response"
	"github.com/sirupsen/logrus"
)

const OpCodeShowBuff uint16 = 0x20
const OpCodeCancelBuff uint16 = 0x21

func WriteCancelBuff(l logrus.FieldLogger) func(stats []Stat) []byte {
	return func(stats []Stat) []byte {
		w := response.NewWriter(l)
		w.WriteShort(OpCodeCancelBuff)
		writeLongMask(w, stats)
		w.WriteByte(1)
		return w.Bytes()
	}
}

type Stat struct {
	first    bool
	buffMask uint64
	amount   uint16
}

func NewBuffStat(first bool, mask uint64, amount uint16) Stat {
	return Stat{first, mask, amount}
}

func (s Stat) Amount() uint16 {
	return s.amount
}

func WriteShowBuff(l logrus.FieldLogger) func(buffId uint32, buffDuration int32, stats []Stat, hasSpecial bool) []byte {
	return func(buffId uint32, buffDuration int32, stats []Stat, hasSpecial bool) []byte {
		w := response.NewWriter(l)
		w.WriteShort(OpCodeShowBuff)
		writeLongMask(w, stats)
		for _, stat := range stats {
			w.WriteShort(stat.Amount())
			w.WriteInt(buffId)
			w.WriteInt32(buffDuration * 1000)
		}
		w.WriteInt(0)
		w.WriteByte(0)
		w.WriteInt(uint32(stats[0].Amount()))

		if hasSpecial {
			w.Skip(3)
		}
		return w.Bytes()
	}
}

func writeLongMask(w *response.Writer, stats []Stat) {
	var fm uint64
	var sm uint64

	for _, stat := range stats {
		if stat.first {
			fm |= stat.buffMask
		} else {
			sm |= stat.buffMask
		}
	}
	w.WriteLong(fm)
	w.WriteLong(sm)
}
