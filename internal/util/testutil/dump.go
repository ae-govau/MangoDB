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

package testutil

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/MangoDB-io/MangoDB/internal/util/hex"
)

func ParseDump(t testing.TB, s string) []byte {
	t.Helper()

	b, err := hex.ParseDump(s)
	require.NoError(t, err)
	return b
}

func ParseDumpFile(t testing.TB, path ...string) []byte {
	t.Helper()

	b, err := os.ReadFile(filepath.Join(path...))
	require.NoError(t, err)
	return ParseDump(t, string(b))
}

func MustParseDump(s string) []byte {
	b, err := hex.ParseDump(s)
	if err != nil {
		panic(err)
	}
	return b
}

func MustParseDumpFile(path ...string) []byte {
	b, err := os.ReadFile(filepath.Join(path...))
	if err != nil {
		panic(err)
	}
	return MustParseDump(string(b))
}
