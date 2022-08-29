package clienterror

type ApiUnauthorizedError struct {
	ApiId string
}

func (e ApiUnauthorizedError) Error() string {
	msg := "Unauthorized error"
	if e.ApiId != "" {
		msg += " for API " + e.ApiId
	}
	return msg
}
