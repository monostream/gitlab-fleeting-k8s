package kubectl

import (
	"errors"
	"io"
	"log"
	"strings"
	"syscall"
)

func isConnectionError(err error) bool {
	retryable := errors.Is(err, syscall.ECONNRESET) ||
		errors.Is(err, syscall.ECONNREFUSED) ||
		errors.Is(err, syscall.ECONNABORTED) ||
		errors.Is(err, syscall.EPIPE) ||
		errors.Is(err, io.ErrUnexpectedEOF) ||
		strings.Contains(err.Error(), "unexpected EOF") ||
		strings.Contains(err.Error(), "TLS handshake timeout")

	if retryable {
		log.Println("connection error, retrying...: " + err.Error())
	}

	return retryable
}
