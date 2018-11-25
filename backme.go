package backme

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"time"
)

const shortForm = "2006-01-02"

var dateRe *regexp.Regexp
var now time.Time

type File struct {
	Name             string
	Path             string
	ParsedTime       time.Time
	ModificationTime time.Time
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
	dir := conf.InputDirs[dirIndex]
	files, err := collectFiles(&dir)
	_ = files
	return err
}

func collectFiles(dir *InputDir) (map[string][]File, error) {
	files := initFilesMap(dir)

	err := filepath.Walk(dir.Path,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			file := filepath.Base(path)
			re, match, err := isBackedUpFile(file, dir.Files)
			if err != nil {
				return err
			}

			if match {
				f := getFile(path, info)
				fmt.Println(f)
				files[re] = append(files[re], f)
			}
			return nil
		})
	return files, err
}

func getFile(path string, info os.FileInfo) File {
	mod := info.ModTime()
	date := parseDate(path)
	return File{Path: path, ParsedTime: date, ModificationTime: mod}
}

func parseDate(path string) time.Time {
	file := filepath.Base(path)
	var timeZero time.Time
	matches := dateRe.FindStringSubmatch(file)

	if len(matches) < 4 {
		return timeZero
	}
	yr, err := strconv.Atoi(matches[1])
	if yr > now.Year() || err != nil {
		return timeZero
	}

	m, err := strconv.Atoi(matches[2])
	if m > 12 || m < 1 || err != nil {
		return timeZero
	}

	d, err := strconv.Atoi(matches[3])
	if d > 31 || d < 1 || err != nil {
		return timeZero
	}

	timeString := fmt.Sprintf("%d-%02d-%02d", yr, m, d)
	t, err := time.ParseInLocation(shortForm, timeString, time.Local)

	if err != nil {
		return timeZero
	}

	return t
}

func initFilesMap(dir *InputDir) map[string][]File {
	res := make(map[string][]File)
	for _, v := range dir.Files {
		var files []File
		res[v] = files
	}
	return res
}

func isBackedUpFile(path string, files []string) (string, bool, error) {
	for _, v := range files {
		re, err := regexp.Compile(v)
		if err != nil {
			return v, false, err
		}

		if re.MatchString(path) {
			return "", true, nil
		}
	}
	return "", false, nil
}
