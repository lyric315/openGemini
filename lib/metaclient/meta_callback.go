/*
Copyright 2022 Huawei Cloud Computing Technologies Co., Ltd.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

 http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package metaclient

import (
	"encoding/json"
	"errors"

	"github.com/openGemini/openGemini/app/ts-meta/meta/message"
	"github.com/openGemini/openGemini/engine/executor/spdy/transport"
	"github.com/openGemini/openGemini/open_src/influx/meta"
)

type BaseCallback struct {
}

func (c *BaseCallback) GetCodec() transport.Codec {
	return &message.MetaMessage{}
}

func (c *BaseCallback) Trans2MetaMsg(data interface{}) (*message.MetaMessage, error) {
	if metaMsg, ok := data.(*message.MetaMessage); !ok {
		return nil, errors.New("data is not a MetaMessage")
	} else {
		return metaMsg, nil
	}
}

type PingCallback struct {
	BaseCallback

	Leader []byte
}

func (c *PingCallback) Handle(data interface{}) error {
	metaMsg, err := c.Trans2MetaMsg(data)
	if err != nil {
		return err
	}
	msg, ok := metaMsg.Data().(*message.PingResponse)
	if !ok {
		return errors.New("data is not a PingResponse")
	}
	c.Leader = msg.Leader
	return nil
}

type PeersCallback struct {
	BaseCallback

	Peers []string
}

func (c *PeersCallback) Handle(data interface{}) error {
	metaMsg, err := c.Trans2MetaMsg(data)
	if err != nil {
		return err
	}
	msg, ok := metaMsg.Data().(*message.PeersResponse)
	if !ok {
		return errors.New("data is not a PeersResponse")
	}
	c.Peers = msg.Peers
	return nil
}

type CreateNodeCallback struct {
	BaseCallback

	NodeStartInfo *meta.NodeStartInfo
}

func (c *CreateNodeCallback) Handle(data interface{}) error {
	metaMsg, err := c.Trans2MetaMsg(data)
	if err != nil {
		return err
	}
	msg, ok := metaMsg.Data().(*message.CreateNodeResponse)
	if !ok {
		return errors.New("data is not a CreateNodeResponse")
	}
	if err = c.NodeStartInfo.UnMarshalBinary(msg.Data); err != nil {
		return err
	}
	return nil
}

type SnapshotCallback struct {
	BaseCallback

	Data []byte
}

func (c *SnapshotCallback) Handle(data interface{}) error {
	metaMsg, err := c.Trans2MetaMsg(data)
	if err != nil {
		return err
	}
	msg, ok := metaMsg.Data().(*message.SnapshotResponse)
	if !ok {
		return errors.New("data is not a SnapshotResponse")
	}
	c.Data = msg.Data
	return nil
}

type JoinCallback struct {
	BaseCallback

	NodeInfo *meta.NodeInfo
}

func (c *JoinCallback) Handle(data interface{}) error {
	metaMsg, err := c.Trans2MetaMsg(data)
	if err != nil {
		return err
	}
	msg, ok := metaMsg.Data().(*message.UpdateResponse)
	if !ok {
		return errors.New("data is not a UpdateResponse")
	}
	if msg.Err != "" {
		return errors.New(msg.Err)
	}
	if err = json.Unmarshal(msg.Data, c.NodeInfo); err != nil {
		return err
	}
	return nil
}

// Execute & Report
type ExecuteAndReportCallback struct {
	BaseCallback

	Typ   uint8
	Index uint64

	ErrCommand *errCommand
}

func (c *ExecuteAndReportCallback) Handle(data interface{}) error {
	metaMsg, err := c.Trans2MetaMsg(data)
	if err != nil {
		return err
	}
	switch c.Typ {
	case message.ExecuteRequestMessage:
		msg, ok := metaMsg.Data().(*message.ExecuteResponse)
		if !ok {
			return errors.New("data is not a ExecuteResponse")
		}
		if msg.Err != "" {
			return errors.New(msg.Err)
		}
		if msg.ErrCommand != "" {
			c.ErrCommand = &errCommand{msg: msg.ErrCommand}
		}
		c.Index = msg.Index
	case message.ReportRequestMessage:
		msg, ok := metaMsg.Data().(*message.ReportResponse)
		if !ok {
			return errors.New("data is not a ReportResponse")
		}
		if msg.Err != "" {
			return errors.New(msg.Err)
		}
		if msg.ErrCommand != "" {
			c.ErrCommand = &errCommand{msg: msg.ErrCommand}
		}
		c.Index = msg.Index
	default:
		panic("not support message type")
	}
	return nil
}

type GetShardInfoCallback struct {
	BaseCallback

	Data []byte
}

func (c *GetShardInfoCallback) Handle(data interface{}) error {
	metaMsg, err := c.Trans2MetaMsg(data)
	if err != nil {
		return err
	}
	msg, ok := metaMsg.Data().(*message.GetShardInfoResponse)
	if !ok {
		return errors.New("data is not a GetShardInfoResponse")
	}
	if msg.Err != "" {
		return errors.New(msg.Err)
	}
	c.Data = msg.Data
	return nil
}
