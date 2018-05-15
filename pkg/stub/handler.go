package stub

import (
	"fmt"

	"github.com/shawn-hurley/starter-pack-operator/pkg/apis/starterpack/v1alpha1"
	"github.com/shawn-hurley/starter-pack-operator/pkg/broker"

	"github.com/operator-framework/operator-sdk/pkg/sdk/handler"
	"github.com/operator-framework/operator-sdk/pkg/sdk/types"
)

func NewHandler() handler.Handler {
	return &Handler{}
}

type Handler struct {
	// Fill me
}

func (h *Handler) Handle(ctx types.Context, event types.Event) error {
	switch o := event.Object.(type) {
	case *v1alpha1.Broker:
		err := broker.Reconcile(o)
		if err != nil {
			fmt.Printf("error reconciling broker - %v", err)
		}
		return err
	}
	return nil
}
