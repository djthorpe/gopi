package tradfri

import (
	"os"
	"path/filepath"
	"testing"
)

////////////////////////////////////////////////////////////////////////////////

func Test_Token_000(t *testing.T) {
	t.Log("Test_Token_000")
}

func Test_Token_001(t *testing.T) {
	tt := new(Token)
	tt.Id = "id"
	tt.Version = "version"
	tt.Token = "token"

	temp := os.TempDir()
	if path, err := tt.CreatePath(temp); err != nil {
		t.Error(err)
	} else {
		t.Log("path=", path)
	}
	defer os.RemoveAll(temp)
	if path, err := tt.CreatePath(filepath.Join(temp, "tradfri")); err != nil {
		t.Error(err)
	} else if err := tt.Write(path); err != nil {
		t.Error(err)
	} else if err := tt.Read(path); err != nil {
		t.Error(err)
	} else if tt.Id != "id" {
		t.Error("Unexpected Id value")
	} else if tt.Version != "version" {
		t.Error("Unexpected Version value")
	} else if tt.Token != "token" {
		t.Error("Unexpected Token value")
	} else if err := os.RemoveAll(path); err != nil {
		t.Error(err)
	}
}
