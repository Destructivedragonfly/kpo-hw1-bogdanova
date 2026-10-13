package main

import (
	"flag"
	"fmt"
	"io"
	"kpo_hw1/entities"
	"kpo_hw1/infrastructure"
	"kpo_hw1/interaction"
	"os"
)

func main() {
	if err := run(os.Args[1:], os.Stdin, os.Stdout); err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		os.Exit(1)
	}
}

func run(args []string, stdin io.Reader, stdout io.Writer) error {
	flagSet := flag.NewFlagSet("vyshkat", flag.ContinueOnError)
	flagSet.SetOutput(os.Stderr)

	script := flagSet.String("script", "", "file input")
	minSimpl := flagSet.Int("min-simpl", 6, "minimal simplicity for a beginner")
	maxEnergy := flagSet.Float64("max-energy", 5.0, "max daily energy consumption")

	err := flagSet.Parse(args)
	if err != nil {
		return err
	}

	in := stdin
	if *script != "" {
		f, err := os.Open(*script)
		if err != nil {
			return err
		}
		defer f.Close()
		in = f
	}

	rec := infrastructure.New(
		infrastructure.NewCurStorage(),
		infrastructure.ServiceCenter{MaxEnergyKWh: *maxEnergy},
		infrastructure.CheckIsSimpl(*minSimpl),
	)

	return interaction.New(interaction.Info{
		In: in, Out: stdout, Reception: rec,
		Vehicles: entities.NewVehicleFactory(), Gear: entities.NewGearFactory(),
		Echo: *script != "",
	}).Run()
}
