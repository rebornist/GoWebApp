package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"
)

func main() {

	r := &router{make(map[string]map[string]HandlerFunc)}

	r.HandleFunc("GET", "/", func(c *Context) {
		t := time.Now()
		fmt.Fprintln(c.ResponseWriter, "welcome!")
		log.Printf("[%s] %q %v\n", c.Request.Method, c.Request.URL.String(), time.Now().Sub(t))
	})

	r.HandleFunc("GET", "/about", func(c *Context) {
		t := time.Now()
		fmt.Fprintln(c.ResponseWriter, "about")
		log.Printf("[%s] %q %v\n", c.Request.Method, c.Request.URL.String(), time.Now().Sub(t))
	})

	r.HandleFunc("GET", "/users/:id", func(c *Context) {
		fmt.Fprintf(c.ResponseWriter, "retrieve user %v\n", c.Params["id"])
	})

	r.HandleFunc("GET", "/users/:user_id/addresses/:address_id", func(c *Context) {
		t := time.Now()
		fmt.Fprintf(c.ResponseWriter, "retrieve user %v's address %v\n", c.Params["user_id"], c.Params["address_id"])
		log.Printf("[%s] %q %v\n", c.Request.Method, c.Request.URL.String(), time.Now().Sub(t))
	})

	r.HandleFunc("POST", "/users", func(c *Context) {
		fmt.Fprintf(c.ResponseWriter, "create user\n")
	})

	r.HandleFunc("POST", "/users/:user_id/addresses", func(c *Context) {
		fmt.Fprintf(c.ResponseWriter, "create user %v's address\n", c.Params["user_id"])
	})

	// 8080 포트로 웹 서버 구동
	http.ListenAndServe(":8080", r)
}

/*
   http 메서드와 URL 패털 별로 핸들러를 등록하고,
   웹 요청이 들어왔을 때 적절한 핸들러로 연결해주는 라우터
*/

type router struct {
	// 키: http 메서드
	// 값: URL 패턴별로 실행할 HandlerFunc
	handlers map[string]map[string]HandlerFunc
}

// 라우터에 핸들러를 등록하기 위한 메서드
func (r *router) HandleFunc(method, pattern string, h HandlerFunc) {
	// http 메서드로 등록된 맵이 있는지 확인
	m, ok := r.handlers[method]
	if !ok {
		// 등록된 맵이 없으면 새 맵을 생성
		m = make(map[string]HandlerFunc)
		r.handlers[method] = m
	}

	// http 메서드로 등록된 맵에 URL 패턴과 핸들러 함수 등록
	m[pattern] = h
}

// 라우터에 등록된 동적 URL 패턴과 실제 URL 경로가 일치하는지 확인하는 함수
func match(pattern, path string) (bool, map[string]string) {
	// 패턴과 패스가 정확히 일치하면 바로 true를 반환
	if pattern == path {
		return true, nil
	}

	// 패턴과 패스를 "/" 단위로 구분
	patterns := strings.Split(pattern, "/")
	paths := strings.Split(path, "/")

	// 패턴과 패스를 "/"로 구분한 후 부분 문자열 집합의 개수가 다르면 false를 반환
	if len(patterns) != len(paths) {
		return false, nil
	}

	// 패턴에 일치하는 URL 매개변수를 담기 위한 params 맵 생성
	params := make(map[string]string)

	// "/"로 구분된 패턴/패스의 각 문자열을 하나씩 비교
	for i := 0; i < len(patterns); i++ {
		switch {
		case patterns[i] == paths[i]:
			// 패턴과 패스의 부분 문자열이 일치하면 바로 다음 루프 수행
		case len(patterns[i]) > 0 && patterns[i][0] == ':':
			// 패턴이 ":" 문자열로 시작하면 params에 URL params를 담은 후 다음 루프 수행
			params[patterns[i][1:]] = paths[i]
		default:
			// 일치하는 경우가 없으면 false를 반환
			return false, nil
		}
	}

	// true와 params를 반환
	return true, params
}

// 웹 요청의 http 메서드와 URL 경로를 분석해서 그에 맞는 핸들러를 찾아 동작시키는 ServHTTP 메서드
func (r *router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	// http 메서드에 맞는 모든 headers를 반복하여 요청 URL에 해당하는 handler를 찾음
	for pattern, handler := range r.handlers[req.Method] {
		if ok, params := match(pattern, req.URL.Path); ok {
			// Context 생성
			c := Context{
				Params:         make(map[string]interface{}),
				ResponseWriter: w,
				Request:        req,
			}
			for k, v := range params {
				c.Params[k] = v
			}
			// 요청 URL에 해당하는 handler 수행
			handler(&c)
			return
		}
	}
	// 요청 URL에 해당하는 handler를 찾지 못하면 NotFound 에러 처리
	http.NotFound(w, req)
	return
}

/*
   URL 패턴에 해당하는 매개변수를 핸들러 함수 내부로 전달하고
   웹 요청의 처리 상태를 저장하는 컨텍스트 작성
*/

// 컨텍스트 타입 정의
type Context struct {
	Params map[string]interface{}

	ResponseWriter http.ResponseWriter
	Request        *http.Request
}

// 핸들러 타입 정의
type HandlerFunc func(*Context)

// 에러 처리 미들웨어 작성
