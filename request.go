package scrago

import (
	"crypto/md5"
	"encoding/hex"
	"io"
	"net/http"
	"net/url"
	"strings"
)

type Request struct {
	id		string		// 唯一值
	Url 	string		// 请求地址
	Method 	string		// 请求方法
	Params 	interface{}	// 请求参数
	Header 	http.Header	// 请求头
	req		*http.Request
}

func NewRequest(method, rawurl string, params interface{}) *Request {
	if method == "" {
		method = http.MethodGet
	} else {
		method = strings.ToUpper(method)
	}

	return &Request{
		Url: rawurl,
		Method: method,
		Params: params,
		Header: http.Header{
			"User-Agent": []string{DefaultUserAgent},
		},
	}
}

func (r *Request) GetRequest() *http.Request {
	if r.req == nil {
		r.getRequest()
	}

	return r.req
}

func (r *Request) getRequest() {
	var params io.Reader = nil

	u, err := url.Parse(r.Url)
	if err != nil {
		return
	}

	if r.Params != nil {
		if r.Method == http.MethodGet {

			values := r.buildParams()

			val := u.Query()
			if len(val) > 0 {
				for k, v := range values {
					if _, is := val[k]; is {
						val[k] = v
					} else {
						for _, v2 := range v {
							val.Add(k, v2)
						}
					}
				}
				u.RawQuery = val.Encode()
			}

		} else {
			body := r.buildParams()
			params = strings.NewReader(body.Encode())
		}
	}

	r.req, _ = http.NewRequest(r.Method, u.String(), params)
	r.req.Header = r.Header
}

func (r *Request) buildParams() url.Values {
	var params url.Values

	if r.Params == nil || r.Params == "" {
		return params
	}

	switch t := r.Params.(type) {
	case url.Values:
		params = t
	case map[string]string:
		sParams := make([]string, 0, len(t))
		for k, v := range t {
			sParams = append(sParams, k+"="+v)
		}
		str := strings.Join(sParams, "&")
		params, _ = url.ParseQuery(str)

	case string:
		params, _ = url.ParseQuery(t)
	}

	return params
}

func (r *Request) Id() string {
	if r.id == "" {
		str := r.Method + r.Url + r.buildParams().Encode()
		ctx := md5.New()
		ctx.Write([]byte(str))
		r.id = hex.EncodeToString(ctx.Sum(nil))
	}
	return r.id
}

func (r *Request) SetId(s string) {
	r.id = s
}
