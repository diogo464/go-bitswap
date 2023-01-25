package metrics

import (
	"context"

	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/metric/instrument"
	"go.opentelemetry.io/otel/metric/instrument/asyncint64"
	"go.opentelemetry.io/otel/metric/instrument/syncfloat64"
	"go.opentelemetry.io/otel/metric/instrument/syncint64"
	"go.opentelemetry.io/otel/metric/unit"
	"go.opentelemetry.io/otel/sdk/instrumentation"
	sdk_metric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/metric/aggregation"
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

	DecisionScope = instrumentation.Scope{
		Name:    "libp2p.io/bitswap/decision",
		Version: Version,
	}

	streamBytesHistogram = sdk_metric.Stream{
		Aggregation: aggregation.ExplicitBucketHistogram{
			Boundaries: []float64{1 << 6, 1 << 10, 1 << 14, 1 << 18, 1<<18 + 15, 1 << 22},
		},
	}

	streamTimeHistogram = sdk_metric.Stream{
		Aggregation: aggregation.ExplicitBucketHistogram{
			Boundaries: []float64{1, 10, 30, 60, 90, 120, 600},
		},
	}

	ViewTimeMetrics = sdk_metric.NewView(sdk_metric.Instrument{
		Kind:  sdk_metric.InstrumentKindSyncHistogram,
		Unit:  unit.Unit("s"),
		Scope: ServerScope,
	}, streamTimeHistogram)

	ViewClientBytesHistograms = sdk_metric.NewView(sdk_metric.Instrument{
		Kind:  sdk_metric.InstrumentKindSyncHistogram,
		Unit:  unit.Bytes,
		Scope: ClientScope,
	}, streamBytesHistogram)

	ViewServerBytesHistograms = sdk_metric.NewView(sdk_metric.Instrument{
		Kind:  sdk_metric.InstrumentKindSyncHistogram,
		Unit:  unit.Bytes,
		Scope: ServerScope,
	}, streamBytesHistogram)

	Views = []sdk_metric.View{
		ViewTimeMetrics,
		ViewClientBytesHistograms,
		ViewServerBytesHistograms,
	}
)

type ClientMetrics struct {
	m metric.Meter

	// Synchronous
	DataReceived    syncint64.Histogram
	DupDataReceived syncint64.Histogram

	// Asynchronous
	BlocksReceived    asyncint64.Counter
	DupBlocksReceived asyncint64.Counter
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

	dataReceived, err := m.SyncInt64().Histogram(
		"data_received",
		instrument.WithUnit(unit.Bytes),
		instrument.WithDescription("Number of data bytes received"),
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

	dupDataReceived, err := m.SyncInt64().Histogram(
		"dup_data_received",
		instrument.WithUnit(unit.Bytes),
		instrument.WithDescription("Number of duplicate data bytes received"),
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
		m: m,

		DataReceived:    dataReceived,
		DupDataReceived: dupDataReceived,

		BlocksReceived:    blocksReceived,
		DupBlocksReceived: dupBlocksReceived,
		MessagesReceived:  messagesReceived,
	}, nil
}

func (m *ClientMetrics) RegisterCallback(cb func(context.Context)) error {
	instruments := []instrument.Asynchronous{
		m.BlocksReceived,
		m.DupBlocksReceived,
		m.MessagesReceived,
	}
	return m.m.RegisterCallback(instruments, cb)
}

type ServerMetrics struct {
	BlocksSent      syncint64.Counter
	DataSent        syncint64.Histogram
	MessageDataSent syncint64.Histogram
	SendTime        syncfloat64.Histogram
}

func NewServerMetrics(meterProvider metric.MeterProvider) (*ServerMetrics, error) {
	m := meterProvider.Meter(ServerScope.Name, metric.WithInstrumentationVersion(ServerScope.Version), metric.WithSchemaURL(ServerScope.SchemaURL))

	blocksSent, err := m.SyncInt64().Counter(
		"blocks_sent",
		instrument.WithUnit(unit.Dimensionless),
		instrument.WithDescription("Number of sent blocks"),
	)
	if err != nil {
		return nil, err
	}

	dataSent, err := m.SyncInt64().Histogram(
		"data_sent",
		instrument.WithUnit(unit.Bytes),
		instrument.WithDescription("Number of data bytes sent"),
	)
	if err != nil {
		return nil, err
	}

	messageDataSent, err := m.SyncInt64().Histogram(
		"message_data_sent",
		instrument.WithUnit(unit.Bytes),
		instrument.WithDescription("Number of payload bytes sent in messages"),
	)
	if err != nil {
		return nil, err
	}

	sendTime, err := m.SyncFloat64().Histogram(
		"send_time",
		instrument.WithUnit(unit.Unit("s")),
		instrument.WithDescription("Time to send a message"),
	)
	if err != nil {
		return nil, err
	}

	return &ServerMetrics{
		BlocksSent:      blocksSent,
		DataSent:        dataSent,
		MessageDataSent: messageDataSent,
		SendTime:        sendTime,
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

type DecisionMetrics struct {
	m metric.Meter

	// Synchronous
	BlockStorePending syncint64.UpDownCounter
	BlockStoreActive  syncint64.UpDownCounter

	// Asynchronous
	PeerQueueActive  asyncint64.Gauge
	PeerQueuePending asyncint64.Gauge
}

func NewDecisionMetrics(meterProvider metric.MeterProvider) (*DecisionMetrics, error) {
	m := meterProvider.Meter(DecisionScope.Name, metric.WithInstrumentationVersion(DecisionScope.Version), metric.WithSchemaURL(DecisionScope.SchemaURL))

	blockStorePending, err := m.SyncInt64().UpDownCounter(
		"blockstore_pending",
		instrument.WithUnit(unit.Dimensionless),
		instrument.WithDescription("Number of pending blockstore tasks"),
	)
	if err != nil {
		return nil, err
	}

	blockStoreActive, err := m.SyncInt64().UpDownCounter(
		"blockstore_active",
		instrument.WithUnit(unit.Dimensionless),
		instrument.WithDescription("Number of active blockstore tasks"),
	)
	if err != nil {
		return nil, err
	}

	peerQueueActive, err := m.AsyncInt64().Gauge(
		"peer_queue_active",
		instrument.WithUnit(unit.Dimensionless),
		instrument.WithDescription("Number of active peer queue tasks"),
	)
	if err != nil {
		return nil, err
	}

	peerQueuePending, err := m.AsyncInt64().Gauge(
		"peer_queue_pending",
		instrument.WithUnit(unit.Dimensionless),
		instrument.WithDescription("Number of pending peer queue tasks"),
	)
	if err != nil {
		return nil, err
	}

	return &DecisionMetrics{
		m: m,

		BlockStorePending: blockStorePending,
		BlockStoreActive:  blockStoreActive,

		PeerQueueActive:  peerQueueActive,
		PeerQueuePending: peerQueuePending,
	}, nil
}

func (m *DecisionMetrics) RegisterCallback(cb func(context.Context)) error {
	insts := []instrument.Asynchronous{
		m.PeerQueueActive,
		m.PeerQueuePending,
	}
	return m.m.RegisterCallback(insts, cb)
}

// func PendingBlocksGauge(ctx context.Context) metrics.Gauge {
// 	return metrics.NewCtx(ctx, "pending_block_tasks", "Total number of pending blockstore tasks").Gauge()
// }

// func ActiveBlocksGauge(ctx context.Context) metrics.Gauge {
// 	return metrics.NewCtx(ctx, "active_block_tasks", "Total number of active blockstore tasks").Gauge()
// }
