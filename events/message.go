package events

import (
	"server/proto/actionpb"
	"server/proto/messagepb"
)

func GetMessageEventPayload(fromUUID, toUUID, text string) *actionpb.Action {
	msg := &messagepb.Message{
		FromUuid: fromUUID,
		ToUuid:   toUUID,
		Type:     "text",
		Text:     text,
	}

	return &actionpb.Action{
		Action: &actionpb.Action_Message{
			Message: msg,
		},
	}
}
