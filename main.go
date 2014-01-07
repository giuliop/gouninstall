package main

import (
	"flag"
	"fmt"
	"os"
	"regexp"
	"runtime"
)

// help shows a helper message if number of parameters is wrong
func help(s string) {
	fmt.Println(s + "\n")
	fmt.Println("Usage: gouninstall <package> [<path>]")
	fmt.Println("If <path> is omitted then the first path in $GOPATH is taken\n")
}

// buildArgs build the list of dirs/files to delete
func buildArgs(pkg, dir string) (string, []string) {
	if dir == "" {
		r := regexp.MustCompile("[^:]*")
		dir = r.FindString(os.ExpandEnv("$GOPATH"))
	}
	pkg_dir := runtime.GOOS + "_" + runtime.GOARCH
	r := regexp.MustCompile("([a-zA-Z0-9_]+)/?$")
	pkgName := r.FindStringSubmatch(pkg)[1]
	dirs := []string{
		dir + "/src/" + pkg,
		dir + "/pkg/" + pkg_dir + "/src/" + pkg,
		dir + "/bin/" + pkgName}

	return pkg, dirs
}

// uninstall delete the dirs/files of the package
func uninstall(pkg string, dirs []string) {
	var okdir = make([]bool, 3)
	for i, d := range dirs {
		var dtype string
		switch i {
		case 0, 1:
			dtype = "Directory"
		case 2:
			dtype = "File"
		default:
			panic("Too many arguments!")
		}
		dinfo, err := os.Stat(d)
		if err != nil {
			if os.IsNotExist(err) {
				fmt.Printf("%s %s not found\n", dtype, d)
			} else {
				fmt.Printf("Error - %s\n", err)
			}
		} else if i < 2 && !dinfo.IsDir() {
			fmt.Printf("Error: %d is a file, not a directory\n")
		} else {
			fmt.Printf("%s %s will be DELETED\n", dtype, d)
			okdir[i] = true
		}
	}
	fmt.Println("\nDo you want to proceed? Type y/n")
	var input string
	fmt.Scanf("%s", &input)
	if input == "y" || input == "yes" {
		for i, d := range dirs {
			if okdir[i] {
				var f func(string) error
				if i < 2 {
					f = os.RemoveAll
				} else {
					f = os.Remove
				}
				err := f(d)
				if err != nil {
					fmt.Println(err)
				}
			}
		}
		fmt.Println("\nDone")
	} else {
		fmt.Println("Action cancelled")
	}
}

func main() {
	flag.Parse()
	args := flag.Args()
	var pkg string
	var dirs []string

	switch len(args) {
	case 0:
		help("No arguments given")
		return
	case 1:
		pkg, dirs = buildArgs(args[0], "")
	case 2:
		pkg, dirs = buildArgs(args[0], args[1])
	default:
		help("Too many arguments given")
		return
	}
	uninstall(pkg, dirs)
}
