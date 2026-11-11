// Copyright 2015 Google Inc. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package fuse

import (
	"sync"

	"github.com/jacobsa/fuse/internal/buffer"
)

var opStatePool = sync.Pool{
	New: func() any {
		return new(opState)
	},
}

func getOpState() *opState {
	return opStatePool.Get().(*opState)
}

func putOpState(s *opState) {
	*s = opState{}
	opStatePool.Put(s)
}

////////////////////////////////////////////////////////////////////////
// buffer.InMessage
////////////////////////////////////////////////////////////////////////

func (c *Connection) getInMessage() *buffer.InMessage {
	x, ok := c.inMessages.Get().(*buffer.InMessage)
	if !ok || x == nil {
		x = buffer.NewInMessage(c.inMessageBufSize)
	}

	return x
}

func (c *Connection) putInMessage(x *buffer.InMessage) {
	c.inMessages.Put(x)
}

////////////////////////////////////////////////////////////////////////
// buffer.OutMessage
////////////////////////////////////////////////////////////////////////

func (c *Connection) getOutMessage() *buffer.OutMessage {
	x, ok := c.outMessages.Get().(*buffer.OutMessage)
	if !ok || x == nil {
		x = new(buffer.OutMessage)
	}
	x.Reset()

	return x
}

func (c *Connection) putOutMessage(x *buffer.OutMessage) {
	c.outMessages.Put(x)
}
