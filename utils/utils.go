package utils

import (
	"encoding/json"
	"strings"

	"github.com/satori/go.uuid"
)

// Spec :
type Spec struct {
	Spec []Relation `json:"spec"`
}

// Relation :
type Relation struct {
	Name         string   `json:"name"`
	Dependencies []string `json:"dependence"`
	Count        int64    `json:"-"`
	K8SInfo      string   `json:"k8sInfo"`
}

// GenerateDependenceMap :
func GenerateDependenceMap(relation string) (isLegal bool, reason string, dependenceMap map[string][]string) {
	isLegal = false
	reason = ""
	dependenceMap = make(map[string][]string, 0)

	if IsNameRepeat(relation) {
		return isLegal, "this relation has repeat name!", dependenceMap
	}

	if isLegal, dependenceMap = IsDAG(relation); !isLegal {
		return isLegal, "this relation is not a DAG !", dependenceMap
	}

	return isLegal, reason, dependenceMap
}

// IsNameRepeat : test is the relation string has repeat name
// @return true if has repeat name , false if not
func IsNameRepeat(relation string) (isRepeat bool) {
	isRepeat = false

	relations := new(Spec)
	json.Unmarshal([]byte(relation), &relations)

	nameMap := make(map[string]string)

	for _, v := range relations.Spec {
		if _, ok := nameMap[v.Name]; ok {
			isRepeat = true
			break
		} else {
			nameMap[v.Name] = v.Name
		}
	}

	return isRepeat
}

// IsDAG : test is the relation string is a DAG(Directed Acyclic Graph),if not return dependenceMap
// dependenceMap   map[name][{"form1","form2"},{"to1","to2"}]
func IsDAG(relation string) (isDAG bool, dependenceMap map[string][]string) {
	isDAG = false
	dependenceMap = make(map[string][]string, 0)

	relations := new(Spec)
	json.Unmarshal([]byte(relation), &relations)

	relationMap := make(map[string]*Relation)

	for _, v := range relations.Spec {
		if _, ok := relationMap[v.Name]; !ok {
			relationMap[v.Name] = new(Relation)
			relationMap[v.Name].Name = v.Name
		}
		relationMap[v.Name].Dependencies = append(relationMap[v.Name].Dependencies, v.Dependencies...)
		if _, ok := dependenceMap[v.Name]; !ok {
			dependenceMap[v.Name] = make([]string, 3)
		}
		dependenceMap[v.Name][0] = strings.Join(append(strings.Split(dependenceMap[v.Name][0], ","), v.Dependencies...), ",")
		dependenceMap[v.Name][2] = v.K8SInfo

		for _, in := range v.Dependencies {
			if _, ok := relationMap[in]; !ok {
				relationMap[in] = new(Relation)
				relationMap[in].Name = in
			}
			relationMap[in].Count++

			if _, ok := dependenceMap[in]; !ok {
				dependenceMap[in] = make([]string, 2)
			}
			dependenceMap[in][1] = strings.Join(append(strings.Split(dependenceMap[in][1], ","), v.Name), ",")

		}
	}

	finish := 0
	for true {
		temp := 0
		for _, relation := range relationMap {
			if relation.Count == 0 {
				finish++
				for _, out := range relation.Dependencies {
					relationMap[out].Count--
				}
				relation.Count = -1
			} else if relation.Count == -1 {
				temp++
			}
		}

		if temp == finish || finish == len(relationMap) {
			break
		}
	}

	if finish == len(relationMap) {
		isDAG = true
	}

	return isDAG, dependenceMap
}

// GenerateUUID :
func GenerateUUID() string {
	uid := uuid.NewV1()
	uids := strings.Split(uid.String(), "-")
	return uids[0] + uids[1] + uids[2] + uids[4] + uids[3]
}
