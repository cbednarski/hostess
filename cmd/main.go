package main

import (
	"fmt"
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

const help = `Hostess: an idempotent tool for managing /etc/hosts

Commands will exit 0 or 1 in a sensible way so you can use the exit code for
bash and make scripting. Add -h to any command to learn more about it.

To preview a hostess-managed hostsfile run ` + "`" + `hostess fix -n` + "`" + `

WARNING: This program is BETA and not all commands are implemented.

Commands:

   add     Add (or update) a hosts entry
   del     Delete a hosts entry
   has     Exit 0 if entry exists, 1 if not
   off     Disable a hosts entry (don't delete it)
   on      Enable a hosts entry (if if exists)
   ls      List entries in the hosts file
   list    Alias for ls
   fix     Reformat the hosts file based on hostess' rules
   dump    Dump the hosts file as JSON
   apply   Apply a JSON hosts dict to your hosts file

Options:

   -f   Force write to the hostsfile even if there are errors or conflicts
   -n   No-op. Show changes but don't write them.
   -q   Supress error messages
   -s   Supress success messages (implies -q)
   -h   Show help for a command

Note: You can specify the HOSTESS_FILE environment variable to operate on a
file other than /etc/hosts

Report bugs at https://github.com/cbednarski/hostess
`

func main() {
	app := cli.NewApp()
	app.Name = "hostess"
	app.Usage = help
	app.Version = "0.1.0"

	var err error = nil

	app.Flags = []cli.Flag{
		cli.BoolFlag{
			Name:  "f",
			Usage: "Force",
		},
		cli.BoolFlag{
			Name:  "n",
			Usage: "Noop",
		},
		cli.BoolFlag{
			Name:  "q",
			Usage: "Quiet",
		},
		cli.BoolFlag{
			Name:  "s",
			Usage: "Silent",
		},
	}

	app.Commands = []cli.Command{
		{
			Name: "add",
			// Usage: "add a task to the list",
			Action: hostess.Add,
			Flags:  app.Flags,
		},
		{
			Name: "fix",
			// Usage: "add a task to the list",
			Action: hostess.Fix,
			Flags:  app.Flags,
		},
		{
			Name: "ls",
			// Usage: "add a task to the list",
			Action: hostess.Ls,
			Flags:  app.Flags,
		},
	}

	// switch flag.Arg(0) {
	// case "add":
	// 	err = hostess.Add(args, flags)
	// case "del":
	// 	err = hostess.Del(args, flags)
	// case "has":
	// 	err = hostess.Has(args, flags)
	// case "off":
	// 	err = hostess.Off(args, flags)
	// case "on":
	// 	err = hostess.On(args, flags)
	// case "ls":
	// 	err = hostess.Ls(args, flags)
	// case "fix":
	// 	err = hostess.Fix(args, flags)
	// case "dump":
	// 	err = hostess.Dump(args, flags)
	// case "apply":
	// 	err = hostess.Apply(args, flags)

	// default:
	// 	fmt.Print(help)
	// }

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	app.Run(os.Args)

	os.Exit(0)
}
