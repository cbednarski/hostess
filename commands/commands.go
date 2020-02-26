package commands

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/cbednarski/hostess/hostess"
)

const (
	IPv4 = 1 << iota
	IPv6 = 1 << iota
)

type Options struct {
	IPVersion int
	Preview bool
	Force bool
}

// ErrCantWriteHostFile indicates that we are unable to write to the hosts file
var ErrCantWriteHostFile = ""

// ErrorLn will print an error message unless -s is passed
func ErrorLn(message string) {
	os.Stderr.WriteString(fmt.Sprintf("%s\n", message))
}

// MaybePrintln will print a message unless -q or -s is passed
func MaybePrintln(options *Options, message string) {
	if !AnyBool(c, "q") && !AnyBool(c, "s") {
		fmt.Println(message)
	}
}

// LoadHostfile will try to load, parse, and return a Hostfile. If we
// encounter errors we will terminate, unless -f is passed.
func LoadHostfile(options *Options) (*hostess.Hostfile, error) {
	hosts, errs := hostess.LoadHostfile()

	if len(errs) > 0 && !options.Force {
		for _, err := range errs {
			ErrorLn(err.Error())
		}
		return nil, errors.New("Errors while parsing hostsfile. Try hostess fmt")
	}

	return hosts, nil
}

// SaveOrPreview will display or write the Hostfile
func SaveOrPreview(options *Options, hostfile *hostess.Hostfile) error {
	// If -n is passed, no-op and output the resultant hosts file to stdout.
	// Otherwise it's for real and we're going to write it.
	if options.Preview {
		fmt.Printf("%s", hostfile.Format())
	} else {
		if err := hostfile.Save(); err != nil {
			return fmt.Errorf("Unable to write to %s. Maybe you need to sudo? (error: %s)", hostess.GetHostsPath(), err)
		}
	}
	return nil
}

// StrPadRight adds spaces to the right of a string until it reaches length.
// If the input string is already that long, do nothing.
func StrPadRight(input string, length int) string {
	minimum := len(input)
	if length <= minimum {
		return input
	}
	return input + strings.Repeat(" ", length-minimum)
}

// Add command parses <hostname> <ip> and adds or updates a hostname in the
// hosts file
func Add(options *Options, hostname, ip string) error {
	hostsfile, err := LoadHostfile(options)

	newHostname, err := hostess.NewHostname(hostname, ip, true)
	if err != nil {
		MaybeError(c, fmt.Sprintf("Failed to parse hosts entry: %s", err))
	}
	// If the command is aff instead of add then the entry should be disabled
	if c.Command.Name == "aff" {
		newHostname.Enabled = false
	}

	replace := hostsfile.Hosts.ContainsDomain(newHostname.Domain)
	// Note that Add() may return an error, but they are informational only. We
	// don't actually care what the error is -- we just want to add the
	// hostname and save the file. This way the behavior is idempotent.
	hostsfile.Hosts.Add(newHostname)

	// If the user passes -n then we'll Add and show the new hosts file, but
	// not save it.
	if c.Bool("n") || AnyBool(c, "n") {
		fmt.Printf("%s", hostsfile.Format())
	} else {
		SaveOrPreview(c, hostsfile)
		// We'll give a little bit of information about whether we added or
		// updated, but if the user wants to know they can use has or ls to
		// show the file before they run the operation. Maybe later we can add
		// a verbose flag to show more information.
		if replace {
			fmt.Printf("Updated %s\n", newHostname.FormatHuman())
		} else {
			fmt.Printf("Added %s\n", newHostname.FormatHuman())
		}
	}
}

// Remove command removes any hostname(s) matching <domain> from the hosts file
func Remove(options *Options, hostname string) error {
	hostsfile := LoadHostfile(c)

	found := hostsfile.Hosts.ContainsDomain(hostname)
	if found {
		hostsfile.Hosts.RemoveDomain(hostname)
		if AnyBool(c, "n") {
			fmt.Printf("%s", hostsfile.Format())
		} else {
			SaveOrPreview(c, hostsfile)
			MaybePrintln(c, fmt.Sprintf("Deleted %s", hostname))
		}
	} else {
		MaybePrintln(c, fmt.Sprintf("%s not found in %s", hostname, hostess.GetHostsPath()))
	}
}

