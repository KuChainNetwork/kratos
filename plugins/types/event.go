package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type Event struct {
	Type       string
	Attributes map[string]string
}

func FromSdkEvent(evt sdk.Event) Event {
	res := Event{
		Type:       evt.Type,
		Attributes: make(map[string]string, len(evt.Attributes)),
	}

	for _, attr := range evt.Attributes {
		res.Attributes[string(attr.Key)] = string(attr.Value)
	}

	return res
}
