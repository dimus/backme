package backme

import (
	"log"
	"regexp"
	"time"
)

const shortForm = "2006-01-02"
const LogSep = "------"

var dateRe *regexp.Regexp
var timeZero time.Time
var now time.Time

type File struct {
	Name             string
	Path             string
	SortingTime      time.Time
	ModificationTime time.Time
}

type CategorizedFiles struct {
	GarbageFiles []File
	Recent       []File
	LastMonth    []File
	Years        map[int][]File
}

func Organize(conf *Config) error {
	dateRe = regexp.MustCompile(`([\d]{4})-([\d]{1,2})-([\d]{1,2})[^\d]`)
	now = time.Now()

	for i := range conf.InputDirs {
		err := organizeDir(i, conf)
		if err != nil {
			return err
		}
	}
	return nil
}

func organizeDir(dirIndex int, conf *Config) error {
	log.Printf("Processing directory %s.\n", conf.InputDirs[dirIndex].Path)
	dir := conf.InputDirs[dirIndex]
	filesRegexMap, err := collectFiles(&dir)
	if err != nil {
		return err
	}
	log.Println(LogSep)
	bins := categorizeFiles(filesRegexMap, dir.KeepAllFiles)
	err = moveFiles(&dir, bins, conf)
	return err
}
