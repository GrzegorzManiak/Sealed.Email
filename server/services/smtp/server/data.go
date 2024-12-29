package server

import (
	"bufio"
	"bytes"
	"github.com/GrzegorzManiak/NoiseBackend/services/smtp/headers"
	"go.uber.org/zap"
	"io"
	"strings"
)

func (s *Session) Data(r io.Reader) error {
	zap.L().Debug("Data received")

	var buffer bytes.Buffer
	bufReader := bufio.NewReader(r)

	for {
		line, err := bufReader.ReadBytes('\n')
		if err != nil {
			if err == io.EOF {
				if buffer.Len() > 0 {
					eofBytes := buffer.Bytes()
					s.RawData = append(s.RawData, eofBytes...)
				}
				break
			}
			zap.L().Error("Error reading data", zap.Error(err))
			return err
		}

		// -- Handle headers
		if !s.Headers.Finished {
			buffer.Write(line)
			strData := buffer.String()
			ProcessHeaders(strData, s)
			buffer.Reset()
			continue
		}

		// -- Handle data
		s.RawData = append(s.RawData, line...)
	}

	return nil
}

func ProcessHeaders(data string, session *Session) {

	if len(strings.TrimSpace(data)) == 0 {
		session.Headers.Finished = true
		return
	}

	lastHeader, _ := session.Headers.Data.Get(session.Headers.LastHeader)
	header, value, err := headers.ParseHeader(data, lastHeader)
	if err != nil {
		zap.L().Debug("Failed to parse header", zap.Error(err))
		return
	}

	zap.L().Debug("Header parsed", zap.String("header", header), zap.String("value", value))
	session.Headers.Data.Add(header, value)
	session.Headers.LastHeader = header
}
