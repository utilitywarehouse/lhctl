package cmd

import (
	"bytes"
	"errors"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestReadConfig(t *testing.T) {

	viper.SetConfigType("yaml")

	// Test: read empty config doesn't error
	yamlConfig := []byte(``)
	viper.ReadConfig(bytes.NewBuffer(yamlConfig))

	c, err := readConfig()

	assert.Equal(t, nil, err)
	assert.Equal(t, Config{}, c)

	// Test: Read config
	yamlConfig = []byte(`
contexts:
- name: my-cluster
  url: https://longhorn.test
  user: user
  pass: pass
default: my-cluster
`)

	viper.ReadConfig(bytes.NewBuffer(yamlConfig))

	c, err = readConfig()
	assert.Equal(t, nil, err)
	assert.Equal(t, c.DefaultContext, "my-cluster")
	assert.Equal(t, len(c.Contexts), 1)
	assert.Equal(t, c.Contexts[0].Name, "my-cluster")
	assert.Equal(t, c.Contexts[0].Url, "https://longhorn.test")
	assert.Equal(t, c.Contexts[0].User, "user")
	assert.Equal(t, c.Contexts[0].Pass, "pass")

	// Test: Read config no default
	yamlConfig = []byte(`
contexts:
- name: my-cluster
  url: https://longhorn.test
  user: user
  pass: pass
`)

	viper.ReadConfig(bytes.NewBuffer(yamlConfig))

	c, err = readConfig()
	assert.Equal(t, nil, err)
	assert.Equal(t, c.DefaultContext, "")

	// Test: Read config multiple contexts
	yamlConfig = []byte(`
contexts:
- name: context-1
  url: https://longhorn.test.1
  user: user-1
  pass: pass-1
- name: context-2
  url: https://longhorn.test.2
  user: user-2
  pass: pass-2
`)

	viper.ReadConfig(bytes.NewBuffer(yamlConfig))

	c, err = readConfig()
	assert.Equal(t, nil, err)
	assert.Equal(t, len(c.Contexts), 2)
}

func TestInitGetClientParams(t *testing.T) {

	viper.SetConfigType("yaml")

	// Test: Test empty conmfig and no url returns error
	yamlConfig := []byte(``)

	viper.ReadConfig(bytes.NewBuffer(yamlConfig))

	_, _, _, err := getClientParams()

	expectedErr := errors.New(
		"You need to provide a url via config or using `--url=` flag",
	)

	assert.Equal(t, expectedErr, err)

	// Test: faulty default context results in error for url
	yamlConfig = []byte(`
contexts:
- name: context
  url: https://longhorn.test
  user: user
  pass: pass
default: default
`)

	viper.ReadConfig(bytes.NewBuffer(yamlConfig))

	_, _, _, err = getClientParams()

	assert.Equal(t, expectedErr, err)

	// Test: default context
	yamlConfig = []byte(`
contexts:
- name: context
  url: https://longhorn.test
  user: user
  pass: pass
default: context
`)

	viper.ReadConfig(bytes.NewBuffer(yamlConfig))

	url, user, pass, err := getClientParams()

	assert.Equal(t, nil, err)
	assert.Equal(t, "https://longhorn.test", url)
	assert.Equal(t, "user", user)
	assert.Equal(t, "pass", pass)

	// Test: context flag supersedes default context
	yamlConfig = []byte(`
contexts:
- name: context
  url: https://longhorn.test
  user: user
  pass: pass
default: default
`)

	viper.ReadConfig(bytes.NewBuffer(yamlConfig))

	contextFlag = "context"

	url, user, pass, err = getClientParams()

	assert.Equal(t, nil, err)
	assert.Equal(t, "https://longhorn.test", url)
	assert.Equal(t, "user", user)
	assert.Equal(t, "pass", pass)

	// Test: Flags override context values
	yamlConfig = []byte(`
contexts:
- name: context
  url: https://longhorn.test
  user: user
  pass: pass
default: context
`)

	viper.ReadConfig(bytes.NewBuffer(yamlConfig))

	// set flags
	urlFlag = "urlFlag"
	userFlag = "userFlag"
	passFlag = "passFlag"

	url, user, pass, err = getClientParams()

	assert.Equal(t, nil, err)
	assert.Equal(t, "urlFlag", url)
	assert.Equal(t, "userFlag", user)
	assert.Equal(t, "passFlag", pass)
}
