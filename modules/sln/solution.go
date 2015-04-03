package sln

import (
	"archive/tar"
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path"
	"regexp"
	"strings"
	"time"

	api "github.com/containerops/rudder"
	"gopkg.in/yaml.v2"

	"github.com/containerops/vessel/modules/utils"
)

// Git represents Git repository information.
type Git struct {
	Auth, Path string
}

// Section represents a solution section.
type Service struct {
	sln *Solution

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

// parseGitPath parses Git path string to a cloneable network address.
func parseGitPath(path string) string {
	if strings.HasPrefix(path, "github.com") {
		return fmt.Sprintf("https://%s.git", path)
	}
	return path
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

	if s.Name == "app" {
		buf.WriteString("WORKDIR $GOPATH/src/" + path.Dir(s.Git.Path) + "\n")
		buf.WriteString("RUN git clone " + parseGitPath(s.Git.Path) + "\n")
		buf.WriteString("WORKDIR " + path.Base(s.Git.Path) + "\n\n")
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

var imageIDPattern = regexp.MustCompile("Successfully built ([0-9a-f]+)")

// FIXME: recursively build.
// Build builds services recursively and returns image ID.
func (s *Service) Build() (string, error) {
	log.Println("Building image...")

	data, err := s.Dockerfile()
	if err != nil {
		return "", fmt.Errorf("generate Dockerfile: %v", err)
	}

	log.Println("Dockerfile:")
	fmt.Println(string(data))

	var (
		input bytes.Buffer
		t     = time.Now()
		so    = utils.NewStreamOutput()
	)

	tr := tar.NewWriter(&input)
	if err = tr.WriteHeader(&tar.Header{
		Name:       "Dockerfile",
		Size:       int64(len(data)),
		ModTime:    t,
		AccessTime: t,
		ChangeTime: t}); err != nil {
		if err != nil {
			return "", fmt.Errorf("write tar header: %v", err)
		}
	} else if _, err = tr.Write(data); err != nil {
		return "", fmt.Errorf("write Dockerfile to tar: %v", err)
	} else if err = tr.Close(); err != nil {
		return "", fmt.Errorf("close tar: %v", err)
	}

	s.sln.Output = so
	if err = s.client.BuildImage(api.BuildImageOption{
		Name:           "app",
		RmTmpContainer: true,
		InputStream:    &input,
		OutputStream:   s.sln.Output,
		RawJSONStream:  true,
	}); err != nil {
		return "", fmt.Errorf("build image: %v", err)
	}

	var imageID string
	if len(so.Events) > 0 {
		e := so.Events[len(so.Events)-1]
		if e["stream"] != nil {
			m := imageIDPattern.FindAllStringSubmatch(e["stream"].(string), 1)
			if m != nil {
				imageID = m[0][1]
			}
		}
	}

	log.Printf("Image ID: %s", imageID)
	return imageID, nil
}

// FIXME: recursively start.
func (s *Service) Start(id string) error {
	// c, err := s.client.CreateContainer(api.CreateContainerOptions{
	// 	Config: &api.Config{
	// 		Image: id,
	// 	},
	// })
	// if err != nil {
	// 	fmt.Println("fail to create container")
	// 	return err
	// }
	// if err = s.client.StartContainer(c.ID, &api.HostConfig{}); err != nil {
	// 	return err
	// }
	// s.client.WaitContainer(c.ID)
	// return s.client.RemoveContainer(api.RemoveContainerOptions{
	// 	ID: c.ID,
	// })
	return nil
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
func (sln *Solution) parseService(name string, sections map[string]map[string]interface{}, marked map[string]bool) (s *Service, err error) {
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
		s.Services[i], err = sln.parseService(name, sections, marked)
		if err != nil {
			return nil, fmt.Errorf("parse '%s' section: %v", name, err)
		}
	}

	s.sln = sln
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

	Output io.Writer
}

// NewClient creates and returns a Docker API client by environment variables.
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
		return nil, fmt.Errorf("parse YAML: %v", err)
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

	sln.Service, err = sln.parseService("app", sections, map[string]bool{})
	if err != nil {
		return nil, fmt.Errorf("parse 'app' section: %v", err)
	}

	sln.client, err = NewClient()
	if err != nil {
		return nil, fmt.Errorf("create new client: %v", err)
	}

	sln.Output = os.Stdout
	return sln, nil
}

// NewSolutionFromFile creates and returns a new solution with given YAML file.
func NewSolutionFromFile(file string) (*Solution, error) {
	data, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, fmt.Errorf("read file: %v", err)
	}
	return NewSolutionFromBytes(data)
}

// Run is used for run the solution and implements Job interface.
func (s *Solution) Run() error {
	imageID, err := s.Build()
	if err != nil {
		return fmt.Errorf("build solution: %v", err)
	} else if err = s.Start(imageID); err != nil {
		return fmt.Errorf("start solution: %v", err)
	}
	return nil
}
