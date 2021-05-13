package writer

import "atlas-wcc/socket/response"

const OpCodeShowBuff uint16 = 0x20

type BuffStat struct {
	first    bool
	buffMask uint64
	amount   uint16
}

func NewBuffStat(first bool, mask uint64, amount uint16) BuffStat {
	return BuffStat{first, mask, amount}
}

func (s BuffStat) Amount() uint16 {
	return s.amount
}

func WriteShowBuff(buffId uint32, buffDuration int32, stats []BuffStat, hasSpecial bool) []byte {
	w := response.NewWriter()
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

func writeLongMask(w *response.Writer, stats []BuffStat) {
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
