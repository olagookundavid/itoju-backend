package api

import (
	"fmt"
	"net/http"
	"time"
)

type UserPoint struct {
	ID    int64     `json:"id"`
	Date  time.Time `json:"-"`
	Point int64     `json:"point"`
}

type UserPointReponse struct {
	TotalPoints    int `json:"total_point"`
	TodayPoints    int `json:"today_point"`
	ThisWeekPoints int `json:"week_point"`
}

func (app *Application) GetUserTotalPoints(w http.ResponseWriter, r *http.Request) {
	user := app.contextGetUser(r)

	userPoint := make(chan int)
	userDayPoint := make(chan int)
	userMonthPoint := make(chan int)

	defer close(userPoint)
	defer close(userDayPoint)
	defer close(userMonthPoint)

	app.Background(func() {
		app.Models.UserPoint.GetUserTotalPoint(user.ID, userPoint)
	})
	app.Background(func() {
		app.Models.UserPoint.GetUserTotalPoints(user.ID, userDayPoint, userMonthPoint)
	})

	// if err != nil {
	// 	app.serverErrorResponse(w, r, err)
	// 	return
	// }

	userPointResponse := UserPointReponse{
		TotalPoints:    <-userPoint,
		TodayPoints:    <-userDayPoint,
		ThisWeekPoints: <-userMonthPoint,
	}

	env := envelope{
		"message":    "Retrieved User Total Points",
		"user_point": userPointResponse}

	err := app.writeJSON(w, http.StatusOK, env, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *Application) AddUserTotalPoints(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Point int64 `json:"point"`
	}
	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}
	user := app.contextGetUser(r)
	err = app.Models.UserPoint.InsertPoint(user.ID, "", input.Point)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	pointMsg := fmt.Sprintf("Added %d Points to User", input.Point)
	env := envelope{
		"message": pointMsg,
	}

	err = app.writeJSON(w, http.StatusOK, env, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
