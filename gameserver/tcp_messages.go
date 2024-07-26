package gameserver

import (
	"fmt"
	"server/proto/actionpb"
	"server/types"
)

func (s *TCPClientsState) ProcessReceivedActions(client *types.TCPClient, action *actionpb.Action) {
	switch act := action.Action.(type) {
	case *actionpb.Action_Interact:
		ActionInteract(s.world, client, act.Interact)
	case *actionpb.Action_InteractWith:
		ActionInteractWith(s.world, client, act.InteractWith)
	case *actionpb.Action_Animation:
		ActionAnimation(s.world, client, act.Animation)
	default:
		fmt.Printf("Unknown action type received %+v\n", action)
	}
}
