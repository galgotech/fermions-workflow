package server

import (
	"context"
	"net/http"

	"github.com/centrifugal/centrifuge"
	cloudevents "github.com/cloudevents/sdk-go/v2"

	"github.com/galgotech/fermions-workflow/pkg/bus"
	"github.com/galgotech/fermions-workflow/pkg/concurrency"
	"github.com/galgotech/fermions-workflow/pkg/log"
)

var loggerCF = log.New("live.centrifuge")

func handleLog(msg centrifuge.LogEntry) {
	arr := make([]interface{}, 0)
	for k, v := range msg.Fields {
		if v == nil {
			v = "<nil>"
		} else if v == "" {
			v = "<empty>"
		}
		arr = append(arr, k, v)
	}

	switch msg.Level {
	case centrifuge.LogLevelTrace:
		loggerCF.Trace(msg.Message, arr...)
	case centrifuge.LogLevelDebug:
		loggerCF.Debug(msg.Message, arr...)
	case centrifuge.LogLevelError:
		loggerCF.Error(msg.Message, arr...)
	case centrifuge.LogLevelInfo:
		loggerCF.Info(msg.Message, arr...)
	default:
		loggerCF.Debug(msg.Message, arr...)
	}
}

func NewCentrifuge(busEvent bus.Bus) (*centrifuge.WebsocketHandler, error) {
	node, err := centrifuge.New(centrifuge.Config{
		LogLevel:   centrifuge.LogLevelTrace,
		LogHandler: handleLog,
	})
	if err != nil {
		return nil, err
	}

	// Override default broker which does not use HistoryMetaTTL.
	broker, err := centrifuge.NewMemoryBroker(node, centrifuge.MemoryBrokerConfig{})
	if err != nil {
		return nil, err
	}
	node.SetBroker(broker)

	node.OnConnecting(func(ctx context.Context, event centrifuge.ConnectEvent) (centrifuge.ConnectReply, error) {
		cred, ok := centrifuge.GetCredentials(ctx)
		if !ok {
			return centrifuge.ConnectReply{}, centrifuge.DisconnectServerError
		}
		loggerCF.Debug("connecting", "user", cred.UserID, "transport", event.Name)
		return centrifuge.ConnectReply{}, nil
	})

	node.OnConnect(func(client *centrifuge.Client) {
		ctx, ctxCancel := context.WithCancel(context.Background())
		transport := client.Transport()
		loggerCF.Debug("connected", "user", client.UserID(), "transport", transport.Name())

		client.OnAlive(func() {
			loggerCF.Info("user connection is still active", "userID", client.UserID())
		})

		client.OnSubscribe(func(e centrifuge.SubscribeEvent, cb centrifuge.SubscribeCallback) {
			loggerCF.Info("user subscribes", "userID", client.UserID(), "channel", e.Channel)

			go func() {
				subscribeChan := busEvent.Subscribe(ctx, e.Channel)
				for subscribe := range concurrency.OrDoneCtx(ctx, subscribeChan) {
					if subscribe.Err != nil {
						loggerCF.Debug("error subscribe", "error", err)
						continue
					}

					data, err := subscribe.Event.MarshalJSON()
					if err != nil {
						loggerCF.Debug("error publishing", "error", err)
						continue
					}

					_, err = node.Publish(e.Channel, data)
					if err != nil {
						loggerCF.Debug("error publishing to personal channel", "error", err)
						continue
					}
				}
			}()

			cb(centrifuge.SubscribeReply{}, nil)
		})

		client.OnUnsubscribe(func(e centrifuge.UnsubscribeEvent) {
			loggerCF.Info("user unsubscribed from", "user", client.UserID(), "channel", e.Channel)
			ctxCancel()
		})

		client.OnPublish(func(e centrifuge.PublishEvent, cb centrifuge.PublishCallback) {
			if e.Channel != "fermions-worklow-ui" {
				loggerCF.Warn("publish received invalid channel", "userID", client.UserID(), "channel", e.Channel)
				cb(centrifuge.PublishReply{}, centrifuge.ErrorBadRequest)
				return
			}
			loggerCF.Info("user publishes into channel", "userID", client.UserID(), "channel", e.Channel)

			event := cloudevents.NewEvent()
			err := event.UnmarshalJSON(e.Data)
			if err != nil {
				cb(centrifuge.PublishReply{}, centrifuge.ErrorBadRequest)
				return
			}

			busEvent.Publish(ctx, event)
			cb(centrifuge.PublishReply{}, nil)
		})

		client.OnDisconnect(func(e centrifuge.DisconnectEvent) {
			loggerCF.Debug("user disconnected", "userID", client.UserID(), "disconnection", e.Disconnect)
			ctxCancel()
		})
	})

	// We also start a separate goroutine for centrifuge itself, since we
	// still need to run gin web server.
	go func() {
		if err := node.Run(); err != nil {
			loggerCF.Fatal(err.Error())
		}
	}()

	return centrifuge.NewWebsocketHandler(node, centrifuge.WebsocketConfig{
		ReadBufferSize:     1024,
		UseWriteBufferPool: true,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}), nil

}
