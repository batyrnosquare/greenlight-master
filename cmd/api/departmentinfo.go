package main

import (
	"errors"
	"fmt"
	"github.com/shynggys9219/greenlight/internal/data"
	"net/http"
)

func (app *application) createDepInfoHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		DepartmentName     string `json:"departmentName"`
		DepartmentDirector string `json:"departmentDirector"`
		StaffQuantity      int    `json:"staffQuantity"`
		ModuleID           int    `json:"moduleID"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.errorResponse(w, r, http.StatusBadRequest, err.Error())
	}
	depInfo := &data.DepartmentInfo{
		DepartmentName:     input.DepartmentName,
		DepartmentDirector: input.DepartmentDirector,
		StaffQuantity:      input.StaffQuantity,
		ModuleID:           input.ModuleID,
	}

	err = app.models.DepInfos.Insert(depInfo)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return

	}

	headers := make(http.Header)
	headers.Set("Location", fmt.Sprintf("/v1/department-info/%d", depInfo.ID))

	err = app.writeJSON(w, http.StatusCreated, envelope{"department_info": depInfo}, headers)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) getDepInfoHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
	}

	depInfos, err := app.models.DepInfos.Get(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"department_infos": depInfos}, nil)
	if err != nil {
		return
	}
}
