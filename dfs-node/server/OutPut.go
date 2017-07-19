package server

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"pilicat-core/logs"

	"github.com/kataras/iris/context"
)

func OutputJson(ctx context.Context, data interface{}) {
	w := ctx.ResponseWriter()
	w.WriteHeader(http.StatusOK)
	w.Header().Add("Content-Type", "application/json;charset=utf-8")

	content, err := json.Marshal(data)
	if err == nil {
		w.Header().Add("Content-Length", fmt.Sprint(len(content)))

		io.Copy(w, bytes.NewReader(content))
	} else {
		logs.Error("OutputJson error: ", err)
	}

	return
}

func OutputText(ctx context.Context, txt string) {
	w := ctx.ResponseWriter()
	w.WriteHeader(http.StatusOK)
	w.Header().Add("Content-Type", "text/plain;charset=utf-8")

	var outBuf bytes.Buffer
	outBuf.WriteString(txt)

	io.Copy(w, &outBuf)
	return
}
