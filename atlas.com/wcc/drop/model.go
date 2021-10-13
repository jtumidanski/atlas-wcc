package drop

type Model struct {
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

func (d Model) UniqueId() uint32 {
	return d.uniqueId
}

func (d Model) ItemId() uint32 {
	return d.itemId
}

func (d Model) Meso() uint32 {
	return d.meso
}

func (d Model) DropperUniqueId() uint32 {
	return d.dropperUniqueId
}

func (d Model) DropType() byte {
	return d.dropType
}

func (d Model) OwnerId() uint32 {
	return d.ownerId
}

func (d Model) OwnerPartyId() uint32 {
	return d.ownerPartyId
}

func (d Model) DropTime() uint64 {
	return d.dropTime
}

func (d Model) DropX() int16 {
	return d.dropX
}

func (d Model) DropY() int16 {
	return d.dropY
}

func (d Model) DropperX() int16 {
	return d.dropperX
}

func (d Model) DropperY() int16 {
	return d.dropperY
}

func (d Model) CharacterDrop() bool {
	return d.characterDrop
}

func (d Model) Mod() byte {
	return d.mod
}

type builder struct {
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

func NewBuilder() *builder {
	return &builder{}
}

func (d *builder) SetWorldId(worldId byte) *builder {
	d.worldId = worldId
	return d
}

func (d *builder) SetChannelId(channelId byte) *builder {
	d.channelId = channelId
	return d
}

func (d *builder) SetMapId(mapId uint32) *builder {
	d.mapId = mapId
	return d
}

func (d *builder) SetUniqueId(uniqueId uint32) *builder {
	d.uniqueId = uniqueId
	return d
}

func (d *builder) SetItemId(itemId uint32) *builder {
	d.itemId = itemId
	return d
}
func (d *builder) SetQuantity(quantity uint32) *builder {
	d.quantity = quantity
	return d
}

func (d *builder) SetMeso(meso uint32) *builder {
	d.meso = meso
	return d
}

func (d *builder) SetDropType(dropType byte) *builder {
	d.dropType = dropType
	return d
}

func (d *builder) SetDropX(dropX int16) *builder {
	d.dropX = dropX
	return d
}

func (d *builder) SetDropY(dropY int16) *builder {
	d.dropY = dropY
	return d
}

func (d *builder) SetOwnerId(ownerId uint32) *builder {
	d.ownerId = ownerId
	return d
}

func (d *builder) SetOwnerPartyId(ownerPartyId uint32) *builder {
	d.ownerPartyId = ownerPartyId
	return d
}

func (d *builder) SetDropTime(dropTime uint64) *builder {
	d.dropTime = dropTime
	return d
}

func (d *builder) SetDropperUniqueId(dropperUniqueId uint32) *builder {
	d.dropperUniqueId = dropperUniqueId
	return d
}

func (d *builder) SetDropperX(dropperX int16) *builder {
	d.dropperX = dropperX
	return d
}

func (d *builder) SetDropperY(dropperY int16) *builder {
	d.dropperY = dropperY
	return d
}

func (d *builder) SetCharacterDrop(characterDrop bool) *builder {
	d.characterDrop = characterDrop
	return d
}

func (d *builder) SetMod(mod byte) *builder {
	d.mod = mod
	return d
}

func (d *builder) Build() Model {
	return Model{
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
