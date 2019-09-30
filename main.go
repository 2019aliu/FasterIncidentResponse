/*
Faster Incident Response (FasterIR, FaIR) is to be a stand alone incident response tracking system.
With the use of Security Orchestration, Automation and Response (SOAR) systems, there tracking of incidents no longer should be kept in the Security Information and Event Management (SIEM) system.

To edit the main database, please see the file mongo/mongoFunctions.go
*/
package main

import (
	mongoFunctions "fir/mongo"
	"fir/routes"

	"github.com/gin-gonic/contrib/static"
)

func main() {
	// Setup MongoDB connection
	client := mongoFunctions.GetClient()
	mongoFunctions.CheckDBConnection(client)

	router := routes.SetupRouter()

	router.Use(static.Serve("/", static.LocalFile("./views", true)))
	router.Use(static.Serve("/js", static.LocalFile("./js", true)))
	router.Use(static.Serve("/css", static.LocalFile("./css", true)))

	// Run server on port 8080
	router.Run()
}
