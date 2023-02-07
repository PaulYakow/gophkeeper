package controller

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/PaulYakow/gophkeeper/internal/server/usecase"
	pb "github.com/PaulYakow/gophkeeper/proto"
)

type PairServer struct {
	pb.UnimplementedPairServer
	pairs usecase.IPairsService
}

func NewPairsServer(pairs usecase.IPairsService) *PairServer {
	return &PairServer{
		pairs: pairs,
	}
}

func (s *PairServer) GetAll(ctx context.Context, req *pb.GetAllRequest) (*pb.GetAllResponse, error) {
	var resp pb.GetAllResponse

	userID, ok := ctx.Value(userIDKey).(int)
	if !ok {
		return nil, status.Error(codes.Aborted, "missing user_id")
	}

	pairs, err := s.pairs.ViewAllPairs(userID)
	if err != nil {
		return nil, err
	}

	for _, pair := range pairs {
		resp.Pairs = append(resp.Pairs, &pb.PairMsg{
			Id:       int64(pair.ID),
			Login:    pair.Login,
			Password: pair.Password,
			Metadata: pair.Metadata,
		})
	}

	return &resp, nil
}
