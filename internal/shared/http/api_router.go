package http

import (
	"net/http"
	"reflect"
	"strings"

	"github.com/azahir21/go-backend-boilerplate/pkg/apperr"
	"github.com/gin-gonic/gin"
)

type EndpointSpec struct {
	Method      string
	Path        string
	Middlewares []gin.HandlerFunc
	// Handler supports:
	//   func(*gin.Context)
	//   func(*gin.Context) error
	//   func(*gin.Context, T) error
	//   func(*gin.Context, *T) error
	Handler interface{}
	// BindFrom: "", "json", "query" (default: auto by HTTP method)
	BindFrom string
}

type APIRouterGroup struct {
	group *gin.RouterGroup
}

func NewAPIRouterGroup(group *gin.RouterGroup) *APIRouterGroup {
	return &APIRouterGroup{group: group}
}

func (g *APIRouterGroup) Register(specs ...EndpointSpec) {
	for _, s := range specs {
		final := g.wrap(s)
		handlers := append(s.Middlewares, final)
		g.group.Handle(s.Method, s.Path, handlers...)
	}
}

var ginContextType = reflect.TypeOf((*gin.Context)(nil))

func (g *APIRouterGroup) wrap(s EndpointSpec) gin.HandlerFunc {
	return func(c *gin.Context) {
		hv := reflect.ValueOf(s.Handler)
		ht := hv.Type()
		if ht.Kind() != reflect.Func || ht.NumIn() < 1 || ht.In(0) != ginContextType {
			apperr.RespondInternalServer(c, "invalid handler signature")
			return
		}

		args := []reflect.Value{reflect.ValueOf(c)}
		// If a second parameter exists, bind it as request model.
		if ht.NumIn() >= 2 {
			reqType := ht.In(1)
			var reqPtr reflect.Value
			if reqType.Kind() == reflect.Ptr {
				reqPtr = reflect.New(reqType.Elem())
			} else {
				reqPtr = reflect.New(reqType)
			}

			// build bind list: explicit BindFrom (comma-separated) or default param,query,json
			var bindList []string
			if strings.TrimSpace(s.BindFrom) != "" {
				for _, p := range strings.Split(s.BindFrom, ",") {
					if v := strings.TrimSpace(p); v != "" {
						bindList = append(bindList, v)
					}
				}
			} else {
				method := strings.ToUpper(s.Method)
				if method == http.MethodGet || method == http.MethodDelete {
					bindList = []string{"param", "query", "json"}
				} else {
					bindList = []string{"json", "query", "param"}
				}
			}

			// perform bindings in order; json binder only runs if there is a body
			for _, b := range bindList {
				switch b {
				case "param", "uri":
					if len(c.Params) == 0 {
						continue
					}
					if err := c.ShouldBindUri(reqPtr.Interface()); err != nil {
						apperr.Respond(c, apperr.BadRequest("Invalid URI parameters").WithCause(err))
						return
					}
				case "query":
					// Skip query binding when there are no query values.
					if len(c.Request.URL.Query()) == 0 {
						continue
					}
					if err := c.ShouldBindQuery(reqPtr.Interface()); err != nil {
						apperr.Respond(c, apperr.BadRequest("Invalid query parameters").WithCause(err))
						return
					}
				case "json", "body":
					// only try JSON bind when there's a body / content-type indicates JSON
					if c.Request.ContentLength == 0 && !strings.Contains(strings.ToLower(c.GetHeader("Content-Type")), "json") {
						continue
					}
					if err := c.ShouldBindJSON(reqPtr.Interface()); err != nil {
						apperr.Respond(c, apperr.BadRequest("Invalid JSON body").WithCause(err))
						return
					}
				default:
					// unknown binder -> skip
				}
			}

			// Store for compatibility with GetRequest[T]
			c.Set("req", reqPtr.Interface())

			// Pass value or pointer depending on handler signature.
			if reqType.Kind() == reflect.Ptr {
				args = append(args, reqPtr)
			} else {
				args = append(args, reqPtr.Elem())
			}
		}

		outs := hv.Call(args)
		// Optional error return handling
		if ht.NumOut() == 1 && !outs[0].IsNil() {
			if err, ok := outs[0].Interface().(error); ok && err != nil {
				// If your handler already wrote a response, this will be ignored by Gin.
				apperr.RespondInternalServer(c, err.Error())
				return
			}
		}
	}
}
