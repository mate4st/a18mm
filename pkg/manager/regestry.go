package manager

import (
	"errors"
	"fmt"

	"golang.org/x/sys/windows"
	"golang.org/x/sys/windows/registry"
)

const registryKeyPathM = `SOFTWARE\Microsoft\Windows\CurrentVersion\Uninstall`
const registryKeyPathW = `SOFTWARE\Wow6432Node\Microsoft\\Windows\CurrentVersion\Uninstall`

const keyInstallLocation = "InstallLocation"
const keyDisplayName = "DisplayName"
const keyValue = "Anno 1800"

type registrySearchPath struct {
	path string
	key  registry.Key
}

func DetectGameInstallDir() error {

	searchPaths := []registrySearchPath{
		{
			path: registryKeyPathM,
			key:  registry.CURRENT_USER,
		},

		{
			path: registryKeyPathW,
			key:  registry.CURRENT_USER,
		},

		{
			path: registryKeyPathM,
			key:  registry.LOCAL_MACHINE,
		},
		{
			path: registryKeyPathW,
			key:  registry.LOCAL_MACHINE,
		},
	}

	var installDir *string
	var err error

	for _, searchPath := range searchPaths {
		installDir, err = searchRegistryPath(searchPath.key, searchPath.path)
		if err != nil {
			return err
		}
		if installDir != nil {
			break
		}

	}

	if installDir == nil {
		return errors.New("could not find game installation directory")
	}

	println(*installDir)

	return nil

}

func searchRegistryPath(key registry.Key, path string) (*string, error) {

	searchDirectory, err := registry.OpenKey(key, path, registry.QUERY_VALUE|registry.ENUMERATE_SUB_KEYS)
	if err != nil {
		if err == windows.ERROR_FILE_NOT_FOUND {
			return nil, nil
		}
		return nil, err
	}
	defer searchDirectory.Close()

	keyInfo, err := searchDirectory.Stat()

	subKeys, err := searchDirectory.ReadSubKeyNames(int(keyInfo.SubKeyCount))
	if err != nil {
		return nil, err
	}

	for _, subKey := range subKeys {
		installLocation, err := searchSubDirectory(key, fmt.Sprintf(`%s\%s`, path, subKey))
		if err != nil {
			return nil, err
		}

		if installLocation != nil {
			return installLocation, nil
		}

	}

	return nil, nil

}

func searchSubDirectory(key registry.Key, path string) (*string, error) {
	subDirectory, err := registry.OpenKey(key, path, registry.QUERY_VALUE)
	if err != nil {
		return nil, err
	}

	defer subDirectory.Close()

	displayName, _, err := subDirectory.GetStringValue(keyDisplayName)
	if err != nil {
		if err == windows.ERROR_FILE_NOT_FOUND {
			return nil, nil
		}

		return nil, err
	}

	if displayName == keyValue {

		installLocation, _, err := subDirectory.GetStringValue(keyInstallLocation)
		if err != nil {
			if err == windows.ERROR_FILE_NOT_FOUND {
				return nil, err
			}
		}

		return &installLocation, nil
	}

	return nil, nil
}
