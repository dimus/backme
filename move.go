package backme

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
)

func moveFiles(dir *InputDir, bins *CategorizedFiles, conf *Config) error {
	archivePath := filepath.Join(dir.Path, conf.OutputDir)
	err := checkCreateDir(archivePath)
	if err != nil {
		return err
	}

	if !dir.KeepAllFiles {
		for _, v := range bins.GarbageFiles {
			err := os.Remove(v.Path)
			if err != nil {
				return err
			}
		}
	} else {
		deletePath := filepath.Join(archivePath, "delete-me")
		err = moveFilesShared(bins.GarbageFiles, deletePath)
		if err != nil {
			return err
		}
	}

	recentPath := filepath.Join(archivePath, "recent")
	err = moveFilesShared(bins.Recent, recentPath)
	if err != nil {
		return err
	}

	monthPath := filepath.Join(archivePath, "last-month")
	err = moveFilesShared(bins.LastMonth, monthPath)
	if err != nil {
		return err
	}

	for k, v := range bins.Years {
		yearPath := filepath.Join(archivePath, strconv.Itoa(k))
		err = moveFilesShared(v, yearPath)
		if err != nil {
			return err
		}
	}

	return nil
}

func moveFilesShared(files []File, path string) error {
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
