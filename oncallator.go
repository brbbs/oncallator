package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/websdev/oncallator/schedule"
	"github.com/websdev/oncallator/terraform"
	"github.com/urfave/cli"
)

const (
	FlagIn = "in"
	FlagOut = "out"
	FlagFormat = "format"

	FormatSchedule = "schedule"
	FormatTerraform = "terraform"
)

func main() {
	app := cli.NewApp()
	app.Name = "oncallator"
	app.Flags = []cli.Flag {
		cli.StringFlag{
			Name: FlagIn,
			Usage: "If set, will read the base schedule from this file. Otherwise, reads from stdin.",
		},
		cli.StringFlag{
			Name: FlagOut,
			Usage: "If set, will write the generated schedule to this file. Otherwise, writes to stdout.",
		},
		cli.StringFlag{
			Name: FlagFormat,
			Usage: `Controls the output format of the schedule.

Allowed values:
	"schedule" -- will perform schedule generation on the input Schedule and output the updated JSON
	"terraform" -- will output Terraform pagerduty_schedule layers using the input Schedule`,
		},
	}
	app.Action = action

	app.Run(os.Args)
}

func action(ctx *cli.Context) error {
	s, err := readSchedule(ctx.String(FlagIn))
	if err != nil {
		return err
	}
	// TODO(brb): This feels gross. Can we find a way to make Schedule.Generate()
	// idempotent?
	if ctx.String(FlagFormat) == FormatSchedule {
		ns, err := s.Generate()
		if err != nil {
			return err
		}
		s = ns
	}
	out, err := output(ctx.String(FlagFormat), s)
	if err != nil {
		return err
	}
	return write(ctx.String(FlagOut), out)
}

func readSchedule(in string) (schedule.Schedule, error) {
	s := schedule.Schedule{}
	text := []byte{}
	if in == "" {
		t, err := ioutil.ReadAll(os.Stdin)
		if err != nil {
			return s, err
		}
		text = t
	} else {
		t, err := ioutil.ReadFile(in)
		if err != nil {
			return s, err
		}
		text = t
	}
	if err := json.Unmarshal(text, &s); err != nil {
		return s, err
	}
	if err := s.Validate(); err != nil {
		return s, err
	}
	return s, nil
}

func output(format string, s schedule.Schedule) ([]byte, error) {
	switch format {
	case FormatSchedule:
		return json.MarshalIndent(s, "", "  ")
	case FormatTerraform:
		return json.MarshalIndent(terraform.NewLayers(s.Rotations), "", "  ")
	default:
		return []byte{}, fmt.Errorf("unknown output format: %s", format)
	}
}

func write(out string, text []byte) error {
	if out == "" {
		_, err := fmt.Println(string(text))
		return err
	} else {
		return ioutil.WriteFile(out, text, 0660)
	}
}
