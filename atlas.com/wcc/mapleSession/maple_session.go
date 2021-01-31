package mapleSession

import (
   "atlas-wcc/rest/requests"
   "atlas-wcc/socket/response/writer"
   "github.com/jtumidanski/atlas-socket/crypto"
   "github.com/jtumidanski/atlas-socket/session"
   "math/rand"
   "net"
   "time"
)

type MapleSession interface {
   session.Session
   AccountId() uint32
   SetAccountId(id uint32)
   Announce(response []byte) error
   WorldId() byte
   SetWorldId(id byte)
   ChannelId() byte
   SetChannelId(id byte)
   CharacterId() uint32
   SetCharacterId(id uint32)
   SetGm(gm bool)
}

type mapleSession struct {
   id          int
   accountId   uint32
   worldId     byte
   channelId   byte
   characterId uint32
   gm          bool
   con         net.Conn
   send        crypto.AESOFB
   recv        crypto.AESOFB
   lastPacket  time.Time
}

const (
   version uint16 = 83
)

func NewSession(id int, con net.Conn) MapleSession {
   recvIv := []byte{70, 114, 122, 82}
   sendIv := []byte{82, 48, 120, 115}
   recvIv[3] = byte(rand.Float64() * 255)
   sendIv[3] = byte(rand.Float64() * 255)
   send := crypto.NewAESOFB(sendIv, uint16(65535)-version)
   recv := crypto.NewAESOFB(recvIv, version)
   return &mapleSession{id, 0, 0, 0, 0, false, con, *send, *recv, time.Now()}
}

func (s *mapleSession) SetAccountId(accountId uint32) {
   s.accountId = accountId
}

func (s *mapleSession) SessionId() int {
   return s.id
}

func (s *mapleSession) AccountId() uint32 {
   return s.accountId
}

func (s *mapleSession) Announce(b []byte) error {
   tmp := make([]byte, len(b)+4)
   copy(tmp, b)
   tmp = append([]byte{0, 0, 0, 0}, b...)
   tmp = s.send.Encrypt(tmp, true, true)
   _, err := s.con.Write(tmp)
   return err
}

func (s *mapleSession) announce(b []byte) error {
   _, err := s.con.Write(b)
   return err
}

func (s *mapleSession) WriteHello() {
   s.announce(writer.WriteHello(version, s.send.IV(), s.recv.IV()))
}

func (s *mapleSession) ReceiveAESOFB() *crypto.AESOFB {
   return &s.recv
}

func (s *mapleSession) GetRemoteAddress() net.Addr {
   return s.con.RemoteAddr()
}

func (s *mapleSession) SetWorldId(worldId byte) {
   s.worldId = worldId
}

func (s *mapleSession) SetChannelId(channelId byte) {
   s.channelId = channelId
}

func (s *mapleSession) WorldId() byte {
   return s.worldId
}

func (s *mapleSession) ChannelId() byte {
   return s.channelId
}

func (s *mapleSession) UpdateLastRequest() {
   s.lastPacket = time.Now()
}

func (s *mapleSession) LastRequest() time.Time {
   return s.lastPacket
}

func (s *mapleSession) Disconnect() {
   _ = s.con.Close()
   requests.CreateLogout(s.accountId)
}

func (s *mapleSession) CharacterId() uint32 {
   return s.characterId
}

func (s *mapleSession) SetCharacterId(id uint32) {
   s.characterId = id
}

func (s *mapleSession) SetGm(gm bool) {
   s.gm = gm
}
