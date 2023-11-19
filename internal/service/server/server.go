package server

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/dwnGnL/ddos-pow/config"
	"github.com/dwnGnL/ddos-pow/lib/cache"
	"github.com/dwnGnL/ddos-pow/lib/goerrors"
	"github.com/dwnGnL/ddos-pow/lib/pow"
	challengeResp "github.com/dwnGnL/ddos-pow/lib/protocol/challenge-resp"
	"golang.org/x/exp/rand"
	"strconv"
	"time"
)

// Quotes of word of wisdom book
var Quotes = []string{
	"The only true wisdom is in knowing you know nothing",

	"At the end of the day you are your own lawmaker",

	"True wisdom comes to each of us when we realize how little we understand about life, ourselves, and the world around us",

	"People who are crazy enough to think they can change the world are the ones who do",

	"Personal Development Is A Major Time-Saver. The Better You Become, The Less Time It Takes You To Achieve Your Goals",
}

type Server struct {
	conf  *config.Config
	cache *cache.InMemoryCache
}

func (s Server) Ping() string {
	return "Pong..."
}

// ResponseChallenge send challenge to solve to client
func (s Server) ResponseChallenge(clientIP string) (msg *challengeResp.Message, err error) {
	randValue := rand.Intn(100000)
	err = s.cache.Add(randValue, s.conf.Pow.HashcashDuration)
	if err != nil {
		return nil, fmt.Errorf("err add rand to cache: %w", err)
	}

	// challenge for client
	hashcash := pow.HashcashData{
		Version:    1,
		ZerosCount: s.conf.Pow.HashcashZerosCount,
		Date:       time.Now().Unix(),
		Resource:   clientIP,
		Rand:       base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%d", randValue))),
		Counter:    0,
	}

	hashCashMarshaled, err := json.Marshal(hashcash)
	if err != nil {
		return nil, fmt.Errorf("err marshal hashCash: %v", err)
	}

	msg = &challengeResp.Message{
		Header:  challengeResp.RESPONSE_CHALLENGE,
		Payload: string(hashCashMarshaled),
	}

	return
}

// ResponseResource check challenge solution and send response
func (s Server) ResponseResource(clientIP string, hashCashSolved string) (msg *challengeResp.Message, err error) {
	goerrors.Log().Printf("client %s requests resource with payload %s\n", clientIP, hashCashSolved)
	var hashcash pow.HashcashData
	msg = new(challengeResp.Message)

	err = json.Unmarshal([]byte(hashCashSolved), &hashcash)
	if err != nil {
		return nil, fmt.Errorf("err unmarshal hashcash: %w", err)
	}

	goerrors.Log().Info("hashcash test 4", hashcash)

	if hashcash.Resource != clientIP {
		msg.Header = challengeResp.FAIL
		msg.Payload = "invalid hashcash resource"
		return msg, nil
	}

	randValueBytes, err := base64.StdEncoding.DecodeString(hashcash.Rand)
	if err != nil {
		return nil, fmt.Errorf("err decode rand: %w", err)
	}

	randValue, err := strconv.Atoi(string(randValueBytes))
	if err != nil {
		return nil, fmt.Errorf("err decode rand: %w", err)
	}

	exists, err := s.cache.Get(randValue)
	if err != nil {
		return nil, fmt.Errorf("err get rand from cache: %w", err)
	}

	if !exists {
		msg.Header = challengeResp.FAIL
		msg.Payload = "challenge expired or not sent"
		return msg, nil
	}

	if time.Now().Unix()-hashcash.Date > s.conf.Pow.HashcashDuration {
		msg.Header = challengeResp.FAIL
		msg.Payload = "challenge expired"
		return msg, nil
	}

	maxIter := hashcash.Counter
	if maxIter == 0 {
		maxIter = 1
	}

	_, err = hashcash.ComputeHashcash(maxIter)
	if err != nil {
		return nil, fmt.Errorf("invalid hashcash")
	}
	goerrors.Log().Printf("client %s succesfully computed hashcash %s\n", clientIP, hashCashSolved)

	// good job, send quote to client
	msg = &challengeResp.Message{
		Header:  challengeResp.RESPONSE_RESOURCE,
		Payload: Quotes[rand.Intn(4)],
	}

	return msg, nil
}

func New(conf *config.Config) *Server {
	return &Server{
		conf:  conf,
		cache: cache.InitInMemoryCache(),
	}
}
