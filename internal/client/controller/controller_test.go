package controller_test

import (
	"context"
	"errors"
	"log"
	"net"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/test/bufconn"

	"github.com/PaulYakow/gophkeeper/internal/client/controller"
	"github.com/PaulYakow/gophkeeper/internal/server/mocks"
	pb "github.com/PaulYakow/gophkeeper/proto"
)

var (
	ctrl *gomock.Controller
	srv  *mockUserServer
)

type mockUserServer struct {
	pb.UnimplementedUserServer
	auth *mocks.MockIAuthorizationService
}

func (s *mockUserServer) Register(ctx context.Context, req *pb.RegisterRequest) (*pb.RegisterResponse, error) {
	var resp pb.RegisterResponse
	token, err := s.auth.RegisterUser(req.GetLogin(), req.GetPassword())
	if err != nil {
		return nil, err
	}

	resp.Token = token
	return &resp, nil
}

func (s *mockUserServer) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	var resp pb.LoginResponse
	token, err := s.auth.LoginUser(req.GetLogin(), req.GetPassword())
	if err != nil {
		return nil, err
	}

	resp.Token = token
	return &resp, nil
}

func dialer() func(context.Context, string) (net.Conn, error) {
	listener := bufconn.Listen(1024 * 1024)

	server := grpc.NewServer()
	srv = &mockUserServer{auth: mocks.NewMockIAuthorizationService(ctrl)}
	pb.RegisterUserServer(server, srv)

	go func() {
		if err := server.Serve(listener); err != nil {
			log.Fatal(err)
		}
	}()

	return func(context.Context, string) (net.Conn, error) {
		return listener.Dial()
	}
}

func mockHelper(t testing.TB) {
	t.Helper()
	ctrl = gomock.NewController(t)
}

func TestRegister(t *testing.T) {
	mockHelper(t)
	defer ctrl.Finish()

	ctx := context.Background()
	conn, err := grpc.DialContext(ctx, "bufnet", grpc.WithContextDialer(dialer()), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Fatalf("Failed to dial bufnet: %v", err)
	}
	defer conn.Close()

	client := controller.NewUserClient(conn)

	login, password := "user", "password"
	token := "token"

	t.Run("proper register", func(t *testing.T) {
		srv.auth.EXPECT().RegisterUser(login, password).Return(token, nil)
		resp, err := client.Register(ctx, login, password)
		require.Equal(t, token, resp)
		require.NoError(t, err)
	})

	t.Run("fail register", func(t *testing.T) {
		errFail := errors.New("fail")
		srv.auth.EXPECT().RegisterUser(login, password).Return("", errFail)
		resp, err := client.Register(ctx, login, password)
		require.Empty(t, resp)
		require.Error(t, err)
	})
}

func TestLogin(t *testing.T) {
	mockHelper(t)
	defer ctrl.Finish()

	ctx := context.Background()
	conn, err := grpc.DialContext(ctx, "bufnet", grpc.WithContextDialer(dialer()), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Fatalf("Failed to dial bufnet: %v", err)
	}
	defer conn.Close()

	client := controller.NewUserClient(conn)

	login, password := "user", "password"
	token := "token"

	t.Run("proper login", func(t *testing.T) {
		srv.auth.EXPECT().LoginUser(login, password).Return(token, nil)
		resp, err := client.Login(ctx, login, password)
		require.Equal(t, token, resp)
		require.NoError(t, err)
	})

	t.Run("fail login", func(t *testing.T) {
		errFail := errors.New("fail")
		srv.auth.EXPECT().LoginUser(login, password).Return("", errFail)
		resp, err := client.Login(ctx, login, password)
		require.Empty(t, resp)
		require.Error(t, err)
	})
}
