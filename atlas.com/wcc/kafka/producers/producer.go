package producers

import (
	"atlas-wcc/kafka/topics"
	"atlas-wcc/retry"
	"context"
	"encoding/binary"
	"encoding/json"
	"github.com/segmentio/kafka-go"
	log "github.com/sirupsen/logrus"
	"os"
	"time"
)

func createKey(key int) []byte {
	b := make([]byte, 8)
	binary.LittleEndian.PutUint32(b, uint32(key))
	return b
}

func produceEvent(l log.FieldLogger, topicToken string, key []byte, event interface{}) {
	name := topics.GetRegistry().Get(l, topicToken)
	w := &kafka.Writer{
		Addr:         kafka.TCP(os.Getenv("BOOTSTRAP_SERVERS")),
		Topic:        name,
		Balancer:     &kafka.LeastBytes{},
		BatchTimeout: 50 * time.Millisecond,
	}

	r, err := json.Marshal(event)
	l.WithField("message", string(r)).Debugf("Writing message to topic %s.", name)
	if err != nil {
		l.WithError(err).Fatalf("Unable to marshall event for topic %s.", name)
	}

	writeMessage := func(attempt int) (bool, error) {
		err = w.WriteMessages(context.Background(), kafka.Message{
			Key:   key,
			Value: r,
		})
		if err != nil {
			l.Warnf("Unable to emit event on topic %s, will retry.", name)
			return true, err
		}
		return false, err
	}

	err = retry.Retry(writeMessage, 10)
	if err != nil {
		l.WithError(err).Fatalf("Unable to emit event on topic %s.", name)
	}
}
