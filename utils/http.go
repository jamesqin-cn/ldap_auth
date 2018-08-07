package utils

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"reflect"

	"github.com/julienschmidt/httprouter"
)

const (
	HTTP_CONTEXT_KEY_REQUEST = 0x01
)

var (
	typeOfContext = reflect.TypeOf((*context.Context)(nil)).Elem()
)

type HttpServer struct {
	host   string
	router *httprouter.Router
}

func NewHttpServer(host string) *HttpServer {
	router := httprouter.New()

	return &HttpServer{
		host,
		router,
	}
}

// isRequestOrContext returns true if type t is either *http.Request or Context
func isContext(t reflect.Type) bool {
	if t.Implements(typeOfContext) {
		return true
	}
	return false
}

func HttpHandlerWrapper(method interface{}) func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	mtype := reflect.TypeOf(method)
	mval := reflect.ValueOf(method)
	numIn, numOut := mtype.NumIn(), mtype.NumOut()
	if numIn != 2 || numOut != 1 {
		log.Println("http handle function wrap by Wrap() must be the type: func (ctx context.Context, req interface{}) (reply interface{})")
		return nil
	}

	ctxType := mtype.In(0)
	reqType := mtype.In(1)

	if !isContext(ctxType) {
		log.Println("http handle function wrap by Wrap() must be the type: func (ctx context.Context, req interface{}) (reply interface{})")
		return nil
	}

	getReqValue := func(r *http.Request, reqType reflect.Type) (reqVal reflect.Value, err error) {
		reqVal = reflect.New(reqType.Elem())
		err = json.NewDecoder(r.Body).Decode(reqVal.Interface())
		if err != nil {
			return
		}
		return reqVal, err
	}

	writeReply := func(w http.ResponseWriter, reply reflect.Value) {
		html, err := json.Marshal(reply.Interface())
		if err != nil {
			log.Println(err)
			http.Error(w, "Encode response body failed", 500)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, "%s", html)
	}

	return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		req, err := getReqValue(r, reqType)
		if err != nil {
			log.Println(err)
			http.Error(w, "Decode request body failed", 500)
			return
		}
		ctx := context.WithValue(context.Background(), HTTP_CONTEXT_KEY_REQUEST, r)
		args := []reflect.Value{
			reflect.ValueOf(ctx),
			req,
		}
		replys := mval.Call(args)
		reply := replys[0]
		writeReply(w, reply)
	}
}

func (s *HttpServer) GET(path string, handle httprouter.Handle) {
	s.router.GET(path, handle)
}

func (s *HttpServer) POST(path string, handle httprouter.Handle) {
	s.router.POST(path, handle)
}

func (s *HttpServer) ANY(path string, handle httprouter.Handle) {
	s.router.GET(path, handle)
	s.router.POST(path, handle)
}

func (s *HttpServer) RunHttpServer() {
	log.Fatal(http.ListenAndServe(s.host, s.router))
}

func GetQuery(r *http.Request, key string, defaultVal string) string {
	values, ok := r.URL.Query()[key]
	if ok && len(values) > 0 && len(values[0]) > 0 {
		return values[0]
	}
	return defaultVal
}

func GetPost(r *http.Request, key string, defaultVal string) string {
	r.ParseMultipartForm(32 << 20)
	if r.MultipartForm != nil {
		values := r.MultipartForm.Value[key]
		if len(values) > 0 && len(values[0]) > 0 {
			return values[0]
		}
	}
	return defaultVal
}

func GetCookie(r *http.Request, key string, defaultVal string) string {
	cookie, err := r.Cookie(key)
	if err == nil && len(cookie.Value) > 0 {
		return cookie.Value
	}
	return defaultVal
}
