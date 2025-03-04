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

package main

import (
	"context"
	"flag"
	"fmt"
	"os/signal"

	"go.uber.org/zap"
	"golang.org/x/sys/unix"

	"github.com/MangoDB-io/MangoDB/internal/clientconn"
	"github.com/MangoDB-io/MangoDB/internal/pgconn"
	"github.com/MangoDB-io/MangoDB/internal/util/debug"
	"github.com/MangoDB-io/MangoDB/internal/util/logging"
	"github.com/MangoDB-io/MangoDB/internal/util/version"
)

// Flags are defined there to be visible in the testcover binary help output (bin/mangodb-testcover -h).
//
//nolint:gochecknoglobals
var (
	modeF          = flag.String("mode", string(clientconn.AllModes[0]), fmt.Sprintf("operation mode: %v", clientconn.AllModes))
	listenAddrF    = flag.String("listen-addr", "127.0.0.1:27017", "listen address")
	postgresqlURLF = flag.String("postgresql-url", "postgres://postgres@127.0.0.1:5432/mangodb", "PostgreSQL URL")
	shadowAddrF    = flag.String("shadow-addr", "127.0.0.1:37017", "")
	debugAddrF     = flag.String("debug-addr", "127.0.0.1:8088", "debug address")
	versionF       = flag.Bool("version", false, "show version and exit")
)

func main() {
	logger := logging.Setup(zap.DebugLevel)
	flag.Parse()

	info := version.Get()

	if *versionF {
		logger.Info(info.Version, zap.String("version", info.Version), zap.String("commit", info.Commit), zap.Bool("dirty", info.Dirty))
		return
	}

	logger.Info("Starting MangoDB "+info.Version+"...", zap.String("commit", info.Commit), zap.Bool("dirty", info.Dirty))

	var found bool
	for _, m := range clientconn.AllModes {
		if *modeF == string(m) {
			found = true
			break
		}
	}
	if !found {
		logger.Sugar().Fatalf("Unknown mode %q.", *modeF)
	}

	ctx, stop := signal.NotifyContext(context.Background(), unix.SIGTERM, unix.SIGINT)
	go func() {
		<-ctx.Done()
		logger.Info("Stopping...")
		stop()
	}()

	go debug.RunHandler(ctx, *debugAddrF, logger.Named("debug"))

	pgPool, err := pgconn.NewPool(*postgresqlURLF, logger)
	if err != nil {
		logger.Fatal(err.Error())
	}
	defer pgPool.Close()

	l := clientconn.NewListener(&clientconn.NewListenerOpts{
		Addr:       *listenAddrF,
		ShadowAddr: *shadowAddrF,
		Mode:       clientconn.Mode(*modeF),
		PgPool:     pgPool,
		Logger:     logger.Named("listener"),
	})

	if err = l.Run(ctx); err != nil {
		logger.Error("Listener stopped", zap.Error(err))
	}
}
