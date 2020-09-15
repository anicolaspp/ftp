package ftp

import (
	"fmt"
	"os"
	"strings"
	"path/filepath"
)

type FS struct {
	baseDir          string
	currentDirectory string
}

func NewFS(basePath string) *FS {
	return &FS{baseDir: basePath, currentDirectory: basePath}
}

//ForUser returns virtual file system for specific user
func (fs *FS) ForUser(user string) *FS {

	userDir := fs.baseDir + "/" + user

	_, err := os.Stat(userDir)

	if os.IsNotExist(err) {
		_, err = os.Create(userDir)

		if err != nil {
			panic(err)
		}
	}

	return &FS{currentDirectory: userDir, baseDir: fs.baseDir}
}

func (fs *FS) Pwd() string {
	return virtualPath(fs.currentDirectory, fs)
}

func (fs *FS) Ls() []string {

	fd, err := os.Open(fs.currentDirectory)

	if err != nil {
		fmt.Println(err)
		return []string{}
	}

	infos, err := fd.Readdir(0)

	if err != nil {
		return []string{}
	}

	result := make([]string, len(infos))

	for i, v := range infos {
		result[i] = strInfo(v)
	}

	return result
}

//Cwd navigates to the given path. The resulting path must be a subtree of the user virtual space.
func (fs *FS) Cwd(path string) (string, error) {
	newPath := fs.currentDirectory + "/" + path

	fd, err := os.Open(newPath)

	if err != nil {
		return fs.Pwd(), PathError{path: path, cause: err.Error()}
	}

	defer fd.Close()

	realPath, _ := filepath.Abs(newPath)

	if !strings.HasPrefix(realPath, fs.baseDir) {
		return fs.currentDirectory, PathError{path: path, cause: "Trying to leave user base directory"}
	}

	fileInfo, err := fd.Stat()

	if err != nil {
		return fs.Pwd(), PathError{path: path, cause: fmt.Sprintf("Error getting path information: %v", err)}
	}

	if !fileInfo.IsDir() {
		return fs.Pwd(), PathError{path: path, cause: fmt.Sprintf("Path %v is not a directory", virtualPath(realPath, fs))}
	}

	fs.currentDirectory = newPath

	return fs.Pwd(), nil
}

func strInfo(info os.FileInfo) string {
	var typ string

	if info.IsDir() {
		typ = "DIR"
	} else {
		typ = "FILE"
	}

	return fmt.Sprintf("%v\t%v\t%v", info.Name(), typ, info.Size())
}

func virtualPath(path string, fs *FS) string {
	virtualPath := strings.Replace(path, fs.baseDir, "", 1)

	return virtualPath
}

type PathError struct {
	path  string
	cause string
}

func (p PathError) Error() string {
	return fmt.Sprintf("Error accesing path: %v. %v", p.path, p.cause)
}
