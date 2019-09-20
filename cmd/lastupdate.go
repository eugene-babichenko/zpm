package cmd

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"time"

	"github.com/pkg/errors"
)

func getLastUpdateTime() (t time.Time, err error) {
	filename := filepath.Join(rootDir, ".lastupdate")
	data, err := ioutil.ReadFile(filename)
	if os.IsNotExist(err) {
		return t, nil
	}
	if err != nil {
		return t, errors.Wrap(err, "while reading last update time")
	}
	t, err = time.Parse(time.RFC3339, string(data))
	if err != nil {
		return t, errors.Wrap(err, "failed to parse time")
	}
	return t, nil
}

func setLastUpdateTime(t time.Time) error {
	s := t.Format(time.RFC3339)
	filename := filepath.Join(rootDir, ".lastupdate")
	err := ioutil.WriteFile(filename, []byte(s), os.ModePerm)
	return err
}
