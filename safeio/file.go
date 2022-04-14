package safeio

import (
	"errors"
	"io/fs"
	"os"
)

type FSActionsChecker interface {
	CheckFSActions(path string, actions FSAction) bool
}

type FSAction uint16

const (
	ReadFileAction FSAction = 1 << iota
	WriteFileAction
	CreateFileAction
	DeleteFileAction
	CreateFileInDirAction
)

func OpenFile(r FSActionsChecker, name string, flag int, perm fs.FileMode) (*os.File, error) {
	if !r.CheckFSActions(name, osFlagToFSActions(flag)) {
		return nil, ErrNotAllowed
	}
	return os.OpenFile(name, flag, perm)
}

func TempFile(r FSActionsChecker, dir string, pattern string) (*os.File, error) {
	if !r.CheckFSActions(dir, CreateFileAction) {
		return nil, ErrNotAllowed
	}
	return os.CreateTemp(dir, pattern)
}

func RemoveFile(r FSActionsChecker, name string) error {
	if !r.CheckFSActions(name, DeleteFileAction) {
		return ErrNotAllowed
	}
	return os.Remove(name)
}

func RenameFile(r FSActionsChecker, oldName, newName string) error {
	if !r.CheckFSActions(oldName, DeleteFileAction) || !r.CheckFSActions(newName, CreateFileAction) {
		return ErrNotAllowed
	}
	return os.Rename(oldName, newName)
}

type readFileChecker interface {
	FSActionsChecker
	RequireBytes(int) uint64
}

// ReadFile returns the contents of the file and requires the memory, so is safe
// to use in memory-constrained environments.
func ReadFile(r readFileChecker, name string) ([]byte, error) {
	if !r.CheckFSActions(name, ReadFileAction) {
		return nil, ErrNotAllowed
	}
	fi, err := os.Stat(name)
	if err != nil {
		return nil, err
	}
	r.RequireBytes(int(fi.Size()))
	return os.ReadFile(name)
}

var ErrNotAllowed = errors.New("safeio: operation not allowed")

func osFlagToFSActions(flag int) FSAction {
	var perms FSAction
	switch {
	case flag&os.O_RDONLY != 0:
		perms |= WriteFileAction
	case flag&os.O_WRONLY != 0:
		perms |= ReadFileAction
	case flag&os.O_RDWR != 0:
		perms |= ReadFileAction | WriteFileAction
	}
	if flag&os.O_CREATE != 0 {
		perms |= CreateFileAction
	}
	return perms
}
