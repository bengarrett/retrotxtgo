package samples

import "log"

// EncodeASCII encodes the content of ?.asc to base64 for use as LogoASCII.
func EncodeASCII() (result string) {
	d, err := ReadLine(asciiFile(), "dos")
	if err != nil {
		log.Fatal(err)
	}
	return Base64Encode(d)
}

// EncodeANSI encodes the content of ZII-RTXT.ans to base64 for use as LogoANSI.
func EncodeANSI() (result string) {
	d, err := ReadLine(ansiFile(), "dos")
	if err != nil {
		log.Fatal(err)
	}
	return Base64Encode(d)
}
