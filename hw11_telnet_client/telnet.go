package main

import (
	"bufio"
	"bytes"
	"errors"
	"io"
	"net"
	"time"
)

var (
	ErrorNotConnected     = errors.New("not connected")
	ErrorAlreadyConnected = errors.New("already connected")
)

type TelnetClient interface {
	Connect() error
	Close() error
	Send() error
	Receive() error
}

type telnetClientImpl struct {
	timeout    time.Duration
	address    string
	conn       net.Conn
	in         io.ReadCloser
	out        io.Writer
	connReader *bufio.Reader
	inReader   *bufio.Reader
}

func readLine(reader *bufio.Reader) ([]byte, error) {
	var buffer bytes.Buffer

	for {
		line, isPrefix, err := reader.ReadLine()

		if err != nil {
			return nil, err
		}

		// дочитали до конца строки?
		if isPrefix {
			buffer.Write(line)
			continue
		}

		// собираем строку из нескольких частей?
		if buffer.Len() != 0 {
			buffer.Write(line)
			line = buffer.Bytes()
		}

		return line, nil
	}
}

func (tc *telnetClientImpl) pumpMessage(from *bufio.Reader, to io.Writer) (err error) {
	line, err := readLine(from)
	if err == nil {
		_, err = to.Write(append(line, '\n'))
	}

	return
}

func (tc *telnetClientImpl) Receive() (err error) {
	if tc.conn == nil {
		return ErrorNotConnected
	}
	return tc.pumpMessage(tc.connReader, tc.out)
}

func (tc *telnetClientImpl) Send() (err error) {
	if tc.conn == nil {
		return ErrorNotConnected
	}
	return tc.pumpMessage(tc.inReader, tc.conn)
}

func (tc *telnetClientImpl) Close() (err error) {
	if tc.conn == nil {
		return
	}

	tmpConn := tc.conn
	tc.conn = nil

	err = tmpConn.Close()

	return
}

func (tc *telnetClientImpl) Connect() (err error) {
	if tc.conn != nil {
		return ErrorAlreadyConnected
	}
	tc.conn, err = net.DialTimeout("tcp", tc.address, tc.timeout)

	if err != nil {
		return
	}

	tc.connReader = bufio.NewReader(tc.conn)
	tc.inReader = bufio.NewReader(tc.in)

	return
}

func NewTelnetClient(address string, timeout time.Duration, in io.ReadCloser, out io.Writer) TelnetClient {
	return &telnetClientImpl{
		timeout: timeout,
		address: address,
		in:      in,
		out:     out,
	}
}
