package command

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
	"sync"
	"text/tabwriter"

	"github.com/hashicorp/consul-terraform-sync/api"
	"github.com/hashicorp/consul-terraform-sync/controller"
	"github.com/hashicorp/consul-terraform-sync/version"
	mcli "github.com/mitchellh/cli"
)

// Exit codes are int values that represent an exit code for a particular error.
// Sub-systems may check this unique error to determine the cause of an error
// without parsing the output or help text.
//
// Errors start at 10
const (
	ExitCodeOK int = 0

	ExitCodeError = 10 + iota
	ExitCodeInterrupt
	ExitCodeRequiredFlagsError
	ExitCodeParseFlagsError
	ExitCodeConfigError
	ExitCodeDriverError

	logSystemName = "cli"
)

var _ api.Server = (*controller.TasksManager)(nil)

// CLI is the main entry point.
type CLI struct {
	sync.Mutex

	// outSteam and errStream are the standard out and standard error streams to
	// write messages from the CLI.
	outStream, errStream io.Writer

	// signalCh is the channel where the cli receives signals.
	signalCh chan os.Signal

	// stopCh is an internal channel used to trigger a shutdown of the CLI.
	stopCh chan struct{}
}

var (
	// Common commands are grouped separately to call them out to operators.
	commonCommands = []string{
		"start",
		"task",
	}
)

// NewCLI creates a new CLI object with the given stdout and stderr streams.
func NewCLI(out, err io.Writer) *CLI {
	return &CLI{
		outStream: out,
		errStream: err,
		signalCh:  make(chan os.Signal, 1),
		stopCh:    make(chan struct{}),
	}
}

// Run accepts a slice of arguments and returns an int representing the exit
// status from the command.
func (cli *CLI) Run(args []string) int {
	subcommands := &mcli.CLI{
		Name:                       "consul-terraform-sync",
		Args:                       args[1:],
		Commands:                   Commands(),
		Autocomplete:               true,
		AutocompleteNoDefaultFlags: true,
		HelpFunc:                   help,
		HelpWriter:                 tabwriter.NewWriter(os.Stdout, 0, 2, 4, ' ', tabwriter.AlignRight),
	}

	if subcommands.IsVersion() {
		fmt.Fprintf(cli.outStream, "%s %s\n", version.Name, version.GetHumanVersion())
		fmt.Fprintf(cli.outStream, "Compatible with Terraform %s\n", version.CompatibleTerraformVersionConstraint)
		return ExitCodeOK
	}

	exitCode, err := subcommands.Run()
	if err != nil {
		fmt.Fprintf(cli.errStream, "Error running the CLI command '%s': %s",
			strings.Join(args, " "), err)
	}
	return exitCode
}

func help(commands map[string]mcli.CommandFactory) string {
	c := make(map[string]string)
	for _, v := range commonCommands {
		c[v] = generateCommandHelp(v, commands[v])
	}

	return generateHelp(c, "Usage CLI: consul-terraform-sync <command> [-help] [options]\n", nil)
}

// To support usage by the start command, this function takes a map[string]string and 
// allows for omission of flags. This functionality will not be needed once running CTS 
// without command arguments is deprecated in a future major release.
func generateHelp(commands map[string]string, usage string, omitFlags []string) string {
	var b bytes.Buffer
	tw := tabwriter.NewWriter(&b, 0, 2, 4, ' ', tabwriter.AlignRight)

	fmt.Fprint(tw, usage)

	if len(commands) > 0 {
		fmt.Fprintf(tw, "\nCommands:\n")

		for _, v := range commonCommands {
			fmt.Fprintf(tw, commands[v])
		}
	}

	fmt.Fprintf(tw, "\n")

	rc := startCommand{}
	printFlags(tw, rc.startFlags(), omitFlags)

	tw.Flush()

	return strings.TrimSpace(b.String())
}

func generateCommandHelp(cmdName string, cmdFactory mcli.CommandFactory) string {
	cmd, err := cmdFactory()
	if err != nil {
		panic(fmt.Sprintf("failed to load %q command: %s", cmdName, err))
	}

	return fmt.Sprintf("%s\t%s\n", cmdName, cmd.Synopsis())
}

// printFlags prints out select flags
func printFlags(w io.Writer, f *flag.FlagSet, omitName []string) {
	fmt.Fprintf(w, "Options:\n")
	f.VisitAll(func(f *flag.Flag) {
		switch f.Name {
		case "h", "help":
			// don't print out help flags
			return
		case "client-type":
			// don't print out development-only flags
			return
		}

		isOmitName := false
		for _, n := range omitName {
			if n == f.Name {
				isOmitName = true
			}
		}

		if !isOmitName {
			fmt.Fprintf(w, "\t-%s %s\n", f.Name, templateDefaultValue(f.Value))
			fmt.Fprintf(w, "\t\t%s\n\n", f.Usage)
		}
	})
}

func templateDefaultValue(value flag.Value) string {
	if value.String() != "" {
		return fmt.Sprintf("[Default: %s]", value)
	}

	return ""
}
