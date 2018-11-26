package backme

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
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
	files, err := collectFiles(&dir)
	if err != nil {
		return err
	}
	log.Println(LogSep)
	bins := processFiles(files)
	err = moveFiles(&dir, bins, conf)
	return err
}

func moveFiles(dir *InputDir, bins *CategorizedFiles, conf *Config) error {
	archivePath := filepath.Join(dir.Path, conf.OutputDir)
	err := checkCreateDir(archivePath)
	if err != nil {
		return err
	}

	for _, v := range bins.GarbageFiles {
		err := os.Remove(v.Path)
		if err != nil {
			return err
		}
	}

	recentPath := filepath.Join(archivePath, "recent")
	err = moveEntries(bins.Recent, recentPath)
	if err != nil {
		return err
	}

	monthPath := filepath.Join(archivePath, "last-month")
	err = moveEntries(bins.LastMonth, monthPath)
	if err != nil {
		return err
	}

	for k, v := range bins.Years {
		yearPath := filepath.Join(archivePath, strconv.Itoa(k))
		err = moveEntries(v, yearPath)
		if err != nil {
			return err
		}
	}

	return nil
}

func moveEntries(files []File, path string) error {
	err := checkCreateDir(path)
	if err != nil {
		return err
	}

	for _, v := range files {
		newPath := filepath.Join(path, filepath.Base(v.Path))
		if v.Path == newPath {
			continue
		}
		err := os.Rename(v.Path, newPath)
		if err != nil {
			return err
		}
	}
	return nil
}

func checkCreateDir(path string) error {
	stat, err := os.Stat(path)
	if os.IsNotExist(err) {
		log.Printf("Creating %s", path)
		err = os.Mkdir(path, os.ModePerm|os.ModeDir)
		if err != nil {
			return err
		}
	}

	if stat != nil && !stat.IsDir() {
		return fmt.Errorf("The %s exists, but is not a directory\n", path)
	}

	return nil
}

func processFiles(files map[string][]File) *CategorizedFiles {
	bins := &CategorizedFiles{Years: make(map[int][]File)}
	allEntries := 0
	for k, v := range files {
		allEntries += len(v)
		log.Printf("Processing files matching /%s/.\n", k)
		log.Printf("Found %d entries.", len(v))
		processFileGroup(v, bins)
	}

	log.Printf("Total: %d entries, %d to delete, %d recent, %d from the last month, %d from all years",
		allEntries, len(bins.GarbageFiles), len(bins.Recent),
		len(bins.LastMonth), yearsEntries(allEntries, bins))
	log.Println(LogSep)
	return bins
}

func yearsEntries(allEntries int, bins *CategorizedFiles) int {
	return allEntries - len(bins.GarbageFiles) -
		len(bins.Recent) - len(bins.LastMonth)
}

func processFileGroup(files []File, bins *CategorizedFiles) {
	sort.Slice(files, func(i, j int) bool {
		return files[i].SortingTime.After(files[j].SortingTime)
	})

	twoDaysAgo := now.Add(-48 * time.Hour)
	monthAgo := now.Add(-31 * 24 * time.Hour)
	for _, v := range files {
		t := v.SortingTime
		switch {
		case t.After(twoDaysAgo):
			bins.Recent = append(bins.Recent, v)
		case t.After(monthAgo):
			categorizeLastMonth(v, bins)
		default:
			categorizeYear(v, bins)
		}
	}
}

func categorizeLastMonth(f File, bins *CategorizedFiles) {
	l := len(bins.LastMonth)
	if l > 0 && bins.LastMonth[l-1].SortingTime.Day() == f.SortingTime.Day() {
		bins.GarbageFiles = append(bins.GarbageFiles, f)
		return
	}

	bins.LastMonth = append(bins.LastMonth, f)
}

func categorizeYear(f File, bins *CategorizedFiles) {
	yr := f.SortingTime.Year()
	if _, ok := bins.Years[yr]; !ok {
		var files []File
		bins.Years[yr] = files
	}

	l := len(bins.Years[yr])
	if l > 0 && bins.Years[yr][l-1].SortingTime.Month() == f.SortingTime.Month() {
		bins.GarbageFiles = append(bins.GarbageFiles, f)
		return
	}
	bins.Years[yr] = append(bins.Years[yr], f)
}

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
