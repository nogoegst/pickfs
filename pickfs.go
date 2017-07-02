// pickfs.go - filesystem that contains only picked files
//
// To the extent possible under law, Ivan Markin waived all copyright
// and related or neighboring rights to this module of avant, using the creative
// commons "cc0" public domain dedication. See LICENSE or
// <http://creativecommons.org/publicdomain/zero/1.0/> for full details.

// Package pickfs file provides an implementation of the vfs.FileSystem
// interface that includes only picked files under defined aliases.

package pickfs // import "github.com/nogoegst/pickfs"

import (
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/tools/godoc/vfs"
	"golang.org/x/tools/godoc/vfs/mapfs"
)

// New returns a new FileSystem with set of files picked
// from the provided map.
// Map keys should be forward slash-separated pathnames
// and not contain a leading slash.
func New(fs vfs.FileSystem, m map[string]string) vfs.FileSystem {
	if len(m) == 0 {
		return fs
	}
	return pickfs{fs, m}
}

type pickfs struct {
	fs vfs.FileSystem
	m  map[string]string
}

func (fs pickfs) String() string { return "pickfs" }

func (fs pickfs) Close() error { return nil }

func (fs pickfs) lookup(p string) (string, bool) {
	p = strings.TrimPrefix(p, "/")
	dir, _ := filepath.Split(p)
	for alias, realf := range fs.m {
		if p == alias {
			return realf, true
		}
		if strings.HasPrefix(dir, alias) {
			subpath := strings.TrimPrefix(p, alias)
			return filepath.Join(realf, subpath), true
		}
	}
	return "", false
}

func (fs pickfs) Open(p string) (vfs.ReadSeekCloser, error) {
	realf, ok := fs.lookup(p)
	if !ok {
		return nil, os.ErrNotExist
	}
	return fs.fs.Open(realf)
}

func (fs pickfs) Lstat(p string) (os.FileInfo, error) {
	realf, ok := fs.lookup(p)
	if !ok {
		return mapfs.New(fs.m).Lstat(p)
	}
	return fs.fs.Lstat(realf)
}

func (fs pickfs) Stat(p string) (os.FileInfo, error) {
	realf, ok := fs.lookup(p)
	if !ok {
		return mapfs.New(fs.m).Stat(p)
	}
	return fs.fs.Stat(realf)
}

func (fs pickfs) ReadDir(p string) ([]os.FileInfo, error) {
	realf, ok := fs.lookup(p)
	if !ok {
		return mapfs.New(fs.m).ReadDir(p)
	}
	return fs.fs.ReadDir(realf)
}
