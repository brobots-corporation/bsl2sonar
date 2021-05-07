package finder

import (
	"github.com/stretchr/testify/suite"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
)

var AbsPathTestSrcFolder, _ = filepath.Abs("../tests/test_conf")
var AbsPathFixtureSonarFile, _ = filepath.Abs("../tests/fixture-sonar-project.properties")
var AbsPathFixtureUnicodeSonarFile, _ = filepath.Abs("../tests/fixture-unicode-sonar-project.properties")
var AbsPathFixtureStdoutFile, _ = filepath.Abs("../tests/fixture-stdout")
var AbsPathFixtureUnicodeStdoutFile, _ = filepath.Abs("../tests/fixture-unicode-stdout")
var AbsPathTemplateSonarFile, _ = filepath.Abs("../tests/template-sonar-project.properties")
var AbsPathTestSonarFile, _ = filepath.Abs("../tests/test-sonar-project.properties")
var AbsPathTestSonarUnicodeFile, _ = filepath.Abs("../tests/test-unicode-sonar-project.properties")

var phrases = "рн_ пс_"
var countSubsystemsFilesPaths = 7
var CountGetObjectsNamesFromSubsystem = 2
var CountGetListMetadataName = 24
var CountGetListBslFiles = 3
var CountGetBslFilesPaths = 63
var subsystemFilePath,_ = filepath.Abs(path.Join(AbsPathTestSrcFolder, "Subsystems/рн_Супер.xml"))
var ObjectFolderPath,_ = filepath.Abs(path.Join(AbsPathTestSrcFolder, "Catalogs/Справочник8"))
var pattern = "*.bsl"
var CountLineBslFiles = 3873
var CountLineBslFilesUnicode = 6833

type FinderTestSuite struct {
	suite.Suite
	BaseFinder       *Finder
	BaseFinderStdOut *Finder
	BaseFinderUnicodeStdOut *Finder
	BaseFinderFileOut *Finder
	BaseFinderUnicodeFileOut *Finder
	//BaseFinderStdOutVerbose *Finder
	fsppContent string // fixture-sonar-project.properties
	fusppContent string // fixture-unicode-sonar-project.properties
	fsContent string // fixture-stdout
	fusContent string // fixture-unicode-stdout
}

func copy(scr string, dst string) {
	// Open source file
	srcFile, err := os.Open(scr)
	if err != nil {
		log.Fatal(err)
	}
	defer srcFile.Close()

	// Create distination file
	dstFile, err := os.Create(dst)
	if err != nil {
		log.Fatal(err)
	}
	defer dstFile.Close()

	// Copy content
	_, err = io.Copy(dstFile, srcFile)
	if err != nil {
		log.Fatal(err)
	}
}

func delete(file string) {
	err := os.Remove(file)
	if err != nil {
		log.Fatal(err)
	}
}

// The SetupSuite method will be run by testify once, at the very
// start of the testing suite, before any tests are run
func (suite *FinderTestSuite) SetupSuite() {
	suite.BaseFinder = NewFinder(AbsPathTestSrcFolder, phrases)
	suite.BaseFinder.Logging = true

	suite.BaseFinderStdOut = NewFinder(AbsPathTestSrcFolder, phrases)
	suite.BaseFinderStdOut.Logging = true

	suite.BaseFinderUnicodeStdOut = NewFinder(AbsPathTestSrcFolder, phrases)
	suite.BaseFinderUnicodeStdOut.Logging = true
	suite.BaseFinderUnicodeStdOut.Unicode = true

	suite.BaseFinderFileOut = NewFinder(AbsPathTestSrcFolder, phrases)
	suite.BaseFinderFileOut.Logging = true
	suite.BaseFinderFileOut.Sfile = AbsPathTestSonarFile

	suite.BaseFinderUnicodeFileOut = NewFinder(AbsPathTestSrcFolder, phrases)
	suite.BaseFinderUnicodeFileOut.Logging = true
	suite.BaseFinderUnicodeFileOut.Unicode = true
	suite.BaseFinderUnicodeFileOut.Sfile = AbsPathTestSonarUnicodeFile

	// read fixture-sonar-project.properties file
	fspp, _ := ioutil.ReadFile(AbsPathFixtureSonarFile)
	// content of file
	suite.fsppContent = string(fspp)

	// read fixture-unicode-sonar-project.properties file
	fuspp, _ := ioutil.ReadFile(AbsPathFixtureUnicodeSonarFile)
	// content of file
	suite.fusppContent = string(fuspp)

	// read fixture-stdout file
	fs, _ := ioutil.ReadFile(AbsPathFixtureStdoutFile)
	// content of file
	suite.fsContent = string(fs)

	// read fixture-unicode-stdout file
	fus, _ := ioutil.ReadFile(AbsPathFixtureUnicodeStdoutFile)
	// content of file
	suite.fusContent = string(fus)

	copy(AbsPathTemplateSonarFile, AbsPathTestSonarFile)
	copy(AbsPathTemplateSonarFile, AbsPathTestSonarUnicodeFile)
}

