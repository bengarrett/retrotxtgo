//Package encoding to handle the opening and reading of text files
package encoding

//https://stackoverflow.com/questions/36111777/golang-how-to-read-a-text-file

import (
	"golang.org/x/net/html/charset"
)

//IsUTF8 blah
func IsUTF8(b []byte) bool {
	//todo limit read size?
	_, name, _ := charset.DetermineEncoding(b, "text/plain")
	if name == "utf-8" {
		return true
	}
	return false
}
