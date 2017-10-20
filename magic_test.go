// Copyright (c) 2017 Mark K Mueller, (markmueller.com), All rights reserved.

package magic_test

import (
	"fmt"
	"github.com/mkmueller/magic"
	. "github.com/smartystreets/goconvey/convey"
	"io/ioutil"
	"os"
	"testing"
)

var tiniestGifEver string

func TestMain(m *testing.M) {

	// Create the Tiniest GIF Ever.
	// Horked from http://probablyprogramming.com/2009/03/15/the-tiniest-gif-ever
	tiniestGifEver = createTestFile([]byte{0x47, 0x49, 0x46, 0x38, 0x39, 0x61,
		0x01, 0x00, 0x01, 0x00, 0x00, 0xFF,
		0x00, 0x2C, 0x00, 0x00, 0x00, 0x00,
		0x01, 0x00, 0x01, 0x00, 0x00, 0x02,
		0x00, 0x3B})

	// Wait here for my tests to complete
	returncode := m.Run()

	// Clean up
	os.Exit(returncode)

}

// Example FileBytes
func ExampleFileBytes() {

	// Create a tiny GIF image file
	fi, err := magic.FileBytes(tiniestGifEver)
	if err == nil {
		fmt.Printf("The file type is %s", fi.Description)
	}
	// Output: The file type is GIF89a Image

}

func TestFileBytes(t *testing.T) {

	Convey("Test unknown file", t, func() {

		fi, err := magic.FileBytes(tiniestGifEver)
		So(err, ShouldEqual, nil)
		So(fi.Description, ShouldEqual, "GIF89a Image")
		So(fi.Ext[0], ShouldEqual, "gif")

	})

}

func TestFileBytes_errors(t *testing.T) {

	unknownFile := createTestFile([]byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00})

	emptyFile := createTestFile(nil)

	Convey("Test file name", t, func() {

		_, err := magic.FileBytes("")
		So(err.Error(), ShouldEqual, "File name not supplied")

	})

	Convey("Test non-existent file", t, func() {

		_, err := magic.FileBytes("./hopefully, this file does not really exist")
		So(err.Error(), ShouldEqual, "File does not exist")

	})

	Convey("Test using directory name", t, func() {

		_, err := magic.FileBytes("..")
		So(err.Error(), ShouldEqual, "Directory")

	})

	Convey("Create an empty file", t, func() {

		_, err := magic.FileBytes(emptyFile)
		So(err.Error(), ShouldEqual, "Empty file")

	})

	Convey("Test an unknown file type", t, func() {

		_, err := magic.FileBytes(unknownFile)
		So(err.Error(), ShouldEqual, "Unknown file type")

	})

	// Clean up
	os.Remove(unknownFile)
	os.Remove(emptyFile)

}

func createTestFile(data []byte) string {

	fl, err := ioutil.TempFile("/tmp", "GOTEST_")
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	defer fl.Close()

	_, err = fl.Write(data)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	return fl.Name()

}
