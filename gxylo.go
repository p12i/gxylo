package main

import (
	"fmt"
	"github.com/jessevdk/go-flags"
	"github.com/p12i/gxylo/proc"
	//	"gopkg.in/yaml.v2"
	//	"io/ioutil"
	"os"
)

const (
	gxylo_description = "This program blabl bla bla"
)

var Options struct {
	// ConfigFile string `short:"c" description:"YAML config file" default:"/etc/gxylo.yml"`
	Pid int `short:"p" description:"Process pid number" required:"true"`
}

var Config struct {
}

var Parser = flags.NewParser(&Options, flags.Default)

func die(err error, format string, args ...interface{}) {
	if err != nil {
		fmt.Fprintln(os.Stderr, fmt.Sprintf(format, args...))
		panic(err)
	}
}

func main() {
	Parser.LongDescription = gxylo_description

	if _, err := Parser.Parse(); err != nil {
		if flagsErr, ok := err.(*flags.Error); ok && flagsErr.Type == flags.ErrHelp {
			os.Exit(0)
		} else {
			os.Exit(1)
		}
	}

	// data, err := ioutil.ReadFile(Options.ConfigFile)
	// die(err, "Unable to read file %s\n", Options.ConfigFile)

	// die(yaml.Unmarshal(data, &Config), "")

	p, err := proc.NewProc(Options.Pid)
	die(err, "Unable to get process info %d", Options.Pid)

	err = p.ParseFDS()
	die(err, "Parsing error")

}
