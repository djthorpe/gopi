// +build dvb

package dvb_test

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/djthorpe/gopi/v3/pkg/sys/dvb"
)

var (
	FILES = []string{
		"/usr/share/dvb/dvb-c/de-Berlin",
		"/usr/share/dvb/dvb-t/de-Berlin",
	}
	FOLDERS = []string{
		"/usr/share/dvb/dvb-c",
		"/usr/share/dvb/dvb-t",
	}
)

func Test_Tunetable_000(t *testing.T) {
	for _, file := range FILES {
		if _, err := os.Stat(file); os.IsNotExist(err) {
			t.Log("Skipping", err)
			continue
		}
		fh, err := os.Open(file)
		if err != nil {
			t.Error(err)
			continue
		}
		defer fh.Close()
		scans, err := dvb.ReadTuneParamsTable(fh)
		if err != nil {
			t.Error(err)
			continue
		}
		for _, scan := range scans {
			t.Log("scan=", scan)
		}
	}
}

func Test_Tunetable_001(t *testing.T) {
	for _, path := range FOLDERS {
		if _, err := os.Stat(path); os.IsNotExist(err) {
			t.Log("Skipping", err)
			continue
		}
		files, err := ioutil.ReadDir(path)
		if err != nil {
			t.Error(err)
			continue
		}
		for _, file := range files {
			if file.Mode().IsRegular() == false {
				continue
			}
			if strings.HasPrefix(file.Name(), ".") {
				continue
			}
			fh, err := os.Open(filepath.Join(path, file.Name()))
			if err != nil {
				t.Error(err)
				continue
			}
			_, err = dvb.ReadTuneParamsTable(fh)
			fh.Close()
			if err != nil {
				t.Error(file.Name(), ":", err)
				continue
			}
		}
	}
}
