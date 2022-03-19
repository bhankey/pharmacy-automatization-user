package user

import (
	"context"
	"github.com/bhankey/pharmacy-automatization/internal/entities"
	"github.com/bhankey/pharmacy-automatization/pkg/api/userservice"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (h *GRPCHandler) GetByEmail(ctx context.Context, req *userservice.Email) (*userservice.User, error) {
	if err := req.ValidateAll(); err != nil {
		return nil, err
	}

	user, err := h.userSrv.GetByEmail(ctx, req.GetEmail())
	if err != nil {
		return nil, err
	}

	return &userservice.User{
		Id:                int64(user.ID),
		Name:              user.Name,
		Email:             user.Email,
		Role:              string(user.Role),
		UseIpCheck:        false, // TODO
		IsBlocked:         false,
		DefaultPharmacyId: int64(user.DefaultPharmacyID),
	}, nil
}

func (h *GRPCHandler) GetByID(ctx context.Context, req *userservice.GetUserByIDRequest) (*userservice.User, error) {
	if err := req.ValidateAll(); err != nil {
		return nil, err
	}

	user, err := h.userSrv.GetByID(ctx, int(req.GetId()))
	if err != nil {
		return nil, err
	}

	return &userservice.User{
		Id:                int64(user.ID),
		Name:              user.Name,
		Email:             user.Email,
		Role:              string(user.Role),
		UseIpCheck:        false, // TODO
		IsBlocked:         false,
		DefaultPharmacyId: int64(user.DefaultPharmacyID),
	}, nil
}

func (h *GRPCHandler) CreateUser(ctx context.Context, req *userservice.NewUser) (*emptypb.Empty, error) {
	if err := req.ValidateAll(); err != nil {
		return nil, err
	}

	if err := h.userSrv.Registry(ctx, entities.User{
		ID:                0,
		Name:              req.Name,
		Surname:           req.Surname,
		Email:             req.Email,
		Password:          req.Password,
		Role:              entities.Role(req.Role),
		DefaultPharmacyID: int(req.DefaultPharmacyId),
	}); err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}

func (h *GRPCHandler) RequestToChangePassword(ctx context.Context, req *userservice.Email) (*emptypb.Empty, error) {
	if err := req.ValidateAll(); err != nil {
		return nil, err
	}

	if err := h.userSrv.RequestToResetPassword(ctx, req.GetEmail()); err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}

func (h *GRPCHandler) ChangePassword(ctx context.Context, req *userservice.ChangePasswordRequest) (*emptypb.Empty, error) {
	if err := req.ValidateAll(); err != nil {
		return nil, err
	}

	if err := h.userSrv.ResetPassword(ctx, req.GetEmail(), req.GetCode(), req.GetNewPassword()); err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}

func (h *GRPCHandler) GetUsers(ctx context.Context, req *userservice.PaginationRequest) (*userservice.Users, error) {
	if err := req.ValidateAll(); err != nil {
		return nil, err
	}

	users, err := h.userSrv.GetBatchOfUsers(ctx, int(req.LastId), int(req.Limit))
	if err != nil {
		return nil, err
	}

	res := make([]*userservice.User, 0, len(users))

	for _, user := range users {
		res = append(res, &userservice.User{
			Id:                int64(user.ID),
			Name:              user.Name,
			Email:             user.Email,
			Role:              string(user.Role),
			UseIpCheck:        false, // TODO
			IsBlocked:         false,
			DefaultPharmacyId: int64(user.DefaultPharmacyID),
		})
	}
	return &userservice.Users{
		Users: res,
	}, nil
}
func (h *GRPCHandler) UpdateUser(ctx context.Context, req *userservice.User) (*emptypb.Empty, error) {
	if err := req.ValidateAll(); err != nil {
		return nil, err
	}

	if err := h.userSrv.UpdateUser(ctx, entities.User{
		ID:                int(req.Id),
		Name:              req.Name,
		Surname:           req.Surname,
		Email:             req.Email,
		Role:              entities.Role(req.Role),
		DefaultPharmacyID: int(req.DefaultPharmacyId),
	}); err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}
