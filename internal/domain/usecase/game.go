package usecase

import (
	"context"

	"github.com/vlad-bti/jsonrpcsrv/internal/domain/entity"
	"github.com/vlad-bti/jsonrpcsrv/pkg/logger"
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
}

type Transactor interface {
	// WithinTransaction runs a function within a database transaction.
	//
	// Transaction is propagated in the context,
	// so it is important to propagate it to underlying repositories.
	// Function commits if error is nil, and rollbacks if not.
	// It returns the same error.
	WithinTransaction(context.Context, func(ctx context.Context) error) error
}

type gameUsecase struct {
	balanceService     BalanceService
	playerService      PlayerService
	transactionService TransactionService
	transactor         Transactor
	log                *logger.Logger
}

func NewGameUsecase(log *logger.Logger, balanceService BalanceService, playerService PlayerService, transactionService TransactionService, transactor Transactor) *gameUsecase {
	return &gameUsecase{
		balanceService:     balanceService,
		playerService:      playerService,
		transactionService: transactionService,
		transactor:         transactor,
		log:                log,
	}
}

func (r *gameUsecase) GetBalance(ctx context.Context, dto GetBalanceDTO) (*entity.Balance, error) {
	balance, err := r.balanceService.GetBalance(ctx, dto.PlayerName, dto.Currency)
	if err != nil {
		r.log.Error("GameUsecase - GetBalance - r.balanceService.GetBalance: %v; PlayerName=%v, Currency=%v",
			err,
			dto.PlayerName,
			dto.Currency,
		)
	}
	return balance, err
}

func (r *gameUsecase) WithdrawAndDeposit(ctx context.Context, dto WithdrawAndDepositDTO) error {
	return r.transactor.WithinTransaction(ctx, func(txCtx context.Context) error {
		transaction, err := r.transactionService.GetTransaction(txCtx, dto.TransactionRef)
		if err != nil {
			r.log.Error("GameUsecase - WithdrawAndDeposit - r.transactor.GetTransaction: %v; TransactionRef=%v",
				err,
				dto.TransactionRef,
			)
			return err
		}
		if transaction != nil {
			if transaction.Status == entity.TRANSACTION_REVERTED {
				return entity.ErrTransactionAlreadyReverted
			}
			return entity.NoError
		}
		if dto.Withdraw > 0 {
			err = r.balanceService.Withdraw(txCtx, dto.PlayerName, dto.Currency, dto.Withdraw)
			if err != nil {
				r.log.Error("GameUsecase - WithdrawAndDeposit - r.balanceService.Withdraw: %v; TransactionRef=%v, PlayerName=%v, Currency=%v, Withdraw=%v",
					err,
					dto.TransactionRef,
					dto.PlayerName,
					dto.Currency,
					dto.Withdraw,
				)
				return err
			}
		}
		if dto.Deposit > 0 {
			err = r.balanceService.Deposit(txCtx, dto.PlayerName, dto.Currency, dto.Deposit)
			if err != nil {
				r.log.Error("GameUsecase - WithdrawAndDeposit - r.balanceService.Deposit: %v; TransactionRef=%v, PlayerName=%v, Currency=%v, Deposit=%v",
					err,
					dto.TransactionRef,
					dto.PlayerName,
					dto.Currency,
					dto.Deposit,
				)
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
		err = r.transactionService.AddTransaction(ctx, trx)
		if err != nil {
			r.log.Info("GameUsecase - WithdrawAndDeposit - r.transactionService.AddTransaction: %v; TransactionRef=%v",
				err,
				dto.TransactionRef,
			)
			return err
		}
		return nil
	})
}

func (r *gameUsecase) RollbackTransaction(ctx context.Context, dto RollbackTransactionDTO) error {
	return r.transactor.WithinTransaction(ctx, func(txCtx context.Context) error {
		transaction, err := r.transactionService.GetTransaction(txCtx, dto.TransactionRef)
		if err != nil {
			r.log.Error("GameUsecase - RollbackTransaction - r.transactionService.GetTransaction: %v; TransactionRef=%v",
				err,
				dto.TransactionRef,
			)
			return err
		}
		if transaction != nil {
			if transaction.Status == entity.TRANSACTION_REVERTED {
				return entity.NoError
			}
			if transaction.Deposit > 0 {
				err = r.balanceService.Withdraw(txCtx, transaction.PlayerName, transaction.Currency, transaction.Deposit)
				if err != nil {
					r.log.Error("GameUsecase - RollbackTransaction - r.balanceService.Withdraw: %v; TransactionRef=%v, PlayerName=%v, Currency=%v, Withdraw=%v",
						err,
						dto.TransactionRef,
						transaction.PlayerName,
						transaction.Currency,
						transaction.Deposit,
					)
					return err
				}
			}
			if transaction.Withdraw > 0 {
				err = r.balanceService.Deposit(txCtx, transaction.PlayerName, transaction.Currency, transaction.Withdraw)
				if err != nil {
					r.log.Error("GameUsecase - RollbackTransaction - r.balanceService.Deposit: %v; TransactionRef=%v, PlayerName=%v, Currency=%v, Deposit=%v",
						err,
						dto.TransactionRef,
						transaction.PlayerName,
						transaction.Currency,
						transaction.Withdraw,
					)
					return err
				}
			}
			if transaction.ChargeFreerounds > 0 {
				err = r.playerService.ChangeFreerounds(txCtx, transaction.PlayerName, transaction.ChargeFreerounds)
				if err != nil {
					r.log.Error("GameUsecase - RollbackTransaction - r.playerService.ChangeFreerounds: %v; TransactionRef=%v, PlayerName=%v, ChargeFreerounds=%v",
						err,
						dto.TransactionRef,
						transaction.PlayerName,
						transaction.ChargeFreerounds,
					)
					return err
				}
			}
		} else {
			transaction = &entity.Transaction{
				TransactionRef: dto.TransactionRef,
			}
		}
		err = r.transactionService.RevertTransaction(txCtx, transaction)
		if err != nil {
			r.log.Error("GameUsecase - RollbackTransaction - r.transactionService.RevertTransaction: %v; TransactionRef=%v",
				err,
				dto.TransactionRef,
			)
			return err
		}
		return nil
	})
}
