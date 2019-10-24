# testy

Simple client library for testingt HTTP and REST apis in Go (inspired by Resty)

Doesn't implement assertions.  Use your favorite test library.

```
  api := testy.New(handler)

  response := api.Get("/hello")
  assert.Equal(t, 200, response.StatusCode, "OK response is expected")
  assert.Equal(t, "hello, world!", response.String())
```
