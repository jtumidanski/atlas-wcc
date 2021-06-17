package mapleSession

import (
	"atlas-wcc/socket/response/writer"
	"github.com/jtumidanski/atlas-socket/crypto"
	"math/rand"
	"net"
	"sync"
	"time"
)

type MapleSession struct {
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

func NewSession(id uint32, con net.Conn) *MapleSession {
	recvIv := []byte{70, 114, 122, 82}
	sendIv := []byte{82, 48, 120, 115}
	recvIv[3] = byte(rand.Float64() * 255)
	sendIv[3] = byte(rand.Float64() * 255)
	send := crypto.NewAESOFB(sendIv, uint16(65535)-version)
	recv := crypto.NewAESOFB(recvIv, version)
	return &MapleSession{id, 0, 0, 0, 0, false, con, *send, *recv, time.Now(), sync.RWMutex{}}
}

func (s *MapleSession) SetAccountId(accountId uint32) {
	s.accountId = accountId
}

func (s *MapleSession) SessionId() uint32 {
	return s.id
}

func (s *MapleSession) AccountId() uint32 {
	return s.accountId
}

func (s *MapleSession) Announce(b []byte) error {
	s.mutex.Lock()
	tmp := make([]byte, len(b)+4)
	copy(tmp, b)
	tmp = append([]byte{0, 0, 0, 0}, b...)
	tmp = s.send.Encrypt(tmp, true, true)
	_, err := s.con.Write(tmp)
	s.mutex.Unlock()
	return err
}

func (s *MapleSession) announce(b []byte) error {
	_, err := s.con.Write(b)
	return err
}

func (s *MapleSession) WriteHello() {
	s.announce(writer.WriteHello(version, s.send.IV(), s.recv.IV()))
}

func (s *MapleSession) ReceiveAESOFB() *crypto.AESOFB {
	return &s.recv
}

func (s *MapleSession) GetRemoteAddress() net.Addr {
	return s.con.RemoteAddr()
}

func (s *MapleSession) SetWorldId(worldId byte) {
	s.worldId = worldId
}

func (s *MapleSession) SetChannelId(channelId byte) {
	s.channelId = channelId
}

func (s *MapleSession) WorldId() byte {
	return s.worldId
}

func (s *MapleSession) ChannelId() byte {
	return s.channelId
}

func (s *MapleSession) UpdateLastRequest() {
	s.lastPacket = time.Now()
}

func (s *MapleSession) LastRequest() time.Time {
	return s.lastPacket
}

func (s *MapleSession) Disconnect() {
	_ = s.con.Close()
}

func (s *MapleSession) CharacterId() uint32 {
	return s.characterId
}

func (s *MapleSession) SetCharacterId(id uint32) {
	s.characterId = id
}

func (s *MapleSession) SetGm(gm bool) {
	s.gm = gm
}

func (s *MapleSession) GM() bool {
	return s.gm
}
