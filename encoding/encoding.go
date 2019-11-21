//Package encoding to handle the opening and reading of text files
package encoding

//https://stackoverflow.com/questions/36111777/golang-how-to-read-a-text-file

import (
	"io/ioutil"
	"log"
	"os"

	"golang.org/x/net/html/charset"
)

//IsUTF8 blah
func IsUTF8(path string) bool {
	file, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	//todo limit read size?
	b, err := ioutil.ReadAll(file)
	if err != nil {
		log.Fatal(err)
	}
	_, name, _ := charset.DetermineEncoding(b, "text/plain")
	if name == "utf-8" {
		return true
	}
	return false
}
