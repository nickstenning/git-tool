package templates

import (
	"github.com/SierraSoftworks/git-tool/pkg/config"
	"github.com/SierraSoftworks/git-tool/pkg/repo"
)

type repoContext struct {
	Repo    *repo.Repo
	Service *config.Service
}

var repoTemplates = buildTemplates(map[string]string{
	"repo.qualified": `{{ .Service.Domain}}/{{ .Repo.FullName }}`,
	"repo.short":     `{{ .Service.Domain}}/{{ .Repo.FullName }} ({{ .Service.WebURL .Repo }})`,
	"repo.full": `
Name:       {{ .Repo.Name }}
Namespace:  {{ .Repo.Namespace }}
Service:    {{ .Service.Domain }}
Path:       {{ .Repo.Path }}

URLs:
 - Website: {{ .Service.WebURL .Repo }}
 - Git SSH: {{ .Service.GitURL .Repo }}
`,
})

// RepoQualifiedName gets a template which will format the fully qualified name of a repo
func RepoQualifiedName(r *repo.Repo, s *config.Service) string {
	return toString(repoTemplates, "repo.qualified", &repoContext{
		Repo:    r,
		Service: s,
	})
}

// RepoShortInfo gets a template which renders a detailed summary of a repository's details
func RepoShortInfo(r *repo.Repo, s *config.Service) string {
	return toString(repoTemplates, "repo.short", &repoContext{
		Repo:    r,
		Service: s,
	})
}

// RepoFullInfo gets a template which renders a detailed summary of a repository's details
func RepoFullInfo(r *repo.Repo, s *config.Service) string {
	return toString(repoTemplates, "repo.full", &repoContext{
		Repo:    r,
		Service: s,
	})
}
