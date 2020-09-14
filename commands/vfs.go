package commands

type FS struct {
	currentDirectory string
}

func NewFS(basePath string) *FS {
	return &FS{currentDirectory: basePath}
}

//ForUser returns virtual file system for specific user
func (fs *FS) ForUser(user string) *FS {
	return &FS{currentDirectory: fs.currentDirectory + "/" + user}
}

func (fs *FS) Pwd() string {
	return fs.currentDirectory
}