// Has command indicates whether a hostname is present in the hosts file
func Has(options *Options, hostname string) error {
	if len(c.Args()) != 1 {
		MaybeError(c, "expected <hostname>")
	}
	domain := c.Args()[0]
	hostsfile := LoadHostfile(c)

	found := hostsfile.Hosts.ContainsDomain(domain)
	if found {
		MaybePrintln(c, fmt.Sprintf("Found %s in %s", domain, hostess.GetHostsPath()))
	} else {
		MaybeError(c, fmt.Sprintf("%s not found in %s", domain, hostess.GetHostsPath()))
	}
}

// OnOff enables (uncomments) or disables (comments) the specified hostname in
// the hosts file. Exits code 1 if the hostname is missing.
func OnOff(options *Options, hostname string) error {
	hostsfile := LoadHostfile(c)

	// Switch on / off commands
	success := false
	if c.Command.Name == "on" {
		success = hostsfile.Hosts.Enable(domain)
	} else {
		success = hostsfile.Hosts.Disable(domain)
	}

	if success {
		SaveOrPreview(c, hostsfile)
		if c.Command.Name == "on" {
			MaybePrintln(c, fmt.Sprintf("Enabled %s", domain))
		} else {
			MaybePrintln(c, fmt.Sprintf("Disabled %s", domain))
		}
	} else {
		MaybeError(c, fmt.Sprintf("%s not found in %s", domain, hostess.GetHostsPath()))
	}
}

func Enable(options *Options, hostname string) error {

}

func Disable(options *Options, hostname string) error {

}

// List command shows a list of hostnames in the hosts file
func List(options *Options) error {
	hostsfile := AlwaysLoadHostFile(c)
	widestHostname := 0
	widestIP := 0
	for _, hostname := range hostsfile.Hosts {
		dlen := len(hostname.Domain)
		if dlen > widestHostname {
			widestHostname = dlen
		}
		ilen := len(hostname.IP)
		if ilen > widestIP {
			widestIP = ilen
		}
	}

	for _, hostname := range hostsfile.Hosts {
		fmt.Printf("%s -> %s %s\n",
			StrPadRight(hostname.Domain, widestHostname),
			StrPadRight(hostname.IP.String(), widestIP),
			hostname.FormatEnabled())
	}
}

// Format command removes duplicates and conflicts from the hosts file
func Format(options *Options) error {
	hostsfile := AlwaysLoadHostFile(c)
	if bytes.Equal(hostsfile.GetData(), hostsfile.Format()) {
		MaybePrintln(c, fmt.Sprintf("%s is already formatted and contains no dupes or conflicts; nothing to do", hostess.GetHostsPath()))
		os.Exit(0)
	}
	SaveOrPreview(c, hostsfile)
}

// Dump command outputs hosts file contents as JSON
func Dump(options *Options) error {
	hostsfile := AlwaysLoadHostFile(c)
	jsonbytes, err := hostsfile.Hosts.Dump()
	if err != nil {
		MaybeError(c, err.Error())
	}
	fmt.Println(fmt.Sprintf("%s", jsonbytes))
}

// Apply command adds hostnames to the hosts file from JSON
func Apply(options *Options, filename string) error {
	jsonbytes, err := ioutil.ReadFile(filename)
	if err != nil {
		MaybeError(c, fmt.Sprintf("Unable to read %s: %s", filename, err))
	}

	hostfile := AlwaysLoadHostFile(c)
	err = hostfile.Hosts.Apply(jsonbytes)
	if err != nil {
		MaybeError(c, fmt.Sprintf("Error applying changes to hosts file: %s", err))
	}

	SaveOrPreview(c, hostfile)
	MaybePrintln(c, fmt.Sprintf("%s applied", filename))
}
