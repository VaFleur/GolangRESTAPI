package user

import "REST_API_app/pkg/logging"

type Service struct {
	storage Storage
	logger  *logging.Logger
}
