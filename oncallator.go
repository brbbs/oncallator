package main

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/websdev/oncallator/rotations"
	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()
	app.Name = "oncallator"
	app.Flags = []cli.Flag{
		cli.StringSliceFlag{
			Name: "user",
			Usage: "list of oncall users to schedule",
		},
		cli.StringFlag{
			Name: "start_date",
			Usage: fmt.Sprintf("An RFC3339-compliant date to start the oncall rotation (e.g., %s)", time.RFC3339),
		},
		cli.StringFlag{
			Name: "end_date",
			Usage: "Rotations will be scheduled up to the end_date",
		},
		cli.DurationFlag{
			Name: "duration",
			Usage: "How long a single oncall rotation lasts",
		},
	}
	app.Action = action

	app.Run(os.Args)
}

func action(ctx *cli.Context) error {
	start, end, err := getDates(ctx)
	if err != nil {
		return err
	}
	if len(ctx.StringSlice("user")) == 0 {
		return fmt.Errorf("must provide at least 1 user")
	}
	n := rotations.Num(start, end, ctx.Duration("duration"))
	rotations := rotations.New(ctx.StringSlice("user"), start, ctx.Duration("duration"), n)
	out, err := json.MarshalIndent(rotations.TerraformLayers(), "", "  ")
	if err != nil {
		return err
	}
	fmt.Println(string(out))
	return nil
}

func getDates(ctx *cli.Context) (time.Time, time.Time, error) {
	start, err := time.Parse(time.RFC3339, ctx.String("start_date"))
	if err != nil {
		return start, time.Time{}, fmt.Errorf("bad start_date: %s", err)
	}
	end, err := time.Parse(time.RFC3339, ctx.String("end_date"))
	if err != nil {
		return start, end, fmt.Errorf("bad end_date: %s", err)
	}
	if !start.Before(end) {
		return start, end, fmt.Errorf("provided start_date after end_date")
	}
	return start, end, nil
}
