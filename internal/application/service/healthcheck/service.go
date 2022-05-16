package healthcheck

import (
	healthcheckservice "github.com/icaroribeiro/new-go-code-challenge-template-2/internal/core/ports/application/service/healthcheck"
	"gorm.io/gorm"
)

type Service struct {
	DB *gorm.DB
}

// New is the factory function that encapsulates the implementation related to healthcheck service.
func New(db *gorm.DB) healthcheckservice.IService {
	return &Service{
		DB: db,
	}
}

// GetStatus is the function that verifies if the service has started up correctly and is ready to accept requests.
// As a way of checking, it makes sure if the database connection is alive.
func (s *Service) GetStatus() error {
	db, err := s.DB.DB()
	if err != nil {
		return err
	}

	if err = db.Ping(); err != nil {
		return err
	}

	return nil
}
