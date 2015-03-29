package hostess

import (
	"bytes"
	"fmt"
	"github.com/codegangsta/cli"
	"os"
	"strings"
)

// ErrCantWriteHostFile indicates that we are unable to write to the hosts file
var ErrCantWriteHostFile = fmt.Errorf(
	"Unable to write to %s. Maybe you need to sudo?", GetHostsPath())

// MaybeErrorln will print an error message unless -s is passed
func MaybeErrorln(c *cli.Context, message string) {
	if !c.Bool("s") {
		os.Stderr.WriteString(fmt.Sprintf("%s\n", message))
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
	hostsfile, errs := LoadHostfile()
	if len(errs) > 0 && !c.Bool("f") {
		for _, err := range errs {
			MaybeErrorln(c, err.Error())
		}
		MaybeError(c, "Errors while parsing hostsfile. Try hostess fix")
	}
	return hostsfile
}

// AlwaysLoadHostFile will load, parse, and return a Hostfile. If we encouter
// errors they will be printed to the terminal, but we'll try to continue.
func AlwaysLoadHostFile(c *cli.Context) *Hostfile {
	hostsfile, errs := LoadHostfile()
	if len(errs) > 0 {
		for _, err := range errs {
			MaybeErrorln(c, err.Error())
		}
	}
	return hostsfile
}

// StrPadRight adds spaces to the right of a string until it reaches l length.
// If the input string is already that long, do nothing.
func StrPadRight(s string, l int) string {
	return s + strings.Repeat(" ", l-len(s))
}

// Add command parses <hostname> <ip> and adds or updates a hostname in the
// hosts file
func Add(c *cli.Context) {
	if len(c.Args()) != 2 {
		MaybeError(c, "expected <hostname> <ip>")
	}

	hostsfile := MaybeLoadHostFile(c)
	hostname := NewHostname(c.Args()[0], c.Args()[1], true)

	var err error
	if !hostsfile.Hosts.Contains(hostname) {
		err = hostsfile.Hosts.Add(hostname)
	}

	if err == nil {
		if c.Bool("n") {
			fmt.Println(hostsfile.Format())
		} else {
			MaybePrintln(c, fmt.Sprintf("Added %s", hostname.FormatHuman()))
			hostsfile.Save()
		}
	} else {
		MaybeError(c, err.Error())
	}
}

// Del command removes any hostname(s) matching <domain> from the hosts file
func Del(c *cli.Context) {
	if len(c.Args()) != 1 {
		MaybeError(c, "expected <hostname>")
	}
	domain := c.Args()[0]
	hostsfile := MaybeLoadHostFile(c)

	found := hostsfile.Hosts.ContainsDomain(domain)
	if found {
		hostsfile.Hosts.RemoveDomain(domain)
		if c.Bool("n") {
			fmt.Println(hostsfile.Format())
		} else {
			err := hostsfile.Save()
			if err != nil {
				MaybeErrorln(c, ErrCantWriteHostFile.Error())
				os.Exit(1)
			}
			MaybePrintln(c, fmt.Sprintf("Deleted %s", domain))
		}
	} else {
		MaybePrintln(c, fmt.Sprintf("%s not found in %s", domain, GetHostsPath()))
	}
}

// Has command indicates whether a hostname is present in the hosts file
func Has(c *cli.Context) {
	if len(c.Args()) != 1 {
		MaybeError(c, "expected <hostname>")
	}
	domain := c.Args()[0]
	hostsfile := MaybeLoadHostFile(c)

	found := hostsfile.Hosts.ContainsDomain(domain)
	if found {
		MaybePrintln(c, fmt.Sprintf("Found %s in %s", domain, GetHostsPath()))
	} else {
		MaybeError(c, fmt.Sprintf("%s not found in %s", domain, GetHostsPath()))
	}

}

// Off command disables (comments) the specified hostname in the hosts file
func Off(c *cli.Context) {
	if len(c.Args()) != 1 {
		MaybeError(c, "expected <hostname>")
	}

}

// On command enabled (uncomments) the specified hostname in the hosts file
func On(c *cli.Context) {
	if len(c.Args()) != 1 {
		MaybeError(c, "expected <hostname>")
	}
}

// Ls command shows a list of hostnames in the hosts file
func Ls(c *cli.Context) {
	hostsfile := AlwaysLoadHostFile(c)
	maxdomain := 0
	maxip := 0
	for _, hostname := range hostsfile.Hosts {
		dlen := len(hostname.Domain)
		if dlen > maxdomain {
			maxdomain = dlen
		}
		ilen := len(hostname.IP)
		if ilen > maxip {
			maxip = ilen
		}
	}

	for _, hostname := range hostsfile.Hosts {
		fmt.Printf("%s -> %s %s\n",
			StrPadRight(hostname.Domain, maxdomain),
			StrPadRight(hostname.IP.String(), maxip),
			hostname.FormatEnabled())
	}
}

const fixHelp = `Programmatically rewrite your hostsfile.

Domains pointing to the same IP will be consolidated, sorted, and extra
whitespace and comments will be removed.

   hostess fix      Rewrite the hostsfile
   hostess fix -n   Show the new hostsfile. Don't write it
`

// Fix command removes duplicates and conflicts from the hosts file
func Fix(c *cli.Context) {
	hostsfile := AlwaysLoadHostFile(c)
	if bytes.Equal(hostsfile.GetData(), hostsfile.Format()) {
		MaybePrintln(c, fmt.Sprintf("%s is already formatted and contains no dupes or conflicts; nothing to do", GetHostsPath()))
		os.Exit(0)
	}
	if c.Bool("n") {
		fmt.Printf("%s", hostsfile.Format())
	} else {
		err := hostsfile.Save()
		if err != nil {
			MaybeErrorln(c, ErrCantWriteHostFile.Error())
		}
		MaybePrintln(c, fmt.Sprintf("Fixed %s", GetHostsPath()))
	}
}

// Dump command outputs hosts file contents as JSON
func Dump(c *cli.Context) {

}

// Apply command adds hostnames to the hosts file from JSON
func Apply(c *cli.Context) {

}
