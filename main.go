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


var logger *log.Logger

func usage(program string) {
	fmt.Printf("Usage: %s <src-dir> <dest-dir>\n", program)
	os.Exit(1)
}

func main() {
	logger = log.New(os.Stderr, "", 0)

	params := os.Args[1:]
	if len(params) < 2 {
		fmt.Printf("Missing required args\n")
		usage(os.Args[0])
	}

	src := getDir(os.Args[1])
	dest := getDir(os.Args[2])

	files, err := getFiles(src)
	if err != nil {
		log.Fatal(err)
	}

	for _, f := range files {
		// compute new destination path and check if it already exists
		newFilePath := path.Join(dest.Name(), strings.TrimPrefix(f, path.Clean(src.Name())))

		if _, err = os.Stat(newFilePath); err == nil {
			fmt.Printf("File [%s] already exists, skipping copy..\n", newFilePath)
			continue
		} else if err != nil && (err.(*os.PathError)).Err.Error() != "no such file or directory" {
			fmt.Printf("Error stating file [%s]: %v\n", newFilePath, err)
			continue
		}

		srcFile, err := os.Open(f)
		if err != nil {
			fmt.Printf("Error opening file [%s]: %v\n", f, err)
			continue
		}

		if err := os.MkdirAll(path.Dir(newFilePath), os.ModePerm); err != nil {
			fmt.Printf("Error creating file [%s]: %v\n", newFilePath, err)
			continue
		}
		destFile, err := os.Create(newFilePath)
		if err != nil {
			fmt.Printf("Error creating file [%s]: %v\n", newFilePath, err)
			continue
		}

		if err = fileCopy(srcFile, destFile); err != nil {
			fmt.Printf("Error copying file [%s]: %v\n", newFilePath, err)
		}

		fmt.Printf("Copied [%s] to [%s]\n", f, newFilePath)
	}
}

func getFiles(inputFile *os.File) ([]string, error) {
	s, err := inputFile.Stat()
	if err != nil {
		return nil, err
	}

	result := make([]string, 0)

	if !s.IsDir() {
		return result, nil
	}

	files, err := ioutil.ReadDir(inputFile.Name())
	if err != nil {
		return nil, err
	}

	for _, file := range files {
		f, err := os.Open(path.Join(inputFile.Name(), file.Name()))
		if err != nil {
			fmt.Printf("Couldn't open file: [%v]\n", file.Name())
			continue
		}

		defer f.Close()

		stat, err := f.Stat()
		if err != nil {
			fmt.Printf("Couldn't stat file: [%v]\n", file.Name())
			continue
		}

		if stat.IsDir() {
			contents, err := getFiles(f)
			if err == nil {
				result = append(result, contents...)
			}
		} else {
			result = append(result, f.Name())
		}
	}

	return result, nil
}

func fileCopy(src, dest *os.File) error {
	defer func() {
		src.Close()
		dest.Close()
	}()

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

func getDir(path string) *os.File {
	src, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}

	info, err := src.Stat()
	if err != nil {
		log.Fatal(err)
	}

	if !info.IsDir() {
		fmt.Printf("Invalid argument: %s\n", path)
		usage(os.Args[0])
	}

	return src
}

