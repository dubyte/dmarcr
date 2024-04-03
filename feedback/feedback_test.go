package feedback

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	v := New()
	assert.NotNil(t, v, "New() should return a non-nil Values instance")
}

func TestReadFromFile(t *testing.T) {
	cases := []struct {
		name     string
		file     string
		hasError bool
		count    int
	}{
		{
			name:     "ValidFile",
			file:     "testdata/valid_report.xml",
			hasError: false,
			count:    1,
		},
		{
			name:     "NonExistentFile",
			file:     "testdata/nonexistent_file.xml",
			hasError: true,
			count:    0,
		},
		{
			name:     "InvalidFile",
			file:     "testdata/invalid_report.xml",
			hasError: true,
			count:    0,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			// Setup
			reports := New()

			// Act
			err := reports.ReadFromFile(tc.file)

			// Assert
			if tc.hasError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Len(t, reports.Feedbacks, tc.count, "Expected feedback count mismatch")
			}
		})
	}
}

func TestReadFromStdin(t *testing.T) {
	cases := []struct {
		name     string
		input    string
		hasError bool
		count    int
	}{
		{
			name:     "ValidXML",
			input:    readTestData("testdata/valid_report.xml", t),
			hasError: false,
			count:    1,
		},
		{
			name:     "InvalidXML",
			input:    readTestData("testdata/invalid_report.xml", t),
			hasError: true,
			count:    0,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			// Setup
			reports := New()

			// Redirect stdin for the test
			restoreFunc, err := redirectStdin(tc.input)
			if err != nil {
				t.Fatalf("%v", err)
			}
			defer restoreFunc()

			// Act
			err = reports.ReadFromStdin()

			// Assert
			if tc.hasError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Len(t, reports.Feedbacks, tc.count, "Expected feedback count mismatch")
			}
		})
	}
}

func redirectStdin(input string) (restoreFunc func(), err error) {
	oldStdin := os.Stdin
	r, w, err := os.Pipe()
	if err != nil {
		return nil, fmt.Errorf("error creating pipe: %v", err)
	}
	os.Stdin = r
	_, err = w.WriteString(input)
	if err != nil {
		return nil, fmt.Errorf("error writing to pipe: %v", err)
	}
	w.Close()

	return func() { os.Stdin = oldStdin }, nil
}

func TestString(t *testing.T) {
	// Setup
	reports := New()

	// Act
	err := reports.ReadFromFile("testdata/valid_report.xml")
	require.NoError(t, err, "ReadFromFile() should not return an error with a valid file")
	output := reports.String()

	// Assert
	// Check for table headers
	assert.Contains(t, output, "Source IP", "String() output should contain 'Source IP'")
	assert.Contains(t, output, "Count", "String() output should contain 'Count'")
	assert.Contains(t, output, "Disposition", "String() output should contain 'Disposition'")
	assert.Contains(t, output, "DKIM", "String() output should contain 'DKIM'")
	assert.Contains(t, output, "SPF", "String() output should contain 'SPF'")

	// Check for specific values from the valid report
	assert.Contains(t, output, "192.0.2.1", "String() output should contain the source IP '192.0.2.1'")
	assert.Contains(t, output, "1", "String() output should contain the count '1'")
	assert.Contains(t, output, "none", "String() output should contain the disposition 'none'")
	assert.Contains(t, output, "pass", "String() output should contain the DKIM result 'pass'")
	assert.Contains(t, output, "pass", "String() output should contain the SPF result 'pass'")
}

// readTestData reads the content of a test data file and returns it as a string.
func readTestData(filename string, t *testing.T) string {
	t.Helper() // Mark this function as a test helper
	data, err := os.ReadFile(filename)
	if err != nil {
		t.Fatalf("Error reading test data file %s: %v", filename, err)
	}
	return string(data)
}
