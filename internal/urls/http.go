package urls

import (
	"fmt"
	"github.com/isaquesb/meli-url-shortener/internal/app"
	"github.com/isaquesb/meli-url-shortener/internal/hasher"
	"github.com/isaquesb/meli-url-shortener/internal/ports/input/http"
	"strings"
)

func CreateShortUrl(r http.Request) (http.Response, error) {
	url := r.FormValue("url")
	if len(url) == 0 {
		return http.NewResponse(http.BadRequest, "Missing 'url' field"), nil
	}

	container := app.GetApp()
	dispatcher := container.Api.GetDispatcher()

	short := hasher.GetUrlHash(url)
	completeUrl := fmt.Sprintf("%s/%s", container.Host, short)

	msg := NewCreateEvent(short, url)
	err := dispatcher.Dispatch(r.Ctx(), msg)

	if err != nil {
		return nil, DispatchError{Err: err}
	}

	contentType := "text/plain"
	content := completeUrl
	acceptHeader := r.Header("Accept")

	if strings.Contains(acceptHeader, "application/json") {
		contentType = "application/json"
		content = fmt.Sprintf(`{"short": "%s"}`, completeUrl)
	}

	return http.NewResponseWithHeaders(
		http.Created,
		content,
		map[string]string{"Content-Type": contentType},
	), nil
}

type DispatchError struct {
	Err error
}

func (d DispatchError) Error() string {
	return d.Err.Error()
}
