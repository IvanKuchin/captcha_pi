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

type login_user struct {
	// User login
	// in: body
	// required: true
	// type: string
	Login string
	// Password MD5 hash (NOT the actual password)
	// in: body
	// required: true
	// type: string
	Password_hash string
}

// swagger:parameters loginID
type login_user_wrapper struct {
	// in: body
	Body login_user
}

type login_response struct {
	Result      string
	Description string
	Sessid      string
}
