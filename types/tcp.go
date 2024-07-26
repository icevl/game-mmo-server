package types

import (
	"bufio"
	"encoding/binary"
	"log"
	"net"
	"server/proto/actionpb"

	"google.golang.org/protobuf/proto"
)

type TCPClient struct {
	Conn   *net.Conn
	UUID   string
	Writer *bufio.Writer
	Send   chan *actionpb.Action
}

func (c *TCPClient) ProcessSenderChannel() {
	for params := range c.Send {
		data, err := proto.Marshal(params)
		if err != nil {
			log.Printf("Serialization error: %s\n", err)
			continue
		}

		writer := c.Writer
		buf := make([]byte, 4)
		binary.LittleEndian.PutUint32(buf, uint32(len(data)))

		combined := append(buf, data...)

		written, err := writer.Write(combined)
		if err != nil || written < len(combined) {
			log.Printf("Error writing message length and message to client: %s\n", err)
			// s.removeClient(uuid)
		} else {
			writer.Flush()
		}
	}
}
