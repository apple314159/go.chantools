// Copyright (C) 2013 by Maxim Bublis <b@codemonkey.ru>
//
// Permission is hereby granted, free of charge, to any person obtaining
// a copy of this software and associated documentation files (the
// "Software"), to deal in the Software without restriction, including
// without limitation the rights to use, copy, modify, merge, publish,
// distribute, sublicense, and/or sell copies of the Software, and to
// permit persons to whom the Software is furnished to do so, subject to
// the following conditions:
//
// The above copyright notice and this permission notice shall be
// included in all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND,
// EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF
// MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND
// NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE
// LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION
// OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION
// WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

// Package pubsub implements channel Publish-Subscribe pattern.
package pubsub

import (
	"fmt"
	"reflect"
	"sync"
)

// Broker.
var broker = _broker{subscriptions: make(map[interface{}]map[reflect.Value]bool)}

// Broker structure.
type _broker struct {
	sync.Mutex
	subscriptions map[interface{}]map[reflect.Value]bool
}

// Sub subscribes channel ch to topic notifications.
func Sub(ch interface{}, topic interface{}) error {
	c := reflect.ValueOf(ch)

	if c.Type().Kind() != reflect.Chan {
		return fmt.Errorf("ch should be channel, got %s", c.Type().Kind())
	}

	broker.Lock()
	defer broker.Unlock()

	if subscribers := broker.subscriptions[topic]; subscribers == nil {
		broker.subscriptions[topic] = make(map[reflect.Value]bool)
	}

	broker.subscriptions[topic][c] = true

	return nil
}

// Unsub unsubscribes channel from notifications.
// Causes broker to stop sending messages to channel.
// When Unsub returns, it is guaranteed that channel will not receive
// more messages.
func Unsub(ch interface{}) error {
	c := reflect.ValueOf(ch)

	if c.Type().Kind() != reflect.Chan {
		return fmt.Errorf("ch should be channel, got %s", c.Type().Kind())
	}

	broker.Lock()
	defer broker.Unlock()

	for _, subscribers := range broker.subscriptions {
		delete(subscribers, c)
	}

	return nil
}

// Pub publishes message to subscribed channels.
func Pub(msg interface{}, topic interface{}) {
	broker.Lock()
	defer broker.Unlock()

	for subscriber := range broker.subscriptions[topic] {
		subscriber.Send(reflect.ValueOf(msg))
	}
}
