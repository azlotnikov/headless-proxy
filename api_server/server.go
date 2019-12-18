package api_server

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"proxy/proxy_server"
)

const apiKeyHeader = "x-api-key"

type apiServer struct {
	proxyServer proxy_server.ILocalProxyServer
	apiKey      string
}

type proxy struct {
	Url string `json:"url" required:"true"`
}

func NewApiServer(proxyServer proxy_server.ILocalProxyServer, apiKey string) IApiServer {
	return &apiServer{
		proxyServer: proxyServer,
		apiKey:      apiKey,
	}
}

func (s *apiServer) Serve(port int) error {
	r := mux.NewRouter()
	r.HandleFunc("/init", s.initProxy).Methods("POST")
	r.HandleFunc("/list", s.listProxy).Methods("GET")
	r.HandleFunc("/shutdown", s.shutdownProxy).Methods("POST")
	return http.ListenAndServe(fmt.Sprintf(":%d", port), r)
}

func (s *apiServer) listProxy(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if err := s.assertApiKey(r); err != nil {
		writeError(w, err)
		return
	}
	list := s.proxyServer.List()
	if err := json.NewEncoder(w).Encode(list); err != nil {
		panic(err)
	}
}

func (s *apiServer) shutdownProxy(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if err := s.assertApiKey(r); err != nil {
		writeError(w, err)
		return
	}
	var p proxy
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		writeError(w, err)
		return
	}
	if p.Url == "" {
		writeError(w, fmt.Errorf("empty proxy url"))
		return
	}
	err := s.proxyServer.ShutdownLocalProxy(p.Url)
	if err != nil {
		writeError(w, err)
		return
	}
	if err := json.NewEncoder(w).Encode(&struct {
		Result string `json:"result"`
	}{"ok"}); err != nil {
		panic(err)
	}
}

func (s *apiServer) initProxy(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if err := s.assertApiKey(r); err != nil {
		writeError(w, err)
		return
	}
	var p proxy
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		writeError(w, err)
		return
	}
	if p.Url == "" {
		writeError(w, fmt.Errorf("empty proxy url"))
		return
	}
	port, err := s.proxyServer.InitLocalProxy(p.Url)
	if err != nil {
		writeError(w, err)
		return
	}
	if err := json.NewEncoder(w).Encode(&struct {
		Port int `json:"port"`
	}{port}); err != nil {
		panic(err)
	}
}

func (s *apiServer) assertApiKey(r *http.Request) error {
	if s.apiKey == "" {
		return nil
	}
	if s.apiKey != r.Header.Get(apiKeyHeader) {
		return fmt.Errorf("invalid api key")
	}
	return nil
}

func writeError(w http.ResponseWriter, err error) {
	if err := json.NewEncoder(w).Encode(&struct {
		Error string `json:"error"`
	}{err.Error()}); err != nil {
		panic(err)
	}
}
