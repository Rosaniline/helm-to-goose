package main

import (
	"fmt"
	"github.com/go-yaml/yaml"
	"github.com/jessevdk/go-flags"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

var opts struct {
	ValueFile      string `short:"v" long:"value" description:"Helm value file" required:"true"`
	OutputFile     string `short:"o" long:"output" description:"Output file" required:"true"`
	DisableSSLMode bool   `long:"disable-ssl-mode" description:"SSL mode for db connection" required:"false"`
}

type Values struct {
	Db DB `yaml:"db"`
}

type DB struct {
	Host     string `yaml:"host"`
	Name     string `yaml:"name"`
	Password string `yaml:"password"`
	User     string `yaml:"user"`
}

func readYAML(filename string) (*Values, error) {

	f, err := ioutil.ReadFile(filename)

	if err != nil {
		return nil, err
	}

	m := new(Values)

	if yaml.Unmarshal(f, &m) != nil {
		return nil, err
	}

	return m, nil
}

func SaveYAML(v *Values, env string, disableSSLMode bool, output string) error {
	temp := make(map[string]map[string]string)
	temp[env] = make(map[string]string)
	mode := "require"

	if disableSSLMode {
		mode = "disable"
	}

	temp[env]["driver"] = "postgres"
	temp[env]["open"] = fmt.Sprintf("host=%v user=%v password=%v dbname=%v sslmode=%v", v.Db.Host, v.Db.User, v.Db.Password, v.Db.Name, mode)

	y, err := yaml.Marshal(temp)

	if err != nil {
		return err
	}

	if err = ioutil.WriteFile(output, y, os.ModePerm); err != nil {
		return err
	}

	return nil
}

func main() {
	_, err := flags.Parse(&opts)

	if err != nil {
		if flagsErr, ok := err.(*flags.Error); ok && flagsErr.Type == flags.ErrHelp {
			os.Exit(0)
		} else {
			os.Exit(1)
		}
	}

	values, err := readYAML(opts.ValueFile)

	if err != nil {
		panic(err)
	}

	env := strings.TrimSuffix(filepath.Base(opts.ValueFile), filepath.Ext(opts.ValueFile))

	if err = SaveYAML(values, env, opts.DisableSSLMode, opts.OutputFile); err != nil {
		panic(err)
	}
}
