package logs

import (
	"fmt"
	"os"

	"github.com/alecthomas/chroma/quick"
	"github.com/alecthomas/chroma/styles"
	"github.com/gookit/color"
	"gopkg.in/yaml.v2"
)

// YamlExample is a YAML example
type YamlExample struct {
	Style struct {
		Name    string `yaml:"name"`
		Count   int    `yaml:"count"`
		Default bool   `yaml:"default"`
	}
}

func (s YamlExample) String(flag string) {
	fmt.Println()
	out, _ := yaml.Marshal(s)
	quick.Highlight(os.Stdout, string(out), "yaml", "terminal256", s.Style.Name)
	if flag != "" {
		fmt.Println(color.Secondary.Sprintf("%s=%q", flag, s.Style.Name))
	}
}

// YamlStyles prints out a list of available YAML color styles.
func YamlStyles(cmd string) {
	for i, s := range styles.Names() {
		var styles YamlExample
		styles.Style.Name = s
		styles.Style.Count = i
		if s == "monokai" {
			styles.Style.Default = true
		}
		styles.String(cmd)
	}
}
