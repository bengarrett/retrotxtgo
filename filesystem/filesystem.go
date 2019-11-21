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
	data := make([]byte, 100)
	count, err := file.Read(data)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("read %d bytes: %q\n", count, data[:count])
	return fmt.Sprintf("%s", data[:count])
}
