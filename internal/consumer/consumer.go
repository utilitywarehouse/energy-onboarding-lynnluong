package consumer

import (
	"context"

	"github.com/utilitywarehouse/energy-pkg/substratemessage"
	"github.com/uw-labs/substrate"
)

type Projector interface {
	Begin()
	Commit(context.Context) error
}

func Handler(db Projector) substratemessage.BatchHandlerFunc {
	return func(ctx context.Context, messages []substrate.Message) error {
		db.Begin()
		for _, message := range messages {
			if err := handleMessage(ctx, message); err != nil {
				return err
			}
		}

		return db.Commit(ctx)

	}
}

func handleMessage(ctx context.Context, msg substrate.Message) error {
	// do something here
	return nil
}
