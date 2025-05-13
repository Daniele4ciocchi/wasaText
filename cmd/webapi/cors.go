package main

import (
	"net/http"

	"github.com/gorilla/handlers"
)

// applyCORSHandler applies a CORS policy to the router. CORS stands for Cross-Origin Resource Sharing: it's a security
// feature present in web browsers that blocks JavaScript requests going across different domains if not specified in a
// policy. This function sends the policy of this API server.
func applyCORSHandler(h http.Handler) http.Handler {
	return handlers.CORS(
		handlers.AllowedHeaders([]string{
			"Content-Type",
			"Authorization",
		}),
		handlers.AllowedMethods([]string{"GET", "POST", "OPTIONS", "DELETE", "PUT"}),
		// Do not modify the CORS origin and max age, they are used in the evaluation.
		handlers.AllowedOrigins([]string{"http://localhost:5173"}), // o "*", ma NON in produzione
		handlers.AllowedOrigins([]string{"http://192.168.1.113:5173"}),
		handlers.AllowedOrigins([]string{"http://192.168.1.187:5173"}),
		handlers.AllowedOrigins([]string{"http://https://editor.swagger.io/"}),
		handlers.AllowedOrigins([]string{"*"}),
		handlers.MaxAge(1),
	)(h)
}
