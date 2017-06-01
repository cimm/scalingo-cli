package update

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/Scalingo/cli/config"
	"github.com/Scalingo/cli/io"
	"gopkg.in/errgo.v1"
)

var (
	lastVersionURL = "http://cli-dl.scalingo.io/version"
	lastVersion    = ""
	gotLastVersion = make(chan struct{})
)

func init() {
	go func() {
		var err error
		lastVersion, err = getLastVersion()
		if err != nil {
			config.C.Logger.Println(err)
		}
		close(gotLastVersion)
	}()
}

func Check() error {
	version := config.Version

	if strings.HasSuffix(version, "dev") {
		fmt.Println("\nNo update checking, dev version:", version)
		return nil
	}

	select {
	case <-gotLastVersion:
	case <-time.After(time.Second * 4):
		fmt.Println("Timeout when connecting on scalingo server.")
		return errgo.New("Timeout when connecting on scalingo server.")
	}
	if version == lastVersion {
		return nil
	}

	io.Statusf("Your Scalingo client (%s) is out-of-date: some features may not work correctly.\n", version)
	io.Infof("Please update to '%s' by reinstalling it: http://cli.scalingo.com\n", lastVersion)
	return nil
}

func getLastVersion() (string, error) {
	res, err := http.Get(lastVersionURL)
	if err != nil {
		return "", errgo.Mask(err)
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", errgo.Mask(err)
	}

	return strings.TrimSpace(string(body)), nil
}
