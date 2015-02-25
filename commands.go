package hostess

import (
	"fmt"
	"github.com/codegangsta/cli"
	"os"
	"strings"
)

// MaybeErrorln will print an error message unless -s is passed
func MaybeErrorln(c *cli.Context, message string) {
	if !c.Bool("s") {
		os.Stderr.WriteString(fmt.Sprintf("%s: %s\n", c.Command.Name, message))
	}
}

// MaybeError will print an error message unless -s is passed and then exit
func MaybeError(c *cli.Context, message string) {
	MaybeErrorln(c, message)
	os.Exit(1)
}

// MaybePrintln will print a message unless -q or -s is passed
func MaybePrintln(c *cli.Context, message string) {
	if !c.Bool("q") && !c.Bool("s") {
		fmt.Println(message)
	}
}

// MaybeLoadHostFile will try to load, parse, and return a Hostfile. If we
// encounter errors we will terminate, unless -f is passed.
func MaybeLoadHostFile(c *cli.Context) *Hostfile {
	hostsfile, errs := LoadHostFile()
	if len(errs) > 0 && !c.Bool("f") {
		for _, err := range errs {
			MaybeErrorln(c, err.Error())
		}
		MaybeError(c, "Errors while parsing hostsfile. Try using fix -f")
	}
	return hostsfile
}

// ShowEnabled turns a boolean into a string (On) or (Off)
func ShowEnabled(on bool) string {
	if on {
		return "(On)"
	} else {
		return "(Off)"
	}
}

// ShowHostname turns a Hostname into a string for display
func ShowHostname(hostname Hostname) string {
	return fmt.Sprintf("%s -> %s %s", hostname.Domain, hostname.Ip, ShowEnabled(hostname.Enabled))
}

// StrPadRight adds spaces to the right of a string until it reaches l length.
// If the input string is already that long, do nothing.
func StrPadRight(s string, l int) string {
	return s + strings.Repeat(" ", l-len(s))
}

func Add(c *cli.Context) {
	if len(c.Args()) != 2 {
		MaybeError(c, "expected <hostname> <ip>")
	}

	hostsfile := MaybeLoadHostFile(c)
	hostname := NewHostname(c.Args()[0], c.Args()[1], true)

	var err error
	if !hostsfile.Contains(hostname) {
		err = hostsfile.Add(hostname)
	}

	if err == nil {
		if c.Bool("n") {
			fmt.Println(hostsfile.Format())
		} else {
			MaybePrintln(c, fmt.Sprintf("Added %s", ShowHostname(*hostname)))
			hostsfile.Save()
		}
	} else {
		MaybeError(c, err.Error())
	}
}

func Del(c *cli.Context) {
	if len(c.Args()) != 1 {
		MaybeError(c, "expected <hostname>")
	}
	domain := c.Args()[0]
	hostsfile := MaybeLoadHostFile(c)

	found := hostsfile.ContainsDomain(domain)
	if found {
		hostsfile.Delete(domain)
		if c.Bool("n") {
			fmt.Println(hostsfile.Format())
		} else {
			MaybePrintln(c, fmt.Sprintf("Deleted %s", domain))
			hostsfile.Save()
		}
	} else {
		MaybePrintln(c, fmt.Sprintf("%s not found in %s", domain, GetHostsPath()))
	}
}

func Has(c *cli.Context) {
	if len(c.Args()) != 1 {
		MaybeError(c, "expected <hostname>")
	}
	domain := c.Args()[0]
	hostsfile := MaybeLoadHostFile(c)

	found := hostsfile.ContainsDomain(domain)
	if found {
		MaybePrintln(c, fmt.Sprintf("Found %s in %s", domain, GetHostsPath()))
	} else {
		MaybeError(c, fmt.Sprintf("%s not found in %s", domain, GetHostsPath()))
	}

}

func Off(c *cli.Context) {
	if len(c.Args()) != 1 {
		MaybeError(c, "expected <hostname>")
	}

}

func On(c *cli.Context) {
	if len(c.Args()) != 1 {
		MaybeError(c, "expected <hostname>")
	}
}

func Ls(c *cli.Context) {
	hostsfile := MaybeLoadHostFile(c)
	maxdomain := 0
	maxip := 0
	for _, hostname := range hostsfile.Hosts {
		dlen := len(hostname.Domain)
		if dlen > maxdomain {
			maxdomain = dlen
		}
		ilen := len(hostname.Ip)
		if ilen > maxip {
			maxip = ilen
		}
	}

	for _, domain := range hostsfile.ListDomains() {
		hostname := hostsfile.Hosts[domain]
		fmt.Printf("%s -> %s %s\n",
			StrPadRight(hostname.Domain, maxdomain),
			StrPadRight(hostname.Ip.String(), maxip),
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
	hostsfile := MaybeLoadHostFile(c)
	if c.Bool("n") {
		fmt.Println(hostsfile.Format())
	} else {
		hostsfile.Save()
	}
}

func Dump(c *cli.Context) {

}

func Apply(c *cli.Context) {

}
