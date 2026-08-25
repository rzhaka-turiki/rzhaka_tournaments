package service

import (
	"context"
	"crypto/sha512"
	"encoding/hex"
	"errors"

	"github.com/google/uuid"
	"github.com/rzhaka-turiki/rzhaka_tournaments/internal/client/apexverifier"
	"github.com/rzhaka-turiki/rzhaka_tournaments/internal/model"
	"github.com/rzhaka-turiki/rzhaka_tournaments/internal/repository"
)

func calculateNIDHash(uid string) string {
	hash := sha512.Sum512([]byte(uid))
	return hex.EncodeToString(hash[:])
}

type ApexAccountService interface {
	GetByUserID(ctx context.Context, userID uuid.UUID) (*model.ApexAccount, error)
	Bind(ctx context.Context, userID uuid.UUID, player string, platform string, level int32) (*model.ApexAccount, error)
	Delete(ctx context.Context, userID uuid.UUID) error
}

type apexAccountService struct {
	apexAccountRepository repository.ApexAccountRepository
	apexVerifierClient    apexverifier.Client
}

func NewApexAccountService(apexAccountRepository repository.ApexAccountRepository, apexVerifierClient apexverifier.Client) ApexAccountService {
	return &apexAccountService{
		apexAccountRepository: apexAccountRepository,
		apexVerifierClient:    apexVerifierClient,
	}
}

func (s *apexAccountService) GetByUserID(ctx context.Context, userID uuid.UUID) (*model.ApexAccount, error) {
	account, err := s.apexAccountRepository.GetByUserID(ctx, userID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, ErrApexAccountNotFound
		}
		return nil, err
	}
	return account, nil
}

func (s *apexAccountService) Bind(
	ctx context.Context,
	userID uuid.UUID,
	player string,
	platform string,
	level int32,
) (*model.ApexAccount, error) {
	_, err := s.apexAccountRepository.GetByUserID(ctx, userID)
	if err == nil {
		return nil, ErrApexAccountAlreadyLinked
	}
	if !errors.Is(err, repository.ErrNotFound) {
		return nil, err
	}
	verified, err := s.apexVerifierClient.VerifyAccount(
		ctx,
		player,
		platform,
		level,
	)
	if err != nil {
		return nil, ErrApexAccountVerificationFailed
	}
	nidHash := calculateNIDHash(verified.UID)

	_, err = s.apexAccountRepository.GetByUID(ctx, verified.UID)
	if err == nil {
		return nil, ErrApexAccountAlreadyUsed
	}
	if !errors.Is(err, repository.ErrNotFound) {
		return nil, err
	}
	_, err = s.apexAccountRepository.GetByNIDHash(ctx, nidHash)
	if err == nil {
		return nil, ErrApexAccountAlreadyUsed
	}
	if !errors.Is(err, repository.ErrNotFound) {
		return nil, err
	}
	account := &model.ApexAccount{
		UserID:  userID,
		UID:     verified.UID,
		NIDHash: nidHash,
	}
	if err := s.apexAccountRepository.Create(ctx, account); err != nil {
		return nil, err
	}
	return account, nil
}

func (s *apexAccountService) Delete(ctx context.Context, userID uuid.UUID) error {
	if err := s.apexAccountRepository.DeleteByUserID(ctx, userID); err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return ErrApexAccountNotFound
		}
		return err
	}
	return nil
}
