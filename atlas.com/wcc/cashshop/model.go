package cashshop

type SpecialCashItem struct {
	sn       uint32
	modifier uint32
	info     byte
}

func (s *SpecialCashItem) SN() uint32 {
	return s.sn
}

func (s *SpecialCashItem) Modifier() uint32 {
	return s.modifier
}

func (s *SpecialCashItem) Info() byte {
	return s.info
}

type Character struct {
}