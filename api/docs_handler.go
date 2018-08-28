package api

import (
	"fmt"
	"net/http"

	"github.com/topfreegames/pitaya"

	"github.com/topfreegames/pitaya-admin/usecases"

	"github.com/topfreegames/pitaya-admin/constants"
)

// DocsHandler handle docs
type DocsHandler struct {
	App             *App
	docsRemoteRoute string
}

// NewDocsHandler create new docs handler
func NewDocsHandler(a *App) *DocsHandler {
	m := &DocsHandler{
		App:             a,
		docsRemoteRoute: a.Config.GetString("pitayaadmin.docs.suffix"),
	}
	return m
}

// ServeHTTP method. Docs will only work if game has a remote for docs
func (s *DocsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	v := r.URL.Query()
	param := v.Get("type")
	roh := v.Get("t")
	routeParam := v.Get("route")

	if param != "" {
		if roh != "" {
			if routeParam != "" {
				ret, err := usecases.GetMethodDoc(param, roh, routeParam)
				if err != nil {
					WriteError(w, 500, "failed to get doc by route", err)
					return
				}
				bts, _ := pitaya.GetSerializer().Marshal(ret)
				Write(w, 200, string(bts))
				return
			}

			ret, err := usecases.GetDocForHandlersOrRemotes(param, roh)
			if err != nil {
				errMsg := fmt.Sprintf("failed to get docs %s", roh)
				WriteError(w, 500, errMsg, err)
				return
			}
			bts, _ := pitaya.GetSerializer().Marshal(ret)
			Write(w, 200, string(bts))
			return
		}
		ret, err := usecases.GetDocumentationForServerType(param)
		if err != nil {
			WriteError(w, 500, "failed to get docs", err)
			return
		}
		bts, _ := pitaya.GetSerializer().Marshal(ret)
		Write(w, 200, string(bts))
		return
	}

	WriteError(w, http.StatusBadRequest, "must specify at least server type", constants.ErrNoServerType)
}
