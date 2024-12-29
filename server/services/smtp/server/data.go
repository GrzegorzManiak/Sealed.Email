package server

import (
	"bufio"
	"bytes"
	"fmt"
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
			return err
		}

		// -- Handle headers
		if !s.Headers.Finished {
			buffer.Write(line)
			strData := buffer.String()
			err := ProcessHeaders(strData, s)
			if err != nil {
				return err
			}
			buffer.Reset()
			continue
		}

		// -- Handle data
		s.RawData = append(s.RawData, line...)
	}

	return nil
}

func ProcessHeaders(data string, session *Session) error {

	if len(strings.TrimSpace(data)) == 0 {
		session.Headers.Finished = true
		if !session.Headers.Data.Has(headers.RequiredHeaders) {
			return fmt.Errorf("missing required headers")
		}
		return nil
	}

	lastHeader, _ := session.Headers.Data.Get(session.Headers.LastHeader)
	header, value, err := headers.ParseHeader(data, lastHeader)
	if err != nil {
		return nil
	}

	header = strings.Trim(header, " \n\t")
	value = strings.Trim(value, " \n\t")

	session.Headers.Data.Add(header, value)
	session.Headers.LastHeader = header

	return nil
}
