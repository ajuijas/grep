package cmd

import (
	"bytes"
	"os"
	"reflect"
	"sort"
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

func createTempFile(filePath, fileName, fileContent string) (*os.File) {

	if err := os.MkdirAll(filePath, os.ModePerm); err != nil {
		panic(err)
	}

	tmpFile, err := os.CreateTemp(filePath, fileName)
	if err!= nil{
		panic(err)
	}
	_, err = tmpFile.WriteString(fileContent)
	if err!=nil {
		panic(err)
	}
	return tmpFile
}

func Test_ReadFile(t *testing.T) {

	// Redirect stdout and stderr to have a clean test output.
	restoreStd := redirectStd()
	defer restoreStd()

	tmpFile := createTempFile(".", "foo*.txt", "abcd\nabcdefg\n")
	defer os.Remove(tmpFile.Name())

	tests := []struct {
		fileName, expected string
	}{
		{"no_such_file.txt", "no_such_file.txt: open: No such file or directory"},
		{"no_such_directory", "no_such_directory: open: No such file or directory"},
		{"no_permission", "no_permission: open: No such file or directory"},
		{tmpFile.Name(), ""},
	}

	for _, test := range tests {

		rootCmd.SetArgs([]string{"not_a_matching_pattern", test.fileName})
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
}

func Test_rootCmd(t *testing.T) {

	tests := []struct {
		stringPattern, fileContent, expected string
	} {
		// I'm using the <filename> placeholder to replace the filename while asserting the result.
		{"abcd", "abcdefg\nhijklmn", "<filename> abcdefg"},
		{"aba", "skfjhksh ks sknabafksd lksdjf aabalsj \nnskfnl adbld\nlsjdfl ksjfn kdsjf\nablisjdijalsba", "<filename> skfjhksh ks sknabafksd lksdjf aabalsj"},
		{"a", "bc\nlsjfoj\nlsdjfljs skdnfks ksdjf", ""},
		{"a", "a", "<filename> a"},
		{"a","a\na\na", "<filename> a\n<filename> a\n<filename> a\n"},
		{"a","a\na\na\nb", "<filename> a\n<filename> a\n<filename> a\n"},

	}

	for _, test := range tests{

		// Redirect stdout and stderr to have a clean test output.
		restoreStd := redirectStd()

		tmpFile := createTempFile(".", "temp_*.txt", test.fileContent)

		defer os.Remove(tmpFile.Name())

		rootCmd.SetArgs([]string{test.stringPattern, tmpFile.Name()})
		err := rootCmd.Execute()
		if err!= nil{
			t.Fatalf("Error while Execute rootCmd error: %v", err)
		}

		out := strings.TrimSpace(restoreStd())
		expected := strings.TrimSpace(strings.ReplaceAll(test.expected, "<filename>", tmpFile.Name()))

		if out != expected {
			t.Errorf("Expected: <<%v>> Got: <<%v>>", expected, out)
		}
	}
}

func Test_read_from_directory(t *testing.T){
	restoreStd := redirectStd()

	tempFoleder := "temp/"

	defer os.RemoveAll(tempFoleder)

	file1 := createTempFile(tempFoleder + "folder1", "file1.txt", "a")
	file2 := createTempFile(tempFoleder + "folder1", "file2.txt", "a")
	file3 := createTempFile(tempFoleder + "folder2", "file1.txt", "a")

	rootCmd.SetArgs([]string{"a", tempFoleder})
	err := rootCmd.Execute()
	if err!= nil{
		t.Fatalf("Error while Execute rootCmd error: %v", err)
	}

	expected := file1.Name() + " a\n" + file2.Name() + " a\n" + file3.Name() + " a"

	outLines := strings.Split(strings.TrimSpace(restoreStd()), "\n")
	expectedLines := strings.Split(strings.TrimSpace(expected), "\n")

	sort.Strings(outLines)
	sort.Strings(expectedLines)

	if !reflect.DeepEqual(outLines, expectedLines) {
		// I'm not considering the order here.
		t.Errorf("Expected: <<%v>> Got: <<%v>>", expectedLines, outLines)
}
}

func Test_redirect_output_to_a_file(t *testing.T){

	restoreStd := redirectStd()
	tmpFile := createTempFile(".", "foo*.txt", "aaa\nbnd\naa")
	expected := "aaa\naa"

	resultFileName := "result.txt"

	defer os.Remove(tmpFile.Name())
	defer os.Remove(resultFileName)


	rootCmd.SetArgs([]string{"a", tmpFile.Name(), "-o", resultFileName})
	err := rootCmd.Execute()
	if err != nil {
		t.Fatal(err.Error())
	}

	out := restoreStd()
	if out != "" {
		t.Errorf("Expected: <<>> Got <<%v>>", out)
	}

	data, _ := os.ReadFile(resultFileName)
	got := strings.TrimSpace(string(data))
	if got != expected {
		t.Errorf("Expected: <<%v>> Got: <<%v>>", expected, string(data))
	}
}
