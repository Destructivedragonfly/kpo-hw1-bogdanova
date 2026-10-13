package interaction

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"kpo_hw1/entities"
	"kpo_hw1/infrastructure"
	"maps"
	"slices"
)

type Info struct {
	In        io.Reader
	Out       io.Writer
	Reception *infrastructure.Reception
	Vehicles  *entities.Factory
	Gear      *entities.Factory
	Echo      bool
}

type App struct {
	Info
	cl   *console
	menu []menuItem
}

type menuItem struct {
	title  string
	action func() error
}

func New(info Info) *App {
	app := &App{
		Info: info,
		cl:   &console{in: bufio.NewReader(info.In), out: info.Out, echo: info.Echo},
	}

	app.menu = []menuItem{
		{"Accept vehicle (with inspection)", app.addVehicle},
		{"Register gear", app.addGear},
		{"Fleet report", app.printReport},
		{"Beginner-friendly vehicles", app.printBeginners},
		{"Full inventory", app.printInventory},
	}

	return app
}

func (app *App) print(format string, args ...any) {
	fmt.Fprintf(app.Out, format+"\n", args...)
}

func (app *App) Run() error {
	for {
		app.showMenu()

		choice, err := app.cl.int("Choose an action: ", 0, len(app.menu))
		if err != nil {
			return app.finish(err)
		}

		if choice == 0 {
			app.print("Done")
			return nil
		}

		if err := app.menu[choice-1].action(); err != nil {
			return app.finish(err)
		}
	}
}

func (app *App) showMenu() {
	app.print("\nVyshKat menu:")
	for i, item := range app.menu {
		app.print("%d. %s", i+1, item.title)
	}
	app.print("0. Exit")
}

func (app *App) finish(err error) error {
	if errors.Is(err, io.EOF) {
		app.print("\nThe end!")
		return nil
	}

	return err
}

func (app *App) askItem(f *entities.Factory, title string, vehicle bool) (entities.Item, error) {
	options := f.Options()
	if len(options) == 0 {
		app.print("No available type")
		return nil, nil
	}

	class, err := app.cl.choose(title, options)
	if err != nil {
		return nil, err
	}

	namePrompt := "Name: "
	if !vehicle {
		namePrompt = "Name is empty, use type name: "
	}
	name, err := app.cl.line(namePrompt)
	if err != nil {
		return nil, err
	}

	num, err := app.cl.line("Inventory number: ")
	if err != nil {
		return nil, err
	}

	if app.Reception.CheckHas(num) {
		app.print("Error: %v", entities.ErrDuplicateNum)
		return nil, nil
	}

	info := entities.Info{Name: name, Num: num}
	if vehicle {
		if info.Simpl, err = app.cl.int("Simplicity for beginners (1-10): ", entities.MinSimpl, entities.MaxSimpl); err != nil {
			return nil, err
		}
	}
	if info.EnergyKWh, err = app.cl.float("Daily energy consumption, kWh (0 if none): "); err != nil {
		return nil, err
	}

	item, err := f.Create(class.ID, info)
	if err != nil {
		app.print("Error: %v", err)
		return nil, nil
	}

	return item, nil
}

func (app *App) addVehicle() error {
	item, err := app.askItem(app.Vehicles, "\nVehicle type:", true)
	if item == nil || err != nil {
		return err
	}

	verdict, err := app.cl.int("Exam result (1 - pass, 2 - fail): ", 1, 2)
	if err != nil {
		return err
	}

	result, err := app.Reception.Exam(item, verdict == 1)

	switch {
	case err != nil:
		app.print("Acceptance error: %v", err)
	case !result.IsAccepted:
		app.print("Rejected: %s. Not added to active fleet.", result.Problem)
	default:
		app.print("Added to active fleet: %s.", describe(item))
	}

	return nil
}

func (app *App) addGear() error {
	item, err := app.askItem(app.Gear, "\nGear type:", false)
	if item == nil || err != nil {
		return err
	}

	if err := app.Reception.Register(item); err != nil {
		app.print("Registration error: %v", err)
		return nil
	}

	app.print("Registered: %s.", describe(item))
	return nil
}

func (app *App) printReport() error {
	s := app.Reception.Report()
	app.print("\nFleet report")
	app.print("Total units: %d", s.Total)

	for _, class := range slices.Sorted(maps.Keys(s.AmountInClass)) {
		app.print("  %s: %d", class, s.AmountInClass[class])
	}

	app.print("Total daily energy consumption: %.2f kWh", s.EnergyKWh)
	return nil
}

func (app *App) printBeginners() error {
	var lines []string
	for _, item := range app.Reception.Beginners() {
		lines = append(lines, fmt.Sprintf("%s %s %s simplicity %d/10",
			item.Class(), item.Name(), item.InventoryNum(), item.SimplForBeginners()))
	}

	app.printList("Beginner-friendly vehicles:", "No suitable devices yet.", lines)
	return nil
}

func (app *App) printInventory() error {
	var lines []string
	for _, item := range app.Reception.Inventory() {
		energy := "—"

		if consumer, ok := item.(entities.EnergyConsumer); ok && consumer.DailyEnergyKWh() > 0 {
			energy = fmt.Sprintf("%.2f kWh/day", consumer.DailyEnergyKWh())
		}

		lines = append(lines, fmt.Sprintf("%s %s %s %s",
			item.Class(), item.Name(), item.InventoryNum(), energy))
	}

	app.printList("All vehicles and gear:", "Inventory is empty", lines)
	return nil
}

func (app *App) printList(title, empty string, lines []string) {
	app.print("\n%s", title)

	if len(lines) == 0 {
		app.print("%s", empty)
		return
	}

	for _, line := range lines {
		fmt.Fprintln(app.Out, line)
	}
}

func describe(item entities.Item) string {
	return fmt.Sprintf("%s «%s», inv. no. %s", item.Class(), item.Name(), item.InventoryNum())
}
