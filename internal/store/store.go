package store

import (
	"context"
	"time"

	"github.com/utilitywarehouse/energy-onboarding-lynnluong/internal/models"
	"github.com/utilitywarehouse/energy-onboarding-lynnluong/internal/store/migrations"
	"github.com/utilitywarehouse/energy-pkg/postgres"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/utilitywarehouse/go-operational/op"
	"github.com/utilitywarehouse/go-ops-health-checks/v3/pkg/sqlhealth"
)

type Store struct {
	pool  *pgxpool.Pool
	batch *pgx.Batch
}

func Setup(ctx context.Context, dsn string) (*Store, error) {
	pool, err := postgres.Setup(ctx, dsn, migrations.Source)
	if err != nil {
		return nil, err
	}
	return New(pool), nil
}

func New(pool *pgxpool.Pool) *Store {
	return &Store{
		pool: pool,
	}
}

func (s *Store) Health() func(cr *op.CheckResponse) {
	return sqlhealth.NewCheck(s, "unable to connect to the DB")
}

func (s *Store) Ping() error {
	return s.pool.Ping(context.Background())
}

func (s *Store) Close() {
	s.pool.Close()
}

func (s *Store) Begin() {
	if s.batch != nil {
		panic("attempting to create new batch without cleaning previous")
	}
	s.batch = &pgx.Batch{}
}

func (s *Store) GetServiceById(ctx context.Context, ID string) (models.EnergyService, error) {
	row := s.pool.QueryRow(ctx, `
		SELECT service_id, state, occurred_at, updated_at
		FROM services
		WHERE service_id = $1
	`, ID)

	m := &models.EnergyService{}
	err := row.Scan(&m.ServiceId, &m.State, &m.OccurredAt, &m.UpdatedAt)
	return *m, err
}

func (s *Store) AddToBatch(svc models.EnergyService) error {
	q := `
      INSERT INTO services (
          service_id,
          state,
          occurred_at,
          updated_at
      )
      VALUES ($1, $2, $3, $4)
      ON CONFLICT (service_id) DO UPDATE
      SET state=$2,occurred_at=$3, updated_at=$4
  `
	s.batch.Queue(q, svc.ServiceId, svc.State, svc.OccurredAt, time.Now())
	return nil
}

func (s *Store) Commit(ctx context.Context) error {
	res := s.pool.SendBatch(ctx, s.batch)
	s.batch = nil
	return res.Close()
}