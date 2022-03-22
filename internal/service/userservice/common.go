package userservice

import (
	"context"
	"errors"
	"fmt"

	"github.com/bhankey/go-utils/pkg/apperror"
	"github.com/bhankey/pharmacy-automatization-user/internal/entities"
	"golang.org/x/crypto/bcrypt"
)

func (s *UserService) UpdateUser(ctx context.Context, user entities.User) error {
	errBase := fmt.Sprintf("userservice.UpdateUser(%v)", user)

	if err := s.userStorage.UpdateUser(ctx, user); err != nil {
		return fmt.Errorf("%s: failed to update user: %w", errBase, err)
	}

	return nil
}

func (s *UserService) GetBatchOfUsers(ctx context.Context, lastClientID int, limit int) ([]entities.User, error) {
	errBase := fmt.Sprintf("userservice.GetBatchOfUsers(%d, %d)", lastClientID, limit)

	users, err := s.userStorage.GetBatchOfUsers(ctx, lastClientID, limit)
	if err != nil {
		return nil, fmt.Errorf("%s: failed to get users: %w", errBase, err)
	}

	return users, nil
}

func (s *UserService) GetByEmail(ctx context.Context, email string) (entities.User, error) {
	errBase := fmt.Sprintf("userservice.GetByEmail(%s)", email)

	user, err := s.userStorage.GetUserByEmail(ctx, email)
	if err != nil {
		return entities.User{}, fmt.Errorf("%s: Failed to get user by email : %w", errBase, err)
	}

	return user, nil
}

func (s *UserService) GetByID(ctx context.Context, id int) (entities.User, error) {
	errBase := fmt.Sprintf("userservice.GetByID(%d)", id)

	user, err := s.userStorage.GetUserByID(ctx, id)
	if err != nil {
		return entities.User{}, fmt.Errorf("%s: Failed to get user by id : %w", errBase, err)
	}

	return user, nil
}

func (s *UserService) IsPasswordCorrect(ctx context.Context, email, password string) (bool, error) {
	errBase := fmt.Sprintf("userservice.IsPasswordCorrect(%s, %s)", email, password)

	user, err := s.userStorage.GetUserByEmail(ctx, email)
	if err != nil && !errors.Is(err, apperror.ErrNoEntity) {
		return false, fmt.Errorf("%s :failed to get user: %w", errBase, err)
	}

	if errors.Is(err, apperror.ErrNoEntity) {
		return false, apperror.NewClientError(apperror.WrongAuthorization, err)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return false, apperror.NewClientError(apperror.WrongAuthorization, err)
	}

	return true, nil
}
