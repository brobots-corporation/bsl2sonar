package cmd

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

var AbsPathTestSrcFolder, _ = filepath.Abs("../tests/test_conf")
var AbsPathTestFailFolder, _ = filepath.Abs("../tests/test_con")
var AbsPathTestFailFile, _ = filepath.Abs("../tests/fixture-stdout")
var AbsPathTestNoExistFile, _ = filepath.Abs("../tests/fixture_stdou")
var AbsPathTemplateSonarFile, _ = filepath.Abs("../tests/template-sonar-project.properties")

func TestArgsCount(t *testing.T) {
	testTable := []struct {
		stringArgs     []string
		expectedString string
	}{
		{
			[]string{AbsPathTestSrcFolder},
			"requires only two arguments",
		},
		{
			[]string{AbsPathTestSrcFolder, "рн_", "рн_"},
			"requires only two arguments",
		},
	}

	for _, testCase := range testTable {
		err := checkArgs(rootCmd, testCase.stringArgs)
		if err != nil {
			assert.Contains(t, err.Error(), testCase.expectedString)
		}
	}
}

func TestPhraseLen(t *testing.T) {
	testTable := []struct {
		stringArgs     []string
		expectedString string
	}{
		{
			[]string{AbsPathTestSrcFolder, "р_"},
			"must be at least 3 characters of parsephrases",
		},
		{
			[]string{AbsPathTestSrcFolder, "р"},
			"must be at least 3 characters of parsephrases",
		},
		{
			[]string{AbsPathTestSrcFolder, "рн_"},
			"must be at least 3 characters of parsephrases",
		},
	}

	for _, testCase := range testTable {
		err := checkArgs(rootCmd, testCase.stringArgs)
		if err != nil {
			assert.Contains(t, err.Error(), testCase.expectedString)
		}
	}
}

func TestIsArgsValid(t *testing.T) {
	testTable := []struct {
		stringArgs     []string
		fileFlag       string
		genFlag        bool
		expectedString string
	}{
		{
			[]string{AbsPathTestFailFolder, "рн_"},
			"",
			false,
			"dosn't exist",
		},
		{
			[]string{AbsPathTestFailFile, "рн_"},
			"",
			false,
			"is not directory",
		},
		{
			[]string{AbsPathTestSrcFolder, "рн_"},
			"",
			true,
			"",
		},
		{
			[]string{AbsPathTestSrcFolder, "рн_"},
			AbsPathTestNoExistFile,
			false,
			"file not found",
		},
	}

	for _, testCase := range testTable {
		_, errText := isArgsValid(testCase.stringArgs, testCase.fileFlag, testCase.genFlag)
		assert.Contains(t, errText, testCase.expectedString)
	}
}

func TestExistingFile(t *testing.T) {

	_, errText := isArgsValid([]string{AbsPathTestSrcFolder, "рн_"}, AbsPathTemplateSonarFile, false)
	assert.Equal(t, "", errText)

}
