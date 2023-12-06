// Copyright 2016 Michal Witkowski. All Rights Reserved.
// See LICENSE for licensing terms.

package grpc_prometheus

import (
	"time"

	"google.golang.org/grpc/codes"
)

type serverReporter struct {
	metrics     *ServerMetrics
	rpcType     grpcType
	serviceName string
	methodName  string
	startTime   time.Time
	remoteIp    string
}

func newServerReporter(m *ServerMetrics, rpcType grpcType, fullMethod, remoteIp string) *serverReporter {
	r := &serverReporter{
		metrics:  m,
		rpcType:  rpcType,
		remoteIp: remoteIp,
	}
	if r.metrics.serverHandledHistogramEnabled {
		r.startTime = time.Now()
	}
	r.serviceName, r.methodName = splitMethodName(fullMethod)
	r.metrics.serverStartedCounter.WithLabelValues(string(r.rpcType), r.serviceName, r.methodName, r.remoteIp).Inc()
	return r
}

func (r *serverReporter) ReceivedMessage() {
	r.metrics.serverStreamMsgReceived.WithLabelValues(string(r.rpcType), r.serviceName, r.methodName, r.remoteIp).Inc()
}

func (r *serverReporter) SentMessage() {
	r.metrics.serverStreamMsgSent.WithLabelValues(string(r.rpcType), r.serviceName, r.methodName, r.remoteIp).Inc()
}

func (r *serverReporter) Handled(code codes.Code) {
	r.metrics.serverHandledCounter.WithLabelValues(string(r.rpcType), r.serviceName, r.methodName, code.String(), r.remoteIp).Inc()
	if r.metrics.serverHandledHistogramEnabled {
		r.metrics.serverHandledHistogram.WithLabelValues(string(r.rpcType), r.serviceName, r.methodName, r.remoteIp).Observe(time.Since(r.startTime).Seconds())
	}
}
