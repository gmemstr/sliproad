package files

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"reflect"
	"testing"
)

const DISK_TESTING_GROUNDS = ".testing_data_directory/"

// Test basic file provider.
// Should return, essentially, nothing. Can be used as a template for other providers.
func TestFileProvider(t *testing.T) {
	fp := FileProvider{}

	setup := fp.Setup(map[string]string{}); if setup != false {
		t.Errorf("Default FileProvider Setup() returned %v, expected false.", setup)
	}

	getdirectory := fp.GetDirectory(""); if len(getdirectory.Files) != 0 {
		t.Errorf("Default FileProvider GetDirectory() files returned %v, expected none.", getdirectory.Files)
	}

	var w bytes.Buffer
	fp.ViewFile("", &w); if len(w.Bytes()) > 0 {
		t.Errorf("Default FileProvider ViewFile() %v, expected nothing.", w)
	}

	savefile := fp.SaveFile(nil, "", ""); if savefile != false {
		t.Errorf("Default FileProvider SaveFile() attempted to save a file.")
	}

	determinetype := fp.DetermineType(""); if determinetype != "" {
		t.Errorf("Default FileProvider DetermineType() did not return an empty string.")
	}
}

// Test functions provided by fileutils, which do not return anything.
func TestFileUtils(t *testing.T) {
	ProviderConfig = map[string]FileProvider{
		"test": {
			Name: "test_disk",
			Provider: "disk",
			Location: DISK_TESTING_GROUNDS,
		},
	}
	var i FileProviderInterface
	TranslateProvider("test", &i)
	typeof := reflect.TypeOf(i); if typeof != (reflect.TypeOf(&DiskProvider{})) {
		t.Errorf("TranslateProvider() returned %v, expected DiskProvider{}", typeof)
	}

	SetupProviders()
	if Providers["test"] == nil {
		t.Errorf("SetupProviders() did not setup test provider")
	}
}

func TestDiskProvider(t *testing.T) {
	// Initialize testing data
	t.Cleanup(DiskProviderTestCleanup)
	err := os.Mkdir(DISK_TESTING_GROUNDS, 0755)
	if err != nil {
		t.Fatalf("Failed to create testing directory for DiskProvider: %v", err.Error())
	}
	err = ioutil.WriteFile(DISK_TESTING_GROUNDS + "testing.txt", []byte("testing file!"), 0755)
	if err != nil {
		t.Fatalf("Failed to create testing file for DiskProvider: %v", err.Error())
	}

	dp := DiskProvider{FileProvider{Location: DISK_TESTING_GROUNDS}}
	setup := dp.Setup(map[string]string{}); if setup != true {
		t.Errorf("DiskProvider Setup() returned %v, expected true.", setup)
	}

	getdirectory := dp.GetDirectory(""); if len(getdirectory.Files) != 1 {
		t.Errorf("DiskProvider GetDirectory() files returned %v, expected 1.", getdirectory.Files)
	}

	var w bytes.Buffer
	dp.ViewFile("testing.txt", &w); if w.String() != "testing file!" {
		fmt.Println(w)
		t.Errorf("DiskProvider ViewFile() returned %v, expected \"testing file!\".", w.String())
	}

 	testfile := bytes.NewReader([]byte("second test file!"))
 	savefile := dp.SaveFile(testfile, "second_test.txt", ""); if savefile != true {
		t.Errorf("DiskProvider SaveFile() could not save a file.")
	}

	determinetype := dp.DetermineType("second_test.txt"); if determinetype != "file" {
		t.Errorf("DiskProvider DetermineType() returned %v, expected \"file\".", determinetype)
	}
}

func DiskProviderTestCleanup() {
	err := os.RemoveAll(DISK_TESTING_GROUNDS)
	if err != nil {
		fmt.Printf("Unable to remove testing directory %v, please manully remove it", DISK_TESTING_GROUNDS)
	}
}