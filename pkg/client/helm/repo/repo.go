package repo

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gofrs/flock"
	"gopkg.in/yaml.v2"
	"helm.sh/helm/v3/pkg/cli"
	"helm.sh/helm/v3/pkg/getter"
	"helm.sh/helm/v3/pkg/repo"
)

var deprecatedRepos = map[string]string{
	"//kubernetes-charts.storage.googleapis.com":           "https://charts.helm.sh/stable",
	"//kubernetes-charts-incubator.storage.googleapis.com": "https://charts.helm.sh/incubator",
}

func checkDeprecatedRepos(url string) error {
	for oldURL, newURL := range deprecatedRepos {
		if strings.Contains(url, oldURL) {
			return fmt.Errorf("repo(%s) is deprecated. use %s", url, newURL)
		}
	}

	return nil
}

type Repos struct {
	settings   *cli.EnvSettings
	fileLock   *flock.Flock
	fileLocked bool
	fileInfo   *repo.File
}

func NewRepos(settings *cli.EnvSettings) *Repos {
	return &Repos{settings: settings}
}

func (r *Repos) lockFile(ctx context.Context) error {
	// check repo directory
	err := os.MkdirAll(filepath.Dir(r.settings.RepositoryConfig), os.ModePerm)
	if err != nil && !os.IsExist(err) {
		return err
	}

	// lock repo lock file
	fileExt := filepath.Ext(r.settings.RepositoryConfig)
	var lockPath string
	if len(fileExt) > 0 && len(fileExt) < len(r.settings.RepositoryConfig) {
		lockPath = strings.Replace(r.settings.RepositoryConfig, fileExt, ".lock", 1)
	} else {
		lockPath = r.settings.RepositoryConfig + ".lock"
	}

	r.fileLock = flock.New(lockPath)

	r.fileLocked, err = r.fileLock.TryLockContext(ctx, time.Second)
	if err != nil {
		return err
	}

	return nil
}

func (r *Repos) unlockFile() {
	if r.fileLock != nil && r.fileLocked {
		r.fileLock.Unlock()
	}
}

func (r *Repos) readFile() error {
	if len(r.settings.RepositoryConfig) == 0 {
		return fmt.Errorf("repository file path is emtpy")
	}

	b, err := ioutil.ReadFile(r.settings.RepositoryConfig)
	if err != nil && !os.IsNotExist(err) {
		return err
	}

	var f repo.File
	if err := yaml.Unmarshal(b, &f); err != nil {
		return err
	}
	r.fileInfo = &f

	return nil
}

func (r *Repos) hasRepoEntry(entry *repo.Entry) error {
	if r.fileInfo == nil || entry == nil {
		return fmt.Errorf("file info or entry is empty")
	}

	// repo name already exists
	if r.fileInfo.Has(entry.Name) {
		existingEntry := r.fileInfo.Get(entry.Name)
		if *entry != *existingEntry {
			return fmt.Errorf("repository name(%s) already exists", entry.Name)
		}

		// the same repo already exists
		return nil
	}

	return nil
}

func (r *Repos) updateRepoEntry(entry *repo.Entry) error {
	if r.fileInfo == nil {
		return fmt.Errorf("file info is empty")
	}

	r.fileInfo.Update(entry)

	if err := r.fileInfo.WriteFile(r.settings.RepositoryConfig, 0644); err != nil {
		return err
	}

	return nil
}

func (r *Repos) AddRepo(name, url string) error {
	if err := checkDeprecatedRepos(url); err != nil {
		return err
	}

	// lock for repo file
	lockCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	err := r.lockFile(lockCtx)
	defer r.unlockFile()
	if err != nil {
		return err
	}

	// read repo file
	if err := r.readFile(); err != nil {
		return err
	}

	entry := &repo.Entry{
		Name: name,
		URL:  url,
	}

	// check if an entry exists in the repo file
	if err := r.hasRepoEntry(entry); err != nil {
		return err
	}

	cr, err := repo.NewChartRepository(entry, getter.All(r.settings))
	if err != nil {
		return err
	}

	if r.settings.RepositoryCache != "" {
		cr.CachePath = r.settings.RepositoryCache
	}
	if _, err := cr.DownloadIndexFile(); err != nil {
		return err
	}

	// update and write entry to repo file
	if err := r.updateRepoEntry(entry); err != nil {
		return err
	}

	return nil
}
