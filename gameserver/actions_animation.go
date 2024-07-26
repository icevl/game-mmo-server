package gameserver

import (
	"fmt"
	"server/proto/animationpb"
	"server/types"
)

func ActionAnimation(world *World, client *types.TCPClient, action *animationpb.Animation) {

	source, err := world.getObject(client.UUID)
	if err != nil {
		fmt.Println("Error getting object from world")
		return
	}

	UpdateAnimationChannel <- &types.Animation{Object: source, Name: action.Name, Speed: action.Speed}

	fmt.Println("Animation action: ", client.UUID)
}
