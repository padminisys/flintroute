package runner

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"os"
	"strings"
	"sync"
	"time"
)

// TestResults holds all test results
type TestResults struct {
	Tests     []*TestResult `json:"tests" xml:"testcase"`
	StartTime time.Time     `json:"start_time" xml:"start_time,attr"`
	EndTime   time.Time     `json:"end_time" xml:"end_time,attr"`
	mu        sync.Mutex
}

// TestResult represents a single test result
type TestResult struct {
	Name     string        `json:"name" xml:"name,attr"`
	Status   string        `json:"status" xml:"status,attr"` // "passed", "failed", "skipped"
	Duration time.Duration `json:"duration" xml:"time,attr"`
	Error    string        `json:"error,omitempty" xml:"error,omitempty"`
	Output   string        `json:"output,omitempty" xml:"system-out,omitempty"`
}

// TestStats represents test statistics
type TestStats struct {
	Total    int           `json:"total"`
	Passed   int           `json:"passed"`
	Failed   int           `json:"failed"`
	Skipped  int           `json:"skipped"`
	Duration time.Duration `json:"duration"`
}

// NewTestResults creates a new test results collection
func NewTestResults() *TestResults {
	return &TestResults{
		Tests:     make([]*TestResult, 0),
		StartTime: time.Now(),
	}
}

// AddResult adds a test result to the collection
func (tr *TestResults) AddResult(result *TestResult) {
	tr.mu.Lock()
	defer tr.mu.Unlock()
	tr.Tests = append(tr.Tests, result)
}

// Finalize marks the end time
func (tr *TestResults) Finalize() {
	tr.EndTime = time.Now()
}

// GenerateJSONReport generates a JSON report
func (tr *TestResults) GenerateJSONReport(path string) error {
	tr.mu.Lock()
	defer tr.mu.Unlock()

	data, err := json.MarshalIndent(tr, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal JSON: %w", err)
	}

	if err := os.WriteFile(path, data, 0644); err != nil {
		return fmt.Errorf("failed to write JSON report: %w", err)
	}

	return nil
}

// JUnit XML types
type testSuite struct {
	XMLName   xml.Name    `xml:"testsuite"`
	Name      string      `xml:"name,attr"`
	Tests     int         `xml:"tests,attr"`
	Failures  int         `xml:"failures,attr"`
	Skipped   int         `xml:"skipped,attr"`
	Time      float64     `xml:"time,attr"`
	Timestamp string      `xml:"timestamp,attr"`
	TestCases []*testCase `xml:"testcase"`
}

type testCase struct {
	Name      string   `xml:"name,attr"`
	ClassName string   `xml:"classname,attr"`
	Time      float64  `xml:"time,attr"`
	Failure   *failure `xml:"failure,omitempty"`
	Skipped   *skipped `xml:"skipped,omitempty"`
	SystemOut string   `xml:"system-out,omitempty"`
}

type failure struct {
	Message string `xml:"message,attr"`
	Type    string `xml:"type,attr"`
	Content string `xml:",chardata"`
}

type skipped struct {
	Message string `xml:"message,attr"`
}

// GenerateXMLReport generates a JUnit-style XML report
func (tr *TestResults) GenerateXMLReport(path string) error {
	tr.mu.Lock()
	defer tr.mu.Unlock()

	stats := tr.GetStats()
	suite := testSuite{
		Name:      "FlintRoute Functional Tests",
		Tests:     stats.Total,
		Failures:  stats.Failed,
		Skipped:   stats.Skipped,
		Time:      stats.Duration.Seconds(),
		Timestamp: tr.StartTime.Format(time.RFC3339),
		TestCases: make([]*testCase, 0, len(tr.Tests)),
	}

	for _, test := range tr.Tests {
		tc := &testCase{
			Name:      test.Name,
			ClassName: "functional",
			Time:      test.Duration.Seconds(),
			SystemOut: test.Output,
		}

		if test.Status == "failed" {
			tc.Failure = &failure{
				Message: "Test failed",
				Type:    "AssertionError",
				Content: test.Error,
			}
		} else if test.Status == "skipped" {
			tc.Skipped = &skipped{
				Message: test.Error,
			}
		}

		suite.TestCases = append(suite.TestCases, tc)
	}

	data, err := xml.MarshalIndent(suite, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal XML: %w", err)
	}

	// Add XML header
	xmlData := []byte(xml.Header + string(data))

	if err := os.WriteFile(path, xmlData, 0644); err != nil {
		return fmt.Errorf("failed to write XML report: %w", err)
	}

	return nil
}

