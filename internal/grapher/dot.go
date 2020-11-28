package grapher

import (
	"fmt"
	"os"
	"os/exec"
)

type OutputFormat int8

const (
	_ OutputFormat = iota
	Svg
	Png
)

func (f OutputFormat) String() string {
	if f == Svg {
		return "-Tsvg"
	}
	if f == Png {
		return "-Tpng"
	}
	return ""
}

type Dot struct {
	Format     OutputFormat
	InputFile  string
	OutputName string
}

func (d *Dot) Execute() error {
	cmd := exec.Command("dot", d.Format.String(), d.InputFile, "-o"+d.OutputName)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()

	if err != nil {
		return fmt.Errorf(`%v
  Most likely, the error occurs due to the fact that Graphviz 
  is not installed or its path to it is not registered in the
  Path environment variable. 
  Please check.`, err)
	}

	return nil
}
