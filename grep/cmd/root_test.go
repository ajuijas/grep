package cmd

import (
	"bytes"
	"os"
	"strings"
	"testing"
)

func redirectStd() (func() string) {

	// Redirecting the std out and error and returns a function which will gives the std out strings.

	oldStdout := os.Stdout
	oldStderr := os.Stderr

	rOut, wOut, _ := os.Pipe()
	rErr, wErr, _ := os.Pipe()
	os.Stdout = wOut
	os.Stderr = wErr

	return func() string {
		wOut.Close()

		var buf bytes.Buffer
		_, _ = buf.ReadFrom(rOut)
		out := buf.String()

		rOut.Close()
		rErr.Close()
		wErr.Close()
		os.Stdout = oldStdout
		os.Stderr = oldStderr
		return out
	}
}

func Test_ReadFile(t *testing.T) {

	// Redirect stdout and stderr to have a clean test output.
	restoreStd := redirectStd()

	tmpFile, err := os.CreateTemp(".", "foo*.txt")
	if err != nil {
		t.Fatalf("Failed to create temp file for tests")
	}
	defer os.Remove(tmpFile.Name())
	fileContent := "abcd\nabcdefg\n"
	_, err = tmpFile.Write([]byte(fileContent))
	if err != nil {
		t.Fatalf("Error writing to temp file:")
		return
	}

	tests := []struct {
		fileName, expected string
	}{
		{"no_such_file.txt", "no_such_file.txt: open: No such file or directory"},
		{"no_such_directory", "no_such_directory: open: No such file or directory"},
		{"no_permission", "no_permission: open: No such file or directory"},
		{tmpFile.Name(), ""},
	}

	for _, test := range tests {

		rootCmd.SetArgs([]string{"abcd", test.fileName})
		err := rootCmd.Execute()
		var outString string
		if err != nil {
			outString = err.Error()
		} else {
			outString = "" // No error from rootCmd.Execute()
		}

		expected := strings.TrimSpace(test.expected)

		if outString != expected {
			t.Errorf("Expected: %v Got: %v", expected, outString)
		}
	}
	restoreStd()
}

func Test_rootCmd(t *testing.T) {

	tests := []struct {
		stringPattern, fileContent, expected string
	} {
		{"abcd", "abcdefg\nhijklmn", "abcdefg"},
	}

	for _, test := range tests{

		// Redirect stdout and stderr to have a clean test output.
		restoreStd := redirectStd()

		tmpFile, err := os.CreateTemp(".", "temp_*.txt")
		if err != nil {
			t.Fatalf("Error createing temp file")
		}
		defer os.Remove(tmpFile.Name())
		_, err = tmpFile.WriteString(test.fileContent)
		if err!=nil {
			t.Fatalf("Error while writing to temp file")
		}

		rootCmd.SetArgs([]string{test.stringPattern, tmpFile.Name()})
		rootCmd.Execute()

		out := restoreStd()

		if out != test.expected {
			t.Errorf("Expected: %v Got: %v", test.expected, out)
		}
	}
}