// PrintSummary prints a summary of test results to stdout
func (tr *TestResults) PrintSummary() {
	tr.mu.Lock()
	defer tr.mu.Unlock()

	stats := tr.GetStats()

	fmt.Println("\n" + strings.Repeat("=", 60))
	fmt.Println("Test Results Summary")
	fmt.Println(strings.Repeat("=", 60))
	fmt.Printf("Total Tests:    %d\n", stats.Total)
	fmt.Printf("Passed:         %d (%.1f%%)\n", stats.Passed, float64(stats.Passed)/float64(stats.Total)*100)
	fmt.Printf("Failed:         %d (%.1f%%)\n", stats.Failed, float64(stats.Failed)/float64(stats.Total)*100)
	fmt.Printf("Skipped:        %d (%.1f%%)\n", stats.Skipped, float64(stats.Skipped)/float64(stats.Total)*100)
	fmt.Printf("Total Duration: %s\n", stats.Duration)
	fmt.Println(strings.Repeat("=", 60))

	if stats.Failed > 0 {
		fmt.Println("\nFailed Tests:")
		for _, test := range tr.Tests {
			if test.Status == "failed" {
				fmt.Printf("  ❌ %s\n", test.Name)
				if test.Error != "" {
					fmt.Printf("     Error: %s\n", test.Error)
				}
			}
		}
	}

	if stats.Skipped > 0 {
		fmt.Println("\nSkipped Tests:")
		for _, test := range tr.Tests {
			if test.Status == "skipped" {
				fmt.Printf("  ⊘ %s\n", test.Name)
				if test.Error != "" {
					fmt.Printf("     Reason: %s\n", test.Error)
				}
			}
		}
	}

	fmt.Println()
}

// GetStats calculates and returns test statistics
func (tr *TestResults) GetStats() *TestStats {
	stats := &TestStats{
		Total:    len(tr.Tests),
		Duration: tr.EndTime.Sub(tr.StartTime),
	}

	for _, test := range tr.Tests {
		switch test.Status {
		case "passed":
			stats.Passed++
		case "failed":
			stats.Failed++
		case "skipped":
			stats.Skipped++
		}
	}

	return stats
}

// HasFailures returns true if any tests failed
func (tr *TestResults) HasFailures() bool {
	tr.mu.Lock()
	defer tr.mu.Unlock()

	for _, test := range tr.Tests {
		if test.Status == "failed" {
			return true
		}
	}
	return false
}

// GetFailedTests returns all failed tests
func (tr *TestResults) GetFailedTests() []*TestResult {
	tr.mu.Lock()
	defer tr.mu.Unlock()

	failed := make([]*TestResult, 0)
	for _, test := range tr.Tests {
		if test.Status == "failed" {
			failed = append(failed, test)
		}
	}
	return failed
}

// GetPassedTests returns all passed tests
func (tr *TestResults) GetPassedTests() []*TestResult {
	tr.mu.Lock()
	defer tr.mu.Unlock()

	passed := make([]*TestResult, 0)
	for _, test := range tr.Tests {
		if test.Status == "passed" {
			passed = append(passed, test)
		}
	}
	return passed
}

// GetSkippedTests returns all skipped tests
func (tr *TestResults) GetSkippedTests() []*TestResult {
	tr.mu.Lock()
	defer tr.mu.Unlock()

	skipped := make([]*TestResult, 0)
	for _, test := range tr.Tests {
		if test.Status == "skipped" {
			skipped = append(skipped, test)
		}
	}
	return skipped
}