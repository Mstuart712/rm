package controllers

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/Mstuart712/rm/api/auth"
	"github.com/Mstuart712/rm/api/models"
	"github.com/Mstuart712/rm/api/responses"
	"github.com/Mstuart712/rm/api/utils/formaterror"
	"github.com/gorilla/mux"
)

func (server *Server) CreateCharacter(w http.ResponseWriter, r *http.Request) {

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	character := models.Character{}
	err = json.Unmarshal(body, &character)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	character.Prepare()
	err = character.Validate()
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	uid, err := auth.ExtractTokenID(r)
	if err != nil {
		responses.ERROR(w, http.StatusUnauthorized, errors.New("Unauthorized"))
		return
	}
	if uid != character.OwnerID {
		responses.ERROR(w, http.StatusUnauthorized, errors.New(http.StatusText(http.StatusUnauthorized)))
		return
	}
	characterCreated, err := character.SaveCharacter(server.DB)
	if err != nil {
		formattedError := formaterror.FormatError(err.Error())
		responses.ERROR(w, http.StatusInternalServerError, formattedError)
		return
	}
	w.Header().Set("Location", fmt.Sprintf("%s%s/%d", r.Host, r.URL.Path, characterCreated.ID))
	responses.JSON(w, http.StatusCreated, characterCreated)
}

func (server *Server) GetCharacters(w http.ResponseWriter, r *http.Request) {

	character := models.Character{}

	characters, err := character.FindAllCharacters(server.DB)
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}
	responses.JSON(w, http.StatusOK, characters)
}

func (server *Server) GetCharacter(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	pid, err := strconv.ParseUint(vars["id"], 10, 64)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}
	character := models.Character{}

	characterReceived, err := character.FindCharacterByID(server.DB, pid)
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}
	responses.JSON(w, http.StatusOK, characterReceived)
}

func (server *Server) UpdateCharacter(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)

	// Check if the character id is valid
	pid, err := strconv.ParseUint(vars["id"], 10, 64)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}

	//CHeck if the auth token is valid and  get the user id from it
	uid, err := auth.ExtractTokenID(r)
	if err != nil {
		responses.ERROR(w, http.StatusUnauthorized, errors.New("Unauthorized"))
		return
	}

	// Check if the Character exist
	character := models.Character{}
	err = server.DB.Debug().Model(models.Character{}).Where("id = ?", pid).Take(&character).Error
	if err != nil {
		responses.ERROR(w, http.StatusNotFound, errors.New("Character not found"))
		return
	}

	// If a user attempt to update a Character not belonging to him
	if uid != character.OwnerID {
		responses.ERROR(w, http.StatusUnauthorized, errors.New("Unauthorized"))
		return
	}
	// Read the data posted
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	// Start processing the request data
	characterUpdate := models.Character{}
	err = json.Unmarshal(body, &characterUpdate)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	//Also check if the request user id is equal to the one gotten from token
	if uid != characterUpdate.OwnerID {
		responses.ERROR(w, http.StatusUnauthorized, errors.New("Unauthorized"))
		return
	}

	characterUpdate.Prepare()
	err = characterUpdate.Validate()
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	characterUpdate.ID = character.ID //this is important to tell the model the Character id to update, the other update field are set above

	characterUpdated, err := characterUpdate.UpdateACharacter(server.DB)

	if err != nil {
		formattedError := formaterror.FormatError(err.Error())
		responses.ERROR(w, http.StatusInternalServerError, formattedError)
		return
	}
	responses.JSON(w, http.StatusOK, characterUpdated)
}

func (server *Server) DeleteCharacter(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)

	// Is a valid Character id given to us?
	pid, err := strconv.ParseUint(vars["id"], 10, 64)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}

	// Is this user authenticated?
	uid, err := auth.ExtractTokenID(r)
	if err != nil {
		responses.ERROR(w, http.StatusUnauthorized, errors.New("Unauthorized"))
		return
	}

	// Check if the Character exist
	character := models.Character{}
	err = server.DB.Debug().Model(models.Character{}).Where("id = ?", pid).Take(&character).Error
	if err != nil {
		responses.ERROR(w, http.StatusNotFound, errors.New("Unauthorized"))
		return
	}

	// Is the authenticated user, the owner of this Character?
	if uid != character.OwnerID {
		responses.ERROR(w, http.StatusUnauthorized, errors.New("Unauthorized"))
		return
	}
	_, err = character.DeleteACharacter(server.DB, pid, uid)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}
	w.Header().Set("Entity", fmt.Sprintf("%d", pid))
	responses.JSON(w, http.StatusNoContent, "")
}
