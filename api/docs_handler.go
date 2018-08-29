package api

import (
	"fmt"
	"net/http"

	"github.com/topfreegames/pitaya"

	"github.com/topfreegames/pitaya-admin/constants"
	"github.com/topfreegames/pitaya-admin/docs"
)

// DocsHandler handle documentation routes
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

// ServeHTTP method. Docs will only work if server has a remote for autodoc
func (s *DocsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	v := r.URL.Query()
	serverType := v.Get("type")
	remoteOrHandler := v.Get("methodtype")
	routeParam := v.Get("route")

	if serverType == "" {
		WriteError(w, http.StatusBadRequest, "must specify at least server type", constants.ErrNoServerType)
		return
	}
	if remoteOrHandler != "" {
		if routeParam != "" {
			ret, err := docs.GetMethodDoc(serverType, remoteOrHandler, routeParam, s.App.GetRemoteDocsRoute())
			if err != nil {
				WriteError(w, http.StatusInternalServerError, "failed to get doc by route", err)
				return
			}
			bts, _ := pitaya.GetSerializer().Marshal(ret)
			Write(w, http.StatusOK, string(bts))
			return
		}

		ret, err := docs.GetDocForHandlersOrRemotes(serverType, remoteOrHandler, s.App.GetRemoteDocsRoute())
		if err != nil {
			errMsg := fmt.Sprintf("failed to get docs %s", remoteOrHandler)
			WriteError(w, http.StatusInternalServerError, errMsg, err)
			return
		}
		bts, _ := pitaya.GetSerializer().Marshal(ret)
		Write(w, http.StatusOK, string(bts))
		return
	}
	ret, err := docs.GetDocumentationForServerType(serverType, s.App.GetRemoteDocsRoute())
	if err != nil {
		WriteError(w, http.StatusInternalServerError, "failed to get docs", err)
		return
	}
	bts, _ := pitaya.GetSerializer().Marshal(ret)
	Write(w, http.StatusOK, string(bts))
	return

}
