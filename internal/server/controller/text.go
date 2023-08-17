package controller

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/PaulYakow/gophkeeper/internal/server/usecase"
	pb "github.com/PaulYakow/gophkeeper/proto"
)

// TextServer реализация интерфейса proto.TextServer (описание - gophkeeper/proto/text.proto)
type TextServer struct {
	pb.UnimplementedTextServer
	notes usecase.ITextService
}

// NewTextServer создаёт объект TextServer.
func NewTextServer(notes usecase.ITextService) *TextServer {
	return &TextServer{
		notes: notes,
	}
}

// GetAll - получение всех значений заметок.
func (s *TextServer) GetAll(ctx context.Context, req *pb.GetAllNotesRequest) (*pb.GetAllNotesResponse, error) {
	var resp pb.GetAllNotesResponse

	userID, ok := ctx.Value(userIDKey).(int)
	if !ok {
		return nil, status.Error(codes.Aborted, "missing user_id")
	}

	notes, err := s.notes.ViewAllNotes(userID)
	if err != nil {
		return nil, err
	}

	for _, note := range notes {
		resp.Notes = append(resp.Notes, &pb.NoteMsg{
			Id:       int64(note.ID),
			Note:     note.Note,
			Metadata: note.Metadata,
		})
	}

	return &resp, nil
}
