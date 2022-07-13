package cmd_test

import "testing"

func Test_ConfigErr(t *testing.T) {
	t.Run("config invalid", func(t *testing.T) {
		const invalid = "zxcvbnnm"
		gotB, err := infoT.tester([]string{"--test", invalid})
		if err == nil {
			t.Errorf("using this invalid config command did not return an error: %s", invalid)
			t.Error(gotB)
		}
	})
}
