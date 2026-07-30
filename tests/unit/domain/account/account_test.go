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
	Describe("#NewAccountBuilder", func() {
		Context("success cases", func() {
			It("should not return an error when build new account", func() {
				acc, err := account.NewAccountBuilder().
					WithID("d589965c-1622-4329-98f9-f13354a2e4dc").
					WithOwnerID("d589965c-1622-4329-98f9-f13354a2e4dc").
					WithAccountExternalID("any_external_id").
					WithAccountNumber("123456").
					WithAccountType(string(account.Asset)).
					WithCurrency(string(money.BRL)).
					Build()
				Expect(err).To(BeNil())
				Expect(acc).NotTo(BeNil())
			})
		})

		Context("error cases", func() {
			It("should return an error when account type is invalid", func() {
				_, err := account.NewAccountBuilder().
					WithID("d589965c-1622-4329-98f9-f13354a2e4dc").
					WithOwnerID("d589965c-1622-4329-98f9-f13354a2e4dc").
					WithAccountExternalID("any_external_id").
					WithAccountNumber("123456").
					WithAccountType("any_type").
					WithCurrency(string(money.BRL)).
					Build()
				Expect(err).NotTo(BeNil())
			})

			It("should return an error when external id is empty", func() {
				_, err := account.NewAccountBuilder().
					WithID("d589965c-1622-4329-98f9-f13354a2e4dc").
					WithOwnerID("d589965c-1622-4329-98f9-f13354a2e4dc").
					WithAccountExternalID("").
					WithAccountNumber("123456").
					WithAccountType(string(account.Asset)).
					WithCurrency(string(money.BRL)).
					Build()
				Expect(err).NotTo(BeNil())
			})

			It("should return an error when currency is invalid", func() {
				_, err := account.NewAccountBuilder().
					WithID("d589965c-1622-4329-98f9-f13354a2e4dc").
					WithOwnerID("d589965c-1622-4329-98f9-f13354a2e4dc").
					WithAccountExternalID("any_external_id").
					WithAccountNumber("123456").
					WithAccountType(string(account.Asset)).
					WithCurrency("INVALID").
					Build()
				Expect(err).NotTo(BeNil())
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
