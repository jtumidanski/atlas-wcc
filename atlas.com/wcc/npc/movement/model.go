package movement

type npcAnimationRequest struct {
	objectId uint32
	second   byte
	third    byte
}

func (r npcAnimationRequest) ObjectId() uint32 {
	return r.objectId
}

func (r npcAnimationRequest) Second() byte {
	return r.second
}

func (r npcAnimationRequest) Third() byte {
	return r.third
}

type npcMoveRequest struct {
	movement []byte
}

func (r npcMoveRequest) Movement() []byte {
	return r.movement
}
