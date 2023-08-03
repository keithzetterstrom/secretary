package errors

import errorutil "github.com/keithzetterstrom/secretary/utils/error"

var ErrValidation = errorutil.NewTypedError("validation", "validation")
var ErrAuthFailed = errorutil.NewTypedError("user_auth_failed", "user auth failed")
var ErrUserNotFound = errorutil.NewTypedError("user_not_found", "user not found")
var ErrNotActivated = errorutil.NewTypedError("user_not_activated", "user account not activated")
var ErrPaymentAlreadyExists = errorutil.NewTypedError("payment_already_exists", "payment already exists")
var ErrTariffNotFound = errorutil.NewTypedError("tariff_not_found", "tariff not found")
