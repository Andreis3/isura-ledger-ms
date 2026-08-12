//go:build unit
// +build unit

package outbox_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestOutbox(t *testing.T) {
	suiteConfig, reporterConfig := GinkgoConfiguration()

	suiteConfig.SkipStrings = []string{"SKIPPED", "PENDING", "NEVER-RUN", "SKIP"}
	reporterConfig.FullTrace = true
	reporterConfig.Verbose = true

	RegisterFailHandler(Fail)
	RunSpecs(t, "Outbox Domain Suite", suiteConfig, reporterConfig)
}
