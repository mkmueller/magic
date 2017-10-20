// Copyright (c) 2017 Mark K Mueller, (markmueller.com), All rights reserved.

// Package magic will determine the type of file based on the first "magic
// bytes". The file signatures are not complete, but represent some of the most
// common file types.
//
package magic

import (
	"errors"
	"os"
)

// Byte array
type ba []byte

// String array
type sa []string

// File signature struct
type fs struct {

	// Description of file type
	Description string

	// Comma delimited list of known extensions
	Ext []string

	// Offset where the signature begins
	Offset int

	// Signature of file type
	Signature []byte

}

// Signature table
type st []fs

// FileInfo struct
type FileInfo struct {

	// Name of the file
	Name string

	// Comma delimited list of known extensions for this file type
	Ext []string

	// Description of the file
	Description string

	// Size of the file in bytes
	Size int64
}

var n_bytes int64

func init() {

	// set n_bytes to longest signature and offset
	for _, sig := range signatureTable {
		x := int64(sig.Offset + len(sig.Signature))
		if x > n_bytes {
			n_bytes = x
		}
	}

}

// FileBytes returns a FileInfo struct. The error will be non-nil if the file
// type cannot be determined or if the file cannot be read.
//
func FileBytes(file string) (FileInfo, error) {

	var fi FileInfo
	var err error

	// Check if file name supplied
	if file == "" {
		err = errors.New("File name not supplied")
		return fi, err
	}

	// Check file exists
	fs, err := os.Stat(file)
	if err != nil {
		err = errors.New("File does not exist")
		return fi, err
	}

	// Check if directory
	if fs.IsDir() {
		err = errors.New("Directory")
		return fi, err
	}

	// Check file size
	size := fs.Size()
	if size == 0 {
		err = errors.New("Empty file")
		return fi, err
	}

	// set read_len to the smaller of n_bytes or file size
	read_len := n_bytes
	if size < read_len {
		read_len = size
	}

	// Open file for reading
	f, err := os.Open(file)
	if err != nil {
		return fi, err
	}
	defer f.Close() // always

	// Read first read_len bytes
	byts := make([]byte, read_len)
	_, err = f.Read(byts)
	if err != nil {
		return fi, err
	}

	// Lookup bytes in signature table
	if desc, ext := lookup(byts); desc != "" {
		fi.Description = desc
		fi.Name = file
		fi.Ext = ext
		fi.Size = size
		return fi, err
	}

	err = errors.New("Unknown file type")
	return fi, err
}

func lookup(b []byte) (string, []string) {
	for _, sig := range signatureTable {
		ln := len(sig.Signature)
		offset := sig.Offset
		if equal(b[offset:ln], sig.Signature) {
			return sig.Description, sig.Ext
		}
	}
	return "", nil
}

func equal(a1, a2 []byte) bool {
	if len(a1) != len(a2) {
		return false
	}
	for i, _ := range a1 {
		if a1[i] != a2[i] {
			return false
		}
	}
	return true
}
