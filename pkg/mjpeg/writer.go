package mjpeg

import (
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"
)

func NewWriter(w io.Writer) io.Writer {
	h := w.(http.ResponseWriter).Header()
	h.Set("Content-Type", "multipart/x-mixed-replace; boundary=frame")
	stime := time.Now()
	h.Set("X-StartTime", fmt.Sprint(stime.UnixMilli()))
	return &writer{wr: w, buf: []byte(header), stime: stime}
}

const header = "--frame\r\nContent-Type: image/jpeg\r\nContent-Length: "

type writer struct {
	wr  io.Writer
	buf []byte
	stime  time.Time
}

func (w *writer) Write(p []byte) (n int, err error) {
	w.buf = w.buf[:len(header)]
	w.buf = append(w.buf, strconv.Itoa(len(p))...)
	w.buf = append(w.buf, ("\r\nX-Timestamp: " + fmt.Sprint(time.Since(w.stime).Milliseconds()))...)
	w.buf = append(w.buf, ("\r\nX-CurrentTime: " + fmt.Sprint(time.Now().UnixMilli()))...)
	w.buf = append(w.buf, "\r\n\r\n"...)
	w.buf = append(w.buf, p...)
	w.buf = append(w.buf, "\r\n"...)

	// Chrome bug: mjpeg image always shows the second to last image
	// https://bugs.chromium.org/p/chromium/issues/detail?id=527446
	if n, err = w.wr.Write(w.buf); err == nil {
		w.wr.(http.Flusher).Flush()
	}

	return
}