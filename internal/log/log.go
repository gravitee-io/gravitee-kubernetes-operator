// Copyright (C) 2015 The Gravitee team (http://gravitee.io)
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//         http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package log

import (
	"context"

	"github.com/google/uuid"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/env"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"k8s.io/apimachinery/pkg/runtime/schema"
	ctrl "sigs.k8s.io/controller-runtime"
)

type key int

var ctxLoggerKey key
var ctxValueKey key

var impl *zap.Logger

const (
	kindField        = "kind"
	nsField          = "namespace"
	nameField        = "name"
	reconcileIdField = "reconcile-id"
)

func Info(message string, fields ...zapcore.Field) {
	impl.Info(message, fields...)
}

func Debug(message string, fields ...zapcore.Field) {
	impl.Debug(message, fields...)
}

func Warn(message string, fields ...zapcore.Field) {
	impl.Warn(message, fields...)
}

func Error(message string, fields ...zapcore.Field) {
	impl.Error(message, fields...)
}

func FromReconcileRequest(ctx context.Context, kind schema.GroupVersionKind, req ctrl.Request) *zap.Logger {
	reconcileId := uuid.New().String()

	logger := impl.With(
		zap.String(kindField, kind.Kind),
		zap.String(nsField, req.Namespace),
		zap.String(nameField, req.Name),
		zap.String(reconcileIdField, reconcileId),
	)

	newCtx := context.WithValue(ctx, ctxValueKey, map[string]interface{}{
		kindField:        kind.Kind,
		nsField:          req.Namespace,
		nameField:        req.Name,
		reconcileIdField: reconcileId,
	})

	return FromCtx(
		context.
			WithValue(newCtx, ctxLoggerKey, logger),
	)
}

func FromCtx(ctx context.Context) *zap.Logger {
	if known, ok := ctx.Value(ctxLoggerKey).(*zap.Logger); ok {
		if fields, fOK := ctx.Value(ctxValueKey).(map[string]interface{}); fOK {
			return known.With(
				zap.String(kindField, fields[kindField].(string)),
				zap.String(nsField, fields[nsField].(string)),
				zap.String(nameField, fields[nameField].(string)),
				zap.String(reconcileIdField, fields[reconcileIdField].(string)),
			)
		}
		return known
	}
	impl.Debug("no logger found in context, returning default logger")
	return impl
}

func init() {
	config := zap.Config{
		Level:         getLevel(),
		Encoding:      env.Config.LogFormat,
		EncoderConfig: setUpEncoderConfig(),
		OutputPaths:   []string{"stdout"},
	}
	impl = zap.Must(config.Build(zap.AddCallerSkip(1)))
}

func setUpEncoderConfig() zapcore.EncoderConfig {
	config := zap.NewProductionEncoderConfig()
	config.MessageKey = "message"
	config.TimeKey = "timestamp"
	config.LevelKey = "level"
	config.NameKey = "logger"
	config.CallerKey = ""
	config.StacktraceKey = "stacktrace"
	config.LineEnding = zapcore.DefaultLineEnding
	config.EncodeLevel = zapcore.LowercaseLevelEncoder
	config.EncodeTime = getTimeEncoder(env.Config.LogTimeFormat)
	return config
}

func getLevel() zap.AtomicLevel {
	level, err := zapcore.ParseLevel(env.Config.LogLevel)
	if err != nil {
		return zap.NewAtomicLevelAt(zap.InfoLevel)
	}
	return zap.NewAtomicLevelAt(level)
}

func getTimeEncoder(format string) zapcore.TimeEncoder {
	switch format {
	case "epochSeconds":
		return zapcore.EpochTimeEncoder
	case "epochMillis":
		return zapcore.EpochMillisTimeEncoder
	case "epochNanos":
		return zapcore.EpochNanosTimeEncoder
	case "iso8601":
		return zapcore.ISO8601TimeEncoder
	default:
		return zapcore.ISO8601TimeEncoder
	}
}
