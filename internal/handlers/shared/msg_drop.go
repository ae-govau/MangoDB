// Copyright 2021 Baltoro OÜ.
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

package shared

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v4"

	"github.com/MangoDB-io/MangoDB/internal/handlers/common"
	"github.com/MangoDB-io/MangoDB/internal/types"
	"github.com/MangoDB-io/MangoDB/internal/wire"
)

func (h *Handler) MsgDrop(ctx context.Context, header *wire.MsgHeader, msg *wire.OpMsg) (*wire.OpMsg, error) {
	if len(msg.Documents) != 1 {
		return nil, common.NewError(common.ErrNotImplemented, fmt.Errorf("multiple documents are not supported"), header, msg)
	}
	document := msg.Documents[0]

	m := document.Map()
	collection := m[document.Command()].(string)
	db := m["$db"].(string)

	// TODO probably not CASCADE
	sql := fmt.Sprintf(`DROP TABLE %s CASCADE`, pgx.Identifier{db, collection}.Sanitize())

	_, err := h.pgPool.Exec(ctx, sql)
	if err != nil {
		// TODO check error code
		return nil, common.NewError(common.ErrNamespaceNotFound, fmt.Errorf("ns not found"), header, msg)
	}

	reply := &wire.OpMsg{
		Documents: []types.Document{types.MakeDocument(
			"nIndexesWas", int32(0), // TOOD
			"ns", db+"."+collection,
			"ok", float64(1),
		)},
	}
	return reply, nil
}
