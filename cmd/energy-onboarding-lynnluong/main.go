package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/utilitywarehouse/energy-onboarding-lynnluong/internal/consumer"
	"github.com/utilitywarehouse/energy-onboarding-lynnluong/internal/store"
	"github.com/utilitywarehouse/energy-pkg/substratemessage"
	"github.com/utilitywarehouse/go-ops-health-checks/v2/pkg/substratehealth"
	_ "github.com/uw-labs/substrate/kafka"
	_ "github.com/uw-labs/substrate/proximo"

	log "github.com/sirupsen/logrus"
	cli "github.com/urfave/cli/v2"
	"github.com/utilitywarehouse/energy-pkg/app"
	"github.com/utilitywarehouse/energy-pkg/ops"
	"golang.org/x/sync/errgroup"
)

const (
	appName = "energy-onboarding-lynnluong"
	appDesc = "Energy onboarding task"

	postgresDSN = "postgres-dsn"

	batchSize = "batch-size"
)

var gitHash string // populated at compile time

var application = &cli.App{
	Name: appName,
	Flags: app.DefaultFlags().
		WithKafka().
		WithCustom(
			&cli.StringFlag{Name: postgresDSN, Required: true, EnvVars: []string{"POSTGRES_DSN"}},
			&cli.IntFlag{Name: batchSize, Value: 1, EnvVars: []string{"BATCH_SIZE"}},
		),
	Before: app.Before,
	Action: action,
}

func main() {
	if err := application.Run(os.Args); err != nil {
		log.WithError(err).Panic("unable to run app")
	}
}

func action(c *cli.Context) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	opsServer := ops.Default().
		WithPort(c.Int(app.OpsPort)).
		WithHash(gitHash).
		WithDetails(appName, appDesc)

	log.WithField("git_hash", gitHash).Info("starting app")

	g, ctx := errgroup.WithContext(ctx)
	store, err := store.Setup(ctx, c.String(postgresDSN))
	if err != nil {
		return err
	}
	defer store.Close()
	opsServer.Add("store", store.Health())

	eventSource, err := app.GetKafkaSource(c, c.String(app.KafkaConsumerGroup), c.String(app.ServiceTopic))
	if err != nil {
		return fmt.Errorf("unable to connect to service source: %w", err)
	}
	opsServer.Add("event-source", substratehealth.NewCheck(eventSource, "unable to consume events"))

	g.Go(func() error {
		defer log.Info("consumer finished")
		return substratemessage.BatchConsumer(
			ctx,
			c.Int(batchSize),
			time.Second,
			eventSource,
			consumer.Handler(store),
		)
	})

	g.Go(func() error {
		return opsServer.Start(ctx)
	})

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	g.Go(func() error {
		defer log.Info("signal handler finished")
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-sigChan:
			cancel()
		}
		return nil
	})

	return g.Wait()
}
