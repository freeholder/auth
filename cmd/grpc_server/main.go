// package main
package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"time"

	"github.com/fatih/color"
	desc "github.com/freeholder/auth/pkg/note_v1"

	// "github.com/jackc/pgx"
	"github.com/jackc/pgx/v5"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

const (
	grpcPort = 50051
	dbDSN    = "host=localhost port=5432 dbname=auth user=auth-user password=auth-password sslmode=disable"
)

type User struct {
	id         int64
	name       string
	email      string
	password   string
	role       int32
	created_at time.Time
	updated_at time.Time
}

type server struct {
	desc.UnimplementedNoteV1Server
	db *pgx.Conn
}

func (s *server) Create(ctx context.Context, req *desc.CreateRequest) (*desc.CreateResponse, error) {
	var id int64

	err := s.db.QueryRow(
		ctx,
		"INSERT INTO users (name, email, password, created_at) VALUES ($1, $2, $3, NOW()) RETURNING id",
		req.Name,
		req.Email,
		req.Password,
	).Scan(&id)

	if err != nil {
		fmt.Printf("The error is - %v", err)
	}
	return &desc.CreateResponse{
		Id: id,
	}, nil
}

func (s *server) Get(ctx context.Context, req *desc.GetRequest) (*desc.GetResponse, error) {
	var user User
	err := s.db.QueryRow(
		ctx,
		"SELECT id, name, email, password, role, created_at, updated_at FROM users WHERE id = $1",
		req.Id,
	).Scan(
		&user.id,
		&user.name,
		&user.email,
		&user.password,
		&user.role,
		&user.created_at,
		&user.updated_at,
	)

	if err != nil {
		fmt.Printf("Error with Get handler - %v", err)
	}
	return &desc.GetResponse{
		Id:        user.id,
		Name:      user.name,
		Email:     user.email,
		Role:      desc.Role(user.role),
		CreatedAt: timestamppb.New(user.created_at),
		UpdatedAt: timestamppb.New(user.updated_at),
	}, nil
}

func (s *server) Update(ctx context.Context, req *desc.UpdateRequest) (*emptypb.Empty, error) {
	_, err := s.db.Exec(
		ctx,
		"UPDATE users SET email = $1, name = $3 WHERE id = $2",
		req.Email,
		req.Id,
		req.Name,
	)

	if err != nil {
		fmt.Printf("The Update Error occured: %v", err)
	}
	return &emptypb.Empty{}, nil
}

func (s *server) Delete(ctx context.Context, req *desc.DeleteRequest) (*emptypb.Empty, error) {
	_, err := s.db.Exec(
		ctx,
		"DELETE FROM users WHERE id = $1",
		req.Id,
	)
	if err != nil {
		fmt.Printf("The Delete Error occured: %v", err)
	}
	return &emptypb.Empty{}, nil
}
func main() {

	ctx := context.Background()
	con, err := pgx.Connect(ctx, dbDSN)

	if err != nil {
		log.Fatalf("failed to connect to database %v", err)
	}

	defer con.Close(ctx)

	fmt.Println(color.GreenString("Server has started!"))

	lis, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%d", grpcPort))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	reflection.Register(s)
	desc.RegisterNoteV1Server(s, &server{db: con})

	log.Printf("gRPC server listening on port %s", lis.Addr().String())

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve")
	}
}
