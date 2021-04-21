package service

import (
	"net/http"
)

type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

type Routes []Route

var routes = Routes{
	Route{
		"Login",
		"POST",
		"/login",
		HandleLogin,
	},
	Route{
		"Teacher Registration",
		"POST",
		"/teacher/register",
		HandleRegistration,
	},
	Route{
		"Student Registration",
		"POST",
		"/student/register",
		HandleRegistration,
	},
	Route{
		"ValidateToken",
		"POST",
		"/validateToken",
		HandleValidateToken,
	},
}
