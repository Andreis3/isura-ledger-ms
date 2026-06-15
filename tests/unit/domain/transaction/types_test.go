//go:build unit
// +build unit

package transaction_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/andreis3/isura-ledger-ms/internal/domain/transaction"
)

var _ = Describe("INTERNAL :: DOMAIN :: TRANSACTION :: TYPES", func() {
	Describe("#CanTransition", func() {
		Context("success cases", func() {
			It("should allow transition from PENDING to COMPLETED", func() {
				Expect(transaction.ValidStateMachine.CanTransition(transaction.Pending, transaction.Completed)).To(BeTrue())
			})

			It("should allow transition from PENDING to FAILED", func() {
				Expect(transaction.ValidStateMachine.CanTransition(transaction.Pending, transaction.Failed)).To(BeTrue())
			})
		})

		Context("error cases", func() {
			It("should not allow transition from COMPLETED to PENDING", func() {
				Expect(transaction.ValidStateMachine.CanTransition(transaction.Completed, transaction.Pending)).To(BeFalse())
			})

			It("should not allow transition from FAILED to PENDING", func() {
				Expect(transaction.ValidStateMachine.CanTransition(transaction.Failed, transaction.Pending)).To(BeFalse())
			})

			It("should not allow transition from an unknown state", func() {
				Expect(transaction.ValidStateMachine.CanTransition("unknown", transaction.Pending)).To(BeFalse())
			})
		})
	})
})
