package main

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
}

// PrintErrLn will print to stderr followed by a newline
func PrintErrLn(err error) {
	os.Stderr.WriteString(fmt.Sprintf("%s\n", err))
}

// LoadHostfile will try to load, parse, and return a Hostfile. If we
// encounter errors we will terminate, unless -f is passed.
func LoadHostfile() (*hostess.Hostfile, error) {
	hosts, errs := hostess.LoadHostfile()

	if len(errs) > 0 {
		for _, err := range errs {
			PrintErrLn(err)
		}
		return nil, errors.New("Errors while parsing hostsfile. Please fix any dupes or conflicts and try again.")
	}

	return hosts, nil
}

// SaveOrPreview will display or write the Hostfile
func SaveOrPreview(options *Options, hostfile *hostess.Hostfile) error {
	// If -n is passed, no-op and output the resultant hosts file to stdout.
	// Otherwise it's for real and we're going to write it.
	if options.Preview {
		fmt.Printf("%s", hostfile.Format())
		return nil
	}

	if err := hostfile.Save(); err != nil {
		return fmt.Errorf("Unable to write to %s. Maybe you need to sudo? (error: %s)", hostess.GetHostsPath(), err)
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
	hostsfile, err := LoadHostfile()

	newHostname, err := hostess.NewHostname(hostname, ip, true)
	if err != nil {
		return err
	}

	replaced := hostsfile.Hosts.ContainsDomain(newHostname.Domain)
	// Note that Add() may return an error, but they are informational only. We
	// don't actually care what the error is -- we just want to add the
	// hostname and save the file. This way the behavior is idempotent.
	hostsfile.Hosts.Add(newHostname)

	// If the user passes -n then we'll Add and show the new hosts file, but
	// not save it.
	if err := SaveOrPreview(options, hostsfile); err != nil {
		return err
	}
	// We'll give a little bit of information about whether we added or
	// updated, but if the user wants to know they can use has or ls to
	// show the file before they run the operation. Maybe later we can add
	// a verbose flag to show more information.
	if replaced {
		fmt.Printf("Updated %s\n", newHostname.FormatHuman())
	} else {
		fmt.Printf("Added %s\n", newHostname.FormatHuman())
	}

	return nil
}

// Remove command removes any hostname(s) matching <domain> from the hosts file
func Remove(options *Options, hostname string) error {
	hostsfile, err := LoadHostfile()
	if err != nil {
		return err
	}

	found := hostsfile.Hosts.ContainsDomain(hostname)
	if !found {
		fmt.Printf("%s not found in %s", hostname, hostess.GetHostsPath())
	}

	hostsfile.Hosts.RemoveDomain(hostname)
	if err := SaveOrPreview(options, hostsfile); err != nil {
		return err
	}
	fmt.Printf("Deleted %s\n", hostname)

	return nil
}

// Has command indicates whether a hostname is present in the hosts file
func Has(options *Options, hostname string) error {
	hostsfile, err := LoadHostfile()
	if err != nil {
		return err
	}

	found := hostsfile.Hosts.ContainsDomain(hostname)
	if found {
		fmt.Printf("Found %s in %s\n", hostname, hostess.GetHostsPath())
	} else {
		fmt.Printf("%s not found in %s\n", hostname, hostess.GetHostsPath())
		// Exit 1 for bash scripts to use this as a check
		os.Exit(1)
	}
	return nil
}

func Enable(options *Options, hostname string) error {
	hostsfile, err := LoadHostfile()
	if err != nil {
		return err
	}

	if err := hostsfile.Hosts.Enable(hostname); err != nil {
		return err
	}

	if err := SaveOrPreview(options, hostsfile); err != nil {
		return err
	}

	fmt.Printf("Enabled %s\n", hostname)

	return nil
}

func Disable(options *Options, hostname string) error {
	hostsfile, err := LoadHostfile()
	if err != nil {
		return err
	}

	if err := hostsfile.Hosts.Disable(hostname); err != nil {
		if err == hostess.ErrHostnameNotFound {
			// If the hostname does not exist then we have still achieved the
			// desired result, so we will not exit with an error here. We'll
			// handle the error by displaying it to the user.
			PrintErrLn(err)
			return nil
		}
		return err
	}

	if err := SaveOrPreview(options, hostsfile); err != nil {
		return err
	}

	fmt.Printf("Disabled %s\n", hostname)

	return nil
}

// List command shows a list of hostnames in the hosts file
func List(options *Options) error {
	hostsfile, err := LoadHostfile()
	if err != nil {
		return err
	}

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

	return nil
}

// Format command removes duplicates and conflicts from the hosts file
func Format(options *Options) error {
	hostsfile, err := LoadHostfile()
	if err != nil {
		return err
	}

	if bytes.Equal(hostsfile.GetData(), hostsfile.Format()) {
		fmt.Printf("%s is already formatted and contains no dupes or conflicts; nothing to do\n", hostess.GetHostsPath())
		return nil
	}

	return SaveOrPreview(options, hostsfile)
}

// Dump command outputs hosts file contents as JSON
func Dump(options *Options) error {
	hostsfile, err := LoadHostfile()
	if err != nil {
		return err
	}

	jsonbytes, err := hostsfile.Hosts.Dump()
	if err != nil {
		return err
	}

	fmt.Println(fmt.Sprintf("%s", jsonbytes))
	return nil
}

// Apply command adds hostnames to the hosts file from JSON
func Apply(options *Options, filename string) error {
	jsonbytes, err := ioutil.ReadFile(filename)
	if err != nil {
		return fmt.Errorf("Unable to read JSON from %s: %s", filename, err)
	}

	hostfile, err := LoadHostfile()
	if err != nil {
		return err
	}


	if err := hostfile.Hosts.Apply(jsonbytes); err != nil {
		return fmt.Errorf("Error applying changes to hosts file: %s", err)
	}

	if err := SaveOrPreview(options, hostfile); err != nil {
		return err
	}

	fmt.Printf("%s applied\n", filename)
	return nil
}
