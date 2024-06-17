package usergrpc

import (
	"context"
	"time"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"gitlab.com/robotomize/gb-golang/homework/03-02-umanager/internal/database"
	"gitlab.com/robotomize/gb-golang/homework/03-02-umanager/pkg/pb"
)

var _ pb.UserServiceServer = (*Handler)(nil)

func New(usersRepository usersRepository, timeout time.Duration) *Handler {
	return &Handler{usersRepository: usersRepository, timeout: timeout}
}

type Handler struct {
	pb.UnimplementedUserServiceServer
	usersRepository usersRepository
	timeout         time.Duration
}

func (h Handler) CreateUser(ctx context.Context, in *pb.CreateUserRequest) (*pb.Empty, error) {
	ctx, cancel := context.WithTimeout(ctx, h.timeout)
	defer cancel()

	user := database.CreateUserReq{
		Name:  request.Name,
		Email: request.Email,
	}

	_, err := h.usersRepository.Create(ctx, user)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create user: %v", err)
	}

	return &pb.Empty{}, nil
}

func (h Handler) GetUser(ctx context.Context, in *pb.GetUserRequest) (*pb.User, error) {
	ctx, cancel := context.WithTimeout(ctx, h.timeout)
	defer cancel()

	id, err := uuid.Parse(request.Id)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid user ID: %v", err)
	}

	user, err := h.usersRepository.FindByID(ctx, id)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "user not found: %v", err)
	}

	return &pb.User{
		Id:    user.ID.String(),
		Name:  user.Name,
		Email: user.Email,
	}, nil
}

func (h Handler) UpdateUser(ctx context.Context, in *pb.UpdateUserRequest) (*pb.Empty, error) {
	ctx, cancel := context.WithTimeout(ctx, h.timeout)
	defer cancel()

	id, err := uuid.Parse(request.Id)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid user ID: %v", err)
	}

	user := database.User{
		ID:    id,
		Name:  request.Name,
		Email: request.Email,
	}

	err = h.usersRepository.Update(ctx, user)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to update user: %v", err)
	}

	return &pb.Empty{}, nil
}

func (h Handler) DeleteUser(ctx context.Context, in *pb.DeleteUserRequest) (*pb.Empty, error) {
	ctx, cancel := context.WithTimeout(ctx, h.timeout)
	defer cancel()

	id, err := uuid.Parse(request.Id)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid user ID: %v", err)
	}

	err = h.usersRepository.DeleteByUserID(ctx, id)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to delete user: %v", err)
	}

	return &pb.Empty{}, nil
}

func (h Handler) ListUsers(ctx context.Context, in *pb.Empty) (*pb.ListUsersResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, h.timeout)
	defer cancel()

	users, err := h.usersRepository.FindAll(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to list users: %v", err)
	}

	var pbUsers []*pb.User
	for _, user := range users {
		pbUsers = append(pbUsers, &pb.User{
			Id:    user.ID.String(),
			Name:  user.Name,
			Email: user.Email,
		})
	}

	return &pb.ListUsersResponse{Users: pbUsers}, nil
}
