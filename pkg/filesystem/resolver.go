package filesystem

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/SierraSoftworks/git-tool/pkg/config"
	"github.com/SierraSoftworks/git-tool/pkg/repo"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

// GetRepo attempts to resolve the details of a repository given its name.
func GetRepo(cfg *config.Config, name string) (*config.Service, *repo.Repo, error) {
	if name == "" {
		return GetCurrentDirectoryRepo(cfg)
	}

	dirParts := strings.Split(filepath.ToSlash(name), "/")
	if len(dirParts) < 2 {
		logrus.WithField("path", name).Debug("Not a fully qualified repository name")
		return nil, nil, nil
	}

	serviceName := dirParts[0]
	service := cfg.GetService(serviceName)

	if service != nil {
		r, err := GetRepoForService(cfg, service, filepath.Join(dirParts[1:]...))
		return service, r, err
	}

	service, r, err := GetFullyQualifiedRepo(cfg, name)
	if err != nil {
		return service, r, err
	}

	if r == nil {
		r, err = GetRepoForService(cfg, cfg.GetDefaultService(), name)
		if r != nil {
			return cfg.GetDefaultService(), r, err
		}
	}

	logrus.WithField("path", name).Debug("Could not find a matching repository")
	return nil, nil, nil
}

// GetRepoForService fetches the repo details for the named repository managed by the
// provided service.
func GetRepoForService(cfg *config.Config, service *config.Service, name string) (*repo.Repo, error) {
	dirParts := strings.Split(filepath.ToSlash(name), "/")

	fullNameLength := len(strings.Split(service.NamingPattern, "/"))
	if len(dirParts) < fullNameLength {
		logrus.WithField("path", name).Debug("Not a fully named repository folder within the service's development directory")
		return nil, nil
	}

	return &repo.Repo{
		FullName: strings.Join(dirParts[:fullNameLength], "/"),
		Service:  service.Domain,
		Path:     filepath.Join(cfg.DevelopmentDirectory, service.Domain, filepath.Join(dirParts[:fullNameLength]...)),
	}, nil
}

// GetFullyQualifiedRepo fetches the repo details for the fully qualified named
// repository which has been provided.
func GetFullyQualifiedRepo(cfg *config.Config, name string) (*config.Service, *repo.Repo, error) {
	dirParts := strings.Split(filepath.ToSlash(name), "/")

	if len(dirParts) < 2 {
		// Not within a service's repository
		logrus.WithField("path", name).Debug("Not a repository folder within the development directory")
		return nil, nil, nil
	}

	serviceName := dirParts[0]
	service := cfg.GetService(serviceName)
	if service == nil {
		logrus.WithField("path", name).Debug("No service found to handle repository type")
		return nil, nil, nil
	}

	r, err := GetRepoForService(cfg, service, strings.Join(dirParts[1:], "/"))
	return service, r, err
}

// GetCurrentDirectoryRepo gets the repo details for the repository open in your
// current directory.
func GetCurrentDirectoryRepo(cfg *config.Config) (*config.Service, *repo.Repo, error) {
	dir, err := os.Getwd()
	if err != nil {
		return nil, nil, errors.Wrap(err, "filesystem: failed to get current directory")
	}

	if !strings.HasPrefix(dir, cfg.DevelopmentDirectory) {
		logrus.WithField("path", dir).Debug("Not within the development directory")
		return nil, nil, nil
	}

	localDir := strings.Trim(filepath.ToSlash(dir[len(cfg.DevelopmentDirectory):]), "/")
	return GetFullyQualifiedRepo(cfg, localDir)
}
