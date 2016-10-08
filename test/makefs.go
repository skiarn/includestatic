package main

import (
	"bufio"
	"encoding/base64"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"reflect"
	"strings"
	"time"
)

func WM(w io.Writer, f interface{}) {
	t := reflect.TypeOf(f)
	for i := 0; i < t.NumMethod(); i++ {
		method := t.Method(i)
		//fmt.Println(method.Name, "*********")
		//fmt.Println(method.Type)

		//build function declaration as string. example: func WM(w io.Writer, f os.FileInfo)
		bFunc := func(m reflect.Method, f interface{}) string {

			getParamList := func(fn func(int) reflect.Type, size int) []string {
				var list []string
				for i := 0; i < size; i++ {
					list = append(list, fmt.Sprintf("%v", fn(i)))
				}
				return list
			}
			inParamList := []string{} //getParamList(m.Type.In, m.Type.NumIn())
			outParamList := getParamList(m.Type.Out, m.Type.NumOut())

			return fmt.Sprintf("func(%s)%s(%s)(%s)", reflect.TypeOf(f), m.Name, strings.Join(inParamList, ","), strings.Join(outParamList, ","))
		}
		codeM := bFunc(method, f)
		v := CallMethod(f, method.Name)
		switch r := v.(type) {
		case string:
			fmt.Fprintf(w, "%s{return `%s`}\n", codeM, r)
		case int64:
			fmt.Fprintf(w, "%s{return int64(%v)}\n", codeM, r)
		case bool:
			fmt.Fprintf(w, "%s{return %v}\n", codeM, r)
		case time.Time:
			t, ok := v.(time.Time)
			if !ok {
				panic(ok)
			}
			code := fmt.Sprintf(`t1, _ := time.Parse(time.RFC3339,"%s");return t1`, t.Format(time.RFC3339))
			fmt.Fprintf(w, "%s{%v}\n", codeM, code)
		default:
			//fmt.Printf("Unknown type for fuction %s, unrecogniced type %v\n", method.Name, v)
			//return nil.
			fmt.Fprintf(w, "%s{return nil}\n", codeM)
			//going deeper.
			//WM(w, f)
		}
		//get return type.. if return type not value as defined return nil.
	}
}

func CallMethod(i interface{}, methodName string) interface{} {
	var ptr reflect.Value
	var value reflect.Value
	var finalMethod reflect.Value

	value = reflect.ValueOf(i)

	// if we start with a pointer, we need to get value pointed to
	// if we start with a value, we need to get a pointer to that value
	if value.Type().Kind() == reflect.Ptr {
		ptr = value
		value = ptr.Elem()
	} else {
		ptr = reflect.New(reflect.TypeOf(i))
		temp := ptr.Elem()
		temp.Set(value)
	}

	// check for method on value
	method := value.MethodByName(methodName)
	if method.IsValid() {
		finalMethod = method
	}
	// check for method on pointer
	method = ptr.MethodByName(methodName)
	if method.IsValid() {
		finalMethod = method
	}

	if finalMethod.IsValid() {
		return finalMethod.Call([]reflect.Value{})[0].Interface()
	}
	// return or panic, method not found of either type
	return ""
}

func main() {

	fs, _ := ioutil.ReadDir("./test/fs")
	for _, f := range fs {
		out, err := os.OpenFile("fs.go", os.O_RDWR|os.O_CREATE, 0666)
		if err != nil {
			log.Fatalf("error opening file: %v", err)
		}
		fmt.Println("Reading:", f.Name())

		w := bufio.NewWriter(os.Stdout)
		defer w.Flush()
		WM(w, f)
		continue
		out.Write([]byte("package main \n"))
		//out.Write([]byte("\nconst (\n" + filename + " = `"))
		f, _ := os.Open(f.Name())
		defer f.Close()
		encoder := base64.NewEncoder(base64.StdEncoding, out)
		defer encoder.Close()
		io.Copy(encoder, f)
		out.Write([]byte("`\n"))
		out.Write([]byte(")\n"))
		out.Close()
	}
}

var fs fileSystem = osFS{}

type fileSystem interface {
	Open(name string) (file, error)
	Stat(name string) (fileInfo, error)
}

type file interface {
	io.Closer
	io.Reader
	io.ReaderAt
	io.Seeker
	Stat() (os.FileInfo, error)
}

// osFS implements fileSystem using the local disk.
type osFS struct{}

func (osFS) Open(name string) (file, error)     { return os.Open(name) }
func (osFS) Stat(name string) (fileInfo, error) { return os.Stat(name) }

//fileInfo used for testing.
type fileInfo interface {
	Name() string       // base name of the file
	Size() int64        // length in bytes for regular files; system-dependent for others
	Mode() os.FileMode  // file mode bits
	ModTime() time.Time // modification time
	IsDir() bool        // abbreviation for Mode().IsDir()
	Sys() interface{}   // underlying data source (can return nil)
}
