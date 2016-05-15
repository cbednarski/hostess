package hostess

import (
	"flag"
	"io/ioutil"
	"os"
	"testing"

	"github.com/codegangsta/cli"
	"github.com/stretchr/testify/assert"
)

func TestStrPadRight(t *testing.T) {
	assert.Equal(t, "", StrPadRight("", 0), "Zero-length no padding")
	assert.Equal(t, "          ", StrPadRight("", 10), "Zero-length 10 padding")
	assert.Equal(t, "string", StrPadRight("string", 0), "6-length 0 padding")
}

func captureOutput(f func()) string {
	rescueStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	f()

	w.Close()
	out, _ := ioutil.ReadAll(r)
	os.Stdout = rescueStdout

	return string(out[:])
}

func TestLs(t *testing.T) {
	os.Setenv("HOSTESS_PATH", "test-fixtures/hostfile1")
	defer os.Setenv("HOSTESS_PATH", "")

	app := cli.NewApp()

	context := cli.NewContext(app, &flag.FlagSet{}, nil)
	Ls(context)

	// Test on/off arguments functionality
	os.Setenv("HOSTESS_PATH", "test-fixtures/ls_on_off")
	set := flag.NewFlagSet("test", 0)
	set.Parse([]string{"list", "on"})
	context = cli.NewContext(app, set, nil)
	command := cli.Command{
		Name:        "list",
		Aliases:     []string{"ls"},
		Usage:       "Testing Ls",
		Description: "Testing Ls",
		Action:      Ls,
	}

	output := captureOutput(func() {
		command.Run(context)
	})

	assert.Equal(t, "chocolate.pie.example.com      -> fe:23b3:890e:342e::ef (On)\n", output)
}
