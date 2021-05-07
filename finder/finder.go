/*
Copyright © 2021 ALEKSEY MAKSIMKIN <maximkin@mail.ru>

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program. If not, see <http://www.gnu.org/licenses/>.
*/

package finder

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"github.com/thoas/go-funk"
	"io/fs"
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"text/template"
)

// Finder is a structure for find and out finding data
type Finder struct {
	srcdir             string
	phrases            string
	Sfile              string `json:"path to sonar-project.properties"`
	Abspath            bool
	Logging            bool
	Unicode            bool `json:"convert Cyrillic symbols to unicode"`
	Generate           bool `json:"generate out data to template"`
	keywordLine        string
	rootSubsystemsPath string
	Logger             *log.Logger
}

// NewFinder is the method for create new finder structure
func NewFinder(srcdir string, phrases string) *Finder {
	finder := &Finder{
		srcdir:             srcdir,
		phrases:            phrases,
		keywordLine:        "$inclusions_line",
		rootSubsystemsPath: path.Join(srcdir, "Subsystems"),
		Logger:             log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime),
	}

	return finder
}

func (f *Finder) stringToUnicode(str string) string {
	// Transform cyrillic symbols to unicode ascii
	return strings.Trim(strconv.QuoteToASCII(str), "\"")
}

func (f *Finder) getSubsystemsFilesPaths() []string {

	var subsystemsFilesPaths []string

	prfxs := strings.Split(f.phrases, " ")
	for _, prfx := range prfxs {

		sPattern := prfx + "*.xml"

		err := filepath.Walk(f.rootSubsystemsPath, func(wpath string, info fs.FileInfo, err error) error {
			if info == nil || !info.IsDir() {
				return nil
			}
			sFiles, _ := filepath.Glob(path.Join(wpath, sPattern))
			subsystemsFilesPaths = append(subsystemsFilesPaths, sFiles...)

			return nil
		})
		if err != nil {
			println(err.Error())
			return []string{}
		}

	}

	if f.Logging {
		f.Logger.Printf(">>> Найдено подсистем для анализа: %d", len(subsystemsFilesPaths))
	}

	return subsystemsFilesPaths
}

func (f *Finder) getObjectsNamesFromSubsystem(filename string) []string {

	// structure for unmarshal xml
	type Metadata struct {
		Names []string `xml:"Subsystem>Properties>Content>Item"`
	}

	// slice for collect all metadata names
	var MetadataNames []string

	// mask for exclusion deleted metadata
	mask := "[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}"
	re := regexp.MustCompile(mask)

	// open xml file
	xmlFile, err := os.Open(filename)
	if err != nil {
		println(err.Error())
		return []string{}
	}

	defer func(xmlFile *os.File) {
		err := xmlFile.Close()
		if err != nil {
			println(err.Error())
		}
	}(xmlFile)

	// read content of xml file
	byteValue, _ := ioutil.ReadAll(xmlFile)

	// unmarshalling
	m := Metadata{Names: []string{}}
	err = xml.Unmarshal(byteValue, &m)
	if err != nil {
		println(err.Error())
		return []string{}
	}

	// check metadata (not deleted or empty) and append to slice
	for _, item := range m.Names {
		if len(item) != 0 && !re.Match([]byte(item)) {
			MetadataNames = append(MetadataNames, item)
		}
	}

	return MetadataNames
}

func (f *Finder) getSliceMetadataName() []string {

	// slice for collect all metadata names
	var SliceMetadataNames []string

	// get subsystems
	SubsystemsFilesPaths := f.getSubsystemsFilesPaths()

	// get bsl files by subsystems
	for _, SubPath := range SubsystemsFilesPaths {

		if f.Logging {
			f.Logger.Printf("%s", SubPath)
		}
		SliceMetadataNames = append(SliceMetadataNames, f.getObjectsNamesFromSubsystem(SubPath)...)

	}

	sort.Strings(SliceMetadataNames)

	if f.Logging {
		f.Logger.Printf(">>> Найдено объектов для анализа: %d", len(SliceMetadataNames))
	}

	return funk.UniqString(SliceMetadataNames)
}

