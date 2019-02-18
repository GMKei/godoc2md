package goreadme

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"runtime"
	"strings"

	"github.com/WillAbides/godoc2md"
)

//go:generate ../bin/goreadme github.com/WillAbides/godoc2md/goreadme

//WriteReadme writes a README.md for pkgname to the given path
func WriteReadme(pkgName, readmePath string) (err error) {
	f, err := os.OpenFile(readmePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0640) //nolint:gosec
	if err != nil {
		return
	}
	defer func() {
		err = f.Close()
	}()
	err = ReadmeMD(pkgName, f)
	return
}

//VerifyReadme checks that the file at readmePath has the correct content for pkgName
func VerifyReadme(pkgName, readmePath string) (bool, error) {
	var want bytes.Buffer
	err := ReadmeMD(pkgName, &want)
	if err != nil {
		return false, err
	}

	got, err := ioutil.ReadFile(readmePath) //nolint:gosec
	if err != nil {
		if os.IsNotExist(err) {
			err = nil
		}
		return false, err
	}
	ok := bytes.Equal(want.Bytes(), got)
	return ok, nil
}

//ReadmeMD writes readme content for the given package to writer
func ReadmeMD(pkgName string, writer io.Writer) error {
	var buf bytes.Buffer
	config := &godoc2md.Config{
		TabWidth:          4,
		DeclLinks:         true,
		Goroot:            runtime.GOROOT(),
		SrcLinkHashFormat: "#L%d",
	}
	godoc2md.Godoc2md([]string{pkgName}, &buf, config)
	mdContent := buf.String()
	mdContent = strings.Replace(mdContent, `/src/target/`, `./`, -1)
	mdContent = strings.Replace(mdContent, fmt.Sprintf("/src/%s/", pkgName), `./`, -1)

	_, err := writer.Write([]byte(mdContent))
	return err
}