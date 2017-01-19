package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path"
	"strings"
)

func usage(program string) {
	fmt.Printf("Usage: %s <src-dir> <dest-dir>\n", program)
	os.Exit(1)
}

func main() {
	params := os.Args[1:]
	if len(params) < 2 {
		fmt.Printf("Missing required args\n")
		usage(os.Args[0])
	}

	src, err := getDir(os.Args[1])
	if err != nil {
		fmt.Printf("Arg error: %v\n", err)
		usage(os.Args[0])
	}

	dest, err := getDir(os.Args[2])
	if err != nil {
		fmt.Printf("Arg error: %v\n", err)
		usage(os.Args[0])
	}

	if err := syncFiles(src, src, dest); err != nil {
		fmt.Printf("Error syncing directory %s: %v\n", src.Name(), err)
		os.Exit(1)
	}
}

func syncFiles(input *os.File, src, dest *os.File) error {
	stat, err := input.Stat()
	if err != nil {
		fmt.Printf("Error processing %s: %v\n", input.Name(), err)
		return nil
	} else if !stat.IsDir() {
		return nil
	}

	files, err := ioutil.ReadDir(input.Name())
	if err != nil {
		return err
	}

	for _, fInfo := range files {
		fName := path.Join(input.Name(), fInfo.Name())

		// TODO: doesn't handle symlinks
		f, err := os.Open(fName)
		if err != nil {
			fmt.Printf("Couldn't open file %s: %v\n", fName, err)
			continue
		}
		defer f.Close()

		fStat, err := f.Stat()
		if err != nil {
			fmt.Printf("Couldn't stat file %s: %v\n", fName, err)
			continue
		}

		if fStat.IsDir() {
			if err := syncFiles(f, src, dest); err != nil {
				fmt.Printf("Error syncing directory %s: %v\n", fName, err)
			}
		} else {
			newFName := path.Join(dest.Name(), strings.TrimPrefix(fName, path.Clean(src.Name())))

			// check if file already exists
			if _, err = os.Stat(newFName); err == nil {
				fmt.Printf("File %s already exists, skipping copy..\n", newFName)
				continue
			} else if err != nil && (err.(*os.PathError)).Err.Error() != "no such file or directory" {
				fmt.Printf("Error opening file %s: %v\n", newFName, err)
				continue
			}

			// make sure parent dirs exist
			if err := os.MkdirAll(path.Dir(newFName), os.ModePerm); err != nil {
				fmt.Printf("Error creating file %s: %v\n", newFName, err)
				continue
			}

			destFile, err := os.Create(newFName)
			if err != nil {
				fmt.Printf("Error creating file %s: %v\n", newFName, err)
				continue
			}

			if err := fileCopy(f, destFile); err != nil {
				fmt.Printf("Error copying file %s to %s: %v\n", fName, newFName, err)
			}

			fmt.Printf("Copied %s to %s\n", f.Name(), destFile.Name())

			defer func(fs... *os.File) {
				for _, f := range fs {
					f.Close()
				}
			}(f, destFile)
		}
	}

	return nil
}

func fileCopy(src, dest *os.File) error {
	//TODO: Needs to handle symlinks

	nw, err := io.Copy(dest, src)
	if err != nil {
		log.Fatal(err)
	}

	if err := dest.Sync(); err != nil {
		return err
	}

	stat, err := src.Stat()
	if err != nil {
		return err
	}

	if nw != stat.Size() {
		return fmt.Errorf("Bytes written [%v] is not the same as src size [%v]", nw, stat.Size())
	}

	return nil
}

func getDir(path string) (*os.File, error) {
	src, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	info, err := src.Stat()
	if err != nil {
		return nil, err
	}

	if !info.IsDir() {
		return nil, fmt.Errorf("%s is not a dir!", path)
	}

	return src, nil
}

