package main

import (
	"encoding/base64"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
)

type fileextSlice []string

// flag.Value interface
func (i *fileextSlice) String() string {
	return fmt.Sprintf("%s", *i)
}

var illegalFileext = [...]string{`/`, `\`, `..`}

//hasIllegal returns true if index has illegal extension characters.
func (i *fileextSlice) hasIllegal(index int) bool {
	l := i.Illegal(index)
	return len(l) > 0
}

//Illegal returns list of illegal extension characters on slice index.
func (i *fileextSlice) Illegal(index int) []string {
	l := *i
	item := l[index]
	var found = []string{}
	for _, inv := range illegalFileext {
		if strings.Contains(item, inv) {
			found = append(found, inv)
		}
	}
	return found
}

func (i *fileextSlice) Remove(index int) {
	s := *i
	s = append(s[:index], s[index+1:]...)
	*i = s
}

// Flag.Set(value string) error
func (i *fileextSlice) Set(value string) error {
	fmt.Printf("%s\n", value)
	*i = append(*i, value)
	index := len(*i) - 1
	hasInvalid := i.hasIllegal(index)
	if hasInvalid {
		i.Remove(index)
		return fmt.Errorf("file extension in flag parameter may not include %s", i.Illegal(index))
	}
	return nil
}

func (i *fileextSlice) Include(filename string) bool {
	var extension = filepath.Ext(filename)
	for _, flagExt := range *i {
		if extension == "."+flagExt {
			return true
		}
	}
	return false
}

var flagFileext fileextSlice

//IncludeFile is the created const file, files are included in files ordered by extension.
type IncludeFile struct {
	Name     string
	Ext      string
	ExtNoDot string
}

func NewIncludeFile(f os.FileInfo) IncludeFile {
	//var extension = filepath.Ext(f.Name())
	var extensionNoDot = extension[1:len(extension)]
	var name = "include" + extensionNoDot + ".go"
	return IncludeFile{Name: name, Ext: extensionNoDot}
}
func makeIncludeExtMap(fs []os.FileInfo) map[string][]IncludeFile {
	imap := map[string][]IncludeFile{}
	for _, f := range fs {
		if !flagFileext.Include(f.Name()) {
			fmt.Println(f.Name() + " not included because has wrong extension.")
			continue
		}
		imap[icFile.Ext] = NewIncludeFile(f)
		return imap
}

func (imap) makeGoFiles(){

}
func getFilesExtMap() {
	fs, _ := ioutil.ReadDir(".")
	for _, f := range fs {
		if !flagFileext.Include(f.Name()) {
			fmt.Println(f.Name() + " not included because has wrong extension.")
			continue
		}
		icFile := NewIncludeFile(f)
		if staticFilesMap[extensionNoDot] == nil {
			out, err := os.OpenFile("include"+extensionNoDot+".go", os.O_RDWR|os.O_CREATE, 0666)
			if err != nil {
				log.Fatalf("error opening file: %v", err)
			}
			out.Write([]byte("package main \n"))
			staticFilesMap[extensionNoDot] = append(staticFilesMap[extensionNoDot], out)
		}
	}
}

// Reads all .txt files in the current folder
// and encodes them as strings literals in textfiles.go
func main() {
	flag.Var(&flagFileext, "ext", "List of file extensions")
	flag.Parse()
	if flag.NFlag() == 0 {
		flag.PrintDefaults()
	} else {
		fmt.Println("Here are the values in 'flagFileext'")
		for i := 0; i < len(flagFileext); i++ {
			fmt.Printf("%s\n", flagFileext[i])
		}
	}

	fs, _ := ioutil.ReadDir(".")
	staticFilesMap := map[string][]*os.File{}

	for _, f := range fs {
		if !flagFileext.Include(f.Name()) {
			fmt.Println(f.Name() + " not included because has wrong extension.")
			continue
		}
		var extension = filepath.Ext(f.Name())
		var extensionNoDot = extension[1:len(extension)]
		if staticFilesMap[extensionNoDot] == nil {
			out, err := os.OpenFile("include"+extensionNoDot+".go", os.O_RDWR|os.O_CREATE, 0666)
			if err != nil {
				log.Fatalf("error opening file: %v", err)
			}
			out.Write([]byte("package main \n"))
			staticFilesMap[extensionNoDot] = append(staticFilesMap[extensionNoDot], out)
		}
		filename := f.Name()[0 : len(f.Name())-len(extension)]
		out := staticFilesMap[extensionNoDot]
		out.Write([]byte("\nconst (\n" + filename + " = `"))
		f, _ := os.Open(f.Name())
		defer f.Close()
		encoder := base64.NewEncoder(base64.StdEncoding, out)
		defer encoder.Close()
		io.Copy(encoder, f)
		out.Write([]byte("`\n"))
	}
	for _, outf := range staticFilesMap {
		outf.Write([]byte(")\n"))
		outf.Close()
	}
}
