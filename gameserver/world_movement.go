package gameserver

import (
	"server/events"
	"server/types"
)

func (w *World) onWalkUpdates(object *types.GameObject) {
	prevNeighbors := object.Neighbors

	w.moveObjectTo(object)

	newNeighbors := object.Neighbors
	added, removed := findChanges(prevNeighbors, newNeighbors)

	for _, neighbor := range added {
		if neighbor.Type != types.ObjectTypePlayer && neighbor.Type != types.ObjectTypeNPC && neighbor.Type != types.ObjectTypeMapObject {
			continue
		}

		// Transmit the object to all players that are new neighbors
		if neighbor.Type == types.ObjectTypePlayer {
			msg := events.GetObjectEventPayload(object, &types.EventPayloadOptions{})
			TCPState.sendToClient(neighbor.UUID, msg)
		}

		// Update neighbors to NPCs near player
		if neighbor.Type == types.ObjectTypeNPC {
			w.updateNeighborsNearObject(object)
		}

		// Transmit the new neighbors to the player
		msg := events.GetObjectEventPayload(neighbor, &types.EventPayloadOptions{})
		TCPState.sendToClient(object.UUID, msg)
	}

	// Destroy the object for all players that are no longer neighbors
	for _, neighbor := range removed {
		if neighbor.Type != types.ObjectTypePlayer && neighbor.Type != types.ObjectTypeNPC && neighbor.Type != types.ObjectTypeMapObject {
			continue
		}

		if neighbor.Type == types.ObjectTypePlayer {
			TCPState.sendToClient(neighbor.UUID, events.GetDestroyObjectEventPayload(object.UUID))
		}

		// Update neighbors to NPCs near player
		if neighbor.Type == types.ObjectTypeNPC || neighbor.Type == types.ObjectTypePlayer {
			w.updateNeighborsNearObject(neighbor)
		}

		TCPState.sendToClient(object.UUID, events.GetDestroyObjectEventPayload(neighbor.UUID))
	}
}
