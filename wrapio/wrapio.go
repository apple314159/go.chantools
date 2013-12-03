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

// Package wrapio implements channel I/O functionality.
package wrapio

import (
	"fmt"
	"reflect"
)

// Decoder interface.
type Decoder interface {
	Decode(e interface{}) error
}

// Encoder interface.
type Encoder interface {
	Encode(e interface{}) error
}

func Notify(ch interface{}, d Decoder) error {
	c := reflect.ValueOf(ch)

	if c.Type().Kind() != reflect.Chan {
		return fmt.Errorf("ch should be channel, got %s", c.Type().Kind())
	}

	t := c.Type().Elem()

	go func() {
		defer c.Close()

		for {
			v := reflect.New(t).Interface()
			err := d.Decode(v)
			if err != nil {
				break
			}

			c.Send(reflect.ValueOf(v).Elem())
		}
	}()
	return nil
}

func Listen(ch interface{}, e Encoder) error {
	c := reflect.ValueOf(ch)

	if c.Type().Kind() != reflect.Chan {
		return fmt.Errorf("ch should be channel, got %s", c.Type().Kind())
	}

	go func() {
		for {
			v, ok := c.Recv()
			if !ok {
				break
			}
			err := e.Encode(v.Interface())
			if err != nil {
				break
			}
		}
	}()
	return nil
}
