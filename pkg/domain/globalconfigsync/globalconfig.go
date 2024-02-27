// Copyright 2023 PingCAP, Inc.
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

package globalconfigsync

import (
	"context"
	"github.com/pingcap/kvproto/pkg/meta_storagepb"
	"github.com/pingcap/tidb/pkg/util/logutil"
	pd "github.com/tikv/pd/client"
	"go.uber.org/zap"
)

// GlobalConfigPath as Etcd prefix
const GlobalConfigPath = "/global/config/"

// GlobalConfigSyncer is used to sync pd global config.
type GlobalConfigSyncer struct {
	pd       pd.Client
	NotifyCh chan meta_storagepb.PutRequest
}

// NewGlobalConfigSyncer creates a GlobalConfigSyncer.
func NewGlobalConfigSyncer(p pd.Client) *GlobalConfigSyncer {
	return &GlobalConfigSyncer{
		pd:       p,
		NotifyCh: make(chan meta_storagepb.PutRequest, 8),
	}
}

// StoreGlobalConfig is used to store global config.
func (s *GlobalConfigSyncer) StoreGlobalConfig(ctx context.Context, putReq meta_storagepb.PutRequest) error {
	if s.pd == nil {
		return nil
	}
	_, err := s.pd.Put(ctx, putReq.Key, putReq.Value)
	if err != nil {
		return err
	}

	logutil.BgLogger().Info("store global config", zap.String("name", string(putReq.Key)), zap.String("value", string(putReq.Value)))
	return nil
}

// Notify pushes global config to internal channel and will be sync into pd's GlobalConfig.
func (s *GlobalConfigSyncer) Notify(putReq meta_storagepb.PutRequest) {
	s.NotifyCh <- putReq
}
