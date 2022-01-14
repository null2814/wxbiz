package qywx

import (
	"context"
	"fmt"
	"strings"

	"github.com/valyala/fasthttp"
)

func NetGET(ctx context.Context, url string) ([]byte, error) {
	req := &fasthttp.Request{}
	req.SetRequestURI(url)
	req.Header.SetMethod("GET")
	resp := &fasthttp.Response{}
	client := &fasthttp.Client{}
	if err := client.Do(req, resp); err != nil {
		return nil, err
	}
	return resp.Body(), nil
}

func NetPOSTJson(ctx context.Context, url string, r []byte) ([]byte, error) {
	req := &fasthttp.Request{}
	req.SetRequestURI(url)
	req.SetBody(r)
	req.Header.SetContentType("application/json")
	req.Header.SetMethod("POST")
	resp := &fasthttp.Response{}
	client := &fasthttp.Client{}
	if err := client.Do(req, resp); err != nil {
		return nil, err
	}
	return resp.Body(), nil
}

func stringTrunc(limit int, o ...interface{}) string {
	trim := func(o string) string {
		o1 := strings.TrimPrefix(o, "[")
		o2 := strings.TrimSuffix(o1, "]")
		return o2
	}
	oStr := fmt.Sprintf("%v", o)
	s_oStr := strings.Split(oStr, "")
	if len(s_oStr) > limit {
		s_tStr := append(s_oStr[:17], "...")
		return trim(strings.Join(s_tStr, ""))
	}
	return trim(oStr)
}
