package command

import (
	"droplet/internal/container"

	"github.com/urfave/cli/v2"
)

func commandExec() *cli.Command {
	return &cli.Command{
		Name:      "exec",
		Usage:     "exec a container",
		ArgsUsage: "<container-id> <command>",
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:    "interactive",
				Usage:   "Execute entrypoint in interative mode",
				Aliases: []string{"i"},
			},
		},
		Action: runExec,
	}
}

func runExec(ctx *cli.Context) error {
	// retrieve container id
	containerId := ctx.Args().Get(0)
	// retrieve args
	args := ctx.Args().Slice()
	// options
	// interactive
	interactive := ctx.Bool("interactive")
	entrypoint := args[1:]

	containerExec := container.NewContainerExec()
	err := containerExec.Exec(container.ExecOption{
		ContainerId: containerId,
		Interactive: interactive,
		Entrypoint:  entrypoint,
	})
	if err != nil {
		return err
	}
	return nil
}
