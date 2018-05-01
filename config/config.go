package config

import (
	"bytes"
	"io/ioutil"
	"os"
	"os/user"
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/Songmu/prompter"
	"github.com/pkg/errors"

	"github.com/potsbo/jobcan/types"
)

// Config is command parameters
type Config struct {
	Credential CredentialConfig
}

// CredentialConfig is jobcan credential
type CredentialConfig struct {
	ClientID    string
	LoginID     string
	Password    string
	AccountType string
}

func (c *Config) valid() bool {
	return c.Credential.valid()
}

func (c *CredentialConfig) valid() bool {
	if len(c.ClientID) == 0 {
		return false
	}
	if len(c.LoginID) == 0 {
		return false
	}
	if len(c.Password) == 0 {
		return false
	}
	if len(c.AccountType) == 0 {
		return false
	}
	return true
}

func configPath() string {
	// only OSX
	usr, _ := user.Current()
	return strings.Replace("~/.jobcan", "~", usr.HomeDir, 1)
}

func Init() {
	ac := []string{types.Admin, types.Staff}

	var config Config
	var credentialConfig CredentialConfig
	credentialConfig.AccountType = prompter.Choose("Choose your account type", ac, types.Staff)
	credentialConfig.ClientID = prompter.Prompt("Enter your client ID", "")
	credentialConfig.LoginID = prompter.Prompt("Enter your login ID", "")
	credentialConfig.Password = prompter.Prompt("Enter your password", "")
	config.Credential = credentialConfig

	var buffer bytes.Buffer
	encoder := toml.NewEncoder(&buffer)
	_ = encoder.Encode(config)

	ioutil.WriteFile(configPath(), []byte(buffer.String()), os.ModePerm)
}

func Read() (Config, error) {
	var config Config
	_, err := toml.DecodeFile(configPath(), &config)
	if err != nil {
		return config, errors.New("Config file is broken ;; please try `jobcan init`.")
	}
	return config, nil
}

// Load loads jobcan config from env vars or stored file
func Load() (Config, error) {
	c := Config{
		Credential: CredentialConfig{
			ClientID:    os.Getenv("JOBCAN_CLIENT_ID"),
			LoginID:     os.Getenv("JOBCAN_LOGIN_ID"),
			Password:    os.Getenv("JOBCAN_PASSWORD"),
			AccountType: os.Getenv("JOBCAN_ACCOUNT_TYPE"),
		},
	}

	if c.valid() {
		return c, nil
	}

	c, err := Read()
	if err != nil {
		return nil, errors.Wrap("failed to read stored credential")
	}
	return nil, errors.Errorf("failed to load config from envs nor config file")
}
