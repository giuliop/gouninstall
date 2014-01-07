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
	fmt.Println("<package> is the package name including import path to be uninstalled")
	fmt.Println("<path> is to be used if you have more than one path in your $GOPATH. If omitted the first path in $GOPATH is taken\n")
	fmt.Println("Please note that this utility will not look in your GOROOT and assume packages are installed in your $GOPATH")
}

// buildArgs build the list of dirs/files to delete
func buildArgs(pkg, dir string) (string, []string) {
	if dir == "" {
		r := regexp.MustCompile("[^:]*")
		dir = r.FindString(os.ExpandEnv("$GOPATH"))
	}
	// strip "/" at end of pkg and dir if present to build path correctly
	if dir[len(dir)-1] == '/' {
		dir = dir[:len(dir)-1]
	}
	if pkg[len(pkg)-1] == '/' {
		pkg = pkg[:len(pkg)-1]
	}
	pkg_dir := runtime.GOOS + "_" + runtime.GOARCH
	r := regexp.MustCompile("([^/]+)/?$")
	pkgName := r.FindStringSubmatch(pkg)[1]
	dirs := []string{
		dir + "/src/" + pkg + "/",
		dir + "/pkg/" + pkg_dir + "/" + pkg + "/",
		dir + "/pkg/" + pkg_dir + "/" + pkg + ".a",
		dir + "/bin/" + pkgName}

	return pkg, dirs
}

// uninstall delete the dirs/files of the package
func uninstall(pkg string, dirs []string) {
	var okdir = make([]bool, 4)
	numFileToDelete := 0
	for i, d := range dirs {
		var dtype string
		switch i {
		case 0, 1:
			dtype = "directory"
		case 2, 3:
			dtype = "file"
		default:
			panic("Too many arguments!")
		}
		dinfo, err := os.Stat(d)
		if err != nil {
			if os.IsNotExist(err) {
				fmt.Printf("Looking for %s %s ... does not exist\n", dtype, d)
			} else {
				fmt.Printf("Error - %s\n", err)
			}
		} else if i < 2 && !dinfo.IsDir() {
			fmt.Printf("Error: %d is a file, not a directory\n")
		} else {
			fmt.Printf("Looking for %s %s ... found, will be DELETED\n", dtype, d)
			okdir[i] = true
			numFileToDelete++
		}
	}
	if numFileToDelete == 0 {
		fmt.Println("\nNothing to delete")
		return
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
