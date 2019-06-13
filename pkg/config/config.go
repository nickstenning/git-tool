package config

import (
	"io/ioutil"
	"os"

	"github.com/go-yaml/yaml"
	"github.com/pkg/errors"
)

// The Config is used to configure the behavior of Git Tool
type Config struct {
	DevelopmentDirectory string `json:"directory" yaml:"directory"`

	Services     []*Service `json:"services" yaml:"services"`
	Applications []*App     `json:"apps" yaml:"apps"`
}

// Default gets a simple default configuration for Git Tool
// for environments where you have not defined a configuration
// file.
func Default() *Config {
	return &Config{
		DevelopmentDirectory: os.Getenv("DEV_DIRECTORY"),
		Services: []*Service{
			&Service{
				Domain:         "github.com",
				NamingPattern:  "*/*",
				Default:        true,
				WebURLTemplate: "https://{ .Service.Domain }/{ .Repo.FullName }",
				GitURLTemplate: "git@{ .Service.Domain }:{ .Repo.FullName }.git",
			},
		},
		Applications: []*App{
			&App{
				Name:        "shell",
				CommandLine: "bash",
			},
		},
	}
}

// Load will attempt to load a configuration object from the provided file.
func Load(file string) (*Config, error) {
	bytes, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, errors.Wrap(err, "config: unable to read config file")
	}

	config := &Config{}
	if err := yaml.Unmarshal(bytes, config); err != nil {
		return nil, errors.Wrap(err, "config: unable to parse config file")
	}

	return config, nil
}

// GetDefaultService fetches the default configured service
func (c *Config) GetDefaultService() *Service {
	for _, s := range c.Services {
		if s.Default {
			return s
		}
	}

	if len(c.Services) > 0 {
		return c.Services[0]
	}

	return nil
}

// GetService will retrieve a known service entry for the given service,
// if one exists in the config file, based on its domain name.
func (c *Config) GetService(domain string) *Service {
	if domain == "" {
		return c.GetDefaultService()
	}

	for _, s := range c.Services {
		if s.Domain == domain {
			return s
		}
	}

	return nil
}

// GetDefaultApp fetches the default configured application
func (c *Config) GetDefaultApp() *App {
	for _, a := range c.Applications {
		if a.Default {
			return a
		}
	}

	if len(c.Applications) > 0 {
		return c.Applications[0]
	}

	return nil
}

// GetApp fetches the app whose name matches the provided name.
func (c *Config) GetApp(name string) *App {
	if name == "" {
		return c.GetDefaultApp()
	}

	for _, a := range c.Applications {
		if a.Name == name {
			return a
		}
	}

	return nil
}
