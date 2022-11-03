package config

import "time"

const (
	NATSDurableID                       = "hello-sub"
	ContextCaller                       = "hello-caller"
	DefaultNATSJSHost                   = "nats://localhost:4222"
	DefaultNATSJSRetryOnFailedConnect   = true
	DefaultNATSJSMaxReconnect           = -1
	DefaultNATSJSReconnectWait          = 1 * time.Second
	DefaultNATSJSRetryAttempts          = 3
	DefaultNATSJSRetryInterval          = 2 * time.Second
	DefaultNATSJSSubscribeRetryAttempts = 3
	DefaultNATSJSSubscribeRetryInterval = 2 * time.Second
	DefaultNATSJSStreamMaxAge           = 1 * 24 * time.Hour
	DefaultNATSJSDeliveryTimeInMinute   = -5
	DefaultNATSJSStreamMaxMessages      = 100000
)
