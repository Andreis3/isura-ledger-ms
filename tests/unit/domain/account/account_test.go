//go:build unit
// +build unit

package account_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/andreis3/isura-ledger-ms/internal/domain/account"
	"github.com/andreis3/isura-ledger-ms/internal/domain/money"
)

var _ = Describe("INTERNAL :: DOMAIN :: ACCOUNT :: ACCOUNT", func() {
	Describe("#NewAccount", func() {
		Context("success cases", func() {
			It("should not return an error when build new account", func() {
				acc, err := account.NewAccount("any_id", "any_external_id", account.Asset, money.BRL)
				Expect(err).To(BeNil())
				Expect(acc).NotTo(BeNil())
			})
		})

		Context("error cases", func() {
			It("should return an error when account type is invalid", func() {
				_, err := account.NewAccount("any_id", "any_external_id", "any_type", money.BRL)
				Expect(err).NotTo(BeNil())
				Expect(err).To(Equal(account.ErrInvalidAccountingType))
			})

			It("should return an error when external id is empty", func() {
				_, err := account.NewAccount("any_id", "", account.Asset, money.BRL)
				Expect(err).NotTo(BeNil())
				Expect(err).To(Equal(account.ErrEmptyExternalID))
			})

			It("should return an error when currency is invalid", func() {
				_, err := account.NewAccount("any_id", "any_external_id", account.Asset, "INVALID")
				Expect(err).NotTo(BeNil())
				Expect(err).To(Equal(money.ErrInvalidCurrency))
			})
		})
	})

	Describe("#IsValid", func() {
		Context("success cases", func() {
			It("should return true for ASSET", func() {
				Expect(account.Asset.IsValid()).To(BeTrue())
			})
			It("should return true for LIABILITY", func() {
				Expect(account.Liability.IsValid()).To(BeTrue())
			})
			It("should return true for REVENUE", func() {
				Expect(account.Revenue.IsValid()).To(BeTrue())
			})
			It("should return true for EXPENSE", func() {
				Expect(account.Expense.IsValid()).To(BeTrue())
			})
		})

		Context("error cases", func() {
			It("should return false for invalid account type", func() {
				Expect(account.AccountType("invalid").IsValid()).To(BeFalse())
			})
		})
	})
})
