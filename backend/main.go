package main

import (
	"context"
	"fmt"
	"log"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	directorypb "authority-app/backend/proto/directory/v1"

	// ヘルスチェック関連のパッケージをインポート
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
)

const (
	grpcPort = ":50051"
)

type server struct {
	directorypb.UnimplementedDirectoryServiceServer
}

func (s *server) CreateDirectory(ctx context.Context, req *directorypb.CreateDirectoryRequest) (*directorypb.CreateDirectoryResponse, error) {
	log.Printf("Received CreateDirectory request: %v", req)

	// ここにディレクトリ作成ロジックを実装
	// 現時点ではスタブとして成功を返す

	// 仮のディレクトリIDを生成
	directoryID := fmt.Sprintf("dir-%s-%d", req.GetName(), len(req.GetName()))

	res := &directorypb.CreateDirectoryResponse{
		Directory: &directorypb.Directory{
			Id:       directoryID,
			Name:     req.GetName(),
			ParentId: req.GetParentId(),
		},
	}
	return res, nil
}

func (s *server) ShareDirectory(ctx context.Context, req *directorypb.ShareDirectoryRequest) (*directorypb.ShareDirectoryResponse, error) {
	log.Printf("Received ShareDirectory request: %v", req)
	return nil, status.Errorf(codes.Unimplemented, "method ShareDirectory not implemented")
}

func (s *server) CheckPermission(ctx context.Context, req *directorypb.CheckPermissionRequest) (*directorypb.CheckPermissionResponse, error) {
	log.Printf("Received CheckPermission request: %v", req)
	return nil, status.Errorf(codes.Unimplemented, "method CheckPermission not implemented")
}

func (s *server) ListChildren(ctx context.Context, req *directorypb.ListChildrenRequest) (*directorypb.ListChildrenResponse, error) {
	log.Printf("Received ListChildren request: %v", req)
	return nil, status.Errorf(codes.Unimplemented, "method ListChildren not implemented")
}

func main() {
	lis, err := net.Listen("tcp", grpcPort)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	directorypb.RegisterDirectoryServiceServer(s, &server{})

	// ヘルスチェックサービスを登録
	healthServer := health.NewServer()
	grpc_health_v1.RegisterHealthServer(s, healthServer)

	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
