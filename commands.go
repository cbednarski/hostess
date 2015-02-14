package hostess

import (
	"fmt"
	"github.com/codegangsta/cli"
	"os"
	"strings"
)

func MaybeErrorln(c *cli.Context, message string) {
	if !c.Bool("q") {
		os.Stderr.WriteString(fmt.Sprintf("%s: %s\n", c.Command.Name, message))
	}
}

func MaybeError(c *cli.Context, message string) {
	MaybeErrorln(c, message)
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

func ShowEnabled(on bool) string {
	if on {
		return "(On)"
	} else {
		return "(Off)"
	}
}

func ShowHostname(hostname Hostname) string {
	return fmt.Sprintf("%s -> %s %s", hostname.Domain, hostname.Ip, ShowEnabled(hostname.Enabled))
}

func StrPadRight(s string, l int) string {
	return s + strings.Repeat(" ", l-len(s))
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
			MaybePrintln(c, fmt.Sprintf("Added %s", ShowHostname(hostname)))
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

func Ls(c *cli.Context) {
	hostfile := MaybeLoadHostFile(c)
	maxdomain := 0
	maxip := 0
	for _, hostname := range hostfile.Hosts {
		dlen := len(hostname.Domain)
		if dlen > maxdomain {
			maxdomain = dlen
		}
		ilen := len(hostname.Ip)
		if ilen > maxip {
			maxip = ilen
		}
	}

	for _, domain := range hostfile.ListDomains() {
		hostname := hostfile.Hosts[domain]
		fmt.Printf("%s -> %s %s\n",
			StrPadRight(hostname.Domain, maxdomain),
			StrPadRight(hostname.Ip, maxip),
			ShowEnabled(hostname.Enabled))
	}
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
