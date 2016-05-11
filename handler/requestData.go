package handler

import (
	"github.com/containerops/vessel/models"
	"github.com/containerops/vessel/setting"
	"errors"
)

func checkRequestData(reqData *models.PipelineSpecTemplate) (err error) {
	if reqData.Kind != setting.Global.Kind{
		err = errors.New("Kind is wrong!")
	}else if reqData.ApiVersion != setting.Global.ApiVersion{
		err = errors.New("ApiVersion is wrong!")
	}
	return
}
