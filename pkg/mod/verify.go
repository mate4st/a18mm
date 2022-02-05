package mod

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"path/filepath"
	"unsafe"

	"golang.org/x/sys/windows"

	"github.com/Masterminds/semver"
	"github.com/google/go-github/v42/github"
)

const latestReleaseUrl = "https://api.github.com/repos/xforce/anno1800-mod-loader/releases/latest"
const assetName = "loader.zip"

const loaderFileName = "python35.dll"

// Second parameter access specific language where 0409 = en-US
const verQueryValuePrefix = `\StringFileInfo\040904b0\`
const verQueryValueProductName = "ProductName"
const verQueryValueProductVersion = "ProductVersion"

const productName = "Anno 1800 Mod Loader"
const noVersion = "v0.0.0"

// VerifyModLoader checks for mod loader updates and turies to determine the current installed version based on the file version info.
// If no version data found it uses lastKnownVersion to decide weather to install a new release or not.
// Set lastKnownVersion to nil will force an install of the latest github release version.
func VerifyModLoader(installLocation string, lastKnownVersion *semver.Version) (*semver.Version, error) {

	if lastKnownVersion == nil {
		lastKnownVersion, _ = semver.NewVersion(noVersion)
	}

	installedVersion, err := getInstalledModVersion(installLocation)

	if err != nil {
		fmt.Println(err)
		installedVersion = lastKnownVersion
	}

	println(installedVersion.String())

	release, err := getLatestRelease()
	if err != nil {
		return nil, err
	}

	latestReleaseVersion, err := semver.NewVersion(*release.TagName)
	if err != nil {
		return nil, err
	}

	println(latestReleaseVersion.String())

	if !latestReleaseVersion.GreaterThan(installedVersion) {
		return installedVersion, nil
	}

	var downloadUrl string
	for _, asset := range release.Assets {
		if *asset.Name == assetName {
			downloadUrl = *asset.BrowserDownloadURL
			break
		}
	}

	if downloadUrl == "" {
		return nil, errors.New("could not find asset download url")
	}

	err = installOrUpdateModLoader(installLocation, downloadUrl)
	if err != nil {
		return nil, err
	}

	return latestReleaseVersion, nil

}

func getLatestRelease() (*github.RepositoryRelease, error) {
	resp, err := http.Get(latestReleaseUrl)
	if err != nil {
		return nil, err
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	release := github.RepositoryRelease{}
	err = json.Unmarshal(body, &release)
	if err != nil {
		return nil, err
	}

	return &release, nil

}

func getInstalledModVersion(installLocation string) (*semver.Version, error) {

	var handle windows.Handle
	versionInfoSize, err := windows.GetFileVersionInfoSize(filepath.Join(installLocation, loaderLocation, loaderFileName), &handle)
	if err != nil {
		return nil, fmt.Errorf("cannot get python lib version info size: %w", err)
	}

	versionInfo := make([]byte, versionInfoSize)
	err = windows.GetFileVersionInfo(filepath.Join(installLocation, loaderLocation, loaderFileName), uint32(handle), versionInfoSize, unsafe.Pointer(&versionInfo[0]))
	if err != nil {
		return nil, fmt.Errorf("cannot acces python lib version info: %w", err)
	}

	productNameBuffer := make([]uint16, 32)
	productNameBufferLength := uint32(unsafe.Sizeof(productNameBuffer))
	productNameBufferMaxLength := productNameBufferLength

	err = windows.VerQueryValue(unsafe.Pointer(&versionInfo[0]), verQueryValuePrefix+verQueryValueProductName, (unsafe.Pointer)(&productNameBuffer), &productNameBufferLength)
	if err != nil {
		return nil, fmt.Errorf("cannot acces python lib prdocut name: %w", err)
	}

	productNameUnicodeString := windows.NTUnicodeString{
		Length:        uint16(productNameBufferLength),
		MaximumLength: uint16(productNameBufferMaxLength),
		Buffer:        &productNameBuffer[0],
	}

	// println(productNameUnicodeString.String())
	if productNameUnicodeString.String() != productName {
		version, err := semver.NewVersion(noVersion)
		if err != nil {
			return nil, err
		}

		return version, nil
	}

	productVersionBuffer := make([]uint16, 32)
	productVersionBufferLength := uint32(unsafe.Sizeof(productVersionBuffer))
	productVersionBufferMaxLength := productVersionBufferLength

	err = windows.VerQueryValue(unsafe.Pointer(&versionInfo[0]), verQueryValuePrefix+verQueryValueProductVersion, (unsafe.Pointer)(&productVersionBuffer), &productNameBufferLength)
	if err != nil {
		return nil, fmt.Errorf("cannot acces python lib data: %w", err)
	}

	productVersionUnicodeString := windows.NTUnicodeString{
		Length:        uint16(productVersionBufferLength),
		MaximumLength: uint16(productVersionBufferMaxLength),
		Buffer:        &productVersionBuffer[0],
	}

	// println(productVersionUnicodeString.String())

	version, err := semver.NewVersion(productVersionUnicodeString.String())
	if err != nil {
		return nil, err
	}

	return version, nil

}
