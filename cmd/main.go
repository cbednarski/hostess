package main

import (
	"flag"
	"fmt"
	"github.com/cbednarski/hostess"
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

WARNING: This program is BETA and not all commands are implemented.

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

Note: You can specify the HOSTESS_FILE environment variable to operate on a
file other than /etc/hosts

Report bugs at https://github.com/cbednarski/hostess
`

func main() {
	hostfile := hostess.NewHostfile(hostess.GetHostsPath())
	hostfile.Load()
	hostfile.Parse()

	flags := make(map[string]*bool)

	flags["force"] = flag.Bool("f", false, "Forcibly apply changes, even if there are errors or conflicts")
	flags["noop"] = flag.Bool("n", false, "No-op. Show changes but don't write them.")
	flags["quiet"] = flag.Bool("q", false, "Suppress error and conflict messages.")
	flags["silent"] = flag.Bool("qq", false, "Suppress all output. Check exit codes for success / failure.")
	flags["help"] = flag.Bool("h", false, "Help")

	flag.Parse()

	// Guard against zero arguments
	var args []string
	if len(flag.Args()) > 0 {
		args = flag.Args()[1:]
	}

	var err error = nil

	switch flag.Arg(0) {
	case "add":
		err = hostess.Add(args, flags)
	case "del":
		err = hostess.Del(args, flags)
	case "has":
		err = hostess.Has(args, flags)
	case "off":
		err = hostess.Fix(args, flags)
	case "on":
		err = hostess.Fix(args, flags)
	case "ls":
		err = hostess.Fix(args, flags)
	case "fix":
		err = hostess.Fix(args, flags)
	case "dump":
		err = hostess.Fix(args, flags)
	case "apply":
		err = hostess.Fix(args, flags)

	default:
		fmt.Print(help)
	}

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	os.Exit(0)
}
