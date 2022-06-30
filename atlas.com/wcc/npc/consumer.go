package npc

import (
	"atlas-wcc/kafka"
	"atlas-wcc/session"
	"fmt"
	"github.com/opentracing/opentracing-go"
	"github.com/sirupsen/logrus"
)

const (
	consumerNameTalk       = "npc_talk_command"
	consumerNameTalkNumber = "npc_talk_num_command"
	consumerNameTalkStyle  = "npc_talk_style_command"
	topicTokenTalk         = "TOPIC_NPC_TALK_COMMAND"
	topicTokenTalkNumber   = "TOPIC_NPC_TALK_NUM_COMMAND"
	topicTokenTalkStyle    = "TOPIC_NPC_TALK_STYLE_COMMAND"
)

func TalkConsumer(wid byte, cid byte) func(groupId string) kafka.ConsumerConfig {
	return func(groupId string) kafka.ConsumerConfig {
		return kafka.NewConsumerConfig[talkCommand](consumerNameTalk, topicTokenTalk, groupId, handleTalk(wid, cid))
	}
}

type talkCommand struct {
	CharacterId uint32 `json:"characterId"`
	NPCId       uint32 `json:"npcId"`
	Message     string `json:"message"`
	Type        string `json:"type"`
	Speaker     string `json:"speaker"`
}

func handleTalk(_ byte, _ byte) kafka.HandlerFunc[talkCommand] {
	return func(l logrus.FieldLogger, span opentracing.Span, event talkCommand) {
		session.ForSessionByCharacterId(event.CharacterId, session.Announce(WriteNPCTalk(l)(event.NPCId, getNPCTalkType(event.Type), event.Message, getNPCTalkEnd(event.Type), getNPCTalkSpeaker(event.Speaker))))
	}
}

func getNPCTalkSpeaker(speaker string) byte {
	switch speaker {
	case "NPC_LEFT":
		return 0
	case "NPC_RIGHT":
		return 1
	case "CHARACTER_LEFT":
		return 2
	case "CHARACTER_RIGHT":
		return 3
	}
	panic(fmt.Sprintf("unsupported npc talk speaker %s", speaker))
}

func getNPCTalkEnd(t string) []byte {
	switch t {
	case "NEXT":
		return []byte{00, 01}
	case "PREVIOUS":
		return []byte{01, 00}
	case "NEXT_PREVIOUS":
		return []byte{01, 01}
	case "OK":
		return []byte{00, 00}
	case "YES_NO":
		return []byte{}
	case "ACCEPT_DECLINE":
		return []byte{}
	case "SIMPLE":
		return []byte{}
	}
	panic(fmt.Sprintf("unsupported talk type %s", t))
}

func getNPCTalkType(t string) byte {
	switch t {
	case "NEXT":
		return 0
	case "PREVIOUS":
		return 0
	case "NEXT_PREVIOUS":
		return 0
	case "OK":
		return 0
	case "YES_NO":
		return 1
	case "ACCEPT_DECLINE":
		return 0x0C
	case "SIMPLE":
		return 4
	}
	panic(fmt.Sprintf("unsupported talk type %s", t))
}

func TalkNumberConsumer(wid byte, cid byte) func(groupId string) kafka.ConsumerConfig {
	return func(groupId string) kafka.ConsumerConfig {
		return kafka.NewConsumerConfig[talkNumberCommand](consumerNameTalkNumber, topicTokenTalkNumber, groupId, handleTalkNumber(wid, cid))
	}
}

type talkNumberCommand struct {
	CharacterId  uint32 `json:"characterId"`
	NPCId        uint32 `json:"npcId"`
	Message      string `json:"message"`
	Type         string `json:"type"`
	Speaker      string `json:"speaker"`
	DefaultValue int32  `json:"defaultValue"`
	MinimumValue int32  `json:"minimumValue"`
	MaximumValue int32  `json:"maximumValue"`
}

func handleTalkNumber(_ byte, _ byte) kafka.HandlerFunc[talkNumberCommand] {
	return func(l logrus.FieldLogger, span opentracing.Span, event talkNumberCommand) {
		session.ForSessionByCharacterId(event.CharacterId, session.Announce(WriteNPCTalkNum(l)(event.NPCId, event.Message, event.DefaultValue, event.MinimumValue, event.MaximumValue)))
	}
}

func TalkStyleConsumer(wid byte, cid byte) func(groupId string) kafka.ConsumerConfig {
	return func(groupId string) kafka.ConsumerConfig {
		return kafka.NewConsumerConfig[talkStyleCommand](consumerNameTalkStyle, topicTokenTalkStyle, groupId, handleTalkStyle(wid, cid))
	}
}

type talkStyleCommand struct {
	CharacterId uint32   `json:"characterId"`
	NPCId       uint32   `json:"npcId"`
	Message     string   `json:"message"`
	Type        string   `json:"type"`
	Speaker     string   `json:"speaker"`
	Styles      []uint32 `json:"styles"`
}

func handleTalkStyle(_ byte, _ byte) kafka.HandlerFunc[talkStyleCommand] {
	return func(l logrus.FieldLogger, span opentracing.Span, event talkStyleCommand) {
		session.ForSessionByCharacterId(event.CharacterId, session.Announce(WriteNPCTalkStyle(l)(event.NPCId, event.Message, event.Styles)))
	}
}
