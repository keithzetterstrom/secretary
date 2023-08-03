package usecase

import (
	"context"

	"github.com/keithzetterstrom/secretary/internal/pkg/models"
	"github.com/keithzetterstrom/secretary/internal/pkg/user"
)

type usecase struct {
	repo user.Repository
}

func New(repo user.Repository) user.Usecase {
	return &usecase{repo: repo}
}

func (u *usecase) Registration(ctx context.Context, registration models.UserRegistration) error {
	return u.repo.Registration(ctx, registration)
}

func (u *usecase) Get(ctx context.Context, tgUserName string) (*models.UserRegistration, error) {
	return u.repo.Get(ctx, tgUserName)
}

func (u *usecase) Activate(ctx context.Context, tgUserName string) error {
	return u.repo.Activate(ctx, tgUserName)
}
