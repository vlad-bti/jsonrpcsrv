package usecase

import (
	"context"

	"github.com/vlad-bti/jsonrpcsrv/internal/domain/entity"
)

type BalanceService interface {
	GetBalance(ctx context.Context, playerName string, currency string) (*entity.Balance, error)
	Deposit(ctx context.Context, playerName string, currency string, value int) error
	Withdraw(ctx context.Context, playerName string, currency string, value int) error
}

type PlayerService interface {
	GetPlayer(ctx context.Context, playerName string) (*entity.Player, error)
	ChangeFreerounds(ctx context.Context, playerName string, value int) error
}

type TransactionService interface {
	GetTransaction(ctx context.Context, transactionRef string) (*entity.Transaction, error)
	AddTransaction(ctx context.Context, trx *entity.Transaction) error
	RevertTransaction(ctx context.Context, trx *entity.Transaction) error
	Begin(ctx context.Context)
	Commit(ctx context.Context)
	Rollback(ctx context.Context)
}

type gameUsecase struct {
	balanceService     BalanceService
	playerService      PlayerService
	transactionService TransactionService
}

func NewGameUsecase(balanceService BalanceService, playerService PlayerService, transactionService TransactionService) *gameUsecase {
	return &gameUsecase{
		balanceService:     balanceService,
		playerService:      playerService,
		transactionService: transactionService,
	}
}

func (g *gameUsecase) GetBalance(ctx context.Context, dto GetBalanceDTO) (*entity.Balance, error) {
	g.transactionService.Begin(ctx)
	balance, err := g.balanceService.GetBalance(ctx, dto.PlayerName, dto.Currency)
	g.transactionService.Rollback(ctx)
	return balance, err
}

func (g *gameUsecase) WithdrawAndDeposit(ctx context.Context, dto WithdrawAndDepositDTO) error {
	g.transactionService.Begin(ctx)
	transaction, err := g.transactionService.GetTransaction(ctx, dto.TransactionRef)
	if err != nil {
		g.transactionService.Rollback(ctx)
		return err
	}
	if transaction != nil {
		if transaction.Status == entity.TRANSACTION_REVERTED {
			g.transactionService.Rollback(ctx)
			return entity.ErrTransactionAlreadyReverted
		}
		g.transactionService.Rollback(ctx)
		return nil
	}
	if dto.Withdraw > 0 {
		err = g.balanceService.Withdraw(ctx, dto.PlayerName, dto.Currency, dto.Withdraw)
		if err != nil {
			g.transactionService.Rollback(ctx)
			return err
		}
	}
	if dto.Deposit > 0 {
		err = g.balanceService.Deposit(ctx, dto.PlayerName, dto.Currency, dto.Deposit)
		if err != nil {
			g.transactionService.Rollback(ctx)
			return err
		}
	}
	trx := &entity.Transaction{
		PlayerName:       dto.PlayerName,
		Withdraw:         dto.Withdraw,
		Deposit:          dto.Deposit,
		Currency:         dto.Currency,
		TransactionRef:   dto.TransactionRef,
		ChargeFreerounds: dto.ChargeFreerounds,
		Status:           entity.TRANSACTION_DEFAULT,
	}
	err = g.transactionService.AddTransaction(ctx, trx)
	if err != nil {
		g.transactionService.Rollback(ctx)
		return err
	}
	g.transactionService.Commit(ctx)
	return nil
}

func (g *gameUsecase) RollbackTransaction(ctx context.Context, dto RollbackTransactionDTO) error {
	g.transactionService.Begin(ctx)
	transaction, err := g.transactionService.GetTransaction(ctx, dto.TransactionRef)
	if err != nil {
		g.transactionService.Rollback(ctx)
		return err
	}
	if transaction != nil {
		if transaction.Status == entity.TRANSACTION_REVERTED {
			g.transactionService.Rollback(ctx)
			return nil
		}
		if transaction.Deposit > 0 {
			err = g.balanceService.Withdraw(ctx, transaction.PlayerName, transaction.Currency, transaction.Deposit)
			if err != nil {
				g.transactionService.Rollback(ctx)
				return err
			}
		}
		if transaction.Withdraw > 0 {
			err = g.balanceService.Deposit(ctx, transaction.PlayerName, transaction.Currency, transaction.Withdraw)
			if err != nil {
				g.transactionService.Rollback(ctx)
				return err
			}
		}
		if transaction.ChargeFreerounds > 0 {
			err = g.playerService.ChangeFreerounds(ctx, transaction.PlayerName, transaction.ChargeFreerounds)
			if err != nil {
				g.transactionService.Rollback(ctx)
				return err
			}
		}
	} else {
		transaction = &entity.Transaction{
			PlayerName:     dto.PlayerName,
			TransactionRef: dto.TransactionRef,
		}
	}
	err = g.transactionService.RevertTransaction(ctx, transaction)
	if err != nil {
		g.transactionService.Rollback(ctx)
		return err
	}
	g.transactionService.Commit(ctx)
	return nil
}
