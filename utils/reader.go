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

package utils

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"strings"
)

// ByteReader defines byte reader
type ByteReader interface {
	io.Reader
	io.ByteReader
}

// NewByteReader create a ByteReader from io.Reader
func NewByteReader(r io.Reader) ByteReader {
	if br, ok := r.(ByteReader); ok {
		return br
	}
	return &byteReader{Reader: r}
}

type byteReader struct {
	io.Reader
}

func (r *byteReader) ReadByte() (byte, error) {
	var buff = [1]byte{}
	_, err := r.Read(buff[:])
	return buff[0], err
}

// ToReader wrap message to io.Reader
func ToReader(message interface{}) (io.Reader, error) {

	switch r := message.(type) {
	case []byte:
		return bytes.NewReader(r), nil
	case [][]byte:
		readers := make([]io.Reader, 0, len(r))
		for _, b := range r {
			readers = append(readers, bytes.NewReader(b))
		}
		return io.MultiReader(readers...), nil
	case string:
		return strings.NewReader(r), nil
	case io.Reader:
		return r, nil
	default:
		return nil, fmt.Errorf("unrecognized type: %T", message)
	}
}

// MustToReader any error to panic
func MustToReader(message interface{}) io.Reader {
	r, err := ToReader(message)
	if nil != err {
		Assert(err)
	}
	return r
}

// ToBytes unwrap a message to []byte
func ToBytes(message interface{}) ([]byte, error) {

	switch r := message.(type) {
	case []byte:
		return r, nil
	case [][]byte:
		buffer := bytes.NewBuffer(make([]byte, 0, CountOf(r)))
		for _, b := range r {
			buffer.Write(b)
		}
		return buffer.Bytes(), nil
	case string:
		return []byte(r), nil
	case *bytes.Buffer:
		return r.Bytes(), nil
	case *bytes.Reader:
		return StealBytes(r)
	case *strings.Reader:
		return StealBytes(r)
	case io.WriterTo:
		return StealBytes(r)
	case io.Reader:
		return ioutil.ReadAll(r)
	default:
		return nil, fmt.Errorf("unrecognized type: %T", message)
	}
}

// MustToBytes any error to panic
func MustToBytes(message interface{}) []byte {
	r, err := ToBytes(message)
	if nil != err {
		Assert(err)
	}
	return r
}

// CountOf count bytes of buffs
func CountOf(buffs [][]byte) (n int64) {
	for _, buff := range buffs {
		n += int64(len(buff))
	}
	return
}

// ByteStealer steal from io.Reader
type ByteStealer struct {
	Data []byte
}

func (s *ByteStealer) Write(p []byte) (n int, err error) {
	if nil == s.Data {
		s.Data = p[0:len(p):len(p)]
	} else {
		s.Data = append(s.Data, p...)
	}
	return len(p), nil
}

func StealBytes(reader io.WriterTo) ([]byte, error) {
	var stealer ByteStealer
	n, err := reader.WriteTo(&stealer)
	if nil != err {
		return nil, err
	}

	// check data length
	if n != int64(len(stealer.Data)) {
		return nil, io.ErrShortWrite
	}

	return stealer.Data, nil
}
