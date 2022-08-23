package clienterror

type ApiNotFoundError struct {
	ApiId string
}

func (e *ApiNotFoundError) Error() string {
	return "No API found for ApiId " + e.ApiId
}
