package core

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func newTestReqCtx(method, path string, body io.Reader) (*RequestContext, *httptest.ResponseRecorder) {
	req := httptest.NewRequest(method, path, body)
	rr := httptest.NewRecorder()
	ctx := NewRequestContext(&Context{store: make(map[string]any)}, rr, req)
	return ctx, rr
}

func TestJSON(t *testing.T) {
	ctx, rr := newTestReqCtx("GET", "/", nil)

	err := ctx.JSON(200, map[string]string{"msg": "ok"})
	if err != nil {
		t.Fatalf("JSON returned error: %v", err)
	}

	if rr.Code != 200 {
		t.Fatalf("expected status 200, got %d", rr.Code)
	}

	var data map[string]string
	if err := json.Unmarshal(rr.Body.Bytes(), &data); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}

	if data["msg"] != "ok" {
		t.Fatalf("expected msg=ok, got %v", data)
	}
}

func TestText(t *testing.T) {
	ctx, rr := newTestReqCtx("GET", "/", nil)

	err := ctx.Text(201, "hello")
	if err != nil {
		t.Fatalf("Text returned error: %v", err)
	}

	if rr.Code != 201 {
		t.Fatalf("expected 201, got %d", rr.Code)
	}
	if rr.Body.String() != "hello" {
		t.Fatalf("expected 'hello', got %s", rr.Body.String())
	}
}

func TestBindJSON(t *testing.T) {
	body := bytes.NewBufferString(`{"name":"Lilium"}`)
	ctx, _ := newTestReqCtx("POST", "/", body)

	var data map[string]string
	if err := ctx.BindJSON(&data); err != nil {
		t.Fatalf("BindJSON error: %v", err)
	}

	if data["name"] != "Lilium" {
		t.Fatalf("expected Lilium, got %v", data["name"])
	}
}

func TestQuery(t *testing.T) {
	ctx, _ := newTestReqCtx("GET", "/test?x=10&x=20&y=1", nil)

	if ctx.Query("y") != "1" {
		t.Fatalf("expected y=1")
	}

	vals := ctx.QueryAll("x")
	if len(vals) != 2 || vals[0] != "10" || vals[1] != "20" {
		t.Fatalf("unexpected QueryAll: %v", vals)
	}
}

func TestParams(t *testing.T) {
	ctx, _ := newTestReqCtx("GET", "/", nil)
	ctx.Params["id"] = "123"

	if ctx.Param("id") != "123" {
		t.Fatalf("expected param=123")
	}
}

func TestStore(t *testing.T) {
	ctx, _ := newTestReqCtx("GET", "/", nil)

	ctx.Set("x", 55)

	val, ok := ctx.Get("x")
	if !ok || val.(int) != 55 {
		t.Fatalf("expected store value=55")
	}
}

func TestRedirect(t *testing.T) {
	ctx, rr := newTestReqCtx("GET", "/", nil)

	err := ctx.Redirect(302, "/next")
	if err != nil {
		t.Fatalf("Redirect error: %v", err)
	}

	if rr.Code != 302 {
		t.Fatalf("expected 302, got %d", rr.Code)
	}

	if loc := rr.Header().Get("Location"); loc != "/next" {
		t.Fatalf("expected Location=/next, got %s", loc)
	}
}

func TestBodyBytes(t *testing.T) {
	body := bytes.NewBufferString("hello world")
	ctx, _ := newTestReqCtx("POST", "/", body)

	b, err := ctx.BodyBytes()
	if err != nil {
		t.Fatalf("BodyBytes error: %v", err)
	}

	if string(b) != "hello world" {
		t.Fatalf("BodyBytes mismatch: %s", string(b))
	}
}

func TestClientIP(t *testing.T) {
	ctx, _ := newTestReqCtx("GET", "/", nil)

	ctx.Req.Header.Set("X-Real-IP", "1.1.1.1")
	if ctx.ClientIP() != "1.1.1.1" {
		t.Fatalf("expected X-Real-IP")
	}

	ctx.Req.Header.Del("X-Real-IP")
	ctx.Req.Header.Set("X-Forwarded-For", "2.2.2.2")
	if ctx.ClientIP() != "2.2.2.2" {
		t.Fatalf("expected X-Forwarded-For")
	}
}

func TestHTML(t *testing.T) {
	ctx, rr := newTestReqCtx("GET", "/", nil)

	err := ctx.HTML(200, "<h1>Hi</h1>")
	if err != nil {
		t.Fatalf("HTML error: %v", err)
	}

	if rr.Body.String() != "<h1>Hi</h1>" {
		t.Fatalf("unexpected HTML body")
	}
}

func TestHeaders(t *testing.T) {
	ctx, rr := newTestReqCtx("GET", "/", nil)

	h := http.Header{}
	h.Set("X-Test", "ok")

	ctx.Headers(h)

	if rr.Header().Get("X-Test") != "ok" {
		t.Fatalf("Headers not applied")
	}
}

func TestStream(t *testing.T) {
	ctx, rr := newTestReqCtx("GET", "/", nil)

	err := ctx.Stream(func(w io.Writer) error {
		_, e := w.Write([]byte("data"))
		return e
	})
	if err != nil {
		t.Fatalf("Stream error: %v", err)
	}

	if rr.Body.String() != "data" {
		t.Fatalf("Stream output mismatch")
	}
}

func TestDeadlineDoneErr(t *testing.T) {
	ctx, _ := newTestReqCtx("GET", "/", nil)

	// set a short deadline
	dCtx, cancel := context.WithTimeout(ctx.Req.Context(), 5*time.Millisecond)
	defer cancel()
	ctx.Req = ctx.Req.WithContext(dCtx)

	_, ok := ctx.Deadline()
	if !ok {
		t.Fatalf("expected deadline ok=true")
	}

	<-ctx.Done() // wait for timeout

	if ctx.Err() == nil {
		t.Fatalf("expected ctx.Err() after deadline")
	}
}

func TestFormPostForm(t *testing.T) {
	body := bytes.NewBufferString("x=10&y=20")
	req := httptest.NewRequest("POST", "/", body)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	rr := httptest.NewRecorder()
	ctx := NewRequestContext(&Context{store: make(map[string]any)}, rr, req)

	if ctx.Form("x") != "10" {
		t.Fatalf("expected Form x=10")
	}
	if ctx.PostForm("y") != "20" {
		t.Fatalf("expected PostForm y=20")
	}

	// ensure only parsed once
	if !ctx.formParsed {
		t.Fatalf("formParsed should be true")
	}
}
