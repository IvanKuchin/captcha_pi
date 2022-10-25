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
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	wrapper "github.com/ivankuchin/captcha_pi/captcha-wrapper"
	"github.com/ivankuchin/timecard.ru-api/logs"
)

// swagger:route GET /captcha/verify Captcha verify
// Verify captcha saved from session
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

func VerifyHandler(w http.ResponseWriter, r *http.Request) {
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

	captcha_id, err := Get(sessid.Value, tID)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		logs.Sugar.Errorw(err.Error(), "traceID", tID)
		return
	}

	if captcha_id == "" {
		w.WriteHeader(http.StatusNotFound)
		logs.Sugar.Infow("captcha not found", "sessid", sessid.Value, "traceID", tID)
		return
	}

	vars := mux.Vars(r)
	solution := vars["solution"]
	result, err := wrapper.Verify(captcha_id, solution)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		logs.Sugar.Errorw(err.Error(), "traceID", tID)
		return
	}

	err = Delete(sessid.Value, tID)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		logs.Sugar.Errorw(err.Error(), "traceID", tID)
		return
	}

	logs.Sugar.Debugw("captcha verification result", "sessid ", sessid.Value, "solution", solution, "captcha_id", captcha_id, "result", strconv.FormatBool(result), "captcha store len", len(cs.m), "traceID", tID)

	fmt.Fprint(w, result)
}
