package urls

import (
	"fmt"
	"github.com/goccy/go-json"
	"github.com/isaquesb/url-shortener/internal/app"
	"github.com/isaquesb/url-shortener/internal/hasher"
	"github.com/isaquesb/url-shortener/internal/ports/input/http"
	"github.com/isaquesb/url-shortener/pkg/logger"
	"strings"
)

func RedirectShort(r http.Request) (http.Response, error) {
	short := r.PathValue("short").(string)
	if len(short) == 0 {
		return http.NewResponse(http.BadRequest, "Missing 'short' field"), nil
	}

	container := app.GetApp()
	repository := container.Api.Repository.Get()
	dispatcher := container.Api.Dispatcher.Get()

	url, err := repository.UrlFromShort(r.Ctx(), short)

	if err != nil {
		return nil, err
	}

	if len(url) == 0 {
		return http.NewResponse(http.NotFound, "Not Found URL for "+short), nil
	}

	err = dispatcher.Dispatch(r.Ctx(), NewVisitEvent([]byte(short)))
	if err != nil {
		logger.Error("Failed to dispatch event: %v", err)
	}

	return http.NewRedirectResponse(url), nil
}

func ShowStats(r http.Request) (http.Response, error) {
	short := r.PathValue("short").(string)
	if len(short) == 0 {
		return http.NewResponse(http.BadRequest, "Missing 'short' field"), nil
	}

	container := app.GetApp()
	repository := container.Api.Repository.Get()

	stats, err := repository.StatsFromShort(r.Ctx(), short)

	if err != nil {
		return nil, err
	}

	if nil == stats {
		return http.NewResponse(http.NotFound, "Not Found URL for "+short), nil
	}

	content, _ := json.Marshal(stats)

	return http.NewResponseWithHeaders(
		200,
		string(content),
		map[string]string{"Content-Type": "application/json"},
	), nil
}

func CreateShortUrl(r http.Request) (http.Response, error) {
	url := r.FormValue("url")
	if len(url) == 0 {
		return http.NewResponse(http.BadRequest, "Missing 'url' field"), nil
	}

	container := app.GetApp()
	dispatcher := container.Api.Dispatcher.Get()

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

func DeleteShortUrl(r http.Request) (http.Response, error) {
	short := r.PathValue("short").(string)
	if len(short) == 0 {
		return http.NewResponse(http.BadRequest, "Missing 'short' field"), nil
	}

	container := app.GetApp()
	repository := container.Api.Repository.Get()
	dispatcher := container.Api.Dispatcher.Get()

	url, err := repository.UrlFromShort(r.Ctx(), short)
	if err != nil {
		return nil, err
	}

	if len(url) == 0 {
		return http.NewResponse(http.NotFound, "Not Found URL for "+short), nil
	}

	err = dispatcher.Dispatch(r.Ctx(), NewDeleteEvent([]byte(short)))

	if err != nil {
		return nil, DispatchError{Err: err}
	}

	return http.NewResponse(http.NoContent, ""), nil
}

type DispatchError struct {
	Err error
}

func (d DispatchError) Error() string {
	return d.Err.Error()
}
