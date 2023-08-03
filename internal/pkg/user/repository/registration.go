package repository

import (
	"context"
	"fmt"

	"github.com/pkg/errors"

	intererrors "github.com/keithzetterstrom/secretary/internal/errors"
	"github.com/keithzetterstrom/secretary/internal/pkg/models"
	"github.com/keithzetterstrom/secretary/internal/pkg/user"
	"github.com/keithzetterstrom/secretary/internal/repository/docs"
	repomodels "github.com/keithzetterstrom/secretary/internal/repository/models"
)

type repository struct {
	docsClient *docs.Client
}

func New(docsClient *docs.Client) user.Repository {
	return &repository{docsClient: docsClient}
}

func (r *repository) Registration(ctx context.Context, registration models.UserRegistration) error {
	ur := repomodels.ModelToUserRegistration(&registration)

	err := r.docsClient.Append(ctx, ur)
	if err != nil {
		return errors.Wrap(err, "failed to reg user")
	}

	_, err = r.docsClient.Get(ctx)
	if err != nil {
		return err
	}

	return nil
}

func (r *repository) Get(ctx context.Context, tgUserName string) (*models.UserRegistration, error) {
	values, err := r.docsClient.Get(ctx)
	if err != nil {
		return nil, err
	}

	for _, row := range values {
		if row[2] == tgUserName {
			return &models.UserRegistration{
				Email:      row[0].(string),
				Name:       row[1].(string),
				TgUserName: row[2].(string),
			}, nil
		}
	}

	return nil, intererrors.ErrUserNotFound
}

func (r *repository) Activate(ctx context.Context, tgUserName string) error {
	values, err := r.docsClient.Get(ctx)
	if err != nil {
		return err
	}

	usr, i := searchUser(tgUserName, values)
	if usr == nil {
		return intererrors.ErrUserNotFound
	}

	rng := fmt.Sprintf("D%d", i)

	err = r.docsClient.UpdateValue(ctx, rng, true)
	if err != nil {
		return err
	}

	return nil
}

func searchUser(tgUserName string, values [][]interface{}) (*models.UserRegistration, int) {
	for i, row := range values {
		if row[2] == tgUserName {
			return &models.UserRegistration{
				Email:      row[0].(string),
				Name:       row[1].(string),
				TgUserName: row[2].(string),
			}, i + 1
		}
	}

	return nil, 0
}
