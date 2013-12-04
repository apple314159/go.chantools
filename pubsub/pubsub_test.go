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

package pubsub

import (
	"testing"
)

func TestSub(t *testing.T) {
	queue1 := make(chan int)

	err := Sub(queue1, "i_wanna_integer")
	defer Unsub(queue1)
	if err != nil {
		t.Errorf("Subscription failed: %s", err)
	}

	queue2 := "not a channel"

	err = Sub(queue2, "i_wanna_error")
	if err == nil {
		t.Errorf("Subscription works only for channels, expected error, got %s", err)
	}
}

func TestPubSub(t *testing.T) {
	queue1 := make(chan int, 1)

	err := Sub(queue1, "i_wanna_integer")
	defer Unsub(queue1)
	if err != nil {
		t.Errorf("Subscription failed: %s", err)
	}

	Pub(42, "i_wanna_integer")

	i := <-queue1

	if i != 42 {
		t.Errorf("Expected receive 42, got %d", i)
	}
}

func TestUnsub(t *testing.T) {
	queue1 := make(chan string, 1)

	err := Sub(queue1, "i_wanna_string")
	if err != nil {
		t.Errorf("Subscription failed: %s", err)
	}

	Pub("fourty-two", "i_wanna_string")

	s := <-queue1

	if s != "fourty-two" {
		t.Errorf("Expected receive 'fourty-two', got %s", s)
	}

	err = Unsub(queue1)
	if err != nil {

	}

	Pub("fourty-two-again", "i_wanna_string")

	select {
	case s = <-queue1:
		t.Errorf("Unexpected receive from channel, got %s", s)
	default:
		break
	}

	queue2 := "not a channel"
	err = Unsub(queue2)
	if err == nil {
		t.Errorf("Unsubscription works only for channels, expected error, got %s", err)
	}
}
