package service

import "context"

// Service interface is implemented by various services
type Service interface {
	// Start the service, it should be implemented as a blocking call
	// method should return gracefully when ctx is Done
	Start(ctx context.Context) error

	// GetErrorChannel() will return the error channel for graceful service shutdown
	GetErrorChannel() chan error
}
