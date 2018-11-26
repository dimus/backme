package backme

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"time"
)

func collectFiles(dir *InputDir) (map[string][]File, error) {
	files := initFilesMap(dir)

	err := filepath.Walk(dir.Path,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			file := filepath.Base(path)
			re, match, err := isBackedUpFile(file, dir.FileRegexPatterns)
			if err != nil {
				return err
			}

			if match {
				f := getFile(path, info)
				if f.SortingTime == timeZero {
					f.SortingTime = f.ModificationTime
				}
				files[re] = append(files[re], f)
			}
			return nil
		})
	return files, err
}

func getFile(path string, info os.FileInfo) File {
	mod := info.ModTime()
	date := parseDate(path)
	return File{Path: path, SortingTime: date, ModificationTime: mod}
}

func parseDate(path string) time.Time {
	file := filepath.Base(path)
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
	filesMap := make(map[string][]File)
	for _, v := range dir.FileRegexPatterns {
		var files []File
		filesMap[v] = files
	}
	return filesMap
}

func isBackedUpFile(path string, files []string) (string, bool, error) {
	for _, v := range files {
		re, err := regexp.Compile(v)
		if err != nil {
			return v, false, err
		}

		if re.MatchString(path) {
			return v, true, nil
		}
	}
	return "", false, nil
}
