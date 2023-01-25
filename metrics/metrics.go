package metrics

import (
	"context"

	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/metric/instrument"
	"go.opentelemetry.io/otel/metric/instrument/asyncint64"
	"go.opentelemetry.io/otel/metric/instrument/syncint64"
	"go.opentelemetry.io/otel/metric/unit"
	"go.opentelemetry.io/otel/sdk/instrumentation"
)

const (
	Version = "0.11.0"
)

var (
	ClientScope = instrumentation.Scope{
		Name:    "libp2p.io/bitswap/client",
		Version: Version,
	}

	ServerScope = instrumentation.Scope{
		Name:    "libp2p.io/bitswap/server",
		Version: Version,
	}

	SessionScope = instrumentation.Scope{
		Name:    "libp2p.io/bitswap/session",
		Version: Version,
	}
)

type ClientMetrics struct {
	m metric.Meter

	// Synchronous
	DupMetric syncint64.Histogram
	AllMetric

	// Asynchronous
	BlocksReceived    asyncint64.Counter
	DataReceived      asyncint64.Counter
	DupBlocksReceived asyncint64.Counter
	DupDataReceived   asyncint64.Counter
	MessagesReceived  asyncint64.Counter
}

func NewClientMetrics(meterProvider metric.MeterProvider) (*ClientMetrics, error) {
	m := meterProvider.Meter(ClientScope.Name, metric.WithInstrumentationVersion(ClientScope.Version), metric.WithSchemaURL(ClientScope.SchemaURL))

	blocksReceived, err := m.AsyncInt64().Counter(
		"blocks_received",
		instrument.WithUnit(unit.Dimensionless),
		instrument.WithDescription("Total number of received blocks"),
	)
	if err != nil {
		return nil, err
	}

	dataReceived, err := m.AsyncInt64().Counter(
		"data_received",
		instrument.WithUnit(unit.Bytes),
		instrument.WithDescription("Total number of data bytes received"),
	)
	if err != nil {
		return nil, err
	}

	dupBlocksReceived, err := m.AsyncInt64().Counter(
		"dup_blocks_received",
		instrument.WithUnit(unit.Dimensionless),
		instrument.WithDescription("Total number of duplicate blocks received"),
	)
	if err != nil {
		return nil, err
	}

	dupDataReceived, err := m.AsyncInt64().Counter(
		"dup_data_received",
		instrument.WithUnit(unit.Bytes),
		instrument.WithDescription("Total number of duplicate data bytes received"),
	)
	if err != nil {
		return nil, err
	}

	messagesReceived, err := m.AsyncInt64().Counter(
		"messages_received",
		instrument.WithUnit(unit.Dimensionless),
		instrument.WithDescription("Total number of messages received"),
	)
	if err != nil {
		return nil, err
	}

	return &ClientMetrics{
		BlocksReceived:    blocksReceived,
		DataReceived:      dataReceived,
		DupBlocksReceived: dupBlocksReceived,
		DupDataReceived:   dupDataReceived,
		MessagesReceived:  messagesReceived,
	}, nil
}

func (m *ClientMetrics) RegisterCallback(cb func(context.Context)) error {
	instruments := []instrument.Asynchronous{
		m.BlocksReceived,
		m.DataReceived,
		m.DupBlocksReceived,
		m.DupDataReceived,
		m.MessagesReceived,
	}
	return m.m.RegisterCallback(instruments, cb)
}

type ServerMetrics struct {
	m          metric.Meter
	BlocksSent asyncint64.Counter
	DataSent   asyncint64.Counter
}

func NewServerMetrics(meterProvider metric.MeterProvider) (*ServerMetrics, error) {
	m := meterProvider.Meter(ServerScope.Name, metric.WithInstrumentationVersion(ServerScope.Version), metric.WithSchemaURL(ServerScope.SchemaURL))

	blocksSent, err := m.AsyncInt64().Counter(
		"blocks_sent",
		instrument.WithUnit(unit.Dimensionless),
		instrument.WithDescription("Total number of sent blocks"),
	)
	if err != nil {
		return nil, err
	}

	dataSent, err := m.AsyncInt64().Counter(
		"data_sent",
		instrument.WithUnit(unit.Bytes),
		instrument.WithDescription("Total number of data bytes sent"),
	)
	if err != nil {
		return nil, err
	}

	return &ServerMetrics{
		BlocksSent: blocksSent,
		DataSent:   dataSent,
	}, nil
}

