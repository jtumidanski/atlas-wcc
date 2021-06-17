package session

import (
	"atlas-wcc/socket/response/writer"
	"github.com/jtumidanski/atlas-socket/crypto"
	"math/rand"
	"net"
	"sync"
	"time"
)

type Model struct {
	id          uint32
	accountId   uint32
	worldId     byte
	channelId   byte
	characterId uint32
	gm          bool
	con         net.Conn
	send        crypto.AESOFB
	recv        crypto.AESOFB
	lastPacket  time.Time
	mutex       sync.RWMutex
}

const (
	version uint16 = 83
)

func NewSession(id uint32, con net.Conn) *Model {
	recvIv := []byte{70, 114, 122, 82}
	sendIv := []byte{82, 48, 120, 115}
	recvIv[3] = byte(rand.Float64() * 255)
	sendIv[3] = byte(rand.Float64() * 255)
	send := crypto.NewAESOFB(sendIv, uint16(65535)-version)
	recv := crypto.NewAESOFB(recvIv, version)
	return &Model{id, 0, 0, 0, 0, false, con, *send, *recv, time.Now(), sync.RWMutex{}}
}

func (s *Model) SetAccountId(accountId uint32) {
	s.accountId = accountId
}

func (s *Model) SessionId() uint32 {
	return s.id
}

func (s *Model) AccountId() uint32 {
	return s.accountId
}

func (s *Model) Announce(b []byte) error {
	s.mutex.Lock()
	tmp := make([]byte, len(b)+4)
	copy(tmp, b)
	tmp = append([]byte{0, 0, 0, 0}, b...)
	tmp = s.send.Encrypt(tmp, true, true)
	_, err := s.con.Write(tmp)
	s.mutex.Unlock()
	return err
}

func (s *Model) announce(b []byte) error {
	_, err := s.con.Write(b)
	return err
}

func (s *Model) WriteHello() {
	_ = s.announce(writer.WriteHello(version, s.send.IV(), s.recv.IV()))
}

func (s *Model) ReceiveAESOFB() *crypto.AESOFB {
	return &s.recv
}

func (s *Model) GetRemoteAddress() net.Addr {
	return s.con.RemoteAddr()
}

func (s *Model) SetWorldId(worldId byte) {
	s.worldId = worldId
}

func (s *Model) SetChannelId(channelId byte) {
	s.channelId = channelId
}

func (s *Model) WorldId() byte {
	return s.worldId
}

func (s *Model) ChannelId() byte {
	return s.channelId
}

func (s *Model) UpdateLastRequest() {
	s.lastPacket = time.Now()
}

func (s *Model) LastRequest() time.Time {
	return s.lastPacket
}

func (s *Model) Disconnect() {
	_ = s.con.Close()
}

func (s *Model) CharacterId() uint32 {
	return s.characterId
}

func (s *Model) SetCharacterId(id uint32) {
	s.characterId = id
}

func (s *Model) SetGm(gm bool) {
	s.gm = gm
}

func (s *Model) GM() bool {
	return s.gm
}
