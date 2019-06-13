package repo

import (
	"os"
	"path/filepath"
	"strings"
)

// Repo describes a repository object with any necessary properties required by
// Git-Tool.
type Repo struct {
	FullName string `json:"fullname"`
	Service  string `json:"service"`
	Path     string `json:"path"`
}

// Exists checks whether a repository entry is present on the local filesystem
// at the expected path.
func (r *Repo) Exists() bool {
	s, err := os.Stat(r.Path)
	if err != nil && os.IsNotExist(err) {
		return false
	}

	return s.IsDir()
}

// IsValid checks whether the current repo is initialized correctly.
func (r *Repo) IsValid() bool {
	s, err := os.Stat(filepath.Join(r.Path, ".git"))
	if err != nil && os.IsNotExist(err) {
		return false
	}

	return s.IsDir()
}

// Namespace gets the portion of the repository's full name prior to its final short name segment.
func (r *Repo) Namespace() string {
	parts := strings.Split(r.FullName, "/")
	return strings.Join(parts[:len(parts)-1], "/")
}

// Name gets the short name component of the repository's full name.
func (r *Repo) Name() string {
	parts := strings.Split(r.FullName, "/")
	return parts[len(parts)-1]
}
