package prompt

import (
	"errors"
	"fmt"
	"os"
	"sort"

	"github.com/bengarrett/retrotxtgo/lib/config/internal/get"
	"github.com/bengarrett/retrotxtgo/lib/config/internal/upd"
	"github.com/bengarrett/retrotxtgo/lib/logs"
	"github.com/bengarrett/retrotxtgo/lib/str"
	"github.com/spf13/viper"
)

var ErrSaveType = errors.New("save value type is unsupported")

// Keys list all the available configuration setting names sorted alphabetically.
func Keys() []string {
	keys := make([]string, len(get.Reset()))
	i := 0
	for key := range get.Reset() {
		keys[i] = key
		i++
	}
	sort.Strings(keys)
	return keys
}

// save the value of the named setting to the configuration file.
func save(name string, setup bool, value interface{}) {
	if name == "" {
		logs.FatalSave(fmt.Errorf("save: %w", logs.ErrNameNil))
	}
	if !Validate(name) {
		logs.FatalSave(fmt.Errorf("save %q: %w", name, logs.ErrConfigName))
	}
	if skipSave(name, value) {
		fmt.Print(skipSet(setup))
		return
	}
	switch v := value.(type) {
	case string:
		if v == "-" {
			value = ""
		}
	default:
	}
	viper.Set(name, value)
	if err := upd.UpdateConfig("", false); err != nil {
		logs.FatalSave(err)
	}
	switch v := value.(type) {
	case string:
		if v == "" {
			fmt.Printf("  %s is now unused\n",
				str.ColSuc(name))
			if !setup {
				os.Exit(0)
			}
			return
		}
	default:
	}
	fmt.Printf("  %s is set to \"%v\"\n",
		str.ColSuc(name), value)
	if !setup {
		os.Exit(0)
	}
}

// skipSave returns true if the named value doesn't need updating.
func skipSave(name string, value interface{}) bool {
	switch v := value.(type) {
	case bool:
		if viper.Get(name).(bool) == v {
			return true
		}
	case string:
		if viper.Get(name).(string) == v {
			return true
		}
		if value.(string) == "" {
			return true
		}
	case uint:
		if viper.Get(name).(int) == int(v) {
			return true
		}
		if name == get.Serve && v == 0 {
			return true
		}
	default:
		logs.FatalSave(fmt.Errorf("name: %s, type: %T, %w", name, value, ErrSaveType))
	}
	return false
}

func skipSet(setup bool) string {
	if !setup {
		return ""
	}
	return str.ColSuc("\r  skipped setting")
}

// Validate the existence of the key in a list of settings.
func Validate(key string) (ok bool) {
	keys := Keys()
	// var i must be sorted in ascending order.
	if i := sort.SearchStrings(keys, key); i == len(keys) || keys[i] != key {
		return false
	}
	return true
}
