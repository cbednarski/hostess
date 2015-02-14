package hostess

import (
	// "errors"
	"fmt"
	"github.com/codegangsta/cli"
	"os"
)

func MaybeErrorln(c *cli.Context, message string) {
	if !c.Bool("q") {
		fmt.Printf("%s: %s\n", c.Command.Name, message)
	}
}

func MaybeError(c *cli.Context, message string) {
	if !c.Bool("q") {
		fmt.Printf("%s: %s\n", c.Command.Name, message)
	}
	os.Exit(1)
}

func MaybePrintln(c *cli.Context, message string) {
	if !c.Bool("s") {
		fmt.Println(message)
	}
}

func MaybeLoadHostFile(c *cli.Context) *Hostfile {
	hostsfile, errs := LoadHostFile()
	if len(errs) > 0 && !c.Bool("f") {
		for _, err := range errs {
			MaybeErrorln(c, err.Error())
		}
		MaybeError(c, "Errors while parsing hostsfile")
	}
	return hostsfile
}

func SprintEnabled(on bool) string {
	if on {
		return "(On)"
	} else {
		return "(Off)"
	}
}

func Add(c *cli.Context) {
	if len(c.Args()) != 2 {
		MaybeError(c, "expected <hostname> <ip>")
	}

	hostsfile := MaybeLoadHostFile(c)
	hostname := Hostname{c.Args()[0], c.Args()[1], true}
	err := hostsfile.Add(hostname)
	if err == nil {
		if c.Bool("n") {
			fmt.Println(hostsfile.Format())
		} else {
			MaybePrintln(c, fmt.Sprintf("Added %s -> %s %s", hostname.Domain, hostname.Ip, SprintEnabled(hostname.Enabled)))
			hostsfile.Save()
		}
	} else {
		MaybeError(c, err.Error())
	}
}

func Del(c *cli.Context) error {
	return nil
}

func Has(c *cli.Context) error {
	return nil
}

func Off(c *cli.Context) error {
	return nil
}

func On(c *cli.Context) error {
	return nil
}

func Ls(c *cli.Context) error {
	return nil
}

const fix_help = `Programmatically rewrite your hostsfile.

Domains pointing to the same IP will be consolidated, sorted, and extra
whitespace and comments will be removed.

   hostess fix      Rewrite the hostsfile
   hostess fix -n   Show the new hostsfile. Don't write it
`

func Fix(c *cli.Context) {
	hostfile := MaybeLoadHostFile(c)
	if c.Bool("n") {
		fmt.Println(hostfile.Format())
	} else {
		hostfile.Save()
	}
}

func Dump(c *cli.Context) error {
	return nil
}

func Apply(c *cli.Context) error {
	return nil
}
