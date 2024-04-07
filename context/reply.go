package context

import (
	"encoding/json"
	"io"
	"net/http"
	"os"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
)

func (l *LuxContext) Reply(contentType string, body []byte) error {
	l.Response.Header().Set("Content-Type", contentType)
	_, err := l.Response.Write(body)
	return err
}

func (l *LuxContext) ReplyJSON(data interface{}) error {
	encoder := json.NewEncoder(l.Response)
	if err := encoder.Encode(data); err != nil {
		return err
	}
	return nil
}

func (l *LuxContext) ReplyPlainText(text string) error {
	return l.Reply("text/plain", []byte(text))
}

func (l *LuxContext) ReplyString(text string) error {
	return l.ReplyPlainText(text)
}

func (l *LuxContext) ReplyProtobuf(data protoreflect.ProtoMessage) error {
	value, err := proto.Marshal(data)
	if err != nil {
		return err
	}
	return l.Reply("application/protobuf", value)
}

func (l *LuxContext) ReplyBinary(data []byte) error {
	return l.Reply("application/octet-stream", data)
}

func (l *LuxContext) ReplyWebP(data []byte) error {
	return l.Reply("image/webp", data)
}

func (l *LuxContext) ReplyJPEG(data []byte) error {
	return l.Reply("image/jpeg", data)
}

func (l *LuxContext) ReplyPNG(data []byte) error {
	return l.Reply("image/png", data)
}

func (l *LuxContext) ReplyWebM(data []byte) error {
	return l.Reply("video/webm", data)
}

func (l *LuxContext) ReplyMP4(data []byte) error {
	return l.Reply("video/mp4", data)
}

func (l *LuxContext) ReplyCSV(data []byte) error {
	return l.Reply("text/csv", data)
}

func (l *LuxContext) ReplyHTML(data []byte) error {
	return l.Reply("text/html", data)
}

func (l *LuxContext) ReplyXML(data []byte) error {
	return l.Reply("text/xml", data)
}

func (l *LuxContext) ReplyExcel(data []byte) error {
	return l.Reply("application/vnd.ms-excel", data)
}

func (l *LuxContext) ReplyWord(data []byte) error {
	return l.Reply("application/msword", data)
}

func (l *LuxContext) ReplyPDF(data []byte) error {
	return l.Reply("application/pdf", data)
}

func (l *LuxContext) ReplyPowerpoint(data []byte) error {
	return l.Reply("application/vnd.ms-powerpoint", data)
}

func (l *LuxContext) ReplyZip(data []byte) error {
	return l.Reply("application/zip", data)
}

func (l *LuxContext) ReplyTar(data []byte) error {
	return l.Reply("application/tar", data)
}

func (l *LuxContext) ReplyGZIP(data []byte) error {
	return l.Reply("application/gzip", data)
}

func (l *LuxContext) Reply7Z(data []byte) error {
	return l.Reply("application/x-7z-compressed", data)
}

func (l *LuxContext) ReplyFile(path string) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()
	buf, err := io.ReadAll(f)
	if err != nil {
		return err
	}
	return l.ReplyBinary(buf)
}

func (l *LuxContext) ReplyAuto(data []byte) error {
	return l.Reply(http.DetectContentType(data), data)
}
