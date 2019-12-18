package api_server

type IApiServer interface {
	Serve(port int) error
}
