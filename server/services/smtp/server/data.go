package server

import (
	"bufio"
	"bytes"
	"errors"
	"io"
	"strings"

	"github.com/GrzegorzManiak/NoiseBackend/internal/email"
	"github.com/GrzegorzManiak/NoiseBackend/services/smtp/services"
	"go.uber.org/zap"
)

func (s *Session) Data(r io.Reader) error {
	var buffer, dkimBuffer bytes.Buffer

	multiWriter := io.MultiWriter(&buffer, &dkimBuffer)
	if _, err := io.Copy(multiWriter, r); err != nil {
		zap.L().Debug("Failed to copy data", zap.Error(err))

		return err
	}

	bufReader := bufio.NewReader(&buffer)
	dkimReader := bufio.NewReader(&dkimBuffer)
	s.DkimResult, _ = services.VerifyDkimSignature(dkimReader)

	for {
		//
		// Builds the line buffer so that we are able to read the data
		// line by line.
		//
		line, err := bufReader.ReadBytes('\n')
		if err != nil {
			if err == io.EOF {
				if buffer.Len() > 0 {
					s.RawData = append(s.RawData, buffer.Bytes()...)
				}

				break
			}

			return err
		}

		//
		// Append the line to the raw dat buffer.
		// This is what the user will see when they request the email.
		//
		s.RawData = append(s.RawData, line...)

		//
		// If the headers are not finished, we need to process them.
		//
		if !s.Headers.Finished {
			if err := processHeaders(line, s); err != nil {
				return err
			}

			continue
		}
	}

	return s.prepareInboundEmail()
}

func processHeaders(line []byte, s *Session) error {
	if len(strings.TrimSpace(string(line))) == 0 {
		s.Headers.Finished = true
		if !s.Headers.Data.Has(email.RequiredHeaders) {
			return errors.New("missing required headers")
		}

		return nil
	}

	lastHeader, _ := s.Headers.Data.Get(s.Headers.LastHeader)

	header, value, err := email.ParseHeader(string(line), lastHeader)
	if err != nil {
		return nil
	}

	header = strings.Trim(header, " \n\t")
	value = strings.Trim(value, " \n\t")

	s.Headers.Data.Add(header, value)
	s.Headers.LastHeader = header

	return nil
}
