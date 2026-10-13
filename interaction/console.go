package interaction

import (
	"bufio"
	"fmt"
	"io"
	"kpo_hw1/entities"
	"strconv"
	"strings"
)

type console struct {
	in   *bufio.Reader
	echo bool
	out  io.Writer
}

func (cl *console) line(prompt string) (string, error) {
	fmt.Fprint(cl.out, prompt)
	str, err := cl.in.ReadString('\n')

	if err != nil && str == "" {
		return "", err
	}

	str = strings.TrimSpace(str)
	if cl.echo {
		fmt.Fprintln(cl.out, str)
	}

	return str, nil
}

func (cl *console) int(prompt string, min int, max int) (int, error) {
	for {
		str, err := cl.line(prompt)
		if err != nil {
			return 0, err
		}

		num, err := strconv.Atoi(str)
		if err == nil && num >= min && num <= max {
			return num, nil
		}

		fmt.Fprintf(cl.out, "Write a number from %d to %d.\n", min, max)
	}
}

func (cl *console) choose(title string, options []entities.Option) (entities.Option, error) {
	fmt.Fprintln(cl.out, title)
	for i, option := range options {
		fmt.Fprintf(cl.out, "  %d. %s\n", i+1, option.Title)
	}

	chosen, err := cl.int("\nIdx: ", 1, len(options))
	if err != nil {
		return entities.Option{}, err
	}

	return options[chosen-1], nil
}

func (cl *console) float(prompt string) (float64, error) {
	for {
		str, err := cl.line(prompt)
		if err != nil {
			return 0, err
		}

		num, err := strconv.ParseFloat(strings.ReplaceAll(str, ",", "."), 64)
		if err == nil && num >= 0 {
			return num, nil
		}

		fmt.Fprintln(cl.out, "Write a float >= 0")
	}
}
