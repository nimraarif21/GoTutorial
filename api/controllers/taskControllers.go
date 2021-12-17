// TaskControllers.go
package controllers

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"

	"github.com/gotutorial/api/models"
	"github.com/gotutorial/api/responses"
)

// CreateTask parses request, validates data and saves the new Task
func (a *App) CreateTask(w http.ResponseWriter, r *http.Request) {
    var resp = map[string]interface{}{"status": "success", "message": "Task successfully created"}

    user := r.Context().Value("userID").(float64)
    Task := &models.Task{}
    body, err := ioutil.ReadAll(r.Body)
    if err != nil {
        responses.ERROR(w, http.StatusBadRequest, err)
        return
    }

    err = json.Unmarshal(body, &Task)
    if err != nil {
        responses.ERROR(w, http.StatusBadRequest, err)
        return
    }

    Task.Prepare() // strip away any white spaces

    if err = Task.Validate(); err != nil {
        responses.ERROR(w, http.StatusBadRequest, err)
        return
    }

    if vne, _ := Task.GetTask(a.DB); vne != nil {
        resp["status"] = "failed"
        resp["message"] = "Task already registered, please choose another name"
        responses.JSON(w, http.StatusBadRequest, resp)
        return
    }

    Task.UserID = uint(user)

    TaskCreated, err := Task.Save(a.DB)
    if err != nil {
        responses.ERROR(w, http.StatusBadRequest, err)
        return
    }

    resp["Task"] = TaskCreated
    responses.JSON(w, http.StatusCreated, resp)
    return
}

func (a *App) GetTasks(w http.ResponseWriter, r *http.Request) {
    user := r.Context().Value("userID").(float64)
    userID := uint(user)
    Tasks, err := models.GetTasks(userID, a.DB)
    if err != nil {
        responses.ERROR(w, http.StatusInternalServerError, err)
        return
    }
    responses.JSON(w, http.StatusOK, Tasks)
    return
}

func (a *App) UpdateTask(w http.ResponseWriter, r *http.Request) {
    var resp = map[string]interface{}{"status": "success", "message": "Task updated successfully"}

    vars := mux.Vars(r)

    user := r.Context().Value("userID").(float64)
    userID := uint(user)

    id, _ := strconv.Atoi(vars["id"])

    Task, err := models.GetTaskById(id, a.DB)
    if err != nil {
        responses.ERROR(w, http.StatusBadRequest, err)
        return
    }

    if Task.UserID != userID {
        resp["status"] = "failed"
        resp["message"] = "Unauthorized Task update"
        responses.JSON(w, http.StatusUnauthorized, resp)
        return
    }

    body, err := ioutil.ReadAll(r.Body)
    if err != nil {
        responses.ERROR(w, http.StatusBadRequest, err)
        return
    }

    TaskUpdate := models.Task{}
    if err = json.Unmarshal(body, &TaskUpdate); err != nil {
        responses.ERROR(w, http.StatusBadRequest, err)
        return
    }

    TaskUpdate.Prepare()

    _, err = TaskUpdate.UpdateTask(id, a.DB)
    if err != nil {
        responses.ERROR(w, http.StatusInternalServerError, err)
        return
    }

    responses.JSON(w, http.StatusOK, resp)
    return
}

func (a *App) DeleteTask(w http.ResponseWriter, r *http.Request) {
    var resp = map[string]interface{}{"status": "success", "message": "Task deleted successfully"}

    vars := mux.Vars(r)

    user := r.Context().Value("userID").(float64)
    userID := uint(user)

    id, _ := strconv.Atoi(vars["id"])

    Task, err := models.GetTaskById(id, a.DB)
    if err != nil {
        responses.ERROR(w, http.StatusBadRequest, err)
        return
    }

    if Task.UserID!= userID {
        resp["status"] = "failed"
        resp["message"] = "Unauthorized Task delete"
        responses.JSON(w, http.StatusUnauthorized, resp)
        return
    }

    err = models.DeleteTask(id, a.DB)
    if err != nil {
        responses.ERROR(w, http.StatusInternalServerError, err)
        return
    }
    responses.JSON(w, http.StatusOK, resp)
    return
}