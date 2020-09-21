package ftp

import (
	"fmt"
	"log"
	"os"
	"strings"
	"path/filepath"
)

//FS is a virtual file system
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

	return &FS{currentDirectory: userDir, baseDir: userDir}
}

func (fs *FS) Pwd() string {
	path := virtualPath(fs.currentDirectory, fs)

	if path == "" {
		return "/"
	}

	return path
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

func (fs *FS) goToRoot() {
	fs.currentDirectory = fs.baseDir
}

//Cwd navigates to the given path. The resulting path must be a subtree of the user virtual space.
// Notice that navigation can only target directories. If the resulting path is a file, an PathError is returned and
// the current directory is unchanged and returned along the error.
func (fs *FS) Cwd(path string) (string, error) {

	if path == "/" {
		fs.goToRoot()

		return fs.Pwd(), nil
	}

	// try to navigate from root
	if strings.HasPrefix(path, "/") {
		pathFromRoot := path[1:]

		current := fs.currentDirectory

		fs.goToRoot()
		cwd, err := fs.Cwd(pathFromRoot)

		if err != nil {
			fs.currentDirectory = current

			return fs.Pwd(), PathError{path: path, cause: err.Error()}
		}

		return cwd, nil
	}

	newPath := fs.currentDirectory + "/" + path

	fd, err := os.Open(newPath)

	if err != nil {
		return fs.Pwd(), PathError{path: path, cause: err.Error()}
	}

	defer fd.Close()

	realPath, _ := filepath.Abs(newPath)

	if !strings.HasPrefix(realPath, fs.baseDir) {
		return fs.Pwd(), PathError{path: path, cause: fmt.Sprintf("Trying to leave user base directory: %v", fs.baseDir)}
	}

	fileInfo, err := fd.Stat()

	if err != nil {
		return fs.Pwd(), PathError{path: path, cause: fmt.Sprintf("Error getting information from path: %v", virtualPath(realPath, fs))}
	}

	if !fileInfo.IsDir() {
		return fs.Pwd(), PathError{path: path, cause: fmt.Sprintf("Path %v is not a directory", virtualPath(realPath, fs))}
	}

	fs.currentDirectory = newPath

	return fs.Pwd(), nil
}

//WriteTo writes the stream of data being received on the data channel to the file referenced by fileName.
//
// Returns the total number of bytes written along with any error (or nil if everything happens correctly)
// Channel Close signal is expected
// If the file exists it is truncated.
// If the file is a directory, PathError is returned
// Failure operations on the file are reported using PathError
func (fs *FS) WriteTo(fileName string, data <-chan Transmission) (int, error) {

	filePath := fs.currentDirectory + "/" + fileName

	log.Printf("Trying to write on file %v", virtualPath(filePath, fs))

	info, exists := fileExists(filePath)

	if exists {
		if info.IsDir() {
			log.Printf("%v is a directory. Cannot write on directories.", virtualPath(filePath, fs))

			return 0, PathError{path: virtualPath(filePath, fs), cause: fmt.Sprintf("Given file name %v is a directory.", fileName)}
		}

		err := os.Remove(filePath)

		if err != nil {
			log.Printf("Cannot delete file %v", virtualPath(filePath, fs))

			return 0, PathError{path: fileName, cause: err.Error()}
		}
	}

	fd, err := os.Create(filePath)

	if err != nil {
		return 0, PathError{path: fileName, cause: err.Error()}
	}

	defer fd.Close()

	totalWritten := 0

	for transmitted := range data {
		size := transmitted.size

		s, err := fd.Write(transmitted.data[0:size])

		totalWritten += s

		if err != nil {
			return totalWritten, PathError{path: virtualPath(filePath, fs), cause: err.Error()}
		}
	}

	log.Printf("Write completed for file %v", virtualPath(filePath, fs))

	return totalWritten, nil
}

//ReadFrom reads from the file name and transmits the file content into the channel
// Returns the total number of byte read along with any error that occurs.
func (fs *FS) ReadFrom(fileName string, to chan<- Transmission) (int, error) {
	filePath := fs.currentDirectory + "/" + fileName

	info, exists := fileExists(filePath)

	if !exists {
		return 0, PathError{path: virtualPath(filePath, fs), cause: "File not found"}
	}

	if info.IsDir() {
		return 0, PathError{path: virtualPath(filePath, fs), cause: "Referenced path is not a file, but a directory"}
	}

	fd, err := os.Open(filePath)

	if err != nil {
		return 0, PathError{path: fileName, cause: err.Error()}
	}

	defer fd.Close()

	totalRead := 0

	buffer := make([]byte, 1024)

	for {
		read, _ := fd.Read(buffer)
		totalRead += read

		if read <= 0 {
			return totalRead, nil
		}

		to <- Transmission{data: buffer, size: read}
	}
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

func fileExists(filename string) (os.FileInfo, bool) {
	info, err := os.Stat(filename)

	if os.IsNotExist(err) {
		return nil, false
	}

	return info, true
}

type Transmission struct {
	size int
	data []byte
}
