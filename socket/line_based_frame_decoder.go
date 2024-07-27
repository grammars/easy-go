package socket

import (
	"bufio"
	"errors"
	"io"
	"log/slog"
)

type LineBasedFrameDecoder[VD any] struct {
}

func (decoder *LineBasedFrameDecoder[VD]) Decode(visitor *Visitor[VD], reader io.Reader) (CodecResult, error) {
	slog.Info("LineBasedFrameDecoder准备解码")
	result := CodecResult{}
	var scanner *bufio.Scanner
	if visitor.Ext == nil {
		scanner = bufio.NewScanner(reader)
		visitor.Ext = scanner
	} else {
		b, ok := visitor.Ext.(*bufio.Scanner)
		if ok {
			scanner = b
		} else {
			return result, errors.New("类型转化失败 visitor.Ext.(*bufio.Scanner)")
		}
	}
	result.BodyBytes = []byte(scanner.Text())
	result.FrameLength = len(result.BodyBytes)
	return result, nil
}
