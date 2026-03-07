package analyzer_test

import (
	"testing"

	"github.com/barashF/test-linter/analyzer"
	"golang.org/x/tools/go/analysis/analysistest"
)

func TestLogLint(t *testing.T) {
	analysistest.Run(t, analysistest.TestData(), analyzer.Analyzer, "tests")
}
