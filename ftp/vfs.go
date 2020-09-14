package ftp

import (
	"fmt"
	"os"
)

type FS struct {
	baseDir string
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
	return fs.currentDirectory
}

func (fs *FS) ls() []string {

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

func strInfo(info os.FileInfo) string {
	var typ string

	if info.IsDir() {
		typ = "DIR"
	} else {
		typ = "FILE"
	}

	return fmt.Sprintf("%v\t%v\t%v", info.Name(), typ, info.Size())
}
