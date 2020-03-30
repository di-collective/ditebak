package rest

import (
	"fmt"
	"net/http"
	"path"
	"strconv"
	"strings"

	"github.com/di-collective/ditebak/backend/pkg/exception"
	"github.com/di-collective/ditebak/backend/pkg/queryables"
	"github.com/di-collective/ditebak/backend/pkg/service"
	"github.com/julienschmidt/httprouter"
)

// REST interface
type REST interface {
	WithRouter(router *httprouter.Router)
}

// rest abstraction
type rest struct {
	resource   string
	queryables queryables.Collection
	service    service.Service

	create  func() interface{}            // constructor of HTTP request payload - CREATE
	update  func() interface{}            // constructor of HTTP request payload - UPDATE
	convert func(interface{}) interface{} // convert HTTP request payload to service payload
}

// New REST API
func New(conf *Config) REST {
	api := &rest{
		resource:   conf.Resource,
		queryables: conf.Queryables,
		service:    conf.Service,
		create:     conf.CreatePayload,
		update:     conf.UpdatePayload,
		convert:    conf.Convert,
	}

	return api
}

// NewRouter initialize routes using julienschmidth httprouter
func (api *rest) NewRouter() *httprouter.Router {
	router := httprouter.New()
	api.WithRouter(router)

	return router
}

// WithRouter initialize routes using julienschmidth httprouter
func (api *rest) WithRouter(router *httprouter.Router) {
	res := strings.ToLower(api.resource)
	root := path.Join("/", res)
	withID := path.Join(root, ":id")

	router.GET(root, api.Find)
	router.POST(root, api.Create)
	router.GET(withID, api.Get)
	router.PATCH(withID, api.Update)
	router.DELETE(withID, api.Delete)
	router.Handle("REMOVE", withID, api.Remove)
}

// Find multiple
func (api *rest) Find(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	ctx := r.Context()
	res := NewAPIResponse(w, r)
	page, size := getPageAndSize(r)

	query := api.queryables.Read(r)
	total, result, err := api.service.Find(ctx, page, size, query)
	exc, throw := exception.IsException(err)
	if throw {
		res.Error(exc.Message(), err).Respond(exc.Code())
		return
	} else if err != nil {
		res.Error(fmt.Sprintf("Failed to find [%s]", api.resource), err).
			Respond(http.StatusInternalServerError)
		return
	}

	res.Paging(total, totalPage(total, int64(size))).
		Payload(result).
		Respond(http.StatusOK)
}

// Get one
func (api *rest) Get(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	id := p.ByName("id")
	ctx := r.Context()
	res := NewAPIResponse(w, r)

	result, err := api.service.Get(ctx, id)
	exc, throw := exception.IsException(err)
	if throw {
		res.Error(exc.Message(), err).Respond(exc.Code())
		return
	} else if err != nil {
		res.Error(fmt.Sprintf("Failed to get [%s] with id: %s", api.resource, id), err).
			Respond(http.StatusInternalServerError)
		return
	}

	res.Payload(result).Respond(http.StatusOK)
}

// Create one
func (api *rest) Create(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	ctx := r.Context()
	res := NewAPIResponse(w, r)

	// initialize payload struct and parse HTTP request to it
	httpPayload := api.create()
	err := ParseBody(r, httpPayload)
	if err != nil {
		res.Error("Failed to parse payload", err).Respond(http.StatusBadRequest)
		return
	}

	// convert httpPayload into servicePayload if necessary
	servicePayload := httpPayload
	if api.convert != nil {
		servicePayload = api.convert(httpPayload)
	}

	// call service with appropriate payload and get the result
	result, err := api.service.Create(ctx, servicePayload)
	exc, throw := exception.IsException(err)
	if throw {
		res.Error(exc.Message(), err).Respond(exc.Code())
		return
	} else if err != nil {
		res.Error(fmt.Sprintf("Failed to create [%s]", api.resource), err).
			Respond(http.StatusInternalServerError)
		return
	}

	res.Payload(result).Respond(http.StatusOK)
}

// Update one
func (api *rest) Update(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	id := p.ByName("id")
	ctx := r.Context()
	res := NewAPIResponse(w, r)

	// initialize payload struct and parse HTTP request to it
	httpPayload := api.update()
	err := ParseBody(r, httpPayload)
	if err != nil {
		res.Error("Failed to parse payload", err).Respond(http.StatusBadRequest)
		return
	}

	// convert httpPayload into servicePayload if necessary
	servicePayload := httpPayload
	if api.convert != nil {
		servicePayload = api.convert(httpPayload)
	}

	// call service with appropriate payload and get the result
	result, err := api.service.Update(ctx, id, servicePayload)
	exc, throw := exception.IsException(err)
	if throw {
		res.Error(exc.Message(), err).Respond(exc.Code())
		return
	} else if err != nil {
		res.Error(fmt.Sprintf("Failed to update [%s] with id: %s", api.resource, id), err).
			Respond(http.StatusInternalServerError)
		return
	}

	res.Payload(result).Respond(http.StatusOK)
}

// Delete one
func (api *rest) Delete(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	id := p.ByName("id")
	ctx := r.Context()
	res := NewAPIResponse(w, r)

	err := api.service.Delete(ctx, id)
	exc, throw := exception.IsException(err)
	if throw {
		res.Error(exc.Message(), err).Respond(exc.Code())
		return
	} else if err != nil {
		res.Error(fmt.Sprintf("Failed to delete [%s] with id: %s", api.resource, id), err).
			Respond(http.StatusInternalServerError)
		return
	}

	res.Respond(http.StatusResetContent)
}

// Remove one physically
func (api *rest) Remove(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	id := p.ByName("id")
	ctx := r.Context()
	res := NewAPIResponse(w, r)

	err := api.service.Remove(ctx, id)
	exc, throw := exception.IsException(err)
	if throw {
		res.Error(exc.Message(), err).Respond(exc.Code())
		return
	} else if err != nil {
		res.Error(fmt.Sprintf("Failed to remove [%s] with id: %s", api.resource, id), err).
			Respond(http.StatusInternalServerError)
		return
	}

	res.Respond(http.StatusResetContent)
}

func getPageAndSize(r *http.Request) (page, size int) {
	page, err := strconv.Atoi(r.FormValue("page"))
	if err != nil || page <= 0 {
		page = 1
	}

	size, err = strconv.Atoi(r.FormValue("size"))
	if err != nil || size <= 0 {
		size = 5
	}

	return page, size
}

func totalPage(total, size int64) int64 {
	if size == 0 {
		return 0
	}

	totalPage := total / size
	if total%size != 0 {
		totalPage++
	}
	return int64(totalPage)
}
