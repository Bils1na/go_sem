package linkgrpc

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"gitlab.com/robotomize/gb-golang/homework/03-02-umanager/internal/database"
	"gitlab.com/robotomize/gb-golang/homework/03-02-umanager/pkg/pb"
)

var _ pb.LinkServiceServer = (*Handler)(nil)

func New(linksRepository linksRepository, timeout time.Duration) *Handler {
	return &Handler{linksRepository: linksRepository, timeout: timeout}
}

type Handler struct {
	pb.UnimplementedLinkServiceServer
	linksRepository linksRepository
	timeout         time.Duration
}

func (h Handler) GetLinkByUserID(ctx context.Context, id *pb.GetLinksByUserId) (*pb.ListLinkResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, h.timeout)
	defer cancel()

	links, err := h.linksRepository.FindByUserID(ctx, request.UserId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get links by user ID: %v", err)
	}

	var pbLinks []*pb.Link
	for _, link := range links {
		pbLinks = append(pbLinks, &pb.Link{
			Id:     link.ID.Hex(),
			Url:    link.URL,
			UserId: link.UserID,
			Title:  link.Title,
		})
	}

	return &pb.ListLinkResponse{Links: pbLinks}, nil
}

func (h Handler) mustEmbedUnimplementedLinkServiceServer() {}

func (h Handler) CreateLink(ctx context.Context, request *pb.CreateLinkRequest) (*pb.Empty, error) {
	ctx, cancel := context.WithTimeout(ctx, h.timeout)
	defer cancel()

	link := database.CreateLinkReq{
		URL:    request.Url,
		UserID: request.UserId,
		Title:  request.Title,
	}

	_, err := h.linksRepository.Create(ctx, link)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create link: %v", err)
	}

	return &pb.Empty{}, nil
}

func (h Handler) GetLink(ctx context.Context, request *pb.GetLinkRequest) (*pb.Link, error) {
	ctx, cancel := context.WithTimeout(ctx, h.timeout)
	defer cancel()

	id, err := primitive.ObjectIDFromHex(request.Id)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid link ID: %v", err)
	}

	link, err := h.linksRepository.FindByID(ctx, id)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "link not found: %v", err)
	}

	return &pb.Link{
		Id:     link.ID.Hex(),
		Url:    link.URL,
		UserId: link.UserID,
		Title:  link.Title,
	}, nil
}

func (h Handler) UpdateLink(ctx context.Context, request *pb.UpdateLinkRequest) (*pb.Empty, error) {
	ctx, cancel := context.WithTimeout(ctx, h.timeout)
	defer cancel()

	id, err := primitive.ObjectIDFromHex(request.Id)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid link ID: %v", err)
	}

	link := database.UpdateLinkReq{
		ID:    id,
		URL:   request.Url,
		Title: request.Title,
	}

	_, err = h.linksRepository.Update(ctx, link)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to update link: %v", err)
	}

	return &pb.Empty{}, nil
}

func (h Handler) DeleteLink(ctx context.Context, request *pb.DeleteLinkRequest) (*pb.Empty, error) {
	ctx, cancel := context.WithTimeout(ctx, h.timeout)
	defer cancel()

	id, err := primitive.ObjectIDFromHex(request.Id)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid link ID: %v", err)
	}

	err = h.linksRepository.Delete(ctx, id)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to delete link: %v", err)
	}

	return &pb.Empty{}, nil
}

func (h Handler) ListLinks(ctx context.Context, request *pb.Empty) (*pb.ListLinkResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, h.timeout)
	defer cancel()

	links, err := h.linksRepository.FindAll(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to list links: %v", err)
	}

	var pbLinks []*pb.Link
	for _, link := range links {
		pbLinks = append(pbLinks, &pb.Link{
			Id:     link.ID.Hex(),
			Url:    link.URL,
			UserId: link.UserID,
			Title:  link.Title,
		})
	}

	return &pb.ListLinkResponse{Links: pbLinks}, nil
}
