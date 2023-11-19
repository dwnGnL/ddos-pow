package client

import (
	"bufio"
	"fmt"
	"github.com/dwnGnL/ddos-pow/config"
	"github.com/dwnGnL/ddos-pow/lib/goerrors"
	"github.com/dwnGnL/ddos-pow/lib/pow"
	challengeResp "github.com/dwnGnL/ddos-pow/lib/protocol/challenge-resp"
	"github.com/goccy/go-json"
	"io"
	"log/slog"
	"net"
)

type Client struct {
	conf *config.Config
}

// RequestChallenge get challenge from server
func (s Client) RequestChallenge() (*pow.HashcashData, error) {
	address := fmt.Sprintf("%s:%d", s.conf.Server.Host, s.conf.Server.Port)

	// connect to server
	conn, err := net.Dial("tcp", address)
	if err != nil {
		slog.Warn("tcp dial", err)
		return nil, err
	}

	defer conn.Close()

	reader := bufio.NewReader(conn)

	// send request go get challenge
	err = sendMsg(challengeResp.Message{
		Header: challengeResp.REQUEST_CHALLENGE,
	}, conn)
	if err != nil {
		return nil, fmt.Errorf("err send request: %w", err)
	}

	// read and parse challenge
	msgStr, err := readConnMsg(reader)
	if err != nil {
		return nil, fmt.Errorf("err read msg: %w", err)
	}

	msg, err := challengeResp.ParseMessage(msgStr)
	if err != nil {
		return nil, fmt.Errorf("err parse msg: %w", err)
	}

	var hashcash pow.HashcashData

	err = json.Unmarshal([]byte(msg.Payload), &hashcash)
	if err != nil {
		return nil, fmt.Errorf("err parse hashcash: %w", err)
	}

	return &hashcash, nil
}

// RequestResource get quote from server
func (s Client) RequestResource(hashcash pow.HashcashData) (*challengeResp.Message, error) {
	// solve server challenge
	hashcash, err := hashcash.ComputeHashcash(s.conf.Pow.HashcashMaxIterations)
	if err != nil {
		return nil, fmt.Errorf("err compute hashcash: %w", err)
	}

	goerrors.Log().Println("hashcash computed:", hashcash)
	byteData, err := json.Marshal(hashcash)
	if err != nil {
		return nil, fmt.Errorf("err marshal hashcash: %w", err)
	}

	address := fmt.Sprintf("%s:%d", s.conf.Server.Host, s.conf.Server.Port)

	// connect to server
	conn, err := net.Dial("tcp", address)
	if err != nil {
		slog.Warn("tcp dial", err)
		return nil, err
	}

	defer conn.Close()

	reader := bufio.NewReader(conn)

	// send challenge solution to server
	err = sendMsg(challengeResp.Message{
		Header:  challengeResp.REQUEST_RESOURCE,
		Payload: string(byteData),
	}, conn)
	if err != nil {
		return nil, fmt.Errorf("err send request: %w", err)
	}

	goerrors.Log().Println("challenge sent to server")

	// parse server response
	msgStr, err := readConnMsg(reader)
	if err != nil {
		return nil, fmt.Errorf("err read msg: %w", err)
	}

	msg, err := challengeResp.ParseMessage(msgStr)
	if err != nil {
		return nil, fmt.Errorf("err parse msg: %w", err)
	}

	return msg, nil
}

func readConnMsg(reader *bufio.Reader) (string, error) {
	return reader.ReadString('\n')
}

func sendMsg(msg challengeResp.Message, conn io.Writer) error {
	msgStr := fmt.Sprintf("%s\n", msg.Stringify())
	_, err := conn.Write([]byte(msgStr))
	goerrors.Log().Println("msg = ", msgStr)
	return err
}

func New(conf *config.Config) *Client {
	return &Client{
		conf: conf,
	}
}
