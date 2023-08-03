package user

import (
	"context"

	"github.com/keithzetterstrom/secretary/internal/pkg/models"
)

type Usecase interface {
	Registration(ctx context.Context, registration models.UserRegistration) error
	Get(ctx context.Context, tgUserName string) (*models.UserRegistration, error)
	Activate(ctx context.Context, tgUserName string) error
}

type Repository interface {
	Registration(ctx context.Context, registration models.UserRegistration) error
	Get(ctx context.Context, tgUserName string) (*models.UserRegistration, error)
	Activate(ctx context.Context, tgUserName string) error
}
