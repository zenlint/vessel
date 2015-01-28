package core

import (
	"archive/tar"
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"regexp"
	"strings"
	"time"

	api "github.com/fsouza/go-dockerclient"
	"gopkg.in/yaml.v2"

	"github.com/dockercn/vessel/modules/utils"
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
	Git         Git
	client      *api.Client
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

	if len(s.Cmd) > 0 {
		buf.WriteString(fmt.Sprintf("CMD [\"%s\"]\n", strings.Join(s.Cmd, "\", \"")))
	}

	if s.Name == "app" {
		buf.WriteString("WORKDIR $GOPATH/src/" + path.Dir(s.Git.Path) + "\n")
		buf.WriteString("RUN git clone https://" + s.Git.Path + ".git\n")
		buf.WriteString("WORKDIR " + path.Base(s.Git.Path) + "\n")
	}

	if len(s.Script) > 0 {
		for _, spt := range s.Script {
			buf.WriteString(fmt.Sprintf("RUN %s\n", spt))
		}
		buf.WriteString("\n")
	}

	fmt.Println(buf.String())
	return buf.Bytes(), nil
}

var imageIdPattern = regexp.MustCompile("Successfully built ([0-9a-f]+)")

// FIXME: recursively build.
// Build builds services recursively and returns image ID.
func (s *Service) Build() (string, error) {
	log.Println("Building image...")

	data, err := s.Dockerfile()
	if err != nil {
		return "", err
	}

	var input bytes.Buffer
	so := utils.NewStreamOutput()
	t := time.Now()

	tr := tar.NewWriter(&input)
	tr.WriteHeader(&tar.Header{Name: "Dockerfile", Size: int64(len(data)), ModTime: t, AccessTime: t, ChangeTime: t})
	tr.Write(data)
	tr.Close()
	opts := api.BuildImageOptions{
		Name:           "test",
		RmTmpContainer: true,
		InputStream:    &input,
		OutputStream:   so,
		RawJSONStream:  true,
	}
	if err = s.client.BuildImage(opts); err != nil {
		return "", err
	}

	var imageId string
	if len(so.Events) > 0 {
		e := so.Events[len(so.Events)-1]
		m := imageIdPattern.FindAllStringSubmatch(e["stream"], 1)
		if m != nil {
			imageId = m[0][1]
		}
	}
	return imageId, nil
}

// FIXME: recursively start.
func (s *Service) Start(id string) error {
	c, err := s.client.CreateContainer(api.CreateContainerOptions{
		Config: &api.Config{
			Image: id,
		},
	})
	if err != nil {
		fmt.Println("fail to create container")
		return err
	}
	if err = s.client.StartContainer(c.ID, &api.HostConfig{}); err != nil {
		return err
	}
	s.client.WaitContainer(c.ID)
	return s.client.RemoveContainer(api.RemoveContainerOptions{
		ID: c.ID,
	})
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
	} else {
		gitSec := sec["git"].(map[interface{}]interface{})
		s.Git.Auth = parseString(gitSec["auth"])
		s.Git.Path = parseString(gitSec["path"])
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
	Notify
	*Service
}

// NewClient creates and returns a Docker client.
func NewClient() (c *api.Client, err error) {
	host := os.Getenv("DOCKER_HOST")
	if os.Getenv("DOCKER_TLS_VERIFY") == "1" {
		certPath := os.Getenv("DOCKER_CERT_PATH")
		c, err = api.NewTLSClient(host, path.Join(certPath, "cert.pem"), path.Join(certPath, "key.pem"), path.Join(certPath, "ca.pem"))
	} else {
		c, err = api.NewClient(host)
	}
	return c, err
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

	sln.client, err = NewClient()
	return sln, err
}

// NewSolutionFromFile creates and returns a new solution with given YAML file.
func NewSolutionFromFile(file string) (*Solution, error) {
	data, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}
	return NewSolutionFromBytes(data)
}
