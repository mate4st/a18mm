package mod

import (
	"archive/zip"
	"errors"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
)

const tempEnvName = "temp"

const loaderLocation = `Bin\Win64`

var tempFilePath = filepath.Join(os.Getenv(tempEnvName), string(rune(rand.Int()))+assetName)

func installOrUpdateModLoader(installLocation string, downloadUrl string) error {

	// ensures that no file is left in temp
	defer os.Remove(tempFilePath)

	err := getModLoader(downloadUrl)
	if err != nil {
		return err
	}

	zipReader, err := zip.OpenReader(tempFilePath)
	if err != nil {
		return err
	}

	defer zipReader.Close()

	for _, zipLoaderFile := range zipReader.File {
		err = unzipFile(zipLoaderFile, filepath.Join(installLocation, loaderLocation, zipLoaderFile.Name))
		if err != nil {
			return err
		}
	}

	return nil

}

func unzipFile(zipFile *zip.File, destination string) error {

	fmt.Printf("Writing %s\n", zipFile.Name)
	rc, err := zipFile.Open()
	if err != nil {
		return fmt.Errorf("cannot open file in zip %w", err)
	}
	defer rc.Close()

	loaderFile, err := os.OpenFile(destination, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)
	if err != nil {
		return fmt.Errorf("cannot access mod loader file %w", err)
	}
	defer loaderFile.Close()

	_, err = io.Copy(loaderFile, rc)
	if err != nil {
		return err
	}

	return nil
}

func getModLoader(downloadUrl string) error {
	resp, err := http.Get(downloadUrl)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return errors.New("failed to download new mod loader release")
	}

	tempFile, err := os.Create(tempFilePath)
	if err != nil {
		return err
	}
	defer tempFile.Close()

	_, err = io.Copy(tempFile, resp.Body)

	return err

}
