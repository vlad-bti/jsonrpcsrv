// Package app configures and runs application.
package app

import (
	"context"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/vlad-bti/jrpc2"
	"github.com/vlad-bti/jsonrpcsrv/config"
	"github.com/vlad-bti/jsonrpcsrv/internal/adapters/db/postgresql"
	"github.com/vlad-bti/jsonrpcsrv/internal/controller/json_rpc/v1"
	"github.com/vlad-bti/jsonrpcsrv/internal/domain/service"
	"github.com/vlad-bti/jsonrpcsrv/internal/domain/usecase"
	"github.com/vlad-bti/jsonrpcsrv/pkg/logger"
	"github.com/vlad-bti/jsonrpcsrv/pkg/postgres"
)

// Run creates objects via constructors.
func Run(cfg *config.Config) {
	log := logger.New(cfg.Log.Level)

	// Repository
	pg, err := postgres.New(cfg.PG.URL, postgres.MaxPoolSize(cfg.PG.PoolMax))
	if err != nil {
		log.Fatal("app - Run - postgres.New: %v", err)
	}
	defer pg.Close()

	balanceStorage := postgresql.NewBalanceStorage(pg)
	playerStorage := postgresql.NewPlayerStorage(pg)
	transactionStorage := postgresql.NewTransactionStorage(pg)
	transactor := postgresql.NewTransactor(log, pg)

	// Service
	balanceService := service.NewBalanceService(log, balanceStorage)
	playerService := service.NewPlayerService(log, playerStorage)
	transactionService := service.NewTransactionService(log, transactionStorage)

	// Use case
	gameUsecase := usecase.NewGameUsecase(log, balanceService, playerService, transactionService, transactor)

	wg := sync.WaitGroup{}

	// JSON-RPC Server
	jsonRpcServer := jrpc2.NewServer(cfg.JsonRpc.Port, cfg.JsonRpc.Route, nil)
	gameHandler := json_rpc.NewGameHandler(gameUsecase)
	gameHandler.Register(jsonRpcServer)

	wg.Add(1)
	go func() {
		jsonRpcServer.Start()
		wg.Done()
	}()

	// Waiting signal
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	select {
	case s := <-interrupt:
		log.Info("app - Run - signal: %v", s.String())
	}

	// Shutdown
	err = jsonRpcServer.Shutdown(context.TODO())
	if err != nil {
		log.Error("app - Run - jsonRpcServer.Shutdown: %v", err)
	}

	wg.Wait()
}
