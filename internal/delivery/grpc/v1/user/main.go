package user

import (
	"context"
	"github.com/bhankey/pharmacy-automatization-user/internal/entities"
	"github.com/bhankey/pharmacy-automatization-user/pkg/api/userservice"
)

type GRPCHandler struct {
	userservice.UnimplementedUserServiceServer // Must be

	userSrv UserSrv
}

type UserSrv interface {
	GetByID(ctx context.Context, id int) (entities.User, error)
	GetByEmail(ctx context.Context, email string) (entities.User, error)
	UpdateUser(ctx context.Context, user entities.User) error
	GetBatchOfUsers(ctx context.Context, lastClientID int, limit int) ([]entities.User, error)
	ResetPassword(ctx context.Context, email, code, newPassword string) error
	RequestToResetPassword(ctx context.Context, email string) error
	Registry(ctx context.Context, user entities.User) error
	IsPasswordCorrect(ctx context.Context, email, password string) (bool, error)
}

func NewUserGRPCHandler(userSrv UserSrv) *GRPCHandler {
	return &GRPCHandler{
		UnimplementedUserServiceServer: userservice.UnimplementedUserServiceServer{},
		userSrv:                        userSrv,
	}
}
