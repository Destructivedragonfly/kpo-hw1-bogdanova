package main

import (
	"bytes"
	"strings"
	"testing"

	"kpo_hw1/entities"
)

const demoFile = "demo.txt"

func runDemo(t *testing.T, extraArgs ...string) string {
	t.Helper()

	var out bytes.Buffer
	args := append([]string{"-script", demoFile}, extraArgs...)
	if err := run(args, strings.NewReader(""), &out); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	return out.String()
}

func TestDemoNoDesync(t *testing.T) {
	text := runDemo(t)

	if strings.Contains(text, "from 0 to 5") {
		t.Error("answers not synchronized")
	}
	if !strings.Contains(text, "Done") {
		t.Error("demo have not gone till the end")
	}
}

func TestDemoReport(t *testing.T) {
	text := runDemo(t)

	if !strings.Contains(text, "Total units: 8") {
		t.Error("expected 8 units")
	}
	if !strings.Contains(text, "Total daily energy consumption: 17.40 kWh") {
		t.Error("expected 17.40 kWh")
	}
}

func beginnersSection(t *testing.T, text string) string {
	t.Helper()

	start := strings.Index(text, "Beginner-friendly vehicles:")
	end := strings.Index(text, "All vehicles and gear:")
	if start < 0 || end < start {
		t.Fatal("no beginner-friendly shown")
	}

	return text[start:end]
}

func TestDemoBeginners(t *testing.T) {
	section := beginnersSection(t, runDemo(t))

	for _, num := range []string{"BK-201", "BK-202", "SC-101"} {
		if !strings.Contains(section, num) {
			t.Errorf("%s should be begginer-friendly", num)
		}
	}
	if strings.Contains(section, "EB-301") {
		t.Error("EB-301 (simplicity 5) is not beginner-friendly")
	}
}

func TestDemoRejectionsAndErrors(t *testing.T) {
	text := runDemo(t)

	for _, want := range []string{
		"Rejected: " + entities.ErrRejected.Error(),
		"Rejected: " + entities.ErrEnergyOverLimit.Error(),
		"Error: " + entities.ErrDuplicateNum.Error(),
		"Error: " + entities.ErrEmptyNum.Error(),
	} {
		if !strings.Contains(text, want) {
			t.Errorf("no %q in the output", want)
		}
	}
}

func TestFlagMinSimpl(t *testing.T) {
	section := beginnersSection(t, runDemo(t, "-min-simpl", "9"))

	if !strings.Contains(section, "BK-201") {
		t.Error("BK-201 (simplicity 9) should fit")
	}
	if strings.Contains(section, "SC-101") {
		t.Error("SC-101 (simplicity 7) with min 9 should not fit")
	}
}
