package request

func FakeRequestInfo(method, path string, opts ...CreateOption) *requestInfo {
	req := &requestInfo{
		Method: method,
		Path:   path,
	}
	for _, opt := range opts {
		opt(req)
	}
	return req
}
