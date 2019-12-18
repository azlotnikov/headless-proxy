package proxy_server

import (
	"fmt"
	"github.com/9seconds/httransform"
	"io/ioutil"
	"net"
	"net/url"
)

type proxy struct {
	localPort int
	srv       *httransform.Server
}

type proxyServer struct {
	entries      map[string]*proxy
	caCert       []byte
	caPrivateKey []byte
}

func NewLocalProxyServer() (ILocalProxyServer, error) {
	caCert, err := ioutil.ReadFile("server.pem")
	if err != nil {
		return nil, err
	}
	caPrivateKey, err := ioutil.ReadFile("server.key")
	if err != nil {
		return nil, err
	}
	server := proxyServer{
		caCert:       caCert,
		caPrivateKey: caPrivateKey,
		entries:      make(map[string]*proxy),
	}

	return &server, nil
}

func (s *proxyServer) List() map[string]int {
	result := make(map[string]int)
	for k, p := range s.entries {
		result[k] = p.localPort
	}
	return result
}

func (s *proxyServer) InitLocalProxy(proxyUrl string) (int, error) {
	existing, ok := s.entries[proxyUrl]
	if ok {
		return existing.localPort, nil
	}
	parsedUrl, err := url.Parse(proxyUrl)
	if err != nil {
		return 0, err
	}
	ln, err := net.Listen("tcp", ":0")
	if err != nil {
		return 0, err
	}
	executor, err := httransform.MakeProxyChainExecutor(parsedUrl)
	if err != nil {
		return 0, err
	}
	opts := httransform.ServerOpts{
		CertCA:   s.caCert,
		CertKey:  s.caPrivateKey,
		Executor: executor,
	}
	srv, err := httransform.NewServer(opts)
	if err != nil {
		return 0, err
	}

	go func() {
		_ = srv.Serve(ln)
	}()

	port := ln.Addr().(*net.TCPAddr).Port
	s.entries[proxyUrl] = &proxy{
		localPort: port,
		srv:       srv,
	}

	return port, nil
}

func (s *proxyServer) ShutdownLocalProxy(proxyUrl string) error {
	existing, ok := s.entries[proxyUrl]
	if !ok {
		return fmt.Errorf("proxy %s not found", proxyUrl)
	}
	delete(s.entries, proxyUrl)
	return existing.srv.Shutdown()
}
