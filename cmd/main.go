package main

import (
	"github.com/cbednarski/hostess"
	"github.com/codegangsta/cli"
	"os"
)

func getCommand() string {
	return os.Args[1]
}

func getArgs() []string {
	return os.Args[2:]
}

const help = `an idempotent tool for managing /etc/hosts

 * WARNING: This program is BETA and not all commands are implemented.

 * Commands will exit 0 or 1 in a sensible way so you can use the exit code for
   bash and make scripting. Add -h to any command to learn more about it.

 * You can specify the HOSTESS_FILE environment variable to operate on a
   file other than /etc/hosts

 * To preview a hostess-managed hostsfile run ` + "`" + `hostess fix -n` + "`" + `

 * Report bugs and feedback at https://github.com/cbednarski/hostess`

func main() {
	app := cli.NewApp()
	app.Name = "hostess"
	app.Usage = help
	app.Version = "0.1.0"

	app.Flags = []cli.Flag{
		cli.BoolFlag{
			Name:  "f",
			Usage: "Force write to the hostsfile even if there are errors or conflicts",
		},
		cli.BoolFlag{
			Name:  "n",
			Usage: "No-op. Show changes but don't write them.",
		},
		cli.BoolFlag{
			Name:  "q",
			Usage: "Quiet. Supress error messages",
		},
		cli.BoolFlag{
			Name:  "s",
			Usage: "Silent. Supress success messages (implies -q)",
		},
	}

	app.Commands = []cli.Command{
		{
			Name:   "add",
			Usage:  "Add (or update) a hosts entry",
			Action: hostess.Add,
			Flags:  app.Flags,
		},
		{
			Name:   "del",
			Usage:  "Delete a hosts entry",
			Action: hostess.Del,
			Flags:  app.Flags,
		},
		{
			Name:   "has",
			Usage:  "Exit 0 if entry exists, 1 if not",
			Action: hostess.Has,
			Flags:  app.Flags,
		},
		{
			Name:   "off",
			Usage:  "Disable a hosts entry (don't delete it)",
			Action: hostess.Off,
			Flags:  app.Flags,
		},
		{
			Name:   "on",
			Usage:  "Enable a hosts entry (if if exists)",
			Action: hostess.On,
			Flags:  app.Flags,
		},
		{
			Name:   "ls, list",
			Usage:  "List entries in the hosts file",
			Action: hostess.Ls,
			Flags:  app.Flags,
		},
		{
			Name:   "fix",
			Usage:  "Reformat the hosts file based on hostess' rules",
			Action: hostess.Fix,
			Flags:  app.Flags,
		},
		{
			Name:   "dump",
			Usage:  "Dump the hosts file as JSON",
			Action: hostess.Dump,
			Flags:  app.Flags,
		},
		{
			Name:   "apply",
			Usage:  "Apply a JSON hosts dict to your hosts file",
			Action: hostess.Apply,
			Flags:  app.Flags,
		},
	}

	app.Run(os.Args)
	os.Exit(0)
}
