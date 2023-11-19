package server

import (
	"bufio"
	"fmt"
	"github.com/dwnGnL/ddos-pow/config"
	"github.com/dwnGnL/ddos-pow/internal/application"
	"github.com/dwnGnL/ddos-pow/lib/goerrors"
	challengeResp "github.com/dwnGnL/ddos-pow/lib/protocol/challenge-resp"
	"net"
	"strings"
)

// SetupHandlers open tcp connection
func SetupHandlers(core application.Core, cfg *config.Config) error {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", cfg.Server.Port))
	if err != nil {
		return fmt.Errorf("net listen err: %w", err)
	}

	defer listener.Close()
	for {
		conn, err := listener.Accept()
		if err != nil {
			return fmt.Errorf("error accept connection: %w", err)
		}

		go handleConnection(core, conn)
	}
}

// handleConnection each new tcp connection
func handleConnection(core application.Core, conn net.Conn) {
	defer conn.Close()

	reader := bufio.NewReader(conn)

	for {
		req, err := reader.ReadString('\n')
		if err != nil {
			return
		}
		msg, err := processRequest(core, req, conn.RemoteAddr().String())
		if err != nil {
			return
		}
		if msg != nil {
			err := sendMsg(*msg, conn)
			if err != nil {
				goerrors.Log().Warn("err send message:", err)
			}
		}
	}
}

// processRequest process tcp request
func processRequest(core application.Core, msgStr string, clienInfo string) (*challengeResp.Message, error) {
	msg, err := challengeResp.ParseMessage(msgStr)
	if err != nil {
		return nil, err
	}

	clientIP := strings.Split(clienInfo, ":")[0]

	serverService := core.GetServer()
	command := msg.Header

	switch command {
	case challengeResp.REQUEST_CHALLENGE:
		msg, err = serverService.ResponseChallenge(clientIP)
	case challengeResp.REQUEST_RESOURCE:
		msg, err = serverService.ResponseResource(clientIP, msg.Payload)
	}

	if err != nil {
		goerrors.Log().Warn("err global", err)
		return nil, err
	}

	return msg, err
}

func sendMsg(msg challengeResp.Message, conn net.Conn) error {
	msgStr := fmt.Sprintf("%s\n", msg.Stringify())
	_, err := conn.Write([]byte(msgStr))
	return err
}
