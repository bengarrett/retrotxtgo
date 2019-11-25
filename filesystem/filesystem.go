//Package filesystem to handle the opening and reading of text files
package filesystem

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
)

//Read a text file and display its content
func Read(path string) string {
	file, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	stat, _ := os.Stat(path)
	fmt.Printf("opening: %s\twhich is %v bytes\n", stat.Name(), stat.Size())
	data := make([]byte, stat.Size())
	count, err := file.Read(data)
	file.Close()
	if err != nil {
		log.Fatal(err)
	}
	//fmt.Printf("read %d bytes: %q\n", count, data[:count])
	return fmt.Sprintf("%s", data[:count])
}

//SeekBytes xx
func SeekBytes(name string, offset int64) ([]byte, error) {
	file, err := os.Open(name)
	defer file.Close()
	if err != nil {
		return nil, err
	}
	buffer := make([]byte, 1024)

	//offset=-128
	_, err = file.Seek(offset, 2)
	if err != nil {
		return nil, err
	}

	_, err = file.Read(buffer)
	if err != nil {
		return nil, err
	}

	return buffer, nil
}

//ReadAllBytes ooof
func ReadAllBytes(name string) ([]byte, error) {
	file, err := os.Open(name)
	fmt.Printf("%v", file)
	defer file.Close()
	if err != nil {
		return nil, err
	}
	buffer, err := ioutil.ReadAll(file)
	return buffer, nil
}
