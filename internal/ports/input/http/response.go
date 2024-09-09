package http

const Ok = 200
const Created = 201
const BadRequest = 400

type GenericResponse struct {
	Status  int
	Body    string
	Headers map[string]string
}

func (r *GenericResponse) Header(key string) string {
	return r.Headers[key]
}

func (r *GenericResponse) GetBody() string {
	return r.Body
}

func (r *GenericResponse) GetStatusCode() int {
	return r.Status
}

func (r *GenericResponse) GetHeaders() map[string]string {
	return r.Headers
}

func NewResponse(status int, body string) Response {
	return &GenericResponse{
		Status: status,
		Body:   body,
	}
}

func NewResponseWithHeaders(status int, body string, headers map[string]string) Response {
	return &GenericResponse{
		Status:  status,
		Body:    body,
		Headers: headers,
	}
}
