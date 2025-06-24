package request

type InvokeRequest struct {
	*ExecRequest
	Namespace string
}
