package backme

import (
	"log"
	"sort"
	"time"
)

func categorizeFiles(filesRegexMap map[string][]File) *CategorizedFiles {
	bins := &CategorizedFiles{Years: make(map[int][]File)}
	allEntries := 0
	for k, v := range filesRegexMap {
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
