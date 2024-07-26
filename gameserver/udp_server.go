package gameserver

import (
	"fmt"
	"log"
	"net"
	"server/proto/actionpb"
	"sync"
	"time"

	"google.golang.org/protobuf/proto"
)

type UDPClient struct {
	Addr *net.UDPAddr
	Conn *net.UDPConn
}

type UDPClientsState struct {
	sync.RWMutex
	clients map[string]*UDPClient
	world   *World
}

type UpdateTransform struct {
	transform  *actionpb.Action
	clientUUID string
}

const (
	udpPort = ":8000"
	bufSize = 1024
)

var UpdateTransformChan = make(chan UpdateTransform)
var UDPState = &UDPClientsState{}

func StartUDPServer() {

	UDPState = &UDPClientsState{
		clients: map[string]*UDPClient{},
		world:   W,
	}

	go processTransformsUpdates()

	addr, err := net.ResolveUDPAddr("udp", udpPort)
	if err != nil {
		log.Fatalf("Error in address resolving: %s", err)
	}

	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		log.Fatalf("Udp connection error: %s", err)
	}
	defer conn.Close()

	fmt.Println("UDP server started on port", udpPort)

	for {
		buf := make([]byte, bufSize)
		n, clientAddr, err := conn.ReadFromUDP(buf)
		if err != nil {
			log.Fatalf("Error in connection read: %s", err)
		}

		// fmt.Printf("Received %d bytes from %s\n", n, clientAddr.String())

		action := &actionpb.Action{}

		if err := proto.Unmarshal(buf[:n], action); err != nil {
			log.Printf("Error in unmarshalling: %s", err)
			continue
		}

		switch action.GetAction().(type) {

		case *actionpb.Action_Ping:
			ping := action.GetPing()
			UDPState.addClient(ping.UUID, clientAddr, conn)

		case *actionpb.Action_Transform:
			transform := action.GetTransform()

			obj, err := W.getObject(transform.UUID)
			if err != nil {
				log.Printf("Object not found: %s", transform.UUID)
				continue
			}

			W.Lock()

			obj.Position.X = float64(transform.Position.X)
			obj.Position.Y = float64(transform.Position.Y)
			obj.Position.Z = float64(transform.Position.Z)
			obj.Rotation.X = float64(transform.Rotation.X)
			obj.Rotation.Y = float64(transform.Rotation.Y)
			obj.Rotation.Z = float64(transform.Rotation.Z)
			obj.Speed = float32(transform.Speed)

			W.Unlock()

			nextStepTime := time.Now().Add(40 * time.Millisecond)

			if obj.NextTransformUpdateTime == nil || time.Now().After(*obj.NextTransformUpdateTime) {
				UpdateMovementChannel <- obj
				obj.NextTransformUpdateTime = &nextStepTime
			}

		default:
			log.Printf("Unknown action type")
		}

	}

}

func (c *UDPClientsState) addClient(uuid string, addr *net.UDPAddr, conn *net.UDPConn) {
	c.Lock()
	defer c.Unlock()
	isNewClient := c.clients[uuid] == nil

	if !isNewClient {
		return
	}

	fmt.Printf("*New client connected: %s\n", addr.String())

	c.clients[uuid] = &UDPClient{
		Addr: addr,
		Conn: conn,
	}
}

func (c *UDPClientsState) removeClient(uuid string) {
	c.Lock()
	defer c.Unlock()

	client := c.clients[uuid]
	if client == nil {
		return
	}

	delete(c.clients, uuid)
}

func (c *UDPClientsState) sendToClient(uuid string, event *actionpb.Action) {
	var client = c.clients[uuid]

	if client == nil {
		return
	}

	data, err := proto.Marshal(event)
	if err != nil {
		log.Fatalf("Serialization error: %s", err)
	}

	_, _ = client.Conn.WriteToUDP(data, client.Addr)
}

func processTransformsUpdates() {
	for update := range UpdateTransformChan {
		client := UDPState.clients[update.clientUUID]

		if client == nil {
			continue
		}

		UDPState.sendToClient(update.clientUUID, update.transform)
	}
}
