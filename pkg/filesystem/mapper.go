package filesystem

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/sirupsen/logrus"

	"github.com/SierraSoftworks/git-tool/pkg/config"

	"github.com/SierraSoftworks/git-tool/pkg/repo"
	"github.com/pkg/errors"
)

// A Mapper holds the information about a developer's source code folder which
// may contain multiple repositories.
type Mapper struct {
	Config *config.Config
}

// GetRepos will fetch all of the repositories contained within a developer's dev
// directory which match the required naming scheme.
func (d *Mapper) GetRepos() ([]*repo.Repo, error) {
	files, err := ioutil.ReadDir(d.Config.DevelopmentDirectory)
	if err != nil {
		return nil, errors.Wrapf(err, "filesystem: unable to list directory contents in %s", d.Config.DevelopmentDirectory)
	}

	repos := []*repo.Repo{}

	for _, f := range files {
		if !f.IsDir() {
			continue
		}

		service := d.Config.GetService(f.Name())
		if service == nil {
			logrus.WithField("service", f.Name()).Warn("Could not find a matching service entry in your configuration")
			continue
		}

		childRepos, err := d.getRepos(service)
		if err != nil {
			return nil, errors.Wrapf(err, "filesystem: unable to list directory contents in %s", d.Config.DevelopmentDirectory)
		}

		repos = append(repos, childRepos...)
	}

	return repos, nil
}

// EnsureRepo will ensure that a repository directory has been created at the correct location
// on the filesystem.
func (d *Mapper) EnsureRepo(service *config.Service, r *repo.Repo) error {
	path := filepath.Join(d.Config.DevelopmentDirectory, service.Domain, filepath.FromSlash(r.FullName))

	s, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			if err := os.MkdirAll(path, os.ModePerm); err != nil {
				return errors.Wrapf(err, "filesystem: unable to create repository directory: %s", path)
			}
			return nil
		}

		return errors.Wrapf(err, "filesystem: unable to check directory: %s", path)
	}

	if s.IsDir() {
		return nil
	}

	return errors.Errorf("filesystem: repository name already exists and is not a directory: %s", path)
}

func (d *Mapper) getRepos(service *config.Service) ([]*repo.Repo, error) {
	logrus.WithField("service", service.Domain).Debug("Enumerating repositories for service")

	path := filepath.Join(d.Config.DevelopmentDirectory, service.Domain)

	pattern := filepath.Join(path, service.NamingPattern)

	files, err := filepath.Glob(pattern)
	if err != nil {
		return nil, errors.Wrapf(err, "repo: unable to list directory contents in %s", pattern)
	}

	repos := []*repo.Repo{}
	for _, f := range files {
		logrus.WithField("service", service.Domain).WithField("path", f).Debug("Enumerated possible repository")
		repo := &repo.Repo{
			Service:  service.Domain,
			FullName: strings.Trim(strings.Replace(f[len(path):], string(filepath.Separator), "/", -1), "/"),
			Path:     f,
		}

		if repo.IsValid() {
			repos = append(repos, repo)
		} else {
			logrus.WithField("service", service.Domain).WithField("path", f).Debug("Marked repository as invalid")
		}
	}

	return repos, nil
}
