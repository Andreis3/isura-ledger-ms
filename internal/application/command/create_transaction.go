package command

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/google/uuid"

	"github.com/andreis3/isura-ledger-ms/internal/application"
	"github.com/andreis3/isura-ledger-ms/internal/application/event"
	"github.com/andreis3/isura-ledger-ms/internal/domain/account"
	"github.com/andreis3/isura-ledger-ms/internal/domain/money"
	"github.com/andreis3/isura-ledger-ms/internal/domain/outbox"
	"github.com/andreis3/isura-ledger-ms/internal/domain/transaction"
)

var (
	ErrInsufficientBalance = errors.New("insufficient balance")
	ErrTransactionExists   = errors.New("transaction already exists")
)

type CreateTransactionInput struct {
	IdempotencyKey  string
	DebitAccountID  account.AccountID
	CreditAccountID account.AccountID
	Amount          int64
	Currency        money.Currency
}

type CreateTransaction struct {
	uow                   application.UnitOfWork
	accountRepository     account.Repository
	transactionRepository transaction.Repository
	outboxRepository      outbox.Repository
	tracer                application.Tracer
}

func NewCreateTransaction(uow application.UnitOfWork,
	accountRepository account.Repository,
	transactionRepository transaction.Repository,
	outboxRepository outbox.Repository,
	tracer application.Tracer,
) *CreateTransaction {
	return &CreateTransaction{
		uow:                   uow,
		accountRepository:     accountRepository,
		transactionRepository: transactionRepository,
		outboxRepository:      outboxRepository,
		tracer:                tracer,
	}
}

func (c *CreateTransaction) Execute(ctx context.Context, input CreateTransactionInput) error {
	ctx, span := c.tracer.Start(ctx, "CreateTransaction.Execute")
	defer span.End()

	existingTransaction, err := c.transactionRepository.ExistsByIdempotencyKey(ctx, input.IdempotencyKey)
	if err != nil {
		span.RecordError(err)
		return err
	}

	if existingTransaction {
		return ErrTransactionExists
	}

	return c.uow.WithTransaction(ctx, func(ctxTx context.Context) error {
		// 1. Lock determinístico para evitar deadlock
		firstID, secondID := input.DebitAccountID, input.CreditAccountID
		if firstID > secondID {
			firstID, secondID = secondID, firstID
		}

		// 2. Busca com SELECT FOR UPDATE da TX
		// Aqui garantimos que o saldo lido é a "última verdade" e ninguém mexe até o commit
		_, err := c.accountRepository.FindBalanceForUpdateByID(ctxTx, firstID)
		if err != nil {
			span.RecordError(err)
			return err
		}

		_, err = c.accountRepository.FindBalanceForUpdateByID(ctxTx, secondID)
		if err != nil {
			span.RecordError(err)
			return err
		}

		// 3. Busca os objetos completos para lógica de domínio
		debitAccount, _ := c.accountRepository.FindByID(ctxTx, firstID)
		creditAccount, _ := c.accountRepository.FindByID(ctxTx, secondID)

		if creditAccount == nil {
			span.RecordError(ErrInsufficientBalance)
			return account.ErrAccountNotFound
		}

		amount, err := money.NewMoney(input.Amount, input.Currency)
		if err != nil {
			span.RecordError(err)
			return err
		}

		if !debitAccount.Balance.IsSufficientBalance(amount) {
			return ErrInsufficientBalance
		}

		transactionID := uuid.New().String()
		newTransaction := transaction.NewTransaction(transaction.TransactionID(transactionID), input.IdempotencyKey)

		debitEntryID := uuid.New().String()
		creditEntryID := uuid.New().String()

		newDebitEntry, err := transaction.NewEntry(
			transaction.EntryID(debitEntryID),
			input.IdempotencyKey,
			transaction.Debit,
			amount,
			transaction.AccountID(debitAccount.ID),
			transaction.TransactionID(transactionID),
		)

		if err != nil {
			return err
		}

		newCreditEntry, err := transaction.NewEntry(
			transaction.EntryID(creditEntryID),
			input.IdempotencyKey,
			transaction.Credit,
			amount,
			transaction.AccountID(creditAccount.ID),
			transaction.TransactionID(transactionID),
		)
		if err != nil {
			span.RecordError(err)
			return err
		}

		debitAccount.Balance, err = debitAccount.Balance.Subtract(amount)
		if err != nil {
			span.RecordError(err)
			return err
		}

		creditAccount.Balance, err = creditAccount.Balance.Add(amount)
		if err != nil {
			span.RecordError(err)
			return err
		}

		err = newTransaction.AddEntry(newDebitEntry)
		if err != nil {
			span.RecordError(err)
			return err
		}

		err = newTransaction.AddEntry(newCreditEntry)
		if err != nil {
			span.RecordError(err)
			return err
		}

		err = newTransaction.Complete()
		if err != nil {
			span.RecordError(err)
			return err
		}

		err = c.transactionRepository.Save(ctxTx, newTransaction)
		if err != nil {
			span.RecordError(err)
			return err
		}

		err = c.accountRepository.UpdateBalance(ctxTx, debitAccount.ID, debitAccount.Balance)
		if err != nil {
			span.RecordError(err)
			return err
		}

		err = c.accountRepository.UpdateBalance(ctxTx, creditAccount.ID, creditAccount.Balance)
		if err != nil {
			span.RecordError(err)
			return err
		}

		payload, err := json.Marshal(event.TransactionCreated{
			TransactionID:   string(newTransaction.ID),
			IdempotencyKey:  input.IdempotencyKey,
			DebitAccountID:  string(input.DebitAccountID),
			CreditAccountID: string(input.CreditAccountID),
			Amount:          input.Amount,
			Currency:        string(input.Currency),
			Status:          string(newTransaction.Status),
			OccurredAt:      time.Now(),
		})
		if err != nil {
			span.RecordError(err)
			return err
		}

		outboxID := uuid.New().String()
		newOutbox := outbox.NewOutbox(outbox.OutboxID(outboxID),
			string(newTransaction.ID),
			outbox.Transaction,
			outbox.TransactionCreated,
			payload,
		)

		err = c.outboxRepository.Save(ctxTx, newOutbox)
		if err != nil {
			return err
		}
		return nil
	})
}
