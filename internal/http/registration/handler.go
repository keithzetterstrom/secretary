package registration

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"

	intererrors "github.com/keithzetterstrom/secretary/internal/errors"
	"github.com/keithzetterstrom/secretary/internal/http/models"
	"github.com/keithzetterstrom/secretary/internal/pkg/user"
	"github.com/keithzetterstrom/secretary/internal/validator"
	"github.com/keithzetterstrom/secretary/utils/logger"
)

type Handler struct {
	validate    *validator.Validator
	userUsecase user.Usecase
	log         logger.Logger
}

func New(
	validate *validator.Validator,
	userUsecase user.Usecase,
	log logger.Logger,
) *Handler {
	return &Handler{
		validate:    validate,
		userUsecase: userUsecase,
		log:         log,
	}
}

func (h *Handler) Registration(c echo.Context) error {
	userRegInput := new(models.UserRegistration)
	if err := c.Bind(userRegInput); err != nil {
		return errors.Wrapf(intererrors.ErrValidation, "failed to bind reg: %v", err)
	}

	err := h.validate.ValidateUserRegistration(*userRegInput)
	if err != nil {
		return errors.Wrapf(intererrors.ErrValidation, "invalid user reg: %v", err)
	}

	userReg := models.UserRegistrationToModel(userRegInput)

	err = h.userUsecase.Registration(c.Request().Context(), userReg)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, http.StatusOK)
}
