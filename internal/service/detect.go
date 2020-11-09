package service

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/iantal/btd/internal/config"
	"github.com/iantal/btd/internal/files"
	"github.com/iantal/btd/internal/util"
	"github.com/sirupsen/logrus"
)

// Detector defines a service for detecting build tool types
type Detector struct {
	log      *util.StandardLogger
	basePath string
	store    files.Storage
	bTypes   config.Types
	rmHost   string
}

// NewDetector creates a Detector
func NewDetector(log *util.StandardLogger, basePath, rmHost string, store files.Storage) *Detector {
	conf, err := config.LoadConfig("/go/config.yml")
	if err != nil {
		log.WithField("error", err).Error("Could not load config file")
	}

	return &Detector{log, basePath, store, conf, rmHost}
}

// Detect is downloading, extracting and analyzing a project
func (d *Detector) Detect(projectID, commit string) ([]string, error) {
	projectPath := filepath.Join(d.store.FullPath(projectID), "bundle")

	if _, err := os.Stat(projectPath); os.IsNotExist(err) {
		err := d.downloadRepository(projectID, commit)
		if err != nil {
			d.log.WithFields(logrus.Fields{
				"projectID": projectID,
				"commit":    commit,
				"error":     err,
			}).Error("Could not download bundled repository")
			return []string{}, fmt.Errorf("Could not download bundled repository for %s", projectID)
		}
	}

	bp := commit + ".bundle"
	srcPath := d.store.FullPath(filepath.Join(projectID, commit, "bundle", bp))
	destPath := d.store.FullPath(filepath.Join(projectID, commit, "unbundle"))

	if _, err := os.Stat(destPath); os.IsNotExist(err) {
		err := d.store.Unbundle(srcPath, destPath)
		if err != nil {
			d.log.WithFields(logrus.Fields{
				"projectID": projectID,
				"commit":    commit,
				"error":     err,
			}).Error("Could not unbundle repository")
			return []string{}, fmt.Errorf("Could not unbundle repository for %s", projectID)
		}
	}

	r := []string{}
	dir := filepath.Join(destPath, commit)
	for _, t := range d.bTypes.BuildTypes {
		for _, file := range t.Files {
			if d.search(file, dir) {
				r = append(r, t.Name)
			}
		}
	}

	return r, nil
}

func (d *Detector) downloadRepository(projectID, commit string) error {
	resp, err := http.DefaultClient.Get("http://" + d.rmHost + "/api/v1/projects/" + projectID + "/" + commit + "/download")
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("Expected error code 200 got %d", resp.StatusCode)
	}

	d.log.WithField("file", resp.Header.Get("Content-Disposition")).Info("Content-Dispozition")

	d.save(projectID, commit, resp.Body)
	resp.Body.Close()

	return nil
}

func (d *Detector) save(projectID, commit string, r io.ReadCloser) {
	d.log.WithFields(logrus.Fields{
		"projectID": projectID,
		"commit":    commit,
	}).Info("Save project to storage")

	bp := commit + ".bundle"
	fp := filepath.Join(projectID, commit, "bundle", bp)
	err := d.store.Save(fp, r)

	if err != nil {
		d.log.WithFields(logrus.Fields{
			"projectID": projectID,
			"commit":    commit,
			"error":     err,
		}).Error("Unable to save file")
		return
	}
}

func (d *Detector) search(file, directory string) bool {
	found := false

	err := filepath.Walk(directory,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			if strings.Contains(path, file) {
				d.log.WithField("name", file).Info("Found file")
				found = true
			}

			return nil
		})
	if err != nil {
		log.Println(err)
	}

	return found
}
