package main

import (
	"bufio"
	"context"
	"fmt"
	"github.com/fasthttp/websocket"
	"go.uber.org/zap"
	"libsocket"
	"log"
	"net/url"
	"os"
	"strings"
	"time"
)

func main() {
	// Create logger
	//logger := libsocket.NewTestLogger(os.Stdout)
	logger, _ := zap.NewProduction()
	defer logger.Sync() // flushes buffer, if any
	sugar := logger.Sugar()
	zapLogger := libsocket.NewLogger(sugar)
	// Set up connection parameters
	wsURL, _ := url.Parse("ws://127.0.0.1:8080/ws?token=demo")
	params := libsocket.OpenConnectionParams{
		URL: *wsURL,
	}

	// Create a params repo that will provide connection parameters
	paramsGetter := func(ctx context.Context) (libsocket.OpenConnectionParams, error) {
		return params, nil
	}
	paramsRepo := libsocket.NewOpenConnectionParamsRepo(zapLogger, paramsGetter)

	// Create WebSocket dialer
	dialer := websocket.DefaultDialer

	// Define message handler
	messageHandler := func(client libsocket.Client, msg libsocket.Message) {
		log.Printf("Received message: %s", msg.Data())
		// Process message
	}

	// Define event handler
	eventHandler := func(client libsocket.Client, event libsocket.EventType) {
		switch event {
		case libsocket.EventConnect:
			log.Println("Connected")
		case libsocket.EventReconnect:
			log.Println("Reconnected")
		case libsocket.EventClose:
			log.Println("Closed")
		}
	}

	// Create a connection factory with passive keep-alive support
	connFactory := libsocket.NewWebsocketFactory(
		zapLogger,
		dialer,
		paramsRepo,
		libsocket.ErrorAdapters{},
	)

	// Create the client
	clientFactory := libsocket.NewBasicClientFactory(
		messageHandler,
		eventHandler,
	)

	client := clientFactory(
		zapLogger,
		connFactory,
		libsocket.WithReopenParam(libsocket.NewReOpenParam(true, time.NewTicker(60*time.Second))),
		libsocket.WithHeartBeatParam(libsocket.NewDefaultHeartBeat()),
		libsocket.WithReConnParam(libsocket.NewReConnParam(true, 10, 1*time.Second, libsocket.ExponentialBackoffSeconds)))

	// Connect
	ctx := context.Background()
	if err := client.Open(ctx); err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}

	// Send a message
	//client.Send(libsocket.NewDataMessage([]byte("Hello, WebSocket server!")))
	go openReadCMD(client)
	// Wait for connection to close
	<-client.CloseChan()
}

func openReadCMD(c libsocket.Client) {
	reader := bufio.NewReader(os.Stdin)
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("输入错误:", err)
			break
		}
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		// 直接发到服务端
		c.Send(libsocket.NewEchoMessage([]byte(line)))
	}
}
