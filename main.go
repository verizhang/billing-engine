package main

import (
	"context"
	"fmt"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/verizhang/billing-engine/config"
	loanpb "github.com/verizhang/billing-engine/contracts/pb/loan"
	paymentpb "github.com/verizhang/billing-engine/contracts/pb/payment"
	"github.com/verizhang/billing-engine/src/handlers"
	"github.com/verizhang/billing-engine/src/repositories"
	"github.com/verizhang/billing-engine/src/services"
	"google.golang.org/grpc"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"net"
	"net/http"
	"time"
)

func main() {
	cfg := config.New()
	grpcServer := startGRPCServer(cfg)
	startHTTPServer(cfg)
	grpcServer.GracefulStop()
}

func registerSvc(cfg config.Config, server *grpc.Server) {
	db := InitDB(
		cfg.PostgresHost,
		cfg.PostgresUsername,
		cfg.PostgresPassword,
		cfg.PostgresDatabase,
		cfg.PostgresPort,
		cfg.PostgresSslmode,
		cfg.PostgresTimeZone,
		cfg.PostgresMaxConnections,
		cfg.PostgresMaxIdleConnection,
		cfg.PostgresConnectionMaxIdleTime,
	)

	// Repository
	unitOfWork := repositories.NewUnitOfWork(db)
	loanRepository := repositories.NewLoanRepository(db)
	paymentRepository := repositories.NewPaymentRepository(db)

	// Service
	loanService := services.NewLoanService(unitOfWork, loanRepository, paymentRepository)
	paymentService := services.NewPaymentService(cfg, paymentRepository, loanRepository, unitOfWork)

	// Handler
	loanHandler := handlers.NewLoanHandler(loanService)
	paymentHandler := handlers.NewPaymentHandler(paymentService)

	loanpb.RegisterLoanServer(server, loanHandler)
	paymentpb.RegisterPaymentServer(server, paymentHandler)
}

func startGRPCServer(cfg config.Config) *grpc.Server {
	grpcServer := grpc.NewServer()

	registerSvc(cfg, grpcServer)

	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", cfg.GRPCPort))
	if err != nil {
		panic(fmt.Sprintf("failed to listen: %v", err))
	}

	go func() {
		fmt.Printf("running gRPC server on port %s\n", cfg.GRPCPort)
		if err := grpcServer.Serve(lis); err != nil {
			panic(fmt.Sprintf("failed to serve gRPC server: %v", err))
		}
	}()

	return grpcServer
}

func startHTTPServer(cfg config.Config) {
	ctx := context.Background()
	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithInsecure()}

	err := loanpb.RegisterLoanHandlerFromEndpoint(ctx, mux, fmt.Sprintf(":%s", cfg.GRPCPort), opts)
	if err != nil {
		panic(fmt.Sprintf("failed to register loan gRPC Gateway: %v", err))
	}

	err = paymentpb.RegisterPaymentHandlerFromEndpoint(ctx, mux, fmt.Sprintf(":%s", cfg.GRPCPort), opts)
	if err != nil {
		panic(fmt.Sprintf("failed to register payment gRPC Gateway: %v", err))
	}

	fmt.Printf("running REST server on port %s\n", cfg.RESTPort)
	err = http.ListenAndServe(fmt.Sprintf(":%s", cfg.RESTPort), mux)
	if err != nil {
		panic(fmt.Sprintf("failed to serve HTTP: %v", err))
	}
}

func InitDB(host, user, password, dbname, port, sslmode, timezone string, maxConnections, maxIdleConnections, connectionsMaxIdleTime int) *gorm.DB {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=%s",
		host,
		user,
		password,
		dbname,
		port,
		sslmode,
		timezone,
	)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		TranslateError: true,
	})
	db = db.Debug()
	if err != nil {
		panic(fmt.Sprintf("failed to connect to database: %v", err))
	}

	postgresDB, err := db.DB()
	postgresDB.SetMaxOpenConns(maxConnections)
	postgresDB.SetMaxIdleConns(maxIdleConnections)
	postgresDB.SetConnMaxIdleTime(time.Duration(connectionsMaxIdleTime) * time.Second)

	return db
}
