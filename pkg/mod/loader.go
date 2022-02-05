package mod

import (
	"encoding/json"
	"io"
	"net/http"
	"unsafe"

	"golang.org/x/sys/windows"

	"github.com/Masterminds/semver"
	"github.com/google/go-github/v42/github"
)

const installPath = "D:\\Gaming\\steamapps\\common\\Anno 1800\\Bin\\Win64\\"

const releaseUrl = "https://api.github.com/repos/xforce/anno1800-mod-loader/releases/latest"
const assetName = "loader.zip"

const loaderFileName = "python35.dll"

// Second parameter access specific language where 0409 = en-US
const verQueryValuePrefix = `\StringFileInfo\040904b0\`
const verQueryValueProductName = "ProductName"
const verQueryValueProductVersion = "ProductVersion"

const productName = "Anno 1800 Mod Loader"
const noVersion = "v0.0.0"

func VerifyModLoader() error {

	installedVersion, err := getInstalledModVersion()
	if err != nil {
		return err
	}

	println(installedVersion.String())

	resp, err := http.Get(releaseUrl)
	if err != nil {
		return err
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	release := github.RepositoryRelease{}
	err = json.Unmarshal(body, &release)

	gitHubLatestReleaseVersion, err := semver.NewVersion(*release.TagName)
	if err != nil {
		return err
	}

	println(gitHubLatestReleaseVersion.String())

	for _, asset := range release.Assets {

		if *asset.Name == assetName {
			println(*asset.BrowserDownloadURL)

		}
	}

	return nil

}

func getInstalledModVersion() (*semver.Version, error) {

	var handle windows.Handle
	versionInfoSize, err := windows.GetFileVersionInfoSize(installPath+loaderFileName, &handle)
	if err != nil {
		return nil, err
	}

	versionInfo := make([]byte, versionInfoSize)
	err = windows.GetFileVersionInfo(installPath+loaderFileName, uint32(handle), versionInfoSize, unsafe.Pointer(&versionInfo[0]))
	if err != nil {
		return nil, err
	}

	productNameBuffer := make([]uint16, 32)
	productNameBufferLength := uint32(unsafe.Sizeof(productNameBuffer))
	productNameBufferMaxLength := productNameBufferLength

	err = windows.VerQueryValue(unsafe.Pointer(&versionInfo[0]), verQueryValuePrefix+verQueryValueProductName, (unsafe.Pointer)(&productNameBuffer), &productNameBufferLength)
	if err != nil {
		return nil, err
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
		return nil, err
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
