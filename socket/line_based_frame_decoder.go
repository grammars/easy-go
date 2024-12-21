package socket

import (
	"bufio"
	"errors"
	"io"
	"log/slog"
)

type LineBasedFrameDecoder struct {
}

func (decoder *LineBasedFrameDecoder) Decode(visitor VisitorSupport, reader io.Reader) (CodecResult, error) {
	slog.Info("LineBasedFrameDecoder准备解码")
	result := CodecResult{}
	var scanner *bufio.Scanner
	if visitor.GetExt() == nil {
		scanner = bufio.NewScanner(reader)
		visitor.SetExt(scanner)
	} else {
		b, ok := visitor.GetExt().(*bufio.Scanner)
		if ok {
			scanner = b
		} else {
			return result, errors.New("类型转化失败 visitor.GetExt().(*bufio.Scanner)")
		}
	}
	if scanner.Scan() {
		text := scanner.Text()
		slog.Info("收到文本", "text", text)
		result.BodyBytes = []byte(text)
		result.FrameLength = len(result.BodyBytes)
		return result, nil
	}
	return result, errors.New("no content for Scan")
}