type SessionMetrics struct {
	DiscoverySuccess syncint64.Counter
	DiscoveryFailure syncint64.Counter
	TimeToFirstBlock syncint64.Histogram
}

func NewSessionMetrics(meterProvider metric.MeterProvider) (*SessionMetrics, error) {
	m := meterProvider.Meter(SessionScope.Name, metric.WithInstrumentationVersion(SessionScope.Version), metric.WithSchemaURL(SessionScope.SchemaURL))

	discoverySuccess, err := m.SyncInt64().Counter(
		"discovery_success",
		instrument.WithUnit(unit.Dimensionless),
		instrument.WithDescription("Total number of times discovery succeeded"),
	)
	if err != nil {
		return nil, err
	}

	discoveryFailure, err := m.SyncInt64().Counter(
		"discovery_failure",
		instrument.WithUnit(unit.Dimensionless),
		instrument.WithDescription("Total number of times discovery failed"),
	)
	if err != nil {
		return nil, err
	}

	timeToFirstBlock, err := m.SyncInt64().Histogram(
		"time_to_first_block",
		instrument.WithUnit(unit.Milliseconds),
		instrument.WithDescription("Time to first block"),
	)
	if err != nil {
		return nil, err
	}

	return &SessionMetrics{
		DiscoverySuccess: discoverySuccess,
		DiscoveryFailure: discoveryFailure,
		TimeToFirstBlock: timeToFirstBlock,
	}, nil
}

func (m *ServerMetrics) RegisterCallback(cb func(context.Context)) error {
	instruments := []instrument.Asynchronous{
		m.BlocksSent,
		m.DataSent,
	}
	return m.m.RegisterCallback(instruments, cb)
}

var (
	// the 1<<18+15 is to observe old file chunks that are 1<<18 + 14 in size
	metricsBuckets = []float64{1 << 6, 1 << 10, 1 << 14, 1 << 18, 1<<18 + 15, 1 << 22}

	timeMetricsBuckets = []float64{1, 10, 30, 60, 90, 120, 600}
)

// func DupHist(ctx context.Context) metrics.Histogram {
// 	return metrics.NewCtx(ctx, "recv_dup_blocks_bytes", "Summary of duplicate data blocks recived").Histogram(metricsBuckets)
// }

// func AllHist(ctx context.Context) metrics.Histogram {
// 	return metrics.NewCtx(ctx, "recv_all_blocks_bytes", "Summary of all data blocks recived").Histogram(metricsBuckets)
// }

// func SentHist(ctx context.Context) metrics.Histogram {
// 	return metrics.NewCtx(ctx, "sent_all_blocks_bytes", "Histogram of blocks sent by this bitswap").Histogram(metricsBuckets)
// }

// func SendTimeHist(ctx context.Context) metrics.Histogram {
// 	return metrics.NewCtx(ctx, "send_times", "Histogram of how long it takes to send messages in this bitswap").Histogram(timeMetricsBuckets)
// }

// func PendingEngineGauge(ctx context.Context) metrics.Gauge {
// 	return metrics.NewCtx(ctx, "pending_tasks", "Total number of pending tasks").Gauge()
// }

// func ActiveEngineGauge(ctx context.Context) metrics.Gauge {
// 	return metrics.NewCtx(ctx, "active_tasks", "Total number of active tasks").Gauge()
// }

// func PendingBlocksGauge(ctx context.Context) metrics.Gauge {
// 	return metrics.NewCtx(ctx, "pending_block_tasks", "Total number of pending blockstore tasks").Gauge()
// }

// func ActiveBlocksGauge(ctx context.Context) metrics.Gauge {
// 	return metrics.NewCtx(ctx, "active_block_tasks", "Total number of active blockstore tasks").Gauge()
// }