func (f *Finder) getSliceFiles(PathToFolder string, pattern string) []string {

	var SliceFiles []string

	err := filepath.Walk(PathToFolder, func(wpath string, info fs.FileInfo, err error) error {
		if info == nil || !info.IsDir() {
			return nil
		}
		sFiles, _ := filepath.Glob(path.Join(wpath, pattern))
		SliceFiles = append(SliceFiles, sFiles...)

		return nil
	})
	if err != nil {
		println(err.Error())
		return []string{}
	}

	return SliceFiles
}

func (f *Finder) getBslFilesPaths() []string {

	SliceMetadataName := f.getSliceMetadataName()

	var SliceBslFilesPaths []string

	for _, MetadataName := range SliceMetadataName {

		MetadataTypeName := MetadataName[:strings.Index(MetadataName, ".")] + "s"
		MetadataOnlyName := MetadataName[strings.Index(MetadataName, ".")+1:]
		MetadataRelPath := path.Join(MetadataTypeName, MetadataOnlyName)
		PathToFolder := path.Join(f.srcdir, MetadataRelPath)

		// check folder exist
		_, err := os.Stat(PathToFolder)
		if os.IsNotExist(err) {
			continue
		}

		// get slice of bsl files in folder
		BslFiles := f.getSliceFiles(PathToFolder, "*.bsl")

		if !f.Abspath {
			// transform path to bsl files without basepath
			for idx, file := range BslFiles {
				BslFiles[idx], _ = filepath.Rel(f.srcdir, file)
			}
		}

		SliceBslFilesPaths = append(SliceBslFilesPaths, BslFiles...)

	}

	if f.Logging {
		f.Logger.Printf(">>> Количество bsl модулей для проверки: %d", len(SliceBslFilesPaths))
	}

	return SliceBslFilesPaths
}

func (f *Finder) getBslFilesLine() string {

	SliceBslFilesPaths := f.getBslFilesPaths()

	// convert Cyrillic symbols to unicode ascii
	if f.Unicode {
		for idx := range SliceBslFilesPaths {
			SliceBslFilesPaths[idx] = f.stringToUnicode(SliceBslFilesPaths[idx])
		}
	}

	var LineBslFiles string

	EndLine := ", \\\n"

	// make one line list of bsl files
	for _, BslFilePath := range SliceBslFilesPaths {
		LineBslFiles = LineBslFiles + BslFilePath + EndLine
	}

	// eval last line end index without spec symbols
	LastLineEndIndex := len(LineBslFiles) - len(EndLine)
	LineBslFiles = LineBslFiles[:LastLineEndIndex]

	return LineBslFiles
}

func (f *Finder) writeBslLineToFile() {

	LineBslFiles := f.getBslFilesLine()

	var spfContent string

	// write
	if f.Generate {

		// read template
		ts, err := template.ParseFiles("./template/sonar-project.properties")
		if err != nil {
			fmt.Print(err)
		}

		buf := bytes.NewBufferString("")

		err = ts.Execute(buf, LineBslFiles)
		if err != nil {
			fmt.Print(err)
		}

		spfContent = buf.String()

	} else {

		// read sonar-project.properties file
		spf, err := ioutil.ReadFile(f.Sfile)
		if err != nil {
			fmt.Print(err)
		}
		// content of file
		spfContent = string(spf)
		spfContent = strings.Replace(spfContent, f.keywordLine, LineBslFiles, -1)

	}

	var writeErr error
	// write sonar properties content to file

	if f.Generate {
		writeErr = ioutil.WriteFile(f.Sfile, []byte(spfContent), fs.ModePerm)
	} else {
		writeErr = ioutil.WriteFile(f.Sfile, []byte(spfContent), fs.ModeExclusive)
	}

	if writeErr != nil {
		fmt.Println(writeErr)
	}
}

func (f *Finder) writeBslLineToSTDOUT() {

	LineBslFiles := f.getBslFilesPaths()

	for idx := range LineBslFiles {
		// convert Cyrillic symbols to unicode ascii and print
		if f.Unicode {
			fmt.Println(f.stringToUnicode(LineBslFiles[idx]))
		} else {
			fmt.Println(LineBslFiles[idx])
		}
	}
}

// DataToSonarQube is a method for output data
func (f *Finder) DataToSonarQube() {

	if len(f.Sfile) != 0 {
		f.writeBslLineToFile()
	} else {
		f.writeBslLineToSTDOUT()
	}
}
