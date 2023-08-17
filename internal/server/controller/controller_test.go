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

	"github.com/PaulYakow/gophkeeper/internal/server/controller"
	"github.com/PaulYakow/gophkeeper/internal/server/mocks"
	pb "github.com/PaulYakow/gophkeeper/proto"
)

const (
	bufSize = 1024 * 1024
)

var grpcMock = struct {
	ctrl    *gomock.Controller
	service *mocks.MockIService
}{}

func mockHelper(t testing.TB) {
	t.Helper()
	grpcMock.ctrl = gomock.NewController(t)
}

func dialer() func(context.Context, string) (net.Conn, error) {
	listener := bufconn.Listen(bufSize)
	server := grpc.NewServer()

	grpcMock.service = mocks.NewMockIService(grpcMock.ctrl)
	pb.RegisterUserServer(server, controller.NewUserServer(grpcMock.service))

	go func() {
		if err := server.Serve(listener); err != nil {
			log.Fatal(err)
		}
	}()

	return func(context.Context, string) (net.Conn, error) {
		return listener.Dial()
	}
}

func TestRegister(t *testing.T) {
	mockHelper(t)
	defer grpcMock.ctrl.Finish()

	ctx := context.Background()
	conn, err := grpc.DialContext(ctx, "bufnet", grpc.WithContextDialer(dialer()), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Fatalf("Failed to dial bufnet: %v", err)
	}
	defer conn.Close()

	client := pb.NewUserClient(conn)

	login, password := "user", "password"
	token := "token"

	t.Run("proper register", func(t *testing.T) {
		grpcMock.service.EXPECT().RegisterUser(login, password).Return(token, nil)
		resp, err := client.Register(ctx, &pb.RegisterRequest{Login: login, Password: password})
		require.NoError(t, err)
		require.Equal(t, token, resp.Token)
		require.Empty(t, resp.Error)
	})

	t.Run("fail register", func(t *testing.T) {
		errFail := errors.New("fail")
		grpcMock.service.EXPECT().RegisterUser(login, password).Return("", errFail)
		resp, err := client.Register(ctx, &pb.RegisterRequest{Login: login, Password: password})
		require.Error(t, err)
		require.Empty(t, resp)
	})
}

func TestLogin(t *testing.T) {
	mockHelper(t)
	defer grpcMock.ctrl.Finish()

	ctx := context.Background()
	conn, err := grpc.DialContext(ctx, "bufnet", grpc.WithContextDialer(dialer()), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Fatalf("Failed to dial bufnet: %v", err)
	}
	defer conn.Close()

	client := pb.NewUserClient(conn)

	login, password := "user", "password"
	token := "token"

	t.Run("proper login", func(t *testing.T) {
		grpcMock.service.EXPECT().LoginUser(login, password).Return(token, nil)
		resp, err := client.Login(ctx, &pb.LoginRequest{Login: login, Password: password})
		require.NoError(t, err)
		require.Equal(t, token, resp.Token)
		require.Empty(t, resp.Error)
	})

	t.Run("fail login", func(t *testing.T) {
		errFail := errors.New("fail")
		grpcMock.service.EXPECT().LoginUser(login, password).Return("", errFail)
		resp, err := client.Login(ctx, &pb.LoginRequest{Login: login, Password: password})
		require.Error(t, err)
		require.Empty(t, resp)
	})
}
