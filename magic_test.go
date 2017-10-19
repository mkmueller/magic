// Copyright (c) 2017 Mark K Mueller, (markmueller.com), All rights reserved.

package magic_test

import (
	"os"
	"fmt"
	"testing"
	"io/ioutil"
	"github.com/mkmueller/magic"
    . "github.com/smartystreets/goconvey/convey"
)

var tiniestGifEver string

func TestMain ( m *testing.M ) {

	tiniestGifEver = createTestFile([]byte{0x47, 0x49, 0x46, 0x38, 0x39, 0x61,
								0x01, 0x00, 0x01, 0x00, 0x00, 0xFF, 0x00,
								0x2C, 0x00, 0x00, 0x00, 0x00, 0x01, 0x00,
								0x01, 0x00, 0x00, 0x02, 0x00, 0x3B})

	returncode := m.Run()

	os.Remove(tiniestGifEver)
	os.Exit(returncode)

}

// Example FileBytes
func ExampleFileBytes () {

	// Thanks to ProbablyProgramming.com for The Tiniest Gif Ever
	fi, err := magic.FileBytes(tiniestGifEver)
	if err == nil {
		fmt.Printf( "The file type is %s and the extension should probably be %s",  fi.Description, fi.Ext )
	}
	// Output: The file type is GIF89a Image and the extension should probably be gif

}

func TestFileBytes ( t *testing.T ) {

    Convey("Test The Tiniest GIF Ever", t, func() {

		fi,err := magic.FileBytes(tiniestGifEver)
		So( err, 			ShouldEqual, nil )
		So( fi.Description,	ShouldEqual, "GIF89a Image" )

    })

}

func TestFileBytes_errors ( t *testing.T ) {

    Convey("Test file name", t, func() {

		_,err := magic.FileBytes("")
		So( err.Error(),	ShouldEqual, "File name not supplied" )

    })

    Convey("Test non-existent file", t, func() {

		_,err := magic.FileBytes("./ thils file should not exist . really")
		So( err.Error(),	ShouldEqual, "File does not exist" )

    })

    Convey("Test using directory name", t, func() {

		_,err := magic.FileBytes("..")
		So( err.Error(),	ShouldEqual, "Directory" )

    })

    Convey("Test using a file too small", t, func() {

		testfile := createTestFile([]byte{0x47, 0x49, 0x46, 0x38, 0x39})

		_,err := magic.FileBytes(testfile)
		So( err.Error(),	ShouldEqual, "File too small to test" )

		os.Remove(testfile)

    })

    Convey("Test an unknown file type", t, func() {

		testfile := createTestFile([]byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
										  0x00, 0x00, 0x00, 0x00, 0x00, 0x00})

		_,err := magic.FileBytes(testfile)
		So( err.Error(),	ShouldEqual, "Unknown" )

		os.Remove(testfile)

    })

}




func createTestFile ( data []byte ) string {

	fl, err := ioutil.TempFile("/tmp", "GOTEST_")
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	defer fl.Close()

	_,err = fl.Write(data)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	return fl.Name()

}
