//go:build unit
// +build unit

package transaction_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestTransaction(t *testing.T) {
	suiteConfig, reporterConfig := GinkgoConfiguration()

	suiteConfig.SkipStrings = []string{"SKIPPED", "PENDING", "NEVER-RUN", "SKIP"}
	reporterConfig.FullTrace = true
	reporterConfig.Verbose = true

	RegisterFailHandler(Fail)
	RunSpecs(t, "Transaction Domain Suite", suiteConfig, reporterConfig)
}
