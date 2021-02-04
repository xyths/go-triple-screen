package utils
import 	"github.com/urfave/cli/v2"
var (
	ConfigFlag = &cli.StringFlag{
		Name:    "config",
		Aliases: []string{"c"},
		Value:   "config.json",
		Usage:   "load configuration from `file`",
	}
	DryRunFlag = &cli.BoolFlag{
		Name:  "dry-run",
		Value: false,
		Usage: "do not run, just print the result",
	}
)
