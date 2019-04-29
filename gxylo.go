package main

import (
	"fmt"
	"github.com/jessevdk/go-flags"
	"github.com/p12i/gxylo/connections"
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
	Pid             int  `short:"p" description:"Process pid number" required:"true"`
	FullConnections bool `short:"l" description:"Show all network connections"`
}

var Config struct {
}

var Parser = flags.NewParser(&Options, flags.Default)

func die(err error, doPanic bool, format string, args ...interface{}) {
	if err != nil {
		fmt.Fprintln(os.Stderr, fmt.Sprintf(format, args...))
		if doPanic {
			panic(err)
		} else {
			os.Exit(1)
		}
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

	p, err := proc.NewProc(uint64(Options.Pid))
	die(err, false, "Unable to get process info %d", Options.Pid)
	err = p.ParseFDS()
	die(err, false, "Parsing error")

	connections_list := connections.ConnectionList{}
	if err := connections_list.ParseConnections(); err != nil {
		die(err, false, "Error during parsing connections")
	}

	fmt.Println(p.Info(&connections_list))

	if Options.FullConnections {
		fmt.Println(connections_list.String())
	}
}
