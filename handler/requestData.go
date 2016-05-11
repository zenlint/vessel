package handler

import (
	"github.com/containerops/vessel/models"
	"errors"
)

func checkRequestData(reqData *models.PipelineSpecTemplate) (err error) {
	if reqData.Kind != "CCloud"{
		err = errors.New("Kind is wrong!")
	}else if reqData.ApiVersion != "v1"{
		err = errors.New("ApiVersion is wrong!")
	}
	return
}
