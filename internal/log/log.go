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
	"io"
	"os"
	"strings"

	stdLog "log"

	"github.com/go-logr/zapr"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/core"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/env"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	kLog "k8s.io/klog/v2"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
	kZap "sigs.k8s.io/controller-runtime/pkg/log/zap"
)

var sink *zap.Logger

type raw struct {
	sink *zap.Logger
}

func (w *raw) Write(p []byte) (int, error) {
	message := strings.TrimSpace(string(p))
	w.sink.Info(message)
	return len(p), nil
}

var Global raw

func (w *raw) Debug(message string) {
	w.sink.Debug(message)
}

func (w *raw) Info(message string) {
	w.sink.Info(message)
}

func (w *raw) Infof(message string, args ...any) {
	w.sink.Info(fmt.Sprintf(message, args...))
}

func (w *raw) Debugf(message string, args ...any) {
	w.sink.Debug(fmt.Sprintf(message, args...))
}

func (w *raw) Warn(message string) {
	w.sink.Warn(message)
}

func (w *raw) Error(err error, message string) {
	w.sink.Error(message, zap.Error(err))
}

func (w *raw) Errorf(err error, message string, args ...any) {
	w.sink.Error(fmt.Sprintf(message, args...), zap.Error(err))
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

func InfoInitReconcile(ctx context.Context, obj client.Object) {
	log.FromContext(ctx).Info(
		fmt.Sprintf(
			"Reconciling resource [%s]",
			obj.GetObjectKind().GroupVersionKind().GroupKind(),
		),
		KeyValues(obj)...,
	)
}

func InfoEndReconcile(ctx context.Context, obj client.Object) {
	log.FromContext(ctx).Info(
		fmt.Sprintf(
			"Resource [%s] has been successfully reconciled",
			obj.GetObjectKind().GroupVersionKind().GroupKind(),
		),
		KeyValues(obj)...,
	)
}

func ErrorRequeuingReconcile(ctx context.Context, err error, obj client.Object) {
	log.FromContext(ctx).Error(
		err,
		fmt.Sprintf(
			"Requeuing reconcile of resource [%s] due to an error",
			obj.GetObjectKind().GroupVersionKind().GroupKind(),
		),
		KeyValues(obj)...,
	)
}

func ErrorAbortingReconcile(ctx context.Context, err error, obj client.Object) {
	log.FromContext(ctx).Error(
		err,
		fmt.Sprintf(
			"Aborting reconcile of resource [%s] due to an unrecoverable error",
			obj.GetObjectKind(),
		),
		KeyValues(obj)...,
	)
}

func KeyValues(obj client.Object, keyValues ...any) []any {
	props := []any{}
	if cta, ok := obj.(core.ContextAwareObject); ok {
		if cta.HasContext() {
			props = append(props, "contextName")
			props = append(props, cta.ContextRef().GetName())
			props = append(props, "contextNamespace")
			props = append(props, cta.ContextRef().GetNamespace())
		}
		props = append(props, "resourceID")
		props = append(props, cta.GetID())
		props = append(props, "environmentID")
		props = append(props, cta.GetEnvID())
		props = append(props, "organizationID")
		props = append(props, cta.GetOrgID())
	}
	props = append(props, keyValues...)
	return props
}

func init() {
	config := zap.Config{
		Level:         getLogLevel(),
		Encoding:      getEncoding(),
		EncoderConfig: setUpEncoderConfig(),
		OutputPaths:   []string{"stdout"},
	}
	sink = zap.Must(config.Build())
	if isSilent() {
		logger := kZap.New(kZap.WriteTo(io.Discard), kZap.UseDevMode(true))
		ctrl.SetLogger(logger)
		log.SetLogger(logger)
		kLog.SetLogger(logger)
		stdLog.SetOutput(io.Discard)
		Global = raw{sink: sink}
	} else {
		logger := zapr.NewLogger(sink)
		ctrl.SetLogger(logger)
		log.SetLogger(logger)
		kLog.SetLogger(logger)
		Global = raw{sink: sink}
		stdLog.SetOutput(&Global)
	}
}

func getEncoding() string {
	if env.Config.Development || env.Config.LogsFormat != "json" {
		return "console"
	}
	return "json"
}

func setUpEncoderConfig() zapcore.EncoderConfig {
	config := zap.NewProductionEncoderConfig()
	config.MessageKey = "message"
	config.TimeKey = env.Config.LogsTimestampField
	config.LevelKey = "level"
	config.NameKey = "logger"
	config.CallerKey = ""
	config.StacktraceKey = "stacktrace"
	config.LineEnding = zapcore.DefaultLineEnding
	config.EncodeLevel = getLevelCase()
	config.EncodeTime = getTimeEncoder()
	return config
}

func getLevelCase() zapcore.LevelEncoder {
	if env.Config.LogsLevelCase == "lower" {
		return zapcore.LowercaseLevelEncoder
	}
	return zapcore.CapitalLevelEncoder
}

func getLogLevel() zap.AtomicLevel {
	level, err := zapcore.ParseLevel(env.Config.LogsLevel)
	if err != nil {
		return zap.NewAtomicLevelAt(zap.InfoLevel)
	}
	return zap.NewAtomicLevelAt(level)
}

func getTimeEncoder() zapcore.TimeEncoder {
	switch env.Config.LogsTimestampFormat {
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

func isSilent() bool {
	silent := os.Getenv("GKO_MANAGER_SILENT_LOG")
	return silent == env.TrueString
}
