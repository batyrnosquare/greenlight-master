package main

import (
	"errors"
	"fmt"
	"github.com/shynggys9219/greenlight/internal/data"
	"net/http"
	"time"
)

func (app *application) createModuleInfo(w http.ResponseWriter, r *http.Request) {
	var input struct {
		ModuleName     string        `json:"moduleName"`
		ModuleDuration time.Duration `json:"moduleDuration"`
		ExamType       string        `json:"examType"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.errorResponse(w, r, http.StatusBadRequest, err.Error())
	}
	moduleInfo := &data.ModuleInfo{
		ModuleName:     input.ModuleName,
		ModuleDuration: input.ModuleDuration,
		ExamType:       input.ExamType,
	}

	err = app.models.ModuleInfos.Insert(moduleInfo)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return

	}

	headers := make(http.Header)
	headers.Set("Location", fmt.Sprintf("/v1/module-info/%d", moduleInfo.ID))

	err = app.writeJSON(w, http.StatusCreated, envelope{"module_info": moduleInfo}, headers)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}

}

func (app *application) getModuleInfo(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
	}

	moduleInfos, err := app.models.ModuleInfos.Get(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"module_infos": moduleInfos}, nil)
	if err != nil {
		return
	}
}

func (app *application) getLastFiftyModuleInfo(w http.ResponseWriter, r *http.Request) {
	moduleInfos := app.models.ModuleInfos.GetLatestFifty()

	err := app.writeJSON(w, http.StatusOK, envelope{"module_infos": moduleInfos}, nil)
	if err != nil {
		return
	}
}

func (app *application) editModuleInfo(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
	}

	moduleInfo, err := app.models.ModuleInfos.Get(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	var input struct {
		ModuleName     string        `json:"moduleName"`
		ModuleDuration time.Duration `json:"moduleDuration"`
		ExamType       string        `json:"examType"`
	}

	err = app.readJSON(w, r, &input)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	moduleInfo.ModuleName = input.ModuleName
	moduleInfo.ModuleDuration = input.ModuleDuration
	moduleInfo.ExamType = input.ExamType

	err = app.models.ModuleInfos.Update(moduleInfo)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"module_info": moduleInfo}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) deleteModuleInfo(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
	}

	err = app.models.ModuleInfos.Delete(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}
	err = app.writeJSON(w, http.StatusOK, envelope{"message": "module info successfully deleted"}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
