package store

import (
    "context"
    "github.com/utilitywarehouse/energy-onboarding-lynnluong/internal/models"
)

// Projector interface definition
type Projector interface {
    Begin()
    AddToBatch(svc models.EnergyService) error
    Commit(ctx context.Context) error
}