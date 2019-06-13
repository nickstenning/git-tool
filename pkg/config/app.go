package config

import (
	"bytes"
	"os"
	"os/exec"
	"text/template"

	"github.com/SierraSoftworks/git-tool/pkg/repo"
	"github.com/pkg/errors"
)

// An App is something which may be executed within the context of your
// repository.
type App struct {
	Name        string   `json:"name" yaml:"name"`
	Default     bool     `json:"default,omitempty" yaml:"default,omitempty"`
	CommandLine string   `json:"command" yaml:"command"`
	Arguments   []string `json:"args,omitempty" yaml:"args,omitempty"`
	Environment []string `json:"environment,omitempty" yaml:"environment,omitempty"`
}

// GetCmd will fetch the *exec.Cmd used to start this application within
// the context of a specific service and repository.
func (a *App) GetCmd(s *Service, r *repo.Repo) (*exec.Cmd, error) {
	args := make([]string, len(a.Arguments))

	for i, arg := range a.Arguments {
		at, err := a.getTemplate(arg, s, r)
		if err != nil {
			return nil, errors.Wrap(err, "config: failed to construct application command line")
		}

		args[i] = at
	}

	env := make([]string, len(a.Environment))

	for i, arg := range a.Environment {
		at, err := a.getTemplate(arg, s, r)
		if err != nil {
			return nil, errors.Wrap(err, "config: failed to construct application environment variables")
		}

		env[i] = at
	}

	cmd := exec.Command(a.CommandLine, args...)

	cmd.Dir = r.Path
	cmd.Env = append(os.Environ(), env...)

	return cmd, nil
}

func (a *App) getTemplate(tmpl string, s *Service, r *repo.Repo) (string, error) {
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
