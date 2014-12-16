// Package df generates corresponding Dockerfiles by YAML.
package df

import (
	"bytes"
	"fmt"
	"strings"

	"gopkg.in/yaml.v2"
)

// Config represents a YAML config section.
type Config struct {
	Image       string
	Environment []string
	Script      []string
	Cmd         []string
}

func generate(cfg *Config) ([]byte, error) {
	var buf bytes.Buffer

	buf.WriteString(fmt.Sprintf("FROM %s\n\n", cfg.Image))

	for _, env := range cfg.Environment {
		buf.WriteString(fmt.Sprintf("ENV %s\n", strings.Replace(env, "=", " ", 1)))
	}
	buf.WriteString("\n")

	for _, spt := range cfg.Script {
		buf.WriteString(fmt.Sprintf("RUN %s\n", spt))
	}
	buf.WriteString("\n")

	buf.WriteString(fmt.Sprintf("CMD [\"%s\"]\n", strings.Join(cfg.Cmd, "\", \"")))

	return buf.Bytes(), nil
}

// Generate takes a YAML file and generate Dockerfiles to given directory.
func Generate(from []byte) (_ map[string][]byte, err error) {
	var cfgs map[string]*Config
	if err = yaml.Unmarshal(from, &cfgs); err != nil {
		return nil, err
	}

	dfs := make(map[string][]byte)
	for name, cfg := range cfgs {
		if dfs[name], err = generate(cfg); err != nil {
			return nil, err
		}
	}

	return dfs, nil
}
