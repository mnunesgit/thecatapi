package controllers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"../../api/models"
	"../formaterror"
	"../responses"
	"github.com/gorilla/mux"
)

func (server *Server) CreateCat(w http.ResponseWriter, r *http.Request) {

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	cat := models.Cat{}
	err = json.Unmarshal(body, &cat)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	cat.Prepare()
	err = cat.Validate()
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	catCreated, err := cat.SaveCat(server.DB)
	if err != nil {
		formattedError := formaterror.FormatError(err.Error())
		responses.ERROR(w, http.StatusInternalServerError, formattedError)
		return
	}
	w.Header().Set("Lacation", fmt.Sprintf("%s%s/%d", r.Host, r.URL.Path, catCreated.ID))
	responses.JSON(w, http.StatusCreated, catCreated)
}

func (server *Server) GetCat(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	pid := vars["name"]
	cat := models.Cat{}

	catReceived, err := cat.FindCatByID(server.DB, pid)
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}
	responses.JSON(w, http.StatusOK, catReceived)
}

func (server *Server) GetCats(w http.ResponseWriter, r *http.Request) {

	cat := models.Cat{}

	cats, err := cat.FindAllCats(server.DB)
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}
	responses.JSON(w, http.StatusOK, cats)
}
