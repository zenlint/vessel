package dependence

import (
	"testing"
	"github.com/containerops/vessel/models"
	"encoding/json"
	"log"
)

func Test_ParsePipelineTemp(t *testing.T){
	str := jsonStr()
	pipelineTemp := &models.PipelineSpecTemplate{}
	err := json.Unmarshal([]byte(str),pipelineTemp)
	if err != nil {
		log.Println(err)
	}
	log.Println(pipelineTemp)
	stageMap,err := ParsePipelineTemplate(pipelineTemp)
	if err != nil {
		log.Println(err)
	}
	log.Println(stageMap)
	bytes,err := json.Marshal(stageMap)
	if err != nil {
		log.Println(err)
	}
	log.Println(string(bytes))
}

func jsonStr() string {
	return `{
	    "kind": "CCloud",
	    "apiVersion": "v1",
	    "metadata": {
		"name": "guestbook",
		"namespace": "guestbook",
		"timeoutDuration": 60
	    },
	    "spec": [{
		"name": "A",
		"replicas": 1,
		"dependence": "",
		"statusCheckLink": "/health",
		"statusCheckInterval": 0,
		"statusCheckCount": 0,
		"image": "containerops.me:8080/opensource/redis",
		"port": 6379,
		"envName": "",
		"envValue": ""
	    }, {
		"name": "B",
		"replicas": 2,
		"dependence": "",
		"statusCheckLink": "/health",
		"statusCheckInterval": 0,
		"statusCheckCount": 0,
		"image": "containerops.me:8080/opensource/redisslavev1",
		"port": 6379,
		"envName": "",
		"envValue": ""
	    }, {
		"name": "C",
		"replicas": 3,
		"dependence": "",
		"statusCheckLink": "/health",
		"statusCheckInterval": 0,
		"statusCheckCount": 0,
		"image": "containerops.me:8080/opensource/aaaaaaa",
		"port": 80,
		"envName": "GET_HOSTS_FROM",
		"envValue": "dns"
	    },{
		"name": "D",
		"replicas": 1,
		"dependence": "A",
		"statusCheckLink": "/health",
		"statusCheckInterval": 0,
		"statusCheckCount": 0,
		"image": "containerops.me:8080/opensource/redis",
		"port": 6379,
		"envName": "",
		"envValue": ""
	    }, {
		"name": "E",
		"replicas": 2,
		"dependence": "A",
		"statusCheckLink": "/health",
		"statusCheckInterval": 0,
		"statusCheckCount": 0,
		"image": "containerops.me:8080/opensource/redisslavev1",
		"port": 6379,
		"envName": "",
		"envValue": ""
	    }, {
		"name": "F",
		"replicas": 3,
		"dependence": "B,I",
		"statusCheckLink": "/health",
		"statusCheckInterval": 0,
		"statusCheckCount": 0,
		"image": "containerops.me:8080/opensource/aaaaaaa",
		"port": 80,
		"envName": "GET_HOSTS_FROM",
		"envValue": "dns"
	    },{
		"name": "G",
		"replicas": 1,
		"dependence": "C,I",
		"statusCheckLink": "/health",
		"statusCheckInterval": 0,
		"statusCheckCount": 0,
		"image": "containerops.me:8080/opensource/redis",
		"port": 6379,
		"envName": "",
		"envValue": ""
	    }, {
		"name": "H",
		"replicas": 2,
		"dependence": "C",
		"statusCheckLink": "/health",
		"statusCheckInterval": 0,
		"statusCheckCount": 0,
		"image": "containerops.me:8080/opensource/redisslavev1",
		"port": 6379,
		"envName": "",
		"envValue": ""
	    }, {
		"name": "I",
		"replicas": 3,
		"dependence": "B,C,D,E,F",
		"statusCheckLink": "/health",
		"statusCheckInterval": 0,
		"statusCheckCount": 0,
		"image": "containerops.me:8080/opensource/aaaaaaa",
		"port": 80,
		"envName": "GET_HOSTS_FROM",
		"envValue": "dns"
	    }]
	}`
}