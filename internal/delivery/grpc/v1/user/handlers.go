package user

import (
	"context"
	"fmt"

	"github.com/bhankey/go-utils/pkg/apperror"
	"github.com/bhankey/pharmacy-automatization-user/internal/entities"
	"github.com/bhankey/pharmacy-automatization-user/pkg/api/userservice"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (h *GRPCHandler) GetByEmail(ctx context.Context, req *userservice.Email) (*userservice.User, error) {
	errBase := fmt.Sprintf("user.GetByEmail(%v)", req)

	if err := req.ValidateAll(); err != nil {
		return nil, apperror.NewClientError(apperror.WrongRequest, fmt.Errorf("%s: %w", errBase, err))
	}

	user, err := h.userSrv.GetByEmail(ctx, req.GetEmail())
	if err != nil {
		return nil, fmt.Errorf("%s: %w", errBase, err)
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
	errBase := fmt.Sprintf("user.GetByID(%v)", req)

	if err := req.ValidateAll(); err != nil {
		return nil, apperror.NewClientError(apperror.WrongRequest, fmt.Errorf("%s: %w", errBase, err))
	}

	user, err := h.userSrv.GetByID(ctx, int(req.GetId()))
	if err != nil {
		return nil, fmt.Errorf("%s: %w", errBase, err)
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
	errBase := fmt.Sprintf("user.CreateUser(%v)", req)

	if err := req.ValidateAll(); err != nil {
		return nil, apperror.NewClientError(apperror.WrongRequest, fmt.Errorf("%s: %w", errBase, err))
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
		return nil, fmt.Errorf("%s: %w", errBase, err)
	}

	return &emptypb.Empty{}, nil
}

func (h *GRPCHandler) RequestToChangePassword(ctx context.Context, req *userservice.Email) (*emptypb.Empty, error) {
	errBase := fmt.Sprintf("user.RequestToChangePassword(%v)", req)

	if err := req.ValidateAll(); err != nil {
		return nil, apperror.NewClientError(apperror.WrongRequest, fmt.Errorf("%s: %w", errBase, err))
	}

	if err := h.userSrv.RequestToResetPassword(ctx, req.GetEmail()); err != nil {
		return nil, fmt.Errorf("%s: %w", errBase, err)
	}

	return &emptypb.Empty{}, nil
}

func (h *GRPCHandler) ChangePassword(
	ctx context.Context,
	req *userservice.ChangePasswordRequest,
) (*emptypb.Empty, error) {
	errBase := fmt.Sprintf("user.ChangePassword(%v)", req)

	if err := req.ValidateAll(); err != nil {
		return nil, apperror.NewClientError(apperror.WrongRequest, fmt.Errorf("%s: %w", errBase, err))
	}

	if err := h.userSrv.ResetPassword(ctx, req.GetEmail(), req.GetCode(), req.GetNewPassword()); err != nil {
		return nil, fmt.Errorf("%s: %w", errBase, err)
	}

	return &emptypb.Empty{}, nil
}

func (h *GRPCHandler) GetUsers(ctx context.Context, req *userservice.PaginationRequest) (*userservice.Users, error) {
	errBase := fmt.Sprintf("user.GetUsers(%v)", req)

	if err := req.ValidateAll(); err != nil {
		return nil, apperror.NewClientError(apperror.WrongRequest, fmt.Errorf("%s: %w", errBase, err))
	}

	users, err := h.userSrv.GetBatchOfUsers(ctx, int(req.LastId), int(req.Limit))
	if err != nil {
		return nil, fmt.Errorf("%s: %w", errBase, err)
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
	errBase := fmt.Sprintf("user.UpdateUser(%v)", req)

	if err := req.ValidateAll(); err != nil {
		return nil, apperror.NewClientError(apperror.WrongRequest, fmt.Errorf("%s: %w", errBase, err))
	}

	if err := h.userSrv.UpdateUser(ctx, entities.User{
		ID:                int(req.Id),
		Name:              req.Name,
		Surname:           req.Surname,
		Email:             req.Email,
		Role:              entities.Role(req.Role),
		DefaultPharmacyID: int(req.DefaultPharmacyId),
	}); err != nil {
		return nil, fmt.Errorf("%s: %w", errBase, err)
	}

	return &emptypb.Empty{}, nil
}

func (h *GRPCHandler) IsPasswordCorrect(
	ctx context.Context,
	req *userservice.EmailAndPassword,
) (*userservice.IsCorrect, error) {
	errBase := fmt.Sprintf("user.IsPasswordCorrect(%v)", req)

	if err := req.ValidateAll(); err != nil {
		return nil, apperror.NewClientError(apperror.WrongRequest, fmt.Errorf("%s: %w", errBase, err))
	}

	isCorrect, err := h.userSrv.IsPasswordCorrect(ctx, req.GetEmail(), req.GetPassword())
	if err != nil {
		return nil, fmt.Errorf("%s: %w", errBase, err)
	}

	return &userservice.IsCorrect{
		IsCorrect: isCorrect,
	}, nil
}
