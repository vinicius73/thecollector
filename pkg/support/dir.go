package support

import (
	"io/fs"
	"os"
	"path/filepath"
)

func GetBinDirPath() string {
	execPath, _ := os.Executable()
	return filepath.Dir(execPath)
}

func EnsureDir(dir string) error {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err = os.MkdirAll(dir, 0o700)

		if err != nil {
			return err
		}
	}

	return nil
}

// ListDirs lists directories recursively up to a maximum depth and returns the list of directories
func ListDirs(dir string, depth int, maxDepth int) ([]string, error) {
	var folders []string

	// Check if the current depth is less than or equal to the maximum depth
	if depth > maxDepth {
		return folders, nil
	}

	// Read the contents of the current directory
	files, err := os.ReadDir(dir)
	if err != nil {
		return folders, err
	}

	// Iterate over all files and directories in the current directory
	for _, file := range files {
		// Check if the file is a directory
		if file.IsDir() {
			// Add the current directory path to the list of folders
			folders = append(folders, filepath.Join(dir, file.Name()))

			// Call the ListDirs function recursively for the found directory,
			// increasing the depth by 1 and adding the results to the list of folders
			subFolders, err := ListDirs(filepath.Join(dir, file.Name()), depth+1, maxDepth)
			if err != nil {
				return folders, err
			}
			folders = append(folders, subFolders...)
		}
	}

	return folders, nil
}

func IsEmptyDir(path string) (bool, error) {
	// Open the directory
	dir, err := os.Open(path)
	if err != nil {
		return false, err
	}
	defer dir.Close()

	// Read the directory contents
	_, err = dir.Readdirnames(1)
	if err == nil {
		// The directory is not empty
		return false, nil
	}

	// Check if the error was due to the directory being empty
	if err == fs.ErrNotExist || err == os.ErrNotExist || err == os.ErrPermission {
		// The directory doesn't exist or we don't have permission to read it
		return false, err
	} else {
		// The directory is empty
		return true, nil
	}
}
