package manager

import (
	"encoding/json"
	"os"
	"path/filepath"
)

const settingsDirName = "a18mm"
const settingsFileName = "settings.json"

const appData = "appdata"

type Configuration struct {
	InstallLocation string `json:"installLocation"`
}

var settingsDirPath = filepath.Join(os.Getenv(appData), settingsDirName)
var settingsFilePath = filepath.Join(settingsDirPath, settingsFileName)

func (c *Configuration) save() error {
	settingsByte, err := json.MarshalIndent(c, "", "    ")
	if err != nil {
		return err
	}

	return os.WriteFile(settingsFilePath, settingsByte, 0644)

}

func Load() (*Configuration, error) {

	configFile, err := os.ReadFile(settingsFilePath)

	if err != nil {
		if !os.IsNotExist(err) {
			return nil, err
		}

		settings, err := initConfiguration()
		if err != nil {
			return nil, err
		}

		return settings, nil

	}

	conf := Configuration{}

	err = json.Unmarshal(configFile, &conf)
	if err != nil {
		return nil, err
	}

	return &conf, nil

}

// initConfiguration creates a default Configuration file
// It does not check for the existing of a Configuration file and will override existing ones
func initConfiguration() (*Configuration, error) {

	err := ensureDirectory(settingsDirPath)
	if err != nil {
		return nil, err
	}

	println("create default Configuration")

	installLocation, err := DetectGameInstallDir()
	if err != nil {
		return nil, err
	}

	defaultConfiguration := Configuration{InstallLocation: *installLocation}

	err = defaultConfiguration.save()
	if err != nil {
		return nil, err
	}

	return &defaultConfiguration, nil
}

// ensureDirectory ensures that configuration base dir exists
func ensureDirectory(settingsDirPath string) error {

	_, err := os.Stat(settingsDirPath)
	if err != nil {

		if !os.IsNotExist(err) {
			return err
		}

		// println("create Configuration dir")
		err = os.MkdirAll(settingsDirPath, 0766)
		if err != nil {
			return err
		}
	}

	return nil

}
