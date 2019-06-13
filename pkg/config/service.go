package config

import (
	"bytes"
	"text/template"

	"github.com/SierraSoftworks/git-tool/pkg/repo"
)

// A Service is a descriptor used to describe a remote git host.
type Service struct {
	Domain string `json:"domain" yaml:"domain"`

	WebURLTemplate string `json:"webUrl" yaml:"webUrl"`
	GitURLTemplate string `json:"gitUrl" yaml:"gitUrl"`

	Default bool `json:"default,omitempty" yaml:"default,omitempty"`

	NamingPattern string `json:"pattern" yaml:"pattern"`
}

// WebURL fetches the HTTP(S) URL which may be used to view a web based
// representation of a repository.
func (s *Service) WebURL(r *repo.Repo) (string, error) {
	return s.getTemplate(s.WebURLTemplate, r)
}

// GitURL fetches the git+ssh URL which may be used to fetch or push
// the repository's code.
func (s *Service) GitURL(r *repo.Repo) (string, error) {
	return s.getTemplate(s.GitURLTemplate, r)
}

func (s *Service) getTemplate(tmpl string, r *repo.Repo) (string, error) {
	t, err := template.New("gitURL").Parse(tmpl)
	if err != nil {
		return "", err
	}

	buf := bytes.NewBuffer([]byte{})
	if err := t.Execute(buf, struct {
		Service *Service
		Repo    *repo.Repo
	}{
		Service: s,
		Repo:    r,
	}); err != nil {
		return "", err
	}

	return buf.String(), nil
}
