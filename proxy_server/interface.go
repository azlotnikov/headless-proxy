package proxy_server

type ILocalProxyServer interface {
	InitLocalProxy(proxyUrl string) (int, error)
	ShutdownLocalProxy(proxyUrl string) error
	List() map[string]int
}
