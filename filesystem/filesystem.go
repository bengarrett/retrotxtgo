//Package filesystem to handle the opening and reading of text files
package filesystem

import (
	"fmt"
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
