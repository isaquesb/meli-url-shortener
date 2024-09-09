package api

import (
	"fmt"
	"github.com/isaquesb/meli-url-shortener/config"
	"github.com/isaquesb/meli-url-shortener/internal/hasher"
	"github.com/isaquesb/meli-url-shortener/internal/ports/input/http"
	"github.com/isaquesb/meli-url-shortener/internal/ports/output/events"
	"strings"
)

func CreateShortUrl(r http.Request) (http.Response, error) {
	url := r.FormValue("url")
	if len(url) == 0 {
		return http.NewResponse(http.BadRequest, "Missing 'url' field"), nil
	}

	container := config.GetApp()
	dispatcher := container.Api.GetDispatcher()

	short := hasher.GetUrlHash(url)
	completeUrl := fmt.Sprintf("%s/%s", container.Host, short)

	msg := &events.Message{
		ContentType: events.TypePlain,
		Key:         short,
		Body:        []byte(fmt.Sprintf("%s%s", short, url)),
	}
	err := dispatcher.Dispatch(r.Ctx(), container.Events["urls.created"], *msg)

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
