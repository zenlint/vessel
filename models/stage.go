package models

type Stage struct {
	Name                string
	Namespace           string
	Replicas            uint64
	Image               string
	Port                uint64
	StatusCheckLink     string
	StatusCheckInterval uint64
	StatusCheckCount    uint64
	EnvName             string
	EnvValue            string
	Dependence          []string
	Status              string
}