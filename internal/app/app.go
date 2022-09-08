// Package app configures and runs application.
package app

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/vlad-bti/jrpc2"
	"github.com/vlad-bti/jsonrpcsrv/config"
	"github.com/vlad-bti/jsonrpcsrv/internal/adapters/db/fakedb"
	"github.com/vlad-bti/jsonrpcsrv/internal/controller/json_rpc/v1"
	"github.com/vlad-bti/jsonrpcsrv/internal/domain/service"
	"github.com/vlad-bti/jsonrpcsrv/internal/domain/usecase"
	"github.com/vlad-bti/jsonrpcsrv/pkg/logger"
)

// Run creates objects via constructors.
func Run(cfg *config.Config) {
	l := logger.New(cfg.Log.Level)

	// Repository
	fakeDB := fakedb.NewFakeDB()
	balanceStorage := fakedb.NewBalanceStorage()
	playerStorage := fakedb.NewPlayerStorage()
	transactionStorage := fakedb.NewTransactionStorage()

	// Service
	balanceService := service.NewBalanceService(balanceStorage)
	playerService := service.NewPlayerService(playerStorage)
	transactionService := service.NewTransactionService(transactionStorage, fakeDB)

	// Use case
	gameUsecase := usecase.NewGameUsecase(balanceService, playerService, transactionService)

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
		l.Info("app - Run - signal: " + s.String())
	}

	// Shutdown
	err := jsonRpcServer.Shutdown(context.TODO())
	if err != nil {
		l.Error(fmt.Errorf("app - Run - jsonRpcServer.Shutdown: %w", err))
	}

	wg.Wait()
}