// The TearDownSuite method will be run by testify once, at the very
// end of the testing suite, after all tests have been run
func (suite *FinderTestSuite) TearDownSuite() {
	delete(AbsPathTestSonarFile)
	delete(AbsPathTestSonarUnicodeFile)
}

func (suite *FinderTestSuite) TestNewFinder() {
	suite.IsType(&Finder{}, suite.BaseFinder)
}

func (suite *FinderTestSuite) TestStringToUnicode() {
	testString := "Проверка преобразования в символы unicode"
	caseString := suite.BaseFinder.stringToUnicode(testString)
	StandardString := strings.Trim(strconv.QuoteToASCII(testString), "\"")
	suite.Equal(StandardString, caseString)
}

func (suite *FinderTestSuite) TestGetSubsystemsFilesPaths() {
	subsystemsFilesPaths := suite.BaseFinder.getSubsystemsFilesPaths()
	suite.Equal(countSubsystemsFilesPaths, len(subsystemsFilesPaths))
}

func (suite *FinderTestSuite) TestGetObjectsNamesFromSubsystem()  {
	metadataNames := suite.BaseFinder.getObjectsNamesFromSubsystem(subsystemFilePath)
	suite.Equal(CountGetObjectsNamesFromSubsystem, len(metadataNames))
}

func (suite *FinderTestSuite) TestGetSliceMetadataName() {
	sliceMetadataNames := suite.BaseFinder.getSliceMetadataName()
	suite.Equal(CountGetListMetadataName, len(sliceMetadataNames))
}

func (suite *FinderTestSuite) TestGetSliceFiles() {
	sliceFiles := suite.BaseFinder.getSliceFiles(ObjectFolderPath, pattern)
	suite.Equal(CountGetListBslFiles, len(sliceFiles))
}

func (suite *FinderTestSuite) TestGetBslFilesPaths() {
	sliceBslFilesPaths := suite.BaseFinder.getBslFilesPaths()
	suite.Equal(CountGetBslFilesPaths, len(sliceBslFilesPaths))
}

func (suite *FinderTestSuite) TestGetBslFilesLine() {
	lineBslFiles := suite.BaseFinder.getBslFilesLine()
	lineBslFilesUnicode := suite.BaseFinderUnicodeStdOut.getBslFilesLine()

	suite.Equal(CountLineBslFiles, len(lineBslFiles))
	suite.Equal(CountLineBslFilesUnicode, len(lineBslFilesUnicode))
}

func (suite *FinderTestSuite) TestWriteBslLineToFile() {
	// write to AbsPathTestSonarFile
	suite.BaseFinderFileOut.writeBslLineToFile()
	// read from AbsPathTestSonarFile
	tsf, _ := ioutil.ReadFile(AbsPathTestSonarFile)
	suite.Equal(len(suite.fsppContent), len(string(tsf)))

	// write to AbsPathTestSonarUnicodeFile
	suite.BaseFinderUnicodeFileOut.writeBslLineToFile()
	// read from AbsPathTestSonarUnicodeFile
	tusf, _ := ioutil.ReadFile(AbsPathTestSonarUnicodeFile)
	suite.Equal(len(suite.fusppContent), len(string(tusf)))
}

func (suite *FinderTestSuite) TestWriteBslLineToSTDOUT() {
	// test without unicode
	r, w, _ := os.Pipe()
	os.Stdout = w

	suite.BaseFinderStdOut.writeBslLineToSTDOUT()

	w.Close()

	stdout, _ := ioutil.ReadAll(r)

	suite.Equal(len(suite.fsContent), len(string(stdout)))

	// test with unicode
	r, w, _ = os.Pipe()
	os.Stdout = w

	suite.BaseFinderUnicodeStdOut.writeBslLineToSTDOUT()

	w.Close()

	stdout, _ = ioutil.ReadAll(r)

	suite.Equal(len(suite.fusContent), len(string(stdout)))

}

func TestSuite(t *testing.T) {
	suite.Run(t, new(FinderTestSuite))
}
