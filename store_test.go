package lightdb

import (
	"os"
	"path"
	"testing"
)

type A struct {
	x string
	Y int
	T bool
	Z byte
}

func TestStore(t *testing.T) {
	obj := A{
		x: "hello",
		Y: 123,
		T: false,
		Z: 0x1,
	}

	test_file_path := path.Join(os.TempDir(), "test_store_write.data")

	writeError := writeObject(test_file_path, obj)
	if writeError != nil {
		t.Fatal(writeError)
	}

	_obj := A{}
	readError := readObject(test_file_path, &_obj)
	if readError != nil {
		t.Fatal(readError)
	}
}
