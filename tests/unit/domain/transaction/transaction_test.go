//go:build unit
// +build unit

package transaction_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/andreis3/isura-ledger-ms/internal/domain/money"
	"github.com/andreis3/isura-ledger-ms/internal/domain/transaction"
)

var _ = Describe("INTERNAL :: DOMAIN :: TRANSACTION :: TRANSACTION", func() {
	Describe("#NewTransaction", func() {
		Context("success cases", func() {
			It("should not return an error when build new transaction", func() {
				trans := transaction.NewTransaction("any_id", "any_idempotency_key")
				Expect(trans).NotTo(BeNil())
			})
		})
	})

	Describe("#AddEntry", func() {
		Context("success cases", func() {
			It("should not return an error when add new entry", func() {
				trans := transaction.NewTransaction("any_id", "any_idempotency_key")
				amount, _ := money.NewMoney(100, money.BRL)
				entry, _ := transaction.NewEntry("any_id", "any_idempotency_key", transaction.Credit, amount, "any_account_id", "any_transaction_id")
				err := trans.AddEntry(entry)
				Expect(err).To(BeNil())
			})
		})

		Context("error cases", func() {
			It("should return an error when add more than two entries", func() {
				trans := transaction.NewTransaction("any_id", "any_idempotency_key")
				amount, _ := money.NewMoney(100, money.BRL)
				entry, _ := transaction.NewEntry("any_id", "any_idempotency_key", transaction.Credit, amount, "any_account_id", "any_transaction_id")
				trans.AddEntry(entry)
				entry2, _ := transaction.NewEntry("any_id2", "any_idempotency_key2", transaction.Debit, amount, "any_account_id2", "any_transaction_id")
				trans.AddEntry(entry2)
				entry3, _ := transaction.NewEntry("any_id3", "any_idempotency_key3", transaction.Debit, amount, "any_account_id3", "any_transaction_id")
				err := trans.AddEntry(entry3)
				Expect(err).NotTo(BeNil())
				Expect(err).To(Equal(transaction.ErrInvalidMaxEntries))
			})

			It("should return an error when add two entries with same direction", func() {
				trans := transaction.NewTransaction("any_id", "any_idempotency_key")
				amount, _ := money.NewMoney(100, money.BRL)
				entry, _ := transaction.NewEntry("any_id", "any_idempotency_key", transaction.Credit, amount, "any_account_id", "any_transaction_id")
				trans.AddEntry(entry)
				entry2, _ := transaction.NewEntry("any_id2", "any_idempotency_key2", transaction.Credit, amount, "any_account_id2", "any_transaction_id")
				err := trans.AddEntry(entry2)
				Expect(err).NotTo(BeNil())
				Expect(err).To(Equal(transaction.ErrDuplicateEntryDirection))
			})

			It("should return an error when add two entries with different amount", func() {
				trans := transaction.NewTransaction("any_id", "any_idempotency_key")
				amount, _ := money.NewMoney(100, money.BRL)
				entry, _ := transaction.NewEntry("any_id", "any_idempotency_key", transaction.Credit, amount, "any_account_id", "any_transaction_id")
				trans.AddEntry(entry)
				amount2, _ := money.NewMoney(200, money.BRL)
				entry2, _ := transaction.NewEntry("any_id2", "any_idempotency_key2", transaction.Debit, amount2, "any_account_id2", "any_transaction_id")
				err := trans.AddEntry(entry2)
				Expect(err).NotTo(BeNil())
				Expect(err).To(Equal(transaction.ErrInvalidDifferentAmount))
			})
		})
	})

	Describe("#Complete", func() {
		Context("success cases", func() {
			It("should complete a transaction", func() {
				trans := transaction.NewTransaction("any_id", "any_idempotency_key")
				err := trans.Complete()
				Expect(err).To(BeNil())
				Expect(trans.Status).To(Equal(transaction.Completed))
			})
		})

		Context("error cases", func() {
			It("should return an error when transaction is already completed", func() {
				trans := transaction.NewTransaction("any_id", "any_idempotency_key")
				trans.Complete()
				err := trans.Complete()
				Expect(err).NotTo(BeNil())
				Expect(err).To(Equal(transaction.ErrInvalidTransactionStatus))
			})
		})
	})

	Describe("#Fail", func() {
		Context("success cases", func() {
			It("should fail a transaction", func() {
				trans := transaction.NewTransaction("any_id", "any_idempotency_key")
				err := trans.Fail()
				Expect(err).To(BeNil())
				Expect(trans.Status).To(Equal(transaction.Failed))
			})
		})

		Context("error cases", func() {
			It("should return an error when transaction is already completed", func() {
				trans := transaction.NewTransaction("any_id", "any_idempotency_key")
				trans.Complete()
				err := trans.Fail()
				Expect(err).NotTo(BeNil())
				Expect(err).To(Equal(transaction.ErrInvalidTransactionStatus))
			})
		})
	})

	Describe("#IsValid", func() {
		Context("success cases", func() {
			It("should return true for valid status", func() {
				Expect(transaction.Pending.IsValid()).To(BeTrue())
				Expect(transaction.Completed.IsValid()).To(BeTrue())
				Expect(transaction.Failed.IsValid()).To(BeTrue())
			})
		})

		Context("error cases", func() {
			It("should return false for invalid status", func() {
				Expect(transaction.TransactionStatus("invalid").IsValid()).To(BeFalse())
			})
		})
	})

	Describe("#NewEntry", func() {
		Context("success cases", func() {
			It("should create a new entry", func() {
				amount, _ := money.NewMoney(100, money.BRL)
				entry, err := transaction.NewEntry("any_id", "any_idempotency_key", transaction.Credit, amount, "any_account_id", "any_transaction_id")
				Expect(err).To(BeNil())
				Expect(entry).NotTo(BeNil())
			})
		})

		Context("error cases", func() {
			It("should return an error when direction is invalid", func() {
				amount, _ := money.NewMoney(100, money.BRL)
				_, err := transaction.NewEntry("any_id", "any_idempotency_key", "invalid", amount, "any_account_id", "any_transaction_id")
				Expect(err).NotTo(BeNil())
				Expect(err).To(Equal(transaction.ErrInvalidDirection))
			})

			It("should return an error when amount is zero", func() {
				amount, _ := money.NewMoney(0, money.BRL)
				_, err := transaction.NewEntry("any_id", "any_idempotency_key", transaction.Credit, amount, "any_account_id", "any_transaction_id")
				Expect(err).NotTo(BeNil())
				Expect(err).To(Equal(transaction.ErrAmountEqualZero))
			})

			It("should return an error when amount is negative", func() {
				amount1, _ := money.NewMoney(100, money.BRL)
				amount2, _ := money.NewMoney(200, money.BRL)
				negativeAmount, _ := amount1.Subtract(amount2)
				_, err := transaction.NewEntry("any_id", "any_idempotency_key", transaction.Credit, negativeAmount, "any_account_id", "any_transaction_id")
				Expect(err).NotTo(BeNil())
				Expect(err).To(Equal(transaction.ErrNegativeAmountValue))
			})
		})
	})

	Describe("#Direction.IsValid", func() {
		Context("success cases", func() {
			It("should return true for valid directions", func() {
				Expect(transaction.Credit.IsValid()).To(BeTrue())
				Expect(transaction.Debit.IsValid()).To(BeTrue())
			})
		})

		Context("error cases", func() {
			It("should return false for invalid direction", func() {
				Expect(transaction.Direction("invalid").IsValid()).To(BeFalse())
			})
		})
	})
})
