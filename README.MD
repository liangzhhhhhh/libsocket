## Hierarchy Overview
1. Connection (Lowest Layer)

Responsible for directly interacting with WebSocket/TCP (sending/receiving messages, managing socket lifecycle).

Operates on the level of byte streams and channels.

2. ConnectionHandler (Middle Layer)

Event dispatching (emitter → EventHandler)

Message handling (MessageHandler)

May include reconnection logic / keepalive mechanisms

Operates on the level of message streams and the event system.

3. Client (Top Layer)

Further wraps the ConnectionHandler, providing a simplified interface:

Open(): automatically establishes the connection and registers handlers

Send(): directly callable by the business/application layer

CloseChan(): convenient for select to wait for shutdown

Operates on the level of business/application developers.

## Usage Examples

### Basic WebSocket Client

```go
func main() {
// Create logger
logger := libsocket.NewTestLogger(os.Stdout)

// Set up connection parameters
wsURL, _ := url.Parse("ws://127.0.0.1:8080/ws?token=demo")
params := libsocket.OpenConnectionParams{
URL: *wsURL,
}

// Create a params repo that will provide connection parameters
paramsGetter := func(ctx context.Context) (libsocket.OpenConnectionParams, error) {
return params, nil
}
paramsRepo := libsocket.NewOpenConnectionParamsRepo(logger, paramsGetter)

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
logger,
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
logger,
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
```

### Using Active Keep-Alive

```go
// Create a connection param with active keep-alive
libsocket.WithReopenParam(libsocket.NewReOpenParam(true, time.NewTicker(60*time.Second)))
```

### Using Periodic Reconnection

```go
// Create a connection param that reconnects every 5 minutes
libsocket.WithReConnParam(libsocket.NewReConnParam(true, 10, 1*time.Second, libsocket.ExponentialBackoffSeconds)))
```

## Architecture

The library is composed of several layers:

1. **Client**: The high-level API that applications interact with
2. **Connection Handlers**: Middleware components that add behavior to connections
3. **Connection**: Low-level component that manages the actual WebSocket connection
4. **Messages**: Structured data types for communication
5. **Events**: Notification system for connection state changes

## Connection Handler Types

- **Basic Connection**: Simple pass-through connection handler
- **Backoff Connection**: Reconnects with exponential backoff on failure
- **Reopen Interval Connection**: Periodically creates a new connection
- **Active Keep-Alive**: Sends periodic ping messages
- **Passive Keep-Alive**: Responds to ping messages with pongs

## License

MIT License

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.