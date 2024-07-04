package main

import (
	"net/http"
	"fmt"
)

func (app *application) healthcheckHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "status: available")
	fmt.Fprintf(w, "environemnt: %s\n", app.config.env)
	fmt.Fprintf(w, "version: %s\n", version)
}
