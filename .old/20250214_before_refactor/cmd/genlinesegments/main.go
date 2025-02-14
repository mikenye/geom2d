package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/mikenye/geom2d/linesegment"
	"github.com/urfave/cli/v3"
	"log"
	"math/rand/v2"
	"os"
)

func main() {
	cmd := &cli.Command{
		Name:      "genlinesegments",
		Usage:     "Generates random line segments in a plane and outputs results to stdout as JSON",
		UsageText: "genlinesegments --number <value> --maxx <value> --minx <value> --maxy <value> --miny <value>",
		Flags: []cli.Flag{
			&cli.IntFlag{
				Name:     "number",
				Usage:    "The number of segments to create",
				Value:    3,
				Aliases:  []string{"n"},
				OnlyOnce: true,
				Validator: func(u int64) error {
					if u <= 0 {
						return fmt.Errorf("number must be greater than zero")
					}
					return nil
				},
			},
			&cli.IntFlag{
				Name:     "maxx",
				Usage:    "The maximum X value of the plane",
				OnlyOnce: true,
				Value:    10,
			},
			&cli.IntFlag{
				Name:     "minx",
				Usage:    "The minimum X value of the plane",
				OnlyOnce: true,
				Value:    0,
			},
			&cli.IntFlag{
				Name:     "maxy",
				Usage:    "The maximum Y value of the plane",
				OnlyOnce: true,
				Value:    10,
			},
			&cli.IntFlag{
				Name:     "miny",
				Usage:    "The minimum Y value of the plane",
				OnlyOnce: true,
				Value:    0,
			},
		},
		HideVersion: true,
		Action:      app,
		Authors:     []any{"https://github.com/mikenye"},
	}
	if err := cmd.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}
}

func randomIntInRange(min, max int64) int64 {
	return min + rand.Int64N(max-min+1)
}

func app(_ context.Context, cmd *cli.Command) error {

	minx := cmd.Int("minx")
	maxx := cmd.Int("maxx")
	miny := cmd.Int("miny")
	maxy := cmd.Int("maxy")
	n := cmd.Int("number")

	// sanity checks
	if minx >= maxx {
		return fmt.Errorf("maxx must be greater than minx")
	}
	if miny >= maxy {
		return fmt.Errorf("maxy must be greater than miny")
	}

	// prep output slice
	output := make([]linesegment.LineSegment[int64], n)

	// fill output slice
	for i := int64(0); i < n; i++ {
		for {
			output[i] = linesegment.New[int64](
				randomIntInRange(minx, maxx), // x1
				randomIntInRange(miny, maxy), // y1
				randomIntInRange(minx, maxx), // x2
				randomIntInRange(miny, maxy), // y2
			)

			// skip degenerate segments
			if !output[i].Start().Eq(output[i].End()) {
				break
			}
		}
	}
	b, err := json.Marshal(output)
	if err != nil {
		return err
	}
	fmt.Print(string(b))
	return nil
}
