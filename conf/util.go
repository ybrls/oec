package conf

import (
	"encoding/json"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	fpath "path/filepath"
	"runtime"
	"strings"
)

const unknownFileExtErrMessage = "Unknown configuration file extension[%s]. Only \".json\" and \".yml(.yaml)\" types are allowed."

func checkFileExtension(filepath string) error {

	extension := fpath.Ext(strings.ToLower(filepath))

	switch extension {
	case ".json", ".yml", ".yaml":
		return nil
	default:
		return errors.Errorf(unknownFileExtErrMessage, extension)
	}
}

func readConfigurationFromFile(filepath string) (*Configuration, error) {

	file, err := ioutil.ReadFile(filepath)
	if err != nil {
		return nil, err
	}

	configuration := &Configuration{}
	extension := fpath.Ext(strings.ToLower(filepath))

	switch extension {
	case ".json":
		return configuration, json.Unmarshal(file, configuration)
	case ".yml", ".yaml":
		return configuration, yaml.Unmarshal(file, configuration)
	default:
		return nil, errors.Errorf(unknownFileExtErrMessage, extension)
	}
}

func addHomeDirPrefix(filepath string) string {
	tildePrefix := "~" + string(os.PathSeparator)

	if strings.HasPrefix(filepath, tildePrefix) {
		return fpath.Join(homeDir(), strings.TrimPrefix(filepath, tildePrefix))
	}

	return fpath.Clean(filepath)
}

func homeDir() string {
	if runtime.GOOS == "windows" {
		return os.Getenv("USERPROFILE")
	}
	return os.Getenv("HOME")
}

func addHomeDirPrefixToLocalActionFilepaths(mappings *ActionMappings) {
	for index, action := range *mappings {
		if action.SourceType == LocalSourceType {
			action.Filepath = addHomeDirPrefix(action.Filepath)
			(*mappings)[index] = action
		}
	}
}

func chmodLocalActions(mappings *ActionMappings, mode os.FileMode) {
	for _, action := range *mappings {
		if action.SourceType == LocalSourceType {
			err := os.Chmod(action.Filepath, mode)
			if err != nil {
				logrus.Warn(err)
			}
		}
	}
}

func addHomeDirPrefixToPrivateKeyFilepaths(mappings *ActionMappings) {
	for index, action := range *mappings {
		if action.SourceType == GitSourceType {
			if action.GitOptions.PrivateKeyFilepath == "" {
				continue
			}
			action.GitOptions.PrivateKeyFilepath = addHomeDirPrefix(action.GitOptions.PrivateKeyFilepath)
			(*mappings)[index] = action

		}
	}
}

func copyActionMappings(mappings ActionMappings) ActionMappings {

	copyActionMappings := make(map[ActionName]MappedAction, len(mappings))
	for k, v := range mappings {
		copyActionMappings[k] = v
	}
	return copyActionMappings
}
