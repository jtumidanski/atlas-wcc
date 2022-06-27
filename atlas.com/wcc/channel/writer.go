package channel

import (
	"atlas-wcc/socket/response"
	"github.com/sirupsen/logrus"
	"strconv"
	"strings"
)

const OpCodeChangeChannel uint16 = 0x10

func WriteChangeChannel(l logrus.FieldLogger) func(ipAddress string, port uint16) []byte {
	return func(ipAddress string, port uint16) []byte {
		w := response.NewWriter(l)
		w.WriteShort(OpCodeChangeChannel)
		w.WriteByte(1)
		ob := ipAsByteArray(ipAddress)
		w.WriteByteArray(ob)
		w.WriteShort(port)
		return w.Bytes()
	}
}

func ipAsByteArray(ipAddress string) []byte {
	var ob = make([]byte, 0)
	os := strings.Split(ipAddress, ".")
	for _, x := range os {
		o, err := strconv.ParseUint(x, 10, 8)
		if err == nil {
			ob = append(ob, byte(o))
		}
	}
	return ob
}
