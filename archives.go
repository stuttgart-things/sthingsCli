/*
Copyright © 2023 Patrick Hermann patrick.hermann@sva.de
*/

package cli

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

func UnzipArchive(source, destination string) error {

	// 1. OPEN THE ZIP FILE
	reader, err := zip.OpenReader(source)
	if err != nil {
		return err
	}
	defer reader.Close()

	// 2. GET THE ABSOLUTE DESTINATION PATH
	destination, err = filepath.Abs(destination)
	if err != nil {
		return err
	}

	// 3. ITERATE OVER ZIP FILES INSIDE THE ARCHIVE AND UNZIP EACH OF THEM
	for _, f := range reader.File {
		err := unzipFile(f, destination)
		if err != nil {
			return err
		}
	}

	return nil
}

func unzipFile(f *zip.File, destination string) error {

	// 4. CHECK IF FILE PATHS ARE NOT VULNERABLE TO ZIP SLIP
	filePath := filepath.Join(destination, f.Name)
	if !strings.HasPrefix(filePath, filepath.Clean(destination)+string(os.PathSeparator)) {
		return fmt.Errorf("invalid file path: %s", filePath)
	}

	// 5. CREATE DIRECTORY TREE
	if f.FileInfo().IsDir() {
		if err := os.MkdirAll(filePath, os.ModePerm); err != nil {
			return err
		}
		return nil
	}

	if err := os.MkdirAll(filepath.Dir(filePath), os.ModePerm); err != nil {
		return err
	}

	// 6. CREATE A DESTINATION FILE FOR UNZIPPED CONTENT
	destinationFile, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
	if err != nil {
		return err
	}
	defer destinationFile.Close()

	// 7. UNZIP THE CONTENT OF A FILE AND COPY IT TO THE DESTINATION FILE
	zippedFile, err := f.Open()
	if err != nil {
		return err
	}
	defer zippedFile.Close()

	if _, err := io.Copy(destinationFile, zippedFile); err != nil {
		return err
	}
	return nil
}
