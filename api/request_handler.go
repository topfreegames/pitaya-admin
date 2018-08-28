package api

// import (
// 	"encoding/json"
// 	"errors"
// 	"net/http"
// )

// // RequestMsg struct
// type RequestMsg struct {
// 	Route   string
// 	Payload []byte
// }

// // RequestHandler handler
// type RequestHandler struct {
// 	App *App
// }

// // NewRequestHandler creates a new request handler
// func NewRequestHandler(a *App) *RequestHandler {
// 	m := &RequestHandler{App: a}
// 	return m
// }

// func (s *RequestHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
// 	var request RequestMsg
// 	err := json.NewDecoder(r.Body).Decode(&request)

// 	if err != nil {
// 		WriteError(w, http.StatusBadRequest, "failed to decode request", err)
// 		return
// 	}

// 	if !pClient.Connected {
// 		WriteError(w, http.StatusBadRequest, "you should be connected", errors.New("client is not connected"))
// 		return
// 	}

// 	_, err = pClient.SendRequest(request.Route, request.Payload)

// 	if err != nil {
// 		WriteError(w, 500, "failed to send request", err)
// 		return
// 	}

// 	m := <-pClient.IncomingMsgChan
// 	WriteSuccessWithJSON(w, 200, m.Data, "Successfully sent a request")

// }

// // ConnectHandler handler
// type ConnectHandler struct {
// 	App *App
// }

// // NewConnectHandler creates a new connect handler
// func NewConnectHandler(a *App) *ConnectHandler {
// 	m := &ConnectHandler{App: a}
// 	return m
// }

// func (s *ConnectHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
// 	v := r.URL.Query()
// 	addr := v.Get("address")

// 	if addr == "" {
// 		WriteError(w, 500, "address should be provided", errors.New("No pitaya address provided"))
// 		return
// 	}

// 	if pClient.Connected {
// 		WriteError(w, http.StatusBadRequest, "already connected", errors.New("client is already connected"))
// 		return
// 	}

// 	err := pClient.ConnectTo(addr)
// 	if err != nil {
// 		WriteError(w, 500, "failed to connect to pitaya", err)
// 		return
// 	}

// 	if err != nil {
// 		WriteError(w, 500, "failed to stablish websocket", err)
// 	}
// 	disconnectedCh = make(chan bool, 1)
// 	//WriteSuccessWithJSON(w, 200, []byte("connected"), "client connected")

// }

// // NotifyHandler handler
// type NotifyHandler struct {
// 	App *App
// }

// // NewNotifyHandler creates a new notify handler
// func NewNotifyHandler(a *App) *NotifyHandler {
// 	m := &NotifyHandler{App: a}
// 	return m
// }

// func (s *NotifyHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
// 	var notifyMsg RequestMsg
// 	err := json.NewDecoder(r.Body).Decode(&notifyMsg)

// 	if err != nil {
// 		WriteError(w, http.StatusBadRequest, "failed to decode notify", err)
// 		return
// 	}

// 	if !pClient.Connected {
// 		WriteError(w, http.StatusBadRequest, "you should be connected", errors.New("fail"))
// 		return
// 	}

// 	err = pClient.SendNotify(notifyMsg.Route, notifyMsg.Payload)

// 	if err != nil {
// 		WriteError(w, 500, "failed to send request", err)
// 		return
// 	}

// 	WriteSuccessWithJSON(w, 200, []byte("server notifyied successfully"), "success on notify")
// }
