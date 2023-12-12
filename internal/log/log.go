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
	"fmt"

	"github.com/go-logr/zapr"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/env"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

var sink *zap.Logger

type raw struct {
	sink *zap.Logger
}

var Global raw

func (w raw) Debug(message string) {
	w.sink.Debug(message)
}

func (w raw) Info(message string) {
	w.sink.Info(message)
}

func (w raw) Infof(message string, args ...any) {
	w.sink.Info(fmt.Sprintf(message, args...))
}

func (w raw) Debugf(message string, args ...any) {
	w.sink.Debug(fmt.Sprintf(message, args...))
}

func (w raw) Warn(message string) {
	w.sink.Warn(message)
}

func (w raw) Error(err error, message string) {
	w.sink.Error(message, zap.Error(err))
}

func Debug(ctx context.Context, message string, keysAndValues ...any) {
	log.FromContext(ctx).V(1).Info(message, keysAndValues...)
}

func Info(ctx context.Context, message string, keysAndValues ...any) {
	log.FromContext(ctx).Info(message, keysAndValues...)
}

func Error(ctx context.Context, err error, message string, keysAndValues ...any) {
	log.FromContext(ctx).Error(err, message, keysAndValues...)
}

func InfoInitReconcile(crx context.Context) {
	log.FromContext(crx).Info("Initializing reconcile")
}

func InfoEndReconcile(crx context.Context) {
	log.FromContext(crx).Info("Reconcile done")
}

func ErrorRequeuingReconcile(crx context.Context, err error) {
	log.FromContext(crx).Error(err, "Requeuing reconcile due to error")
}

func ErrorAbortingReconcile(crx context.Context, err error) {
	log.FromContext(crx).Error(err, "Aborting reconcile due to an unrecoverable error")
}

func init() {
	config := zap.Config{
		Level:         getLogLevel(),
		Encoding:      env.Config.LogFormat,
		EncoderConfig: setUpEncoderConfig(),
		OutputPaths:   []string{"stdout"},
	}
	sink = zap.Must(config.Build())
	ctrl.SetLogger(zapr.NewLogger(sink))
	Global = raw{sink: sink}
}

func setUpEncoderConfig() zapcore.EncoderConfig {
	config := zap.NewProductionEncoderConfig()
	config.MessageKey = "message"
	config.TimeKey = env.Config.LogTimestampField
	config.LevelKey = "level"
	config.NameKey = "logger"
	config.CallerKey = ""
	config.StacktraceKey = "stacktrace"
	config.LineEnding = zapcore.DefaultLineEnding
	config.EncodeLevel = zapcore.LowercaseLevelEncoder
	config.EncodeTime = getTimeEncoder(env.Config.LogTimestampFormat)
	return config
}

func getLogLevel() zap.AtomicLevel {
	level, err := zapcore.ParseLevel(env.Config.LogLevel)
	if err != nil {
		return zap.NewAtomicLevelAt(zap.InfoLevel)
	}
	return zap.NewAtomicLevelAt(level)
}

func getTimeEncoder(format string) zapcore.TimeEncoder {
	switch format {
	case "epoch-second":
		return zapcore.EpochTimeEncoder
	case "epoch-millis":
		return zapcore.EpochMillisTimeEncoder
	case "epoch-nano":
		return zapcore.EpochNanosTimeEncoder
	case "iso-8601":
		return zapcore.ISO8601TimeEncoder
	default:
		return zapcore.ISO8601TimeEncoder
	}
}
