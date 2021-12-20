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
	"github.com/gotutorial/utils"
)

func (a *App) AssignTask(w http.ResponseWriter, r *http.Request) {
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
        resp["message"] = "Unauthorized Task assignment"
        responses.JSON(w, http.StatusUnauthorized, resp)
        return
    }
    body, err := ioutil.ReadAll(r.Body)
    if err != nil {
        responses.ERROR(w, http.StatusBadRequest, err)
        return
    }
		
		assigntask := models.PendingTask{}
		if err = json.Unmarshal(body, &assigntask); err!= nil{
			responses.ERROR(w, http.StatusInternalServerError, err)
			return
		}
		if assigntask.Email != ""{
			assigned_user, err := models.Getuserbyemail(assigntask.Email, a.DB)
			if err == nil {
					TaskUpdate := models.Task{}
					TaskUpdate.AssignedTo = *assigned_user
					a.DB.Save(&TaskUpdate)
					_, err = TaskUpdate.UpdateTask(id, a.DB)
					if err != nil {
							responses.ERROR(w, http.StatusInternalServerError, err)
							return
					}
				} else{
				go utils.SendEmail(assigntask.Email)
				assigntask.Task = *Task
				assigntask.Save(a.DB)
			}
			responses.JSON(w, http.StatusOK, resp)
      return
			}
   		responses.ERROR(w, http.StatusInternalServerError, err)
			return
}
