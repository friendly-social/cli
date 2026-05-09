package profile

import (
	"context"

	sdk "github.com/friendly-social/golang-sdk"
)

type Service struct {
	client *sdk.Client
}

func NewService(client *sdk.Client) *Service {
	return &Service{
		client: client,
	}
}

func (s *Service) GetDetails(auth *sdk.Authorization) (*sdk.UserDetails, error) {
	details, err := s.client.GetSelfDetails(context.Background(), auth)
	if err != nil {
		return nil, err
	}

	return details, nil
}
