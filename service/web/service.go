package web

import "context"

type Service struct {
	errChan chan error
}

func NewService() (*Service, error) {
	return &Service{
		errChan: make(chan error),
	}, nil
}

func (svc *Service) GetErrorChannel() chan error {
	return svc.errChan
}

func (svc *Service) Start(ctx context.Context) error {
	return nil
}
