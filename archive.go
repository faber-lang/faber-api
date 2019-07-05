package main

import (
	"archive/tar"
	"io/ioutil"
)

func CreateSourceArchive(source string) (string, error) {
	file, err := ioutil.TempFile("", "faber.tar")
	if err != nil {
		return "", err
	}
	defer file.Close()

	t := tar.NewWriter(file)
	defer t.Close()

	content := []byte(source)

	header := &tar.Header{
		Name: "/faber.fab",
		Size: int64(len(content)),
	}

	if err := t.WriteHeader(header); err != nil {
		return "", err
	}

	if _, err := t.Write(content); err != nil {
		return "", err
	}

	return file.Name(), nil
}
