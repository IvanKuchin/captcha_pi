// Package handlers for the RESTful Server
//
// Documentation for REST API
//
//  Schemes: http, https
//  BasePath: /
//  Version: 1.0.7
//
//  Consumes:
//  - application/json
//
//  Produces:
//  - application/json
//
// swagger:meta
package apihandlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	wrapper "github.com/ivankuchin/captcha_pi/captcha-wrapper"
	"github.com/ivankuchin/timecard.ru-api/logs"
)

// swagger:route GET /captcha/getnew Captcha newCaptcha
// Create new captcha
//
// Consumes:
// - application/json
//
// Produces:
// - text/plain
//
// Schemes: http, https
//
// responses:
// 200: captchaID
// 404: notFoundWrapper
// 400: badRequestWrapper

func GetNewHandler(w http.ResponseWriter, r *http.Request) {
	tID := generateTraceID()

	sessid, err := r.Cookie("sessid")
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		logs.Sugar.Errorw(err.Error(), "traceID", tID)
		return
	}

	if sessid.Value == "" {
		w.WriteHeader(http.StatusBadRequest)
		logs.Sugar.Infow("cookie is empty", "traceID", tID)
		return
	}

	captcha_id, err := wrapper.MakeNew(tID)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		logs.Sugar.Errorw(err.Error(), "traceID", tID)

		return
	}

	err = Put(sessid.Value, captcha_id)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		logs.Sugar.Errorw(err.Error(), "traceID", tID)
	}

	logs.Sugar.Debugw("captcha "+captcha_id+" created for sessid "+sessid.Value, "captcha store len", len(cs.m), "traceID", tID)

	response := map[string]string{"captcha_id": captcha_id}
	js, err := json.Marshal(response)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		logs.Sugar.Errorw("json.marshal", "traceID", tID)
	}

	fmt.Fprint(w, string(js))
}
