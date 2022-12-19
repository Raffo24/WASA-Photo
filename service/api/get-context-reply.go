package api

import (
	"encoding/json"
	"net/http"
	"wasaPhoto/service/api/reqcontext"

	"github.com/julienschmidt/httprouter"
)

// getContextReply is an example of HTTP endpoint that returns "Hello World!" as a plain text. The signature of this
// handler accepts a reqcontext.RequestContext (see httpRouterHandler).
func (rt *_router) getContextReply(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	w.Header().Set("Content-Type", "application/json")
	if json.NewEncoder(w).Encode("Hello World") != nil {
		w.WriteHeader(500)
		logerr(w.Write([]byte("Internal Server Error")))
	}
}
