package interaction

import (
	"bufio"
	"bytes"
	"io"
	"kpo_hw1/entities"
	"kpo_hw1/infrastructure"
	"strings"
	"testing"
)

func newConsole(input string) *console {
	return &console{in: bufio.NewReader(strings.NewReader(input)), out: &bytes.Buffer{}}
}

func TestConsoleInt(t *testing.T) {
	n, err := newConsole("abc\n15\n7\n").int("N: ", 1, 10)

	if err != nil || n != 7 {
		t.Errorf("got %d, %v; want 7 after two invalid answers", n, err)
	}
}

func TestConsoleFloat(t *testing.T) {
	x, err := newConsole("-1\n2,5\n").float("X: ")

	if err != nil || x != 2.5 {
		t.Errorf("got %v, %v; want 2.5 (negative rejected, comma accepted)", x, err)
	}
}

func TestConsoleEndOfInput(t *testing.T) {
	if _, err := newConsole("").line("Q: "); err != io.EOF {
		t.Errorf("got %v, want io.EOF", err)
	}
}

func newApp(input string) (*App, *infrastructure.Reception, *bytes.Buffer) {
	vehicles := entities.NewFactory()
	_ = vehicles.Register("scooter", "Scooter", func(in entities.Info) (entities.Item, error) {
		return entities.NewVehicle("Scooter", in)
	})
	gearKinds := entities.NewFactory()
	_ = gearKinds.Register("helmet", "Helmet", func(in entities.Info) (entities.Item, error) {
		return entities.NewGear("Helmet", in)
	})

	rec := infrastructure.New(
		infrastructure.NewCurStorage(),
		infrastructure.ServiceCenter{MaxEnergyKWh: 5},
		infrastructure.CheckIsSimpl(6),
	)
	out := &bytes.Buffer{}
	app := New(Info{In: strings.NewReader(input), Out: out, Reception: rec, Vehicles: vehicles, Gear: gearKinds})

	return app, rec, out
}

func TestAppExit(t *testing.T) {
	app, _, out := newApp("0\n")

	if err := app.Run(); err != nil || !strings.Contains(out.String(), "Done") {
		t.Errorf("err = %v, output must contain Done", err)
	}
}

func TestAppEndOfInput(t *testing.T) {
	app, _, _ := newApp("")

	if err := app.Run(); err != nil {
		t.Errorf("end of input must not be an error, got %v", err)
	}
}

func TestAppReport(t *testing.T) {
	app, _, out := newApp("3\n0\n")
	_ = app.Run()

	if !strings.Contains(out.String(), "Total units: 0") {
		t.Errorf("report is missing, output: %q", out.String())
	}
}

func TestAppAddVehicle(t *testing.T) {
	app, rec, _ := newApp("1\n1\nXiaomi\nSC-1\n8\n1,5\n1\n0\n")
	_ = app.Run()
	if !rec.CheckHas("SC-1") {
		t.Error("accepted vehicle must be in the fleet")
	}

	app, rec, _ = newApp("1\n1\nXiaomi\nSC-1\n8\n1,5\n2\n0\n")
	_ = app.Run()
	if rec.CheckHas("SC-1") {
		t.Error("rejected vehicle must not be in the fleet")
	}
}

func TestAppLists(t *testing.T) {
	app, _, out := newApp("1\n1\nXiaomi\nSC-1\n8\n1,5\n1\n4\n5\n0\n")
	_ = app.Run()

	if n := strings.Count(out.String(), "SC-1"); n != 3 {
		t.Errorf("SC-1 appears %d times, want 3", n)
	}
}

func TestAppAddGear(t *testing.T) {
	app, rec, _ := newApp("2\n1\nMy helmet\nH-1\n0\n0\n")
	_ = app.Run()

	if !rec.CheckHas("H-1") {
		t.Error("gear must be registered")
	}
}
