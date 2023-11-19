package service

import (
	"github.com/dwnGnL/ddos-pow/config"
	"github.com/dwnGnL/ddos-pow/internal/service/client"
	"github.com/dwnGnL/ddos-pow/internal/service/server"
	"github.com/dwnGnL/ddos-pow/lib/pow"
	challenge_resp "github.com/dwnGnL/ddos-pow/lib/protocol/challenge-resp"
)

type ServiceImpl struct {
	conf   *config.Config
	server ServerService
	client ClientService
}

type ServerService interface {
	Ping() string
	ResponseChallenge(clientIP string) (msg *challenge_resp.Message, err error)
	ResponseResource(clientIP string, hashCashSolved string) (msg *challenge_resp.Message, err error)
}

type ClientService interface {
	RequestChallenge() (*pow.HashcashData, error)
	RequestResource(data pow.HashcashData) (*challenge_resp.Message, error)
}

type Option func(*ServiceImpl)

func New(conf *config.Config, opts ...Option) *ServiceImpl {
	s := ServiceImpl{
		conf:   conf,
		server: server.New(conf),
		client: client.New(conf),
	}

	for _, opt := range opts {
		opt(&s)
	}

	return &s
}

func (s ServiceImpl) GetServer() ServerService {
	return s.server
}

func (s ServiceImpl) GetClient() ClientService {
	return s.client
}
