package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/websdev/oncallator/schedule"
	"github.com/urfave/cli"
)

const (
	FlagIn = "in"
	FlagOut = "out"
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
	}
	app.Action = action

	app.Run(os.Args)
}

func action(ctx *cli.Context) error {
	s, err := readSchedule(ctx.String(FlagIn))
	if err != nil {
		return err
	}
	ns, err := s.Generate()
	if err != nil {
		return err
	}
	return writeSchedule(ctx.String(FlagOut), ns)
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
	return s, nil
}

func writeSchedule(out string, s schedule.Schedule) error {
	text, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		return err
	}
	if out == "" {
		_, err := fmt.Println(string(text))
		return err
	} else {
		return ioutil.WriteFile(out, text, 0660)
	}
}
