package consumer

import (
	"context"
	"fmt"
	"time"

	envelope "github.com/utilitywarehouse/energy-contracts/pkg/generated"
	"github.com/utilitywarehouse/energy-contracts/pkg/generated/platform"
	"github.com/utilitywarehouse/energy-onboarding-lynnluong/internal/models"
	"github.com/utilitywarehouse/energy-onboarding-lynnluong/internal/store"
	"github.com/utilitywarehouse/energy-pkg/domain"
	"github.com/uw-labs/substrate"
	"google.golang.org/protobuf/proto"
)

type Consumer struct {
	Db store.Projector
}

func (c *Consumer) HandleBatch(ctx context.Context, messages []substrate.Message) error {
	c.Db.Begin()
	for _, message := range messages {
		if err := c.handleMessage(message); err != nil {
			return err
		}
	}
	return c.Db.Commit(ctx)
}

func (c *Consumer) handleMessage(msg substrate.Message) error {
	var env envelope.Envelope
	err := proto.Unmarshal(msg.Data(), &env)
	if err != nil {
		return fmt.Errorf("%w failed to unmarshal", err)
	}

	if env.Message == nil {
		return fmt.Errorf("message is empty")
	}

	payload, err := env.GetMessage().UnmarshalNew()
	if err != nil {
		if env.Message.TypeUrl == "type.googleapis.com/IgnitionRegistrationSubmittedEvent" {
			return nil // Skipping specific messages as per the original logic
		}
		return fmt.Errorf("%w failed to unmarshal new, message: %+v", err, env.GetMessage())
	}

	// Use EventToServiceState to get the appropriate service state
	serviceState, err := domain.EventToServiceState(payload)
	if err != nil {
		return fmt.Errorf("failed to map event to service state: %w", err)
	}

	// Determine the service ID from the payload
	serviceId := getServiceIdFromPayload(payload)
	if serviceId == "" {
		return fmt.Errorf("failed to extract service ID from payload")
	}

	// Update the service in the database using the mapped service state
	err = c.Db.AddToBatch(models.EnergyService{
		ServiceId:  serviceId,
		UpdatedAt:  time.Now(),
		OccurredAt: time.Now(),
		State:      serviceState.String(),
	})
	if err != nil {
		return fmt.Errorf("%w failed to add to batch", err)
	}

	return nil
}

// getServiceIdFromPayload is a helper function to extract the service ID from the event payload
func getServiceIdFromPayload(payload proto.Message) string {
	switch ev := payload.(type) {
	case *platform.ElectricityServiceRequestReceivedEvent:
		return ev.ServiceId
	case *platform.ElectricityServiceRequestReleasedEvent:
		return ev.ServiceId
	default:
		return ""
	}
}