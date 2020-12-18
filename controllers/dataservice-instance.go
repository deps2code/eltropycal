package controllers

import (
	dataservices "github.com/eltropycal/dataservices"
)

// DBClient - An instance of the Postgres client
var (
	DataService dataservices.IPostgresClient
)
