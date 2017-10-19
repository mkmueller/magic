// Copyright (c) 2017 Mark K Mueller, (markmueller.com), All rights reserved.

// Package magic will determine the type of file based on it's magic-bytes. This
// package will identify the most common types of files. There are a kbutt load of
// signatures at http://www.garykessler.net/library/file_sigs.html
//
package magic

import (
	"os"
	"errors"
)

// File signature struct
type fs struct {

	// Description of file type
	Description string

	// Comma delimited list of known extensions
	Ext			string

	// Offset where the signature begins
	Offset		int

	// Signature of file type
	Signature	[]byte

}

// FileInfo struct
type FileInfo struct {

	// Name of the file
	Name string

	// Comma delimited list of known extensions for this file type
	Ext			string

	// Description of the file
	Description string

	// Size of the file in bytes
	Size int64

}

const n_bytes = 8

// FileBytes returns a FileInfo struct. The error will be non-nil if the file
// type cannot be determined.
//
func FileBytes ( file string ) (FileInfo, error) {

	var fi FileInfo
	var err error

	// Check if file name supplied
	if file == "" {
		err = errors.New("File name not supplied")
		return fi, err
	}

	// Check file exists
	fs,err := os.Stat(file)
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
	if size < n_bytes {
		err = errors.New("File too small to test")
		return fi, err
	}

	// Open file for reading
	f,err := os.Open(file)
	if err != nil {
		return fi, err
	}
	defer f.Close() // always

	// Read first n_bytes
	byts := make([]byte, n_bytes)
	_,err = f.Read(byts)
	if err != nil {
		return fi, err
	}

	// Lookup bytes in signature table
	if desc,ext := lookup(byts); desc != "" {
		fi.Description = desc
		fi.Name = file
		fi.Ext  = ext
		fi.Size = size
		return fi, err
	}

	err = errors.New("Unknown")
	return fi, err
}

func lookup ( b []byte ) (string,string) {
	for _,sig := range SignatureTable {
		ln := len(sig.Signature)
		offset := sig.Offset
		if equal( b[offset:ln], sig.Signature ) {
			return sig.Description,sig.Ext
		}
	}
	return "",""
}

func equal ( a1, a2 []byte ) bool {
	for i,_ := range a1 {
		if a1[i] != a2[i] {
			return false
		}
	}
	return true
}

