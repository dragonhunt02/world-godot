package main

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestConfigFile(t *testing.T) {
	cfgFileContent := []byte(`{"store-options": {"/path/to/store/":{"uncompressed": true}}}`)
	f, err := ioutil.TempFile("", "")
	require.NoError(t, err)
	f.Close()
	defer os.Remove(f.Name())
	require.NoError(t, ioutil.WriteFile(f.Name(), cfgFileContent, 0644))

	// Set the global config file name
	cfgFile = f.Name()

	// Call init, this should use the custom config file and global "cfg" should contain the
	// values
	initConfig()

	// If everything worked, the options should be set according to the config file created above
	opt, err := cfg.GetStoreOptionsFor("/path/to/store")
	require.NoError(t, err)
	require.True(t, opt.Uncompressed)

	// The options for a non-matching store should be default
	opt, err = cfg.GetStoreOptionsFor("/path/other-store")
	require.NoError(t, err)
	require.False(t, opt.Uncompressed)
}

func TestConfigFileMultipleMatches(t *testing.T) {
	cfgFileContent := []byte(`{"store-options": {"/path/to/store/":{"uncompressed": true}, "/path/to/store":{"uncompressed": false}}}`)
	f, err := ioutil.TempFile("", "")
	require.NoError(t, err)
	f.Close()
	defer os.Remove(f.Name())
	require.NoError(t, ioutil.WriteFile(f.Name(), cfgFileContent, 0644))

	// Set the global config file name
	cfgFile = f.Name()

	// Call init, this should use the custom config file and global "cfg" should contain the
	// values
	initConfig()

	// We expect this to fail because both "/path/to/store/" and "/path/to/store" matches the
	// provided location
	_, err = cfg.GetStoreOptionsFor("/path/to/store")
	require.Error(t, err)
}
