//go:build unit
// +build unit

package outbox_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/andreis3/isura-ledger-ms/internal/domain/outbox"
)

var _ = Describe("INTERNAL :: DOMAIN :: OUTBOX :: OUTBOX", func() {
	Describe("#NewOutbox", func() {
		Context("success cases", func() {
			It("should create a new outbox entry", func() {
				ob := outbox.NewOutbox("any_id", "aggregate_id", outbox.Transaction, outbox.TransactionCreated, []byte(`{"key":"value"}`))
				Expect(ob).NotTo(BeNil())
				Expect(ob.Status).To(Equal(outbox.Pending))
				Expect(ob.Attempts).To(Equal(0))
			})
		})
	})

	Describe("#Publish", func() {
		Context("success cases", func() {
			It("should transition from PENDING to SUCCESS", func() {
				ob := outbox.NewOutbox("any_id", "aggregate_id", outbox.Transaction, outbox.TransactionCreated, []byte(`{"key":"value"}`))
				err := ob.Publish()
				Expect(err).To(BeNil())
				Expect(ob.Status).To(Equal(outbox.Success))
				Expect(ob.PublishedAt).NotTo(BeNil())
			})
		})

		Context("error cases", func() {
			It("should return error when transitioning from SUCCESS to SUCCESS", func() {
				ob := outbox.NewOutbox("any_id", "aggregate_id", outbox.Transaction, outbox.TransactionCreated, []byte(`{"key":"value"}`))
				ob.Publish() // now SUCCESS
				err := ob.Publish()
				Expect(err).NotTo(BeNil())
				Expect(err).To(Equal(outbox.ErrTransitionStatus))
			})
		})
	})

	Describe("#MarkFailed", func() {
		Context("success cases", func() {
			It("should transition from PENDING to FAILED", func() {
				ob := outbox.NewOutbox("any_id", "aggregate_id", outbox.Transaction, outbox.TransactionCreated, []byte(`{"key":"value"}`))
				err := ob.MarkFailed()
				Expect(err).To(BeNil())
				Expect(ob.Status).To(Equal(outbox.Failed))
				Expect(ob.Attempts).To(Equal(1))
				Expect(ob.LastAttemptAt).NotTo(BeNil())
			})
		})

		Context("error cases", func() {
			It("should return error when transitioning from SUCCESS to FAILED", func() {
				ob := outbox.NewOutbox("any_id", "aggregate_id", outbox.Transaction, outbox.TransactionCreated, []byte(`{"key":"value"}`))
				ob.Publish() // now SUCCESS
				err := ob.MarkFailed()
				Expect(err).NotTo(BeNil())
				Expect(err).To(Equal(outbox.ErrTransitionStatus))
			})
		})
	})

	Describe("#Retry", func() {
		Context("success cases", func() {
			It("should transition from FAILED to PENDING", func() {
				ob := outbox.NewOutbox("any_id", "aggregate_id", outbox.Transaction, outbox.TransactionCreated, []byte(`{"key":"value"}`))
				ob.MarkFailed() // now FAILED, attempt 1
				err := ob.Retry()
				Expect(err).To(BeNil())
				Expect(ob.Status).To(Equal(outbox.Pending))
			})
		})

		Context("error cases", func() {
			It("should return error when max attempts exceeded", func() {
				ob := outbox.NewOutbox("any_id", "aggregate_id", outbox.Transaction, outbox.TransactionCreated, []byte(`{"key":"value"}`))
				ob.MarkFailed() // attempt 1
				ob.Retry()
				ob.MarkFailed() // attempt 2
				ob.Retry()
				ob.MarkFailed() // attempt 3

				err := ob.Retry()
				Expect(err).NotTo(BeNil())
				Expect(err).To(Equal(outbox.ErrMaxAttemptsExceeded))
			})

			It("should return error when transitioning from SUCCESS to PENDING", func() {
				ob := outbox.NewOutbox("any_id", "aggregate_id", outbox.Transaction, outbox.TransactionCreated, []byte(`{"key":"value"}`))
				ob.Publish() // now SUCCESS
				err := ob.Retry()
				Expect(err).NotTo(BeNil())
				Expect(err).To(Equal(outbox.ErrTransitionStatus))
			})
		})
	})

	Describe("#IsValid", func() {
		Context("success cases", func() {
			It("should return true for valid status", func() {
				Expect(outbox.Pending.IsValid()).To(BeTrue())
				Expect(outbox.Failed.IsValid()).To(BeTrue())
				Expect(outbox.Success.IsValid()).To(BeTrue())
			})
		})

		Context("error cases", func() {
			It("should return false for invalid status", func() {
				Expect(outbox.StatusOutbox("invalid").IsValid()).To(BeFalse())
			})
		})
	})
})
