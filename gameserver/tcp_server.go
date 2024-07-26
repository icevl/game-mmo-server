package gameserver

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net"
	"os"
	"server/entity"
	"server/events"
	"server/proto/actionpb"
	"server/proto/objectpb"
	"server/types"
	"sync"
	"time"

	"github.com/google/uuid"
	"google.golang.org/protobuf/proto"
)

const (
	tcpPort = ":8001"
)

type TCPClientsState struct {
	sync.RWMutex
	clients map[string]*types.TCPClient
	world   *World
}

type SenderParams struct {
	UUID   string
	Action *actionpb.Action
}

var TCPState = &TCPClientsState{}
var SenderChannel = make(chan *SenderParams)

func StartTCPServer() {

	TCPState = &TCPClientsState{
		clients: map[string]*types.TCPClient{},
		world:   W,
	}

	listener, err := net.Listen("tcp", tcpPort)
	if err != nil {
		fmt.Println("Error in TCP listen:", err)
		os.Exit(1)
	}

	defer listener.Close()

	fmt.Println("TCP server started on port", tcpPort)

	go processSenderChannel()

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error in TCP incoming connection:", err)
			continue
		}

		go TCPState.handleConnection(conn)
	}

}

func (s *TCPClientsState) handleConnection(conn net.Conn) {
	connectionId := conn.RemoteAddr().String()
	uuid, _ := uuid.NewUUID()
	clientUUID := uuid.String()

	defer func() {
		fmt.Println("Client disconnected", connectionId, clientUUID)
		conn.Close()
		s.removeClient(clientUUID)
	}()

	fmt.Println("New connection", connectionId, clientUUID)

	s.addClient(clientUUID, &conn)
	go s.spawnPlayer(clientUUID)

	client := s.getClient(clientUUID)

	for {

		sizeBytes := make([]byte, 4)
		_, err := io.ReadFull(conn, sizeBytes)
		if err != nil {
			if err != io.EOF {
				fmt.Println("Error in reading message size:", err)
			}
			break
		}

		size := binary.LittleEndian.Uint32(sizeBytes)

		data := make([]byte, size)
		_, err = io.ReadFull(conn, data)
		if err != nil {
			fmt.Println("E2:", err)
			if err != io.EOF {
				fmt.Println("Error in TCP data reading:", err)
			}
			break
		}

		action := &actionpb.Action{}
		if err := proto.Unmarshal(data, action); err != nil {
			fmt.Println("Error unmarshaling Protobuf message:", err)
			continue
		}

		s.ProcessReceivedActions(client, action)
	}
}

func (s *TCPClientsState) addClient(uuid string, conn *net.Conn) {
	s.Lock()
	defer s.Unlock()
	isNewClient := s.clients[uuid] == nil

	if !isNewClient {
		return
	}

	writer := bufio.NewWriter(*conn)

	client := &types.TCPClient{Conn: conn, UUID: uuid, Writer: writer, Send: make(chan *actionpb.Action)}
	s.clients[uuid] = client

	go client.ProcessSenderChannel()
}

func (s *TCPClientsState) getClient(uuid string) *types.TCPClient {
	s.RLock()
	defer s.RUnlock()
	return s.clients[uuid]
}

func (s *TCPClientsState) removeClient(uuid string) {
	s.Lock()
	defer s.Unlock()

	client := s.clients[uuid]
	if client == nil {
		return
	}

	s.world.removeObject(client.UUID)

	close(client.Send)
	client.Writer.Flush()
	(*client.Conn).Close()

	event := events.GetDestroyObjectEventPayload(client.UUID)
	for id := range s.clients {
		if id == uuid {
			continue
		}
		s.sendToClient(id, event)
	}

	delete(s.clients, uuid)
}

