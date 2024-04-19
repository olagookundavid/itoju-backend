package api

import (
	"errors"
	"net/http"
	"time"

	"github.com/olagookundavid/itoju/internal/models"
)

func (app *Application) GetUserFoodMetrics(w http.ResponseWriter, r *http.Request) {
	dateString, err := app.readStringParam(r, "date")
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}
	date, err := time.Parse("2006-01-02", dateString)
	if err != nil {
		app.badRequestResponse(w, r, errors.New("invalid date format"))
		return
	}
	user := app.contextGetUser(r)
	foodMetric, err := app.Models.FoodMetric.GetUserFoodMetric(user.ID, date)

	if err != nil {
		switch {
		case errors.Is(err, models.ErrRecordNotFound):
			env := envelope{
				"message":    "Retrieved All Food Metrics for user",
				"foodMetric": foodMetric,
			}
			err = app.writeJSON(w, http.StatusOK, env, nil)
			if err != nil {
				app.serverErrorResponse(w, r, err)
			}

		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	env := envelope{
		"message":    "Retrieved All Food Metrics for user",
		"foodMetric": foodMetric,
	}

	err = app.writeJSON(w, http.StatusOK, env, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *Application) UpdateUserFoodMetrics(w http.ResponseWriter, r *http.Request) {
	user := app.contextGetUser(r)
	dateString, err := app.readStringParam(r, "date")
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}
	date, err := time.Parse("2006-01-02", dateString)
	if err != nil {
		app.badRequestResponse(w, r, errors.New("invalid date format"))
		return
	}

	foodMetric, err := app.Models.FoodMetric.GetUserFoodMetric(user.ID, date)
	if err != nil {
		switch {
		case errors.Is(err, models.ErrRecordNotFound):
			var input struct {
				BreakfastMeal  *string   `json:"breakfast_meal"`
				LunchMeal      *string   `json:"lunch_meal"`
				DinnerMeal     *string   `json:"dinner_meal"`
				BreakfastExtra *string   `json:"breakfast_extra"`
				LunchExtra     *string   `json:"lunch_extra"`
				DinnerExtra    *string   `json:"dinner_extra"`
				BreakfastFruit *string   `json:"breakfast_fruit"`
				LunchFruit     *string   `json:"lunch_fruit"`
				DinnerFruit    *string   `json:"dinner_fruit"`
				BreakfastTags  *[]string `json:"breakfast_tags"`
				LunchTags      *[]string `json:"lunch_tags"`
				DinnerTags     *[]string `json:"dinner_tags"`
				SnackName      *string   `json:"snack_name"`
				SnackTags      *[]string `json:"snack_tags"`
				GlassNo        *int      `json:"glass_no"`
			}

			err := app.readJSON(w, r, &input)
			if err != nil {
				app.badRequestResponse(w, r, err)
				return
			}
			foodMetric := &models.FoodMetric{
				UserID: user.ID,
				Date:   date,
			}

			if input.BreakfastMeal != nil {
				foodMetric.BreakfastMeal = *input.BreakfastMeal
			} else {
				foodMetric.BreakfastMeal = ""
			}

			if input.LunchMeal != nil {
				foodMetric.LunchMeal = *input.LunchMeal
			} else {
				foodMetric.LunchMeal = ""
			}

			if input.DinnerMeal != nil {
				foodMetric.DinnerMeal = *input.DinnerMeal
			} else {
				foodMetric.DinnerMeal = ""
			}

			if input.BreakfastExtra != nil {
				foodMetric.BreakfastExtra = *input.BreakfastExtra
			} else {
				foodMetric.BreakfastExtra = ""
			}

			if input.LunchExtra != nil {
				foodMetric.LunchExtra = *input.LunchExtra
			} else {
				foodMetric.LunchExtra = ""
			}

			if input.DinnerExtra != nil {
				foodMetric.DinnerExtra = *input.DinnerExtra
			} else {
				foodMetric.DinnerExtra = ""
			}

			if input.BreakfastFruit != nil {
				foodMetric.BreakfastFruit = *input.BreakfastFruit
			} else {
				foodMetric.BreakfastFruit = ""
			}

			if input.LunchFruit != nil {
				foodMetric.LunchFruit = *input.LunchFruit
			} else {
				foodMetric.LunchFruit = ""
			}

			if input.DinnerFruit != nil {
				foodMetric.DinnerFruit = *input.DinnerFruit
			} else {
				foodMetric.DinnerFruit = ""
			}

			if input.BreakfastTags != nil {
				foodMetric.BreakfastTags = *input.BreakfastTags
			} else {
				foodMetric.BreakfastTags = []string{}
			}

			if input.LunchTags != nil {
				foodMetric.LunchTags = *input.LunchTags
			} else {
				foodMetric.LunchTags = []string{}
			}

			if input.DinnerTags != nil {
				foodMetric.DinnerTags = *input.DinnerTags
			} else {
				foodMetric.DinnerTags = []string{}
			}

			if input.SnackName != nil {
				foodMetric.SnackName = *input.SnackName
			} else {
				foodMetric.SnackName = ""
			}

			if input.SnackTags != nil {
				foodMetric.SnackTags = *input.SnackTags
			} else {
				foodMetric.SnackTags = []string{}
			}

			if input.GlassNo != nil {
				foodMetric.GlassNo = *input.GlassNo
			} else {
				foodMetric.GlassNo = 0
			}

			err = app.Models.FoodMetric.InsertFoodMetric(foodMetric)

			if err != nil {
				app.serverErrorResponse(w, r, err)
				return
			}
			env := envelope{
				"message": "Successfully updated User Sleep Metrics!",
			}

			err = app.writeJSON(w, http.StatusOK, env, nil)
			if err != nil {
				app.serverErrorResponse(w, r, err)
			}
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	var input struct {
		BreakfastMeal  *string   `json:"breakfast_meal"`
		LunchMeal      *string   `json:"lunch_meal"`
		DinnerMeal     *string   `json:"dinner_meal"`
		BreakfastExtra *string   `json:"breakfast_extra"`
		LunchExtra     *string   `json:"lunch_extra"`
		DinnerExtra    *string   `json:"dinner_extra"`
		BreakfastFruit *string   `json:"breakfast_fruit"`
		LunchFruit     *string   `json:"lunch_fruit"`
		DinnerFruit    *string   `json:"dinner_fruit"`
		BreakfastTags  *[]string `json:"breakfast_tags"`
		LunchTags      *[]string `json:"lunch_tags"`
		DinnerTags     *[]string `json:"dinner_tags"`
		SnackName      *string   `json:"snack_name"`
		SnackTags      *[]string `json:"snack_tags"`
		GlassNo        *int      `json:"glass_no"`
	}
	err = app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if input.BreakfastMeal != nil {
		foodMetric.BreakfastMeal = *input.BreakfastMeal
	}

	if input.LunchMeal != nil {
		foodMetric.LunchMeal = *input.LunchMeal
	}

	if input.DinnerMeal != nil {
		foodMetric.DinnerMeal = *input.DinnerMeal
	}

	if input.BreakfastExtra != nil {
		foodMetric.BreakfastExtra = *input.BreakfastExtra
	}
	if input.LunchExtra != nil {
		foodMetric.LunchExtra = *input.LunchExtra
	}
	if input.DinnerExtra != nil {
		foodMetric.DinnerExtra = *input.DinnerExtra
	}

	if input.BreakfastFruit != nil {
		foodMetric.BreakfastFruit = *input.BreakfastFruit
	}

	if input.LunchFruit != nil {
		foodMetric.LunchFruit = *input.LunchFruit
	}

	if input.DinnerFruit != nil {
		foodMetric.DinnerFruit = *input.DinnerFruit
	}

	if input.BreakfastTags != nil {
		foodMetric.BreakfastTags = *input.BreakfastTags
	}

	if input.LunchTags != nil {
		foodMetric.LunchTags = *input.LunchTags
	}

	if input.DinnerTags != nil {
		foodMetric.DinnerTags = *input.DinnerTags
	}
	if input.SnackName != nil {
		foodMetric.SnackName = *input.SnackName
	}

	if input.SnackTags != nil {
		foodMetric.SnackTags = *input.SnackTags
	}

	if input.GlassNo != nil {
		foodMetric.GlassNo = *input.GlassNo
	}

	err = app.Models.FoodMetric.UpdateFoodMetric(foodMetric)
	if err != nil {
		switch {
		case errors.Is(err, models.ErrEditConflict):
			app.editConflictResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}
	env := envelope{
		"message": "Successfully updated User Food Metrics",
	}
	err = app.writeJSON(w, http.StatusOK, env, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
