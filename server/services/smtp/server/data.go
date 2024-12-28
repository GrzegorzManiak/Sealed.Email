package server

import (
	"bufio"
	"bytes"
	"github.com/GrzegorzManiak/NoiseBackend/services/smtp/headers"
	"go.uber.org/zap"
	"io"
	"strings"
)

func ProcessData(reader io.Reader, session *Session) error {
	zap.L().Debug("Data received")

	var buffer bytes.Buffer
	bufReader := bufio.NewReader(reader)

	for {
		line, err := bufReader.ReadBytes('\n')
		if err != nil {
			if err == io.EOF {
				if buffer.Len() > 0 {
					eofBytes := buffer.Bytes()
					session.RawData = append(session.RawData, eofBytes...)
				}
				break
			}
			zap.L().Error("Error reading data", zap.Error(err))
			return err
		}

		// -- Handle headers
		if !session.Headers.Finished {
			buffer.Write(line)
			strData := buffer.String()
			strData = strings.TrimSpace(strData)
			ProcessHeaders(strData, session)
			buffer.Reset()
			continue
		}

		// -- Handle data
		session.RawData = append(session.RawData, line...)
	}

	return nil
}

func ProcessHeaders(data string, session *Session) {

	if len(data) == 0 {
		session.Headers.Finished = true
		return
	}

	header, value, err := headers.ParseHeader(data)
	if err != nil {
		zap.L().Debug("Failed to parse header", zap.Error(err))
		return
	}

	zap.L().Debug("Header parsed", zap.String("header", header), zap.String("value", value))
	session.Headers.Data.Add(header, value)
}
