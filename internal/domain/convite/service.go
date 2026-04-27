package convite

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

type ConviteService struct {
	respository ConviteRepository
}

func NewConviteService(respository ConviteRepository) *ConviteService {
	return &ConviteService{respository: respository}
}

func (s *ConviteService) CreateConvite(ctx context.Context, colaboradorId string) error {
	id, err := uuid.Parse(colaboradorId)
	if err != nil {
		return fmt.Errorf("UUID inválido: %v", err)
	}

	_, err = s.respository.Store(ctx, id)
	return err
}

func (s *ConviteService) GetAllConvites(ctx context.Context) ([]Convite, error) {
	return s.respository.FindAll(ctx)
}

func (s *ConviteService) DisableConvite(ctx context.Context, token string) error {
	tk, err := uuid.Parse(token)
	if err != nil {
		return fmt.Errorf("UUID inválido: %v", err)
	}

	return s.respository.Disable(ctx, tk)
}
