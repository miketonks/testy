package testy

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
)

const (
	// MethodGet HTTP method
	MethodGet = "GET"

	// MethodPost HTTP method
	MethodPost = "POST"

	// MethodPut HTTP method
	MethodPut = "PUT"

	// MethodDelete HTTP method
	MethodDelete = "DELETE"

	// MethodPatch HTTP method
	MethodPatch = "PATCH"

	// MethodHead HTTP method
	MethodHead = "HEAD"

	// MethodOptions HTTP method
	MethodOptions = "OPTIONS"
)

// Client ...
type Client struct {
	handler    http.Handler
	QueryParam url.Values
	FormData   url.Values
	Header     http.Header
	Body       []byte
	Result     interface{}
	Error      interface{}
}

// Response ...
type Response struct {
	RawResponse *http.Response
	Body        []byte
	Status      string
	StatusCode  int
	Size        int64
}

// New ...
func New(h http.Handler) *Client {
	return &Client{
		handler:    h,
		QueryParam: url.Values{},
		FormData:   url.Values{},
		Header:     http.Header{},
	}
}

// Get ...
func (c *Client) Get(url string) *Response {
	return c.Execute("GET", url)
}

// Patch ...
func (c *Client) Patch(url string) *Response {
	return c.Execute("PATCH", url)
}

// Post ...
func (c *Client) Post(url string) *Response {
	return c.Execute("POST", url)
}

// Delete ...
func (c *Client) Delete(url string) *Response {
	return c.Execute("DELETE", url)
}

// Execute ...
func (c *Client) Execute(method, url string) *Response {

	if len(c.QueryParam) > 0 {
		url = fmt.Sprintf("%s?%s", url, c.QueryParam.Encode())
	}

	var reader io.Reader
	if c.Body != nil {
		reader = bytes.NewReader(c.Body)
	}
	request, _ := http.NewRequest(method, url, reader)
	request.Header = c.Header

	recorder := httptest.NewRecorder()
	c.handler.ServeHTTP(recorder, request)

	result := recorder.Result()
	response := Response{
		RawResponse: result,
		Status:      result.Status,
		StatusCode:  result.StatusCode,
	}

	var err error
	if response.Body, err = ioutil.ReadAll(result.Body); err != nil {
		panic(err)
	}

	response.Size = int64(len(response.Body))

	if c.Result != nil {
		err = json.Unmarshal(response.Body, c.Result)
		if err != nil {
			panic(err)
		}
	}
	return &response
}

// SetHeader method is to set a single header field and its value in the current request.
//
// For Example: To set `Content-Type` and `Accept` as `application/json`.
// 		client.R().
//			SetHeader("Content-Type", "application/json").
//			SetHeader("Accept", "application/json")
//
// Also you can override header value, which was set at client instance level.
func (c *Client) SetHeader(header, value string) *Client {
	c.Header.Set(header, value)
	return c
}

// SetHeaders method sets multiple headers field and its values at one go in the current request.
//
// For Example: To set `Content-Type` and `Accept` as `application/json`
//
// 		client.R().
//			SetHeaders(map[string]string{
//				"Content-Type": "application/json",
//				"Accept": "application/json",
//			})
// Also you can override header value, which was set at client instance level.
func (c *Client) SetHeaders(headers map[string]string) *Client {
	for h, v := range headers {
		c.SetHeader(h, v)
	}
	return c
}

// SetQueryParam method sets single parameter and its value in the current request.
// It will be formed as query string for the request.
//
// For Example: `search=kitchen%20papers&size=large` in the URL after `?` mark.
// 		client.R().
//			SetQueryParam("search", "kitchen papers").
//			SetQueryParam("size", "large")
// Also you can override query params value, which was set at client instance level.
func (c *Client) SetQueryParam(param, value string) *Client {
	c.QueryParam.Set(param, value)
	return c
}

// SetQueryParams method sets multiple parameters and its values at one go in the current request.
// It will be formed as query string for the request.
//
// For Example: `search=kitchen%20papers&size=large` in the URL after `?` mark.
// 		client.R().
//			SetQueryParams(map[string]string{
//				"search": "kitchen papers",
//				"size": "large",
//			})
// Also you can override query params value, which was set at client instance level.
func (c *Client) SetQueryParams(params map[string]string) *Client {
	for p, v := range params {
		c.SetQueryParam(p, v)
	}
	return c
}

// SetQueryParamsFromValues method appends multiple parameters with multi-value
// (`url.Values`) at one go in the current request. It will be formed as
// query string for the request.
//
// For Example: `status=pending&status=approved&status=open` in the URL after `?` mark.
// 		client.R().
//			SetQueryParamsFromValues(url.Values{
//				"status": []string{"pending", "approved", "open"},
//			})
// Also you can override query params value, which was set at client instance level.
func (c *Client) SetQueryParamsFromValues(params url.Values) *Client {
	for p, v := range params {
		for _, pv := range v {
			c.QueryParam.Add(p, pv)
		}
	}
	return c
}

// SetQueryString method provides ability to use string as an input to set URL query string for the request.
//
// Using String as an input
// 		client.R().
//			SetQueryString("productId=232&template=fresh-sample&cat=resty&source=google&kw=buy a lot more")
func (c *Client) SetQueryString(query string) *Client {
	params, err := url.ParseQuery(strings.TrimSpace(query))
	if err == nil {
		for p, v := range params {
			for _, pv := range v {
				c.QueryParam.Add(p, pv)
			}
		}
	} else {
		//c.client.log.Errorf("%v", err)
	}
	return c
}

// SetBody ...
func (c *Client) SetBody(body []byte) *Client {
	c.Body = body
	return c
}

// SetResult ...
func (c *Client) SetResult(result interface{}) *Client {
	c.Result = result
	return c
}

//
// // SetBody ...
// func (c *Client) SetBody(body interface{}) *Client {
//
// 	var bodyBytes []byte
// 	contentType := r.Header.Get("Content-Type")
// 	kind := kindOf(body)
// 	r.bodyBuf = nil
//
// 	if reader, ok := r.Body.(io.Reader); ok {
// 		if c.setContentLength || r.setContentLength { // keep backward compability
// 			r.bodyBuf = acquireBuffer()
// 			_, err = r.bodyBuf.ReadFrom(reader)
// 			r.Body = nil
// 		} else {
// 			// Otherwise buffer less processing for `io.Reader`, sounds good.
// 			return
// 		}
// 	} else if b, ok := r.Body.([]byte); ok {
// 		bodyBytes = b
// 	} else if s, ok := r.Body.(string); ok {
// 		bodyBytes = []byte(s)
// 	} else if IsJSONType(contentType) &&
// 		(kind == reflect.Struct || kind == reflect.Map || kind == reflect.Slice) {
// 		bodyBytes, err = jsonMarshal(c, r, r.Body)
// 	} else if IsXMLType(contentType) && (kind == reflect.Struct) {
// 		bodyBytes, err = xml.Marshal(r.Body)
// 	}
//
// 	if bodyBytes == nil && r.bodyBuf == nil {
// 		err = errors.New("unsupported 'Body' type/value")
// 	}
//
// 	// if any errors during body bytes handling, return it
// 	if err != nil {
// 		return
// 	}
//
// 	// []byte into Buffer
// 	if bodyBytes != nil && r.bodyBuf == nil {
// 		r.bodyBuf = acquireBuffer()
// 		_, _ = r.bodyBuf.Write(bodyBytes)
// 	}
//
// 	return
// }

func (r *Response) String() string {
	return string(r.Body)
}
