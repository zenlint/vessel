package core

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"strings"

	"gopkg.in/yaml.v2"
)

// Git represents Git repository information.
type Git struct {
	Auth, Path string
}

// Section represents a solution section.
type Service struct {
	Name        string
	Image       string
	Environment []string
	Script      []string
	Cmd         []string
	Ports       []string
	Expose      []string
	Services    []*Service
}

// Dockerfile generates Dockerfile for the service.
func (s *Service) Dockerfile() ([]byte, error) {
	var buf bytes.Buffer

	buf.WriteString(fmt.Sprintf("FROM %s\n\n", s.Image))

	if len(s.Environment) > 0 {
		for _, env := range s.Environment {
			buf.WriteString(fmt.Sprintf("ENV %s\n", strings.Replace(env, "=", " ", 1)))
		}
		buf.WriteString("\n")
	}

	if len(s.Script) > 0 {
		for _, spt := range s.Script {
			buf.WriteString(fmt.Sprintf("RUN %s\n", spt))
		}
		buf.WriteString("\n")
	}

	if len(s.Cmd) > 0 {
		buf.WriteString(fmt.Sprintf("CMD [\"%s\"]\n", strings.Join(s.Cmd, "\", \"")))
	}

	return buf.Bytes(), nil
}

func parseString(v interface{}) string {
	if v == nil {
		return ""
	}
	return v.(string)
}

func parseStrings(v interface{}) []string {
	if v == nil {
		return []string{}
	}
	strs := v.([]interface{})
	list := make([]string, len(strs))
	for i := range strs {
		list[i] = parseString(strs[i])
	}
	return list
}

// parseService parses section to service with no cycle dependency.
func parseService(name string, sections map[string]map[string]interface{}, marked map[string]bool) (s *Service, err error) {
	if marked[name] {
		return nil, fmt.Errorf("cycle dependency found for '%s'", name)
	}
	marked[name] = true

	sec := sections[name]

	s = new(Service)
	s.Name = name
	s.Image = parseString(sec["image"])
	s.Environment = parseStrings(sec["environment"])
	s.Script = parseStrings(sec["script"])
	s.Cmd = parseStrings(sec["cmd"])

	if name != "app" {
		s.Ports = parseStrings(sec["ports"])
		s.Expose = parseStrings(sec["expose"])
	}

	serviceNames := parseStrings(sec["services"])
	s.Services = make([]*Service, len(serviceNames))
	for i, name := range serviceNames {
		s.Services[i], err = parseService(name, sections, marked)
		if err != nil {
			return nil, err
		}
	}
	return s, nil
}

// NotifyEmail represents e-mail notification information.
type NotifyEmail struct {
	Recipients []string
}

// Notify represents notification information.
type Notify struct {
	Email NotifyEmail
}

// Solution represents a build solution with needed information.
type Solution struct {
	Git Git
	Notify
	*Service
}

// NewSolutionFromBytes creates and returns a new solution with given YAML in bytes.
func NewSolutionFromBytes(data []byte) (_ *Solution, err error) {
	sections := make(map[string]map[string]interface{})
	if err = yaml.Unmarshal(data, &sections); err != nil {
		return nil, fmt.Errorf("error parsing YAML: %v", err)
	}

	sln := new(Solution)

	app := sections["app"]
	if app == nil {
		return nil, errors.New("section 'app' not found")
	}

	gitSec := app["git"].(map[interface{}]interface{})
	sln.Git.Auth = parseString(gitSec["auth"])
	sln.Git.Path = parseString(gitSec["path"])

	if app["notify"] != nil {
		recipients := app["notify"].(map[interface{}]interface{})["email"].(map[interface{}]interface{})["recipients"].([]interface{})
		sln.Notify.Email.Recipients = make([]string, len(recipients))
		for i := range recipients {
			sln.Notify.Email.Recipients[i] = recipients[i].(string)
		}
	}

	sln.Service, err = parseService("app", sections, map[string]bool{})
	if err != nil {
		return nil, err
	}

	return sln, nil
}

// NewSolutionFromFile creates and returns a new solution with given YAML file.
func NewSolutionFromFile(file string) (*Solution, error) {
	data, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}
	return NewSolutionFromBytes(data)
}
