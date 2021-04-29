package domain

type Drop struct {
	worldId         byte
	channelId       byte
	mapId           uint32
	uniqueId        uint32
	itemId          uint32
	quantity        uint32
	meso            uint32
	dropType        byte
	dropX           int16
	dropY           int16
	ownerId         uint32
	ownerPartyId    uint32
	dropTime        uint64
	dropperUniqueId uint32
	dropperX        int16
	dropperY        int16
	characterDrop   bool
	mod             byte
}

func (d Drop) UniqueId() uint32 {
	return d.uniqueId
}

func (d Drop) ItemId() uint32 {
	return d.itemId
}

func (d Drop) Meso() uint32 {
	return d.meso
}

func (d Drop) DropperUniqueId() uint32 {
	return d.dropperUniqueId
}

func (d Drop) DropType() byte {
	return d.dropType
}

func (d Drop) OwnerId() uint32 {
	return d.ownerId
}

func (d Drop) OwnerPartyId() uint32 {
	return d.ownerPartyId
}

func (d Drop) DropTime() uint64 {
	return d.dropTime
}

func (d Drop) DropX() int16 {
	return d.dropX
}

func (d Drop) DropY() int16 {
	return d.dropY
}

func (d Drop) DropperX() int16 {
	return d.dropperX
}

func (d Drop) DropperY() int16 {
	return d.dropperY
}

func (d Drop) CharacterDrop() bool {
	return d.characterDrop
}

func (d Drop) Mod() byte {
	return d.mod
}

type dropBuilder struct {
	worldId         byte
	channelId       byte
	mapId           uint32
	uniqueId        uint32
	itemId          uint32
	quantity        uint32
	meso            uint32
	dropType        byte
	dropX           int16
	dropY           int16
	ownerId         uint32
	ownerPartyId    uint32
	dropTime        uint64
	dropperUniqueId uint32
	dropperX        int16
	dropperY        int16
	characterDrop   bool
	mod             byte
}

func NewDropBuilder() *dropBuilder {
	return &dropBuilder{}
}

func (d *dropBuilder) SetWorldId(worldId byte) *dropBuilder {
	d.worldId = worldId
	return d
}

func (d *dropBuilder) SetChannelId(channelId byte) *dropBuilder {
	d.channelId = channelId
	return d
}

func (d *dropBuilder) SetMapId(mapId uint32) *dropBuilder {
	d.mapId = mapId
	return d
}

func (d *dropBuilder) SetUniqueId(uniqueId uint32) *dropBuilder {
	d.uniqueId = uniqueId
	return d
}

func (d *dropBuilder) SetItemId(itemId uint32) *dropBuilder {
	d.itemId = itemId
	return d
}
func (d *dropBuilder) SetQuantity(quantity uint32) *dropBuilder {
	d.quantity = quantity
	return d
}

func (d *dropBuilder) SetMeso(meso uint32) *dropBuilder {
	d.meso = meso
	return d
}

func (d *dropBuilder) SetDropType(dropType byte) *dropBuilder {
	d.dropType = dropType
	return d
}

func (d *dropBuilder) SetDropX(dropX int16) *dropBuilder {
	d.dropX = dropX
	return d
}

func (d *dropBuilder) SetDropY(dropY int16) *dropBuilder {
	d.dropY = dropY
	return d
}

func (d *dropBuilder) SetOwnerId(ownerId uint32) *dropBuilder {
	d.ownerId = ownerId
	return d
}

func (d *dropBuilder) SetOwnerPartyId(ownerPartyId uint32) *dropBuilder {
	d.ownerPartyId = ownerPartyId
	return d
}

func (d *dropBuilder) SetDropTime(dropTime uint64) *dropBuilder {
	d.dropTime = dropTime
	return d
}

func (d *dropBuilder) SetDropperUniqueId(dropperUniqueId uint32) *dropBuilder {
	d.dropperUniqueId = dropperUniqueId
	return d
}

func (d *dropBuilder) SetDropperX(dropperX int16) *dropBuilder {
	d.dropperX = dropperX
	return d
}

func (d *dropBuilder) SetDropperY(dropperY int16) *dropBuilder {
	d.dropperY = dropperY
	return d
}

func (d *dropBuilder) SetCharacterDrop(characterDrop bool) *dropBuilder {
	d.characterDrop = characterDrop
	return d
}

func (d *dropBuilder) SetMod(mod byte) *dropBuilder {
	d.mod = mod
	return d
}

func (d *dropBuilder) Build() Drop {
	return Drop{
		uniqueId:        d.uniqueId,
		worldId:         d.worldId,
		channelId:       d.channelId,
		mapId:           d.mapId,
		itemId:          d.itemId,
		quantity:        d.quantity,
		meso:            d.meso,
		dropType:        d.dropType,
		dropX:           d.dropX,
		dropY:           d.dropY,
		ownerId:         d.ownerId,
		ownerPartyId:    d.ownerPartyId,
		dropTime:        d.dropTime,
		dropperUniqueId: d.dropperUniqueId,
		dropperX:        d.dropperX,
		dropperY:        d.dropperY,
		characterDrop:   d.characterDrop,
		mod:             d.mod,
	}
}
