package view

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/gonzabosio/res-manager/model"
	"github.com/maxence-charriere/go-app/v10/pkg/app"
)

type JoinTeam struct {
	app.Compo

	user        model.User
	teamName    string
	password    string
	errMessage  string
	accessToken string
}

func (j *JoinTeam) OnMount(ctx app.Context) {
	if err := ctx.SessionStorage().Get("user", &j.user); err != nil {
		app.Log("Could not get user from local storage")
	}
	if err := ctx.LocalStorage().Get("access-token", &j.accessToken); err != nil {
		app.Log(err)
		ctx.Navigate("/")
	}
}

func (j *JoinTeam) Render() app.UI {
	return app.Div().Body(
		app.Text("Join Team"),
		app.Form().Body(
			app.Input().Type("text").Value(j.teamName).MaxLength(30).
				Placeholder("Team name").
				AutoFocus(true).
				OnChange(j.ValueTo(&j.teamName)),
			app.Input().Type("password").Value(j.password).
				Placeholder("Password").
				OnChange(j.ValueTo(&j.password)),
		),
		app.Button().Text("Join").OnClick(j.joinAction).Class("global-btn"),
		app.P().Text(j.errMessage).Class("err-message"),
	)
}

func (j *JoinTeam) joinAction(ctx app.Context, e app.Event) {
	if j.teamName == "" || j.password == "" {
		j.errMessage = "Empty team name or password field"
	} else {
		ctx.Async(func() {
			req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("%v/join-team", app.Getenv("BACK_URL")), strings.NewReader(fmt.Sprintf(
				`{"name":"%v","password":"%v"}`,
				j.teamName, j.password)))
			if err != nil {
				app.Log(err)
				j.errMessage = "Could not build the team request"
				return
			}
			req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", j.accessToken))
			req.Header.Add("Content-Type", "application/json")
			client := http.Client{}
			res, err := client.Do(req)
			if err != nil {
				log.Println(err)
				return
			}
			defer res.Body.Close()

			b, err := io.ReadAll(res.Body)
			if err != nil {
				log.Println("Failed to send request:", err)
				return
			}
			if res.StatusCode == http.StatusOK {
				var resBody okResponseBody
				if err := json.Unmarshal(b, &resBody); err != nil {
					app.Log(err)
					j.errMessage = "Failed to parse json"
					return
				}
				if err := ctx.LocalStorage().Set("teamName", j.teamName); err != nil {
					app.Log(err)
				}
				teamIDstr := strconv.FormatInt(resBody.TeamID, 10)
				if err := ctx.LocalStorage().Set("teamID", teamIDstr); err != nil {
					app.Log(err)
				}
				ctx.Async(func() {
					req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("%v/participant", app.Getenv("BACK_URL")), strings.NewReader(fmt.Sprintf(`
					{"admin":%v,"user_id":%v,"team_id":%v}
					`, false, j.user.Id, resBody.TeamID,
					)))
					if err != nil {
						app.Log(err)
						j.errMessage = "Could not build the participant request"
						return
					}
					req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", j.accessToken))
					req.Header.Add("Content-Type", "application/json")
					client := http.Client{}
					res, err := client.Do(req)
					if err != nil {
						app.Log(err)
						j.errMessage = "Failed to add participant"
						return
					}
					defer res.Body.Close()
					b, err := io.ReadAll(res.Body)
					if err != nil {
						app.Log(err)
						j.errMessage = "Failed reading the partipant response"
						return
					}

					if res.StatusCode == http.StatusOK {
						var body participantResponse
						if err = json.Unmarshal(b, &body); err != nil {
							app.Log(err)
							j.errMessage = "Could not parse the participant response"
							return
						}
						ctx.SessionStorage().Set("admin", body.Participant.Admin)
						ctx.Navigate("dashboard")
					} else if res.StatusCode == http.StatusUnauthorized {
						ctx.LocalStorage().Del("access-token")
						ctx.Navigate("/")
					} else {
						var body errResponseBody
						if err = json.Unmarshal(b, &body); err != nil {
							app.Log(err)
							j.errMessage = body.Message
							return
						}
						app.Log(body.Err)
						j.errMessage = body.Message
					}
				})
			} else if res.StatusCode == http.StatusUnauthorized {
				ctx.LocalStorage().Del("access-token")
				ctx.Navigate("/")
			} else {
				var resBody errResponseBody
				if err := json.Unmarshal(b, &resBody); err != nil {
					j.errMessage = "Failed to parse json"
					app.Log(resBody.Err)
					return
				}
				app.Log(resBody.Err)
				j.errMessage = resBody.Message
				ctx.Dispatch(func(ctx app.Context) {})
			}
		})
	}
}
