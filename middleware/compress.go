package middleware

import (
	"bytes"
	"compress/gzip"
	"net/http"
	"strings"

	"github.com/andybalholm/brotli"
	"github.com/golang/snappy"
	"github.com/snowmerak/lux/v3/context"
)

var CompressResponse = compressResponse{}

type compressResponse struct{}

type SnappyResponseMiddleware Response

func (cr compressResponse) Snappy() SnappyResponseMiddleware {
	return func(l *context.LuxContext) (*context.LuxContext, error) {
		acceptEncodings := strings.Split(l.Request.Header.Get("Accept-Encoding"), ", ")
		if len(acceptEncodings) > 0 && acceptEncodings[0] != "snappy" || len(acceptEncodings) == 0 {
			return l, nil
		}
		buf := bytes.NewBuffer(nil)
		writer := snappy.NewBufferedWriter(buf)
		_, err := writer.Write(l.Response.Body)
		if err != nil {
			return l, err
		}
		writer.Flush()
		writer.Close()
		l.Response.Body = buf.Bytes()
		l.Response.Headers.Add("Content-Encoding", "snappy")
		l.Request.Header.Set("Accept-Encoding", strings.Join(acceptEncodings[1:], ", "))
		return l, nil
	}
}

type GzipResponseMiddleware Response

func (cr compressResponse) Gzip() GzipResponseMiddleware {
	return func(l *context.LuxContext) (*context.LuxContext, error) {
		acceptEncodings := strings.Split(l.Request.Header.Get("Accept-Encoding"), ", ")
		if len(acceptEncodings) > 0 && acceptEncodings[0] != "gzip" || len(acceptEncodings) == 0 {
			return l, nil
		}
		buf := bytes.NewBuffer(nil)
		writer := gzip.NewWriter(buf)
		_, err := writer.Write(l.Response.Body)
		if err != nil {
			l.Response.StatusCode = http.StatusInternalServerError
			return l, err
		}
		writer.Flush()
		writer.Close()
		l.Response.Body = buf.Bytes()
		l.Response.Headers.Add("Content-Encoding", "gzip")
		if len(acceptEncodings) >= 2 && acceptEncodings[1] == "defalte" {
			acceptEncodings = acceptEncodings[1:]
		}
		l.Request.Header.Set("Accept-Encoding", strings.Join(acceptEncodings[1:], ", "))
		return l, nil
	}
}

type BrotliResponseMiddleware Response

func (cr compressResponse) Brotli() BrotliResponseMiddleware {
	return func(l *context.LuxContext) (*context.LuxContext, error) {
		acceptEncodings := strings.Split(l.Request.Header.Get("Accept-Encoding"), ", ")
		if len(acceptEncodings) > 0 && acceptEncodings[0] != "br" || len(acceptEncodings) == 0 {
			return l, nil
		}
		buf := bytes.NewBuffer(nil)
		writer := brotli.NewWriter(buf)
		_, err := writer.Write(l.Response.Body)
		if err != nil {
			l.Response.StatusCode = http.StatusInternalServerError
			return l, err
		}
		writer.Flush()
		writer.Close()
		l.Response.Body = buf.Bytes()
		l.Response.Headers.Add("Content-Encoding", "br")
		l.Request.Header.Set("Accept-Encoding", strings.Join(acceptEncodings[1:], ", "))
		return l, nil
	}
}
