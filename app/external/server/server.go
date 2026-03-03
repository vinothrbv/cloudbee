package server

import (
	"context"
	"time"

	"github.com/vinothrbv/cloudbee/app/domain/entity"
	"github.com/vinothrbv/cloudbee/app/domain/repository"
	"github.com/vinothrbv/cloudbee/pb"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// Server implements pb.PostServiceServer.
type Server struct {
	pb.UnimplementedPostServiceServer
	postRepo repository.PostRepository
}

// NewServer returns a new gRPC server with instances to repositories.
func NewServer(repo repository.PostRepository) *Server {
	return &Server{postRepo: repo}
}

func timeFromProto(ts *timestamppb.Timestamp, fallback time.Time) time.Time {
	if ts == nil || !ts.IsValid() {
		return fallback
	}
	return ts.AsTime()
}

// renderData converts the entity into response data.
func renderData(p *entity.Post) *pb.Post {
	if p == nil {
		return nil
	}
	return &pb.Post{
		Id:              p.ID,
		Title:           p.Title,
		Content:         p.Content,
		Author:          p.Author,
		Tags:            p.Tags,
		PublicationDate: timestamppb.New(p.PublicationDate),
		CreatedAt:       timestamppb.New(p.CreatedAt),
		UpdatedAt:       timestamppb.New(p.UpdatedAt),
	}
}

func renderListData(p []entity.Post) []*pb.Post {

	var posts []*pb.Post
	for i := range p {
		outPost := renderData(&p[i])
		posts = append(posts, outPost)
	}
	return posts
}

// CreatePost creates a new post.
// Accepts CreatePostRequest and returns CreatePostResponse.
func (s *Server) CreatePost(ctx context.Context,
	req *pb.CreatePostRequest) (*pb.CreatePostResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request is required")
	}

	if req.Title == "" {
		return nil, status.Error(codes.InvalidArgument, "title is required")
	}

	now := time.Now()
	post := &entity.Post{
		ID:              req.GetId(),
		Title:           req.GetTitle(),
		Content:         req.GetContent(),
		Author:          req.GetAuthor(),
		Tags:            req.GetTags(),
		PublicationDate: timeFromProto(req.PublicationDate, now),
		CreatedAt:       now,
		UpdatedAt:       now,
	}

	err := s.postRepo.Create(post)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.CreatePostResponse{Post: renderData(post)}, nil
}

// GetPost fetches a post based on the id.
// Accepts GetPostRequest and returns GetPostResponse.
func (s *Server) GetPost(ctx context.Context,
	req *pb.GetPostRequest) (*pb.GetPostResponse, error) {
	if req == nil || req.Id == 0 {
		return nil, status.Error(codes.InvalidArgument, "id is required and must be > 0")
	}

	post, err := s.postRepo.Get(req.Id)
	if err != nil {
		if err.Error() == "post not found" {
			return nil, status.Error(codes.NotFound, err.Error())
		}
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.GetPostResponse{Post: renderData(post)}, nil
}

// UpdatePost updates a post based on the id if available.
// Accepts UpdatePostRequest and returns UpdatePostResponse.
func (s *Server) UpdatePost(ctx context.Context, req *pb.UpdatePostRequest) (*pb.UpdatePostResponse, error) {
	if req == nil || req.Id == 0 {
		return nil, status.Error(codes.InvalidArgument, "id is required and must be > 0")
	}

	existing, err := s.postRepo.Get(req.Id)

	if err != nil {
		if err.Error() == "post not found" {
			return nil, status.Error(codes.NotFound, err.Error())
		}
		return nil, status.Error(codes.Internal, err.Error())
	}

	now := time.Now()
	post := &entity.Post{
		ID:              req.Id,
		Title:           req.GetTitle(),
		Content:         req.GetContent(),
		Author:          req.GetAuthor(),
		Tags:            req.GetTags(),
		PublicationDate: timeFromProto(req.PublicationDate, now),
		CreatedAt:       existing.CreatedAt,
		UpdatedAt:       now,
	}

	if err := s.postRepo.Update(post); err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.UpdatePostResponse{Post: renderData(post)}, nil
}

// DeletePost deletes a post based on the id.
// Accepts DeletePostRequest and returns DeletePostResponse.
func (s *Server) DeletePost(ctx context.Context,
	req *pb.DeletePostRequest) (*pb.DeletePostResponse, error) {
	if req == nil || req.Id == 0 {
		return nil, status.Error(codes.InvalidArgument, "id is required and must be > 0")
	}

	err := s.postRepo.Delete(req.Id)

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.DeletePostResponse{Deleted: true}, nil
}

// ListPost lists all posts.
// Returns ListPostResponse.
func (s *Server) ListPost(context.Context,
	*pb.ListPostRequest) (*pb.ListPostResponse, error) {
	return &pb.ListPostResponse{
		Post: renderListData(s.postRepo.List()),
	}, nil
}
