package store

import (
	"context"

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

func (s *Store) Commit(ctx context.Context) error {
	res := s.pool.SendBatch(ctx, s.batch)
	s.batch = nil
	return res.Close()
}
