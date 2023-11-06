/*
 * Copyright 2019 the go-netty project
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *      https://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package tcp

import (
	"context"
	"time"

	"github.com/mijingduI/go-netty/transport"
)

// DefaultOption default tcp options
var DefaultOption = &Options{
	Timeout:         time.Second * 5,
	KeepAlive:       true,
	KeepAlivePeriod: time.Minute,
	Linger:          -1,
	NoDelay:         true,
}

// Options fot tcp transport
type Options struct {
	Timeout         time.Duration `json:"timeout"`
	KeepAlive       bool          `json:"keep-alive"`
	KeepAlivePeriod time.Duration `json:"keep-alive-period"`
	Linger          int           `json:"linger"`
	NoDelay         bool          `json:"nodelay"`
	SockBuf         int           `json:"sockbuf"`
	ReadBufferSize  int           `json:"readBufferSize"`
	WriteBufferSize int           `json:"writeBufferSize"`
}

type contextKey struct{}

// WithOptions to wrap the tcp options
func WithOptions(option *Options) transport.Option {
	return func(options *transport.Options) error {
		options.Context = context.WithValue(options.Context, contextKey{}, option)
		return nil
	}
}

// FromContext to unwrap the tcp options
func FromContext(ctx context.Context, def *Options) *Options {
	if v, ok := ctx.Value(contextKey{}).(*Options); ok {
		return v
	}
	return def
}
