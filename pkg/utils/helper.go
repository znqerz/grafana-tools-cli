package utils

import (
	"fmt"
	"os"
	"os/user"
	"reflect"
	"strings"
)

// PkgRootPath returns object pkg root path.
func PkgRootPath(v interface{}) (string, error) {
	var (
		wd      string
		err     error
		pkgPath = reflect.TypeOf(v).PkgPath()
	)
	if wd, err = os.Getwd(); err != nil {
		return "", err
	}
	dirs := strings.Split(pkgPath, "/")
	suffix := fmt.Sprintf("/%s", strings.Join(dirs[3:], "/"))

	return strings.TrimSuffix(wd, suffix), nil
}

// FileExist check file exist or not, if exist then return true
func FileExist(name string) bool {
	if _, err := os.Stat(name); os.IsNotExist(err) {
		return false
	}

	return true
}

// WriteStr2File write string content to file
func WriteStr2File(name string, content string) error {
	if err := os.WriteFile(name, []byte(content), 0600); err != nil {
		return err
	}

	return nil
}

// ReadFile2Str return file string content
func ReadFile2Str(name string) (string, error) {
	var (
		err     error
		content []byte
	)

	if content, err = os.ReadFile(name); err != nil {
		return "", err
	}

	return string(content), nil
}

// HomeDir return home path
func HomeDir() (string, error) {
	var (
		u   *user.User
		err error
	)
	if u, err = user.Current(); err != nil {
		return "", err
	}

	return u.HomeDir, nil
}
