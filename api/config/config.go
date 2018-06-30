package config

import (
	"bufio"
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"golang.org/x/oauth2"
)

var config *Config

// Config represents a configuration.
type Config struct {
	Project     string
	Instance    string
	Creds       string
	TokenSource oauth2.TokenSource
}

// registerFlags registers a set of standard flags for this config.
func (c *Config) registerFlags() {
	flag.StringVar(&c.Project, "project", c.Project, "project ID, if unset uses gcloud configured project")
	flag.StringVar(&c.Instance, "instance", c.Instance, "Cloud Bigtable instance")
	flag.StringVar(&c.Creds, "creds", c.Creds, "if set, use application credentials in this file")
}

// Load returns initialized configuration
func Load() (*Config, error) {
	filename := filepath.Join(os.Getenv("HOME"), ".cbtrc")
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		// silent fail if the file isn't there
		if os.IsNotExist(err) {
			return &Config{}, nil
		}
		return nil, fmt.Errorf("Reading %s: %v", filename, err)
	}
	s := bufio.NewScanner(bytes.NewReader(data))
	for s.Scan() {
		line := s.Text()
		i := strings.Index(line, "=")
		if i < 0 {
			return nil, fmt.Errorf("Bad line in %s: %q", filename, line)
		}
		key, val := strings.TrimSpace(line[:i]), strings.TrimSpace(line[i+1:])
		switch key {
		default:
			return nil, fmt.Errorf("Unknown key in %s: %q", filename, key)
		case "project":
			config.Project = val
		case "instance":
			config.Instance = val
		case "creds":
			config.Creds = val
		}
	}
	return config, s.Err()
}

type gcloudCredential struct {
	AccessToken string    `json:"access_token"`
	Expiry      time.Time `json:"token_expiry"`
}

func (cred *gcloudCredential) Token() *oauth2.Token {
	return &oauth2.Token{AccessToken: cred.AccessToken, TokenType: "Bearer", Expiry: cred.Expiry}
}

// GcloudConfig configuration fot the gcloud
type GcloudConfig struct {
	Configuration struct {
		Properties struct {
			Core struct {
				Project string `json:"project"`
			} `json:"core"`
		} `json:"properties"`
	} `json:"configuration"`
	Credential gcloudCredential `json:"credential"`
}

// GcloudCmdTokenSource represents gcloud command that returns a token source
type GcloudCmdTokenSource struct {
	Command string
	Args    []string
}

// Token implements the oauth2.TokenSource interface
func (g *GcloudCmdTokenSource) Token() (*oauth2.Token, error) {
	gcloudConfig, err := loadGcloudConfig(g.Command, g.Args)
	if err != nil {
		return nil, err
	}
	return gcloudConfig.Credential.Token(), nil
}

func loadGcloudConfig(gcloudCmd string, gcloudCmdArgs []string) (*GcloudConfig, error) {
	out, err := exec.Command(gcloudCmd, gcloudCmdArgs...).Output()
	if err != nil {
		return nil, fmt.Errorf("Could not retrieve gcloud configuration")
	}

	var gcloudConfig GcloudConfig
	if err := json.Unmarshal(out, &gcloudConfig); err != nil {
		return nil, fmt.Errorf("Could not parse gcloud configuration")
	}

	return &gcloudConfig, nil
}

func (c *Config) setFromGcloud() error {
	if c.Creds == "" {
		c.Creds = os.Getenv("GOOGLE_APPLICATION_CREDENTIALS")
		if c.Creds == "" {
			log.Printf("-creds flag unset, will use gcloud credential")
		}
	} else {
		os.Setenv("GOOGLE_APPLICATION_CREDENTIALS", c.Creds)
	}

	if c.Project == "" {
		log.Printf("-project flag unset, will use gcloud active project")
	}

	if c.Creds != "" && c.Project != "" {
		return nil
	}

	gcloudCmd := "gcloud"
	if runtime.GOOS == "windows" {
		gcloudCmd = gcloudCmd + ".cmd"
	}

	gcloudCmdArgs := []string{"config", "config-helper",
		"--format=json(configuration.properties.core.project,credential)"}

	gcloudConfig, err := loadGcloudConfig(gcloudCmd, gcloudCmdArgs)
	if err != nil {
		return err
	}

	if c.Project == "" && gcloudConfig.Configuration.Properties.Core.Project != "" {
		log.Printf("gcloud active project is \"%s\"",
			gcloudConfig.Configuration.Properties.Core.Project)
		c.Project = gcloudConfig.Configuration.Properties.Core.Project
	}

	if c.Creds == "" {
		c.TokenSource = oauth2.ReuseTokenSource(
			gcloudConfig.Credential.Token(),
			&GcloudCmdTokenSource{Command: gcloudCmd, Args: gcloudCmdArgs})
	}

	return nil
}
