package utils

import (
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"

	"github.com/VKCOM/noverify/src/constfold"
	"github.com/VKCOM/noverify/src/ir"
	"github.com/VKCOM/noverify/src/meta"
)

// ResolveRequirePath resolves the path for the passed expression
// based on the project path and current state.
//
// For the function to resolve the path correctly, the passed expression
// must be constant, otherwise, an empty string and a false flag will be returned.
//
// Example:
//   require_once "./file.php";                     // Correct
//   require_once __DIR__ ."/file.php";             // Correct
//   require_once some_root(__DIR__) ."/file.php";  // Incorrect, some_root function is not constant
//
// Note: For testing using package import, you need to take into account that the
// www folder is added to the path, so for all new tests, all test files must be placed
// in the www folder for golden tests, or the www folder must be added to the path for
// other tests for the imports to work correctly.
func ResolveRequirePath(st *meta.ClassParseState, projectPath string, e ir.Node) (string, bool) {
	res := constfold.Eval(st, e)
	if res.Type == meta.Undefined {
		return "", false
	}

	path, ok := res.ToString()
	if !ok {
		return "", false
	}

	// In order to correctly handle paths in unix-like systems and in windows,
	// we need to bring all slashes to the form as in unix.
	path = filepath.ToSlash(path)

	// If the path is absolute, then we don't need
	// to do anything with it.
	if filepath.IsAbs(path) {
		return filepath.Clean(path), true
	}

	// If the path contains a prefix equal to the project path,
	// then nothing needs to be done with it.
	//
	// This usually happens during golden tests.
	if strings.HasPrefix(path, projectPath) {
		return filepath.Clean(path), true
	}

	// We need to put a slash in the beginning only for unix,
	// in the case of windows, this is not required.
	var pathBegin string
	if os.PathSeparator == '/' {
		pathBegin = `/`
	}

	if !strings.HasSuffix(projectPath, `/`) {
		projectPath = projectPath + `/`
	}

	fullName := pathBegin + projectPath + path
	clean := filepath.Clean(fullName)

	return clean, true
}

func GenIndent(level int) string {
	var res string
	for i := 0; i < level; i++ {
		res += "   "
	}
	return res
}

func IsEmbeddedConstant(name string) bool {
	name = strings.ToLower(name)
	return name == "null" || name == "true" || name == "false"
}

func IsSuperGlobal(name string) bool {
	name = strings.ToLower(name)
	return name == "GLOBALS" ||
		name == "_SERVER" ||
		name == "_GET" ||
		name == "_POST" ||
		name == "_REQUEST" ||
		name == "_COOKIE" ||
		name == "_FILES" ||
		name == "_SESSION" ||
		name == "_ENV"
}

func DefaultCacheDir() string {
	defaultCacheDir, err := os.UserCacheDir()
	if err != nil {
		defaultCacheDir = ""
	} else {
		defaultCacheDir = filepath.Join(defaultCacheDir, "phpstats")
	}
	return defaultCacheDir
}

func DefaultGraphsDir() string {
	defaultCacheDir, err := os.UserCacheDir()
	if err != nil {
		defaultCacheDir = ""
	} else {
		defaultCacheDir = filepath.Join(defaultCacheDir, "graphs")
	}
	return defaultCacheDir
}

func NormalizeSlashes(str string) string {
	return strings.ReplaceAll(str, `\`, `\\`)
}

var NameToIdentifierRegexp = regexp.MustCompile("[^a-zA-Z0-9]")

func NameToIdentifier(str string) string {
	return NameToIdentifierRegexp.ReplaceAllString(str, "_")
}

func OpenFile(url string) error {
	var cmd string
	var args []string

	switch runtime.GOOS {
	case "windows":
		cmd = "cmd"
		args = []string{"/c", "start"}
	case "darwin":
		cmd = "open"
	default: // "linux", "freebsd", "openbsd", "netbsd"
		cmd = "xdg-open"
	}
	args = append(args, url)
	return exec.Command(cmd, args...).Start()
}
