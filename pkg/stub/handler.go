package stub

import (
	"context"
	"fmt"

	"github.com/shawn-hurley/starter-pack-operator/pkg/apis/starterpack/v1alpha1"
	"github.com/shawn-hurley/starter-pack-operator/pkg/broker"

	"github.com/operator-framework/operator-sdk/pkg/sdk"
	log "github.com/sirupsen/logrus"
)

func NewHandler() sdk.Handler {
	return &Handler{}
}

type Handler struct {
	// Fill me
}

func (h *Handler) Handle(ctx context.Context, event sdk.Event) error {
	switch o := event.Object.(type) {
	case *v1alpha1.Broker:
		if !event.Deleted {
			err := broker.Reconcile(o)
			if err != nil {
				fmt.Printf("error reconciling broker - %v", err)
			}
			return err
		}
		log.Infof("deleted event: for %v", o.GetName())
	}
	return nil
}
