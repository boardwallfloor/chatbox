package main

import (
	"chatbox/auth_module/pb"
	"context"
	"flag"
	"log"
	"net"

	"google.golang.org/grpc"
)

//	type DBPool struct {
//		connection chan *sql.DB
//		max        int
//		mu         sync.Mutex
//	}
//
//	func NewDBPool(max int, dsn string) (*DBPool, error) {
//		conn := make(chan *sql.DB, max)
//		for i := 0; i < max; i++ {
//			db, err := sql.Open("sqlite3", dsn)
//			if err != nil {
//				return nil, err
//			}
//			conn <- db
//		}
//		pool := DBPool{
//			connection: conn,
//			max:        max,
//		}
//		return &pool, nil
//	}
//
//	func (pool *DBPool) Get() (*sql.DB, error) {
//		select {
//		case conn := <-pool.connection:
//			return conn, nil
//		default:
//			return nil, fmt.Errorf("Connection pool exhausted")
//		}
//	}
//
//	func (pool *DBPool) Put(conn *sql.DB) {
//		pool.mu.Lock()
//		pool.connection <- conn
//		pool.mu.Unlock()
//	}
type User struct {
	username string
	password string
}

type AuthOperation interface {
	Login(username, password string) (bool, string, error)
	CreateSession(uid string) (string, error)
}

type AuthServices struct {
	pb.UnimplementedAuthServiceServer
	db AuthOperation
}

func (as *AuthServices) AuthenticateUser(ctx context.Context, ar *pb.AuthenticateRequest) (*pb.AuthenticateResponse, error) {
	status, uid, err := as.db.Login(ar.Username, ar.Password)
	if err != nil {
		return nil, err
	}
	if !status {
		return &pb.AuthenticateResponse{Success: false, Token: "", Error: "Login Failed: Invalid Username or Password"}, nil
	}
	ssid, err := as.db.CreateSession(uid)
	if err != nil {
		return nil, err
	}
	return &pb.AuthenticateResponse{Success: true, Token: ssid, Error: ""}, nil
}

type UserBuild struct {
	pb.User
}

func (ub *UserBuild) AddId(id int) *UserBuild {
	ub.Id = int32(id)
	return ub
}

func (ub *UserBuild) AddUsername(username string) *UserBuild {
	ub.Username = username
	return ub
}

func (ub *UserBuild) AddEmail(email string) *UserBuild {
	ub.Email = email
	return ub
}

func initServerProd() {
	log.Println("Running test grpc auth server")
	sessMap := make(map[string]string)
	inmem := Inmem{userList: []InmemUser{
		{User: User{username: "user1", password: "pass1"}, uid: "123456"},
		{User: User{username: "user2", password: "pass2"}, uid: "654321"},
		{User: User{username: "user3", password: "pass3"}, uid: "987654"},
		{User: User{username: "user4", password: "pass4"}, uid: "456789"},
		{User: User{username: "user5", password: "pass5"}, uid: "321987"},
		{User: User{username: "user6", password: "pass6"}, uid: "789456"},
		{User: User{username: "user7", password: "pass7"}, uid: "654789"},
		{User: User{username: "user8", password: "pass8"}, uid: "123987"},
		{User: User{username: "user9", password: "pass9"}, uid: "456123"},
		{User: User{username: "user10", password: "pass10"}, uid: "789321"},
	}, session: sessMap}
	srv := AuthServices{db: &inmem}
	grpcServer := grpc.NewServer()

	pb.RegisterAuthServiceServer(grpcServer, &srv)

	lis, err := net.Listen("tcp", ":8008")
	if err != nil {
		log.Fatalf("TCP listening error, %s", err)
	}
	err = grpcServer.Serve(lis)
	if err != nil {
		log.Fatalf("grpc serving error, %s", err)
	}

	log.Println("Running Auth module on :8008")
}

func testServer(started chan struct{}) {
	log.Println("Running test grpc auth server")
	sessMap := make(map[string]string)
	inmem := Inmem{userList: []InmemUser{
		{User: User{username: "user1", password: "pass1"}, uid: "123456"},
		{User: User{username: "user2", password: "pass2"}, uid: "654321"},
		{User: User{username: "user3", password: "pass3"}, uid: "987654"},
		{User: User{username: "user4", password: "pass4"}, uid: "456789"},
		{User: User{username: "user5", password: "pass5"}, uid: "321987"},
		{User: User{username: "user6", password: "pass6"}, uid: "789456"},
		{User: User{username: "user7", password: "pass7"}, uid: "654789"},
		{User: User{username: "user8", password: "pass8"}, uid: "123987"},
		{User: User{username: "user9", password: "pass9"}, uid: "456123"},
		{User: User{username: "user10", password: "pass10"}, uid: "789321"},
	}, session: sessMap}
	srv := AuthServices{db: &inmem}
	grpcServer := grpc.NewServer()

	pb.RegisterAuthServiceServer(grpcServer, &srv)

	lis, err := net.Listen("tcp", ":8008")
	if err != nil {
		log.Fatalf("TCP listening error, %s", err)
	}
	close(started)
	err = grpcServer.Serve(lis)
	if err != nil {
		log.Fatalf("grpc serving error, %s", err)
	}

	log.Println("Running test Auth module on :8008")
}

func main() {
	modePtr := flag.String("mode", "production", "Operating mode: 'production' or 'test'")
	flag.Parse()

	switch *modePtr {
	case "production":
		initServerProd()
	case "test":
		ch := make(chan struct{})
		testServer(ch)
	default:
		ch := make(chan struct{})
		testServer(ch)
	}
}
