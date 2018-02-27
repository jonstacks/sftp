package sftp

import (
	"os"
	"testing"

	"github.com/spf13/afero"
)

func TestHandleManager(t *testing.T) {
	fs := afero.NewMemMapFs()
	hm := newHandleManager()

	f, err := fs.OpenFile("/etc/passwd", os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		t.Errorf("Error: %s", err)
	}

	handle1 := hm.NextHandle(f)
	if handle1 != "1" {
		t.Errorf("Got bad file handle: %s", handle1)
	}

	handle2 := hm.NextHandle(f)
	if handle2 != "2" {
		t.Errorf("Got bad file handle: %s", handle2)
	}

	f, ok := hm.GetHandle(handle1)
	if !ok {
		t.Errorf("Expected to be able to get handle for 1, but was not able.")
	}

	err = hm.CloseHandle(handle1)
	if err != nil {
		t.Errorf("Expected to be able to close file by handle, but was not able.")
	}

	oh := hm.OpenHandles()
	if len(oh) != 1 {
		t.Errorf("Expected there to only be 1 open handle, but there was %d", len(oh))
	}

	_, ok = oh["2"]
	if !ok {
		t.Errorf("Expected to still have a file handle 2")
	}

	err = hm.CloseHandle("2")
	if err != nil {
		t.Errorf("Expected to be able to CloseHandle '2' but got error: %s", err)
	}

}
