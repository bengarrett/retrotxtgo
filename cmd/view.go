/*
Copyright © 2019 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"bufio"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"golang.org/x/text/encoding/charmap"
)

// viewCmd represents the view command
var viewCmd = &cobra.Command{
	Use:   "view",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		//		name := "textfiles/ZII-RTXT.ans"
		name := "/home/ben/hx_back.ans"
		fmt.Println("view called", name)

		file, err := os.Open(name)

		// todo reverse scan of file looking for SAUCE00 and COMNTT

		//defer close(name)
		Check(ErrorFmt{"file open", name, err})
		defer file.Close()

		decode := charmap.CodePage437.NewDecoder()

		scanner := bufio.NewScanner(file)

		// todo scan for unique color codes like 24-bit color

		// todo scan for new lines or character counts and hard-code the width

		for scanner.Scan() {
			out, _ := decode.Bytes(scanner.Bytes())
			fmt.Printf("%s%s\n", string(out), "\033[0m")
		}

		// buf := make([]byte, 32*1024) //

		// for {
		// 	n, err := file.Read(buf)
		// 	if n > 0 {
		// 		fmt.Print(buf[:n])
		// 	}
		// 	if err == io.EOF {
		// 		break
		// 	}
		// 	if err != nil {
		// 		log.Printf("read %d bytes: %v", n, err)
		// 	}
		// }

		// r := bufio.NewReaderSize(f, 100)
		// b, _ := r.
		// //b, _ := r.Peek(100)

		// decode := charmap.CodePage437.NewDecoder()
		// out, _ := decode.Bytes(b)
		// fmt.Println(out)

		//func DetermineEncodingFromReader(r io.Reader) (e encoding.Encoding, name string, certain bool, err error) {
		// b, err := bufio.NewReader(f).Peek(1024)
		// if err != nil {
		// 	return
		// }

		// e, name, certain := charset.DetermineEncoding(b, "")
		// fmt.Println(e, name, certain)
		// //}

		//encoder := x.text.encoding.NewEncoder()
		//scanner := bufio.NewScanner(f)
		// var buf bytes.Buffer
		// for scanner.Scan() {
		// 	buf.Write(scanner.Bytes())
		// }
		// fmt.Println(buf.String())

	},
}

func init() {
	rootCmd.AddCommand(viewCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// viewCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// viewCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
