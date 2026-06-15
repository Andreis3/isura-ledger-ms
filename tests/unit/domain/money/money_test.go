//go:build unit
// +build unit

package money_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/andreis3/isura-ledger-ms/internal/domain/money"
)

var _ = Describe("INTERNAL :: DOMAIN :: MONEY :: MONEY", func() {
	Describe("#NewMoney", func() {
		Context("success cases", func() {
			It("should not return an error when build new money", func() {
				m, err := money.NewMoney(100, money.BRL)
				Expect(err).To(BeNil())
				Expect(m.Amount()).To(Equal(int64(100)))
				Expect(m.Currency()).To(Equal(money.BRL))
			})
		})

		Context("error cases", func() {
			It("should return an error when amount is negative", func() {
				_, err := money.NewMoney(-1, money.BRL)
				Expect(err).NotTo(BeNil())
				Expect(err).To(Equal(money.ErrNegativeAmount))
			})

			It("should return an error when currency is invalid", func() {
				_, err := money.NewMoney(100, "INVALID")
				Expect(err).NotTo(BeNil())
				Expect(err).To(Equal(money.ErrInvalidCurrency))
			})
		})
	})

	Describe("#Add", func() {
		Context("success cases", func() {
			It("should add two money values", func() {
				m1, _ := money.NewMoney(100, money.BRL)
				m2, _ := money.NewMoney(50, money.BRL)
				result, err := m1.Add(m2)
				Expect(err).To(BeNil())
				Expect(result.Amount()).To(Equal(int64(150)))
			})
		})

		Context("error cases", func() {
			It("should return an error when currencies mismatch", func() {
				m1, _ := money.NewMoney(100, money.BRL)
				m2, _ := money.NewMoney(50, money.USD)
				_, err := m1.Add(m2)
				Expect(err).NotTo(BeNil())
				Expect(err).To(Equal(money.ErrCurrencyMismatch))
			})
		})
	})

	Describe("#Subtract", func() {
		Context("success cases", func() {
			It("should subtract two money values", func() {
				m1, _ := money.NewMoney(100, money.BRL)
				m2, _ := money.NewMoney(50, money.BRL)
				result, err := m1.Subtract(m2)
				Expect(err).To(BeNil())
				Expect(result.Amount()).To(Equal(int64(50)))
			})
		})

		Context("error cases", func() {
			It("should return an error when currencies mismatch", func() {
				m1, _ := money.NewMoney(100, money.BRL)
				m2, _ := money.NewMoney(50, money.USD)
				_, err := m1.Subtract(m2)
				Expect(err).NotTo(BeNil())
				Expect(err).To(Equal(money.ErrCurrencyMismatch))
			})
		})
	})

	Describe("Utility functions", func() {
		It("IsZero", func() {
			m1, _ := money.NewMoney(0, money.BRL)
			Expect(m1.IsZero()).To(BeTrue())
			m2, _ := money.NewMoney(100, money.BRL)
			Expect(m2.IsZero()).To(BeFalse())
		})

		It("IsNegative", func() {
			// Since NewMoney prevents negative, we construct it directly or modify for test if needed,
			// but we can't export fields. The current implementation of NewMoney prevents < 0.
			// Let's test IsNegative using Subtract result.
			m1, _ := money.NewMoney(50, money.BRL)
			m2, _ := money.NewMoney(100, money.BRL)
			m3, _ := m1.Subtract(m2)
			Expect(m3.IsNegative()).To(BeTrue())
		})

		It("IsPositive", func() {
			m1, _ := money.NewMoney(100, money.BRL)
			Expect(m1.IsPositive()).To(BeTrue())
			m2, _ := money.NewMoney(0, money.BRL)
			Expect(m2.IsPositive()).To(BeFalse())
		})

		It("Equal", func() {
			m1, _ := money.NewMoney(100, money.BRL)
			m2, _ := money.NewMoney(100, money.BRL)
			m3, _ := money.NewMoney(50, money.BRL)
			m4, _ := money.NewMoney(100, money.USD)
			Expect(m1.Equal(m2)).To(BeTrue())
			Expect(m1.Equal(m3)).To(BeFalse())
			Expect(m1.Equal(m4)).To(BeFalse())
		})

		It("IsSufficientBalance", func() {
			m1, _ := money.NewMoney(100, money.BRL)
			m2, _ := money.NewMoney(50, money.BRL)
			m3, _ := money.NewMoney(150, money.BRL)
			m4, _ := money.NewMoney(100, money.USD)
			Expect(m1.IsSufficientBalance(m2)).To(BeTrue())
			Expect(m1.IsSufficientBalance(m1)).To(BeTrue())
			Expect(m1.IsSufficientBalance(m3)).To(BeFalse())
			Expect(m1.IsSufficientBalance(m4)).To(BeFalse())
		})

		It("String", func() {
			m1, _ := money.NewMoney(1050, money.BRL)
			Expect(m1.String()).To(Equal("10.50 BRL"))
			m2, _ := money.NewMoney(5, money.BRL)
			Expect(m2.String()).To(Equal("0.05 BRL"))
		})
	})
})
