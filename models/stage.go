package models

type Stage struct {
	Name                string
	Namespace           string
	Replicas            int64
	Image               string
	Port                int64
	StatusCheckLink     string
	StatusCheckInterval int64
	StatusCheckCount    int64
	EnvName             string
	EnvValue            string
	Dependence          []string
	Status              string
}