func processSenderChannel() {
	for params := range SenderChannel {
		client, ok := TCPState.clients[params.UUID]
		if !ok {
			continue
		}

		data, err := proto.Marshal(params.Action)
		if err != nil {
			log.Printf("Serialization error: %s\n", err)
			// s.removeClient(uuid)
			continue
		}

		writer := client.Writer
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

func (s *TCPClientsState) sendToClient(uuid string, event *actionpb.Action) {
	// SenderChannel <- &SenderParams{UUID: uuid, Action: event}

	client, ok := TCPState.clients[uuid]

	if ok {
		client.Send <- event
	}
}

func (c *TCPClientsState) testMessageTick() {
	ticker := time.NewTicker(20 * time.Second)

	for range ticker.C {
		c.world.Lock()

		var messages = []string{"Greetings, traveler! Welcome to our realm. What brings you to these lands?", "Ah, a newcomer! Prepare for a thrilling journey in our mystical world.", "Hey there, adventurer! Care to join us on a quest for glory and treasure?", "What's up?", "Good day!", "Welcome, bold warrior! Let's vanquish foes and uncover hidden secrets together!", "Hail, explorer! Unveil ancient mysteries and be the champion in our epic tale"}

		for _, obj := range c.world.objects {
			randomMessage := messages[rand.Intn(len(messages))]

			for id := range c.clients {
				msg := events.GetMessageEventPayload(obj.UUID, "", randomMessage)
				c.sendToClient(id, msg)
			}
		}

		c.world.Unlock()
	}
}

func (c *TCPClientsState) spawnPlayer(uuid string) {
	connection, ok := c.clients[uuid]
	if !ok {
		return
	}

	fmt.Println("Spawning player", uuid)

	slots := make(map[string]entity.HumanSlot)
	// MaleHair2 MilCut
	slots["Hair"] = entity.HumanSlot{Recipe: "MilCut", Color: "#000000"}
	slots["Beard"] = entity.HumanSlot{Recipe: "MaleBeard1", Color: "#FFFFFF"}
	//slots["Legs"] = types.HumanSlot{Recipe: "MalePants"}
	slots["Legs"] = entity.HumanSlot{Recipe: "MaleSweatPants_Recipe", Color: "#000000"}
	slots["Feet"] = entity.HumanSlot{Recipe: "TallShoes_Black_Recipe"}
	//slots["Chest"] = types.HumanSlot{Recipe: "MaleChallengerTorso"}
	slots["Chest"] = entity.HumanSlot{Recipe: "MaleShirt2", Color: "#CACACA"}
	// slots["Cape"] = &actors.HumanSlot{Recipe: "CapeBasic"}

	teleport := c.world.getTeleport("main")
	if teleport == nil {
		fmt.Println("Teleport not found")
		return
	}

	playerObject := &types.GameObject{
		Entity: entity.Entity{
			Name:      "Player",
			Speed:     2,
			Health:    200,
			MaxHealth: 200,
			HumanCharacter: &entity.HumanCharacter{
				Gender: "male",
				Slots:  slots,
			},
			EquippedItems: &entity.EquippedItems{
				RightHand: entity.GetItem("dragon_axe"),
			},
		},
		UUID:     connection.UUID,
		Type:     types.ObjectTypePlayer,
		Position: types.Vector3{X: teleport.Position.X, Y: teleport.Position.Y, Z: teleport.Position.Z},
		Rotation: types.Vector3{X: 0, Y: teleport.Rotation.Y, Z: 0},
	}

	c.world.addObject(playerObject)
	c.world.updateNeighbors(playerObject)
	c.world.updateNeighborsNearObject(playerObject)

	// Send all objects to the new player
	mapObjectsBatch := &objectpb.ObjectStateBatch{ObjectStates: []*objectpb.ObjectState{}}
	for _, obj := range playerObject.GetMapObjectsNearby() {
		if obj.UUID == playerObject.UUID {
			continue
		}
		mapObjectsBatch.ObjectStates = append(mapObjectsBatch.ObjectStates, &objectpb.ObjectState{
			UUID:         obj.UUID,
			VariantIndex: obj.VariationIndex,
		})
	}

	// Send nearby NPCs, players, loot items
	objectsBatch := &objectpb.ObjectBatch{Object: []*objectpb.Object{}}

	// Send the player itself
	objectsBatch.Object = append(objectsBatch.Object, events.GetObjectEvent(playerObject, &types.EventPayloadOptions{IsSelf: true}))
	for _, obj := range playerObject.Neighbors {
		if obj.Type != types.ObjectTypePlayer && obj.Type != types.ObjectTypeNPC && obj.Type != types.ObjectTypeMapObject {
			continue
		}

		objectsBatch.Object = append(objectsBatch.Object, events.GetObjectEvent(obj, &types.EventPayloadOptions{}))
	}

	c.sendToClient(uuid, &actionpb.Action{
		Action: &actionpb.Action_ObjectBatch{
			ObjectBatch: objectsBatch,
		},
	})

	c.sendToClient(uuid, &actionpb.Action{
		Action: &actionpb.Action_ObjectStateBatch{
			ObjectStateBatch: mapObjectsBatch,
		},
	})

	// Broadcast neighbors players
	for _, obj := range playerObject.GetPlayersNearby() {
		if obj.UUID == playerObject.UUID {
			continue
		}

		c.sendToClient(obj.UUID, events.GetObjectEventPayload(playerObject, &types.EventPayloadOptions{}))

	}

}
