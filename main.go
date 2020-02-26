package main

import (
	"errors"
	"flag"
	"fmt"
	"os"

	"github.com/cbednarski/hostess/commands"
	"github.com/cbednarski/hostess/hostess"
)

const help = `An idempotent tool for managing %s

Commands

    fmt                  Reformat the hosts file

    add <hostname> <ip>  Add or overwrite a hosts entry
    rm <hostname>        Remote a hosts entry
    on <hostname>        Enable a hosts entry
    off <hostname>       Disable a hosts entry

    ls                   List hosts entries
    has                  Exit 0 if entry present in hosts file, 1 if not

    dump                 Export hosts entries as JSON
    apply                Import hosts entries from JSON

    All commands that change the hosts file will implicitly reformat it.

Flags

    -n will preview changes but not rewrite your hosts file
    -f force changes even if there are errors parsing the hosts file
    -4 limit changes to IPv4 entries
    -6 limit changes to IPv6 entries

Configuration

    HOSTESS_FMT may be set to unix or windows to force that platform's syntax
    HOSTESS_PATH may be set to point to a file other than the platform default

About

    Copyright 2015-2020 Chris Bednarski <chris@cbednarski.com>; MIT Licensed 
    Portions Copyright the Go authors, licensed under BSD-style license
    Bugs and updates via https://github.com/cbednarski/hostess
`

var (
	Version           = "dev"
	ErrInvalidCommand = errors.New("invalid command")
)

func ExitWithError(err error) {
	if err != nil {
		os.Stderr.WriteString(err.Error())
		os.Stderr.WriteString("\n")
		os.Exit(1)
	}
}

func CommandUsage(command string) error {
	return fmt.Errorf("Usage: %s %s <hostname>", os.Args[0], command)
}

func wrappedMain() error {
	cli := flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	ipv4 := cli.Bool("4",false, "IPv4")
	ipv6 := cli.Bool("6",false, "IPv6")
	preview := cli.Bool("n",false, "preview")
	force := cli.Bool("f",false, "force")
	cli.Usage = func() {
		fmt.Printf(help, hostess.GetHostsPath())
	}

	if err := cli.Parse(os.Args[2:]); err != nil {
		return err
	}

	options := &commands.Options{
		IPVersion: 0,
		Preview:   *preview,
		Force:     *force,
	}
	if *ipv4 {
		options.IPVersion = options.IPVersion|commands.IPv4
	}
	if *ipv6 {
		options.IPVersion = options.IPVersion|commands.IPv6
	}

	if len(os.Args) == 1 {
		cli.Usage()
		return nil
	}

	command := os.Args[1]
	switch command {

	case "-v", "--version", "version":
		fmt.Println(Version)
		return nil

	case "-h", "--help", "help":
		cli.Usage()
		return nil

	case "add":
		if len(cli.Args()) != 2 {
			return fmt.Errorf("Usage: %s add <hostname> <ip>", cli.Name())
		}
		return commands.Add(options, cli.Arg(0), cli.Arg(1))

	case "rm":
		if cli.Arg(0) == "" {
			return CommandUsage(command)
		}
		return commands.Remove(options, cli.Arg(0))

	case "has":
		if cli.Arg(0) == "" {
			return CommandUsage(command)
		}
		return commands.Has(options, cli.Arg(0))

	case "on":
		if cli.Arg(0) == "" {
			return CommandUsage(command)
		}
		return commands.Enable(options, cli.Arg(0))

	case "off":
		if cli.Arg(0) == "" {
			return CommandUsage(command)
		}
		return commands.Disable(options, cli.Arg(0))

	case "fmt":
		return commands.Format(options)

	case "ls":
		return commands.List(options)

	case "dump":
		return commands.Dump(options)

	case "apply":
		if cli.Arg(0) == "" {
			return fmt.Errorf("Usage: %s apply <filename>", os.Args[0])
		}
		return commands.Apply(options, cli.Arg(0))

	default:
		return ErrInvalidCommand
	}
}

func main() {
	ExitWithError(wrappedMain())
}
