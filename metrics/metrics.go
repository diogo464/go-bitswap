package metrics

import (
	"context"

	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/metric/instrument"
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
		Kind:  sdk_metric.InstrumentKindHistogram,
		Unit:  unit.Unit("s"),
		Scope: ServerScope,
	}, streamTimeHistogram)

	ViewClientBytesHistograms = sdk_metric.NewView(sdk_metric.Instrument{
		Kind:  sdk_metric.InstrumentKindHistogram,
		Unit:  unit.Bytes,
		Scope: ClientScope,
	}, streamBytesHistogram)

	ViewServerBytesHistograms = sdk_metric.NewView(sdk_metric.Instrument{
		Kind:  sdk_metric.InstrumentKindHistogram,
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
	DataReceived    instrument.Int64Histogram
	DupDataReceived instrument.Int64Histogram

	// Asynchronous
	BlocksReceived    instrument.Int64ObservableCounter
	DupBlocksReceived instrument.Int64ObservableCounter
	MessagesReceived  instrument.Int64ObservableCounter
}

func NewClientMetrics(meterProvider metric.MeterProvider) (*ClientMetrics, error) {
	m := meterProvider.Meter(ClientScope.Name, metric.WithInstrumentationVersion(ClientScope.Version), metric.WithSchemaURL(ClientScope.SchemaURL))

	blocksReceived, err := m.Int64ObservableCounter(
		"bitswap.client.blocks_received",
		instrument.WithUnit(unit.Dimensionless),
		instrument.WithDescription("Total number of received blocks"),
	)
	if err != nil {
		return nil, err
	}

	dataReceived, err := m.Int64Histogram(
		"bitswap.client.data_received",
		instrument.WithUnit(unit.Bytes),
		instrument.WithDescription("Number of data bytes received"),
	)
	if err != nil {
		return nil, err
	}

	dupBlocksReceived, err := m.Int64ObservableCounter(
		"bitswap.client.dup_blocks_received",
		instrument.WithUnit(unit.Dimensionless),
		instrument.WithDescription("Total number of duplicate blocks received"),
	)
	if err != nil {
		return nil, err
	}

	dupDataReceived, err := m.Int64Histogram(
		"bitswap.client.dup_data_received",
		instrument.WithUnit(unit.Bytes),
		instrument.WithDescription("Number of duplicate data bytes received"),
	)
	if err != nil {
		return nil, err
	}

	messagesReceived, err := m.Int64ObservableCounter(
		"bitswap.client.messages_received",
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

func (m *ClientMetrics) RegisterCallback(cb func(context.Context, metric.Observer) error) error {
	instruments := []instrument.Asynchronous{
		m.BlocksReceived,
		m.DupBlocksReceived,
		m.MessagesReceived,
	}
	_, err := m.m.RegisterCallback(cb, instruments...)
	return err
}

type ServerMetrics struct {
	BlocksSent      instrument.Int64Counter
	DataSent        instrument.Int64Histogram
	MessageDataSent instrument.Int64Histogram
	SendTime        instrument.Float64Histogram
}

func NewServerMetrics(meterProvider metric.MeterProvider) (*ServerMetrics, error) {
	m := meterProvider.Meter(ServerScope.Name, metric.WithInstrumentationVersion(ServerScope.Version), metric.WithSchemaURL(ServerScope.SchemaURL))

	blocksSent, err := m.Int64Counter(
		"bitswap.server.blocks_sent",
		instrument.WithUnit(unit.Dimensionless),
		instrument.WithDescription("Number of sent blocks"),
	)
	if err != nil {
		return nil, err
	}

	dataSent, err := m.Int64Histogram(
		"bitswap.server.data_sent",
		instrument.WithUnit(unit.Bytes),
		instrument.WithDescription("Number of data bytes sent"),
	)
	if err != nil {
		return nil, err
	}

	messageDataSent, err := m.Int64Histogram(
		"bitswap.server.message_data_sent",
		instrument.WithUnit(unit.Bytes),
		instrument.WithDescription("Number of payload bytes sent in messages"),
	)
	if err != nil {
		return nil, err
	}

	sendTime, err := m.Float64Histogram(
		"bitswap.server.send_time",
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
	DiscoverySuccess instrument.Int64Counter
	DiscoveryFailure instrument.Int64Counter
	TimeToFirstBlock instrument.Int64Histogram
}

func NewSessionMetrics(meterProvider metric.MeterProvider) (*SessionMetrics, error) {
	m := meterProvider.Meter(SessionScope.Name, metric.WithInstrumentationVersion(SessionScope.Version), metric.WithSchemaURL(SessionScope.SchemaURL))

	discoverySuccess, err := m.Int64Counter(
		"bitswap.session.discovery_success",
		instrument.WithUnit(unit.Dimensionless),
		instrument.WithDescription("Total number of times discovery succeeded"),
	)
	if err != nil {
		return nil, err
	}

	discoveryFailure, err := m.Int64Counter(
		"bitswap.session.discovery_failure",
		instrument.WithUnit(unit.Dimensionless),
		instrument.WithDescription("Total number of times discovery failed"),
	)
	if err != nil {
		return nil, err
	}

	timeToFirstBlock, err := m.Int64Histogram(
		"bitswap.session.time_to_first_block",
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
	BlockStorePending instrument.Int64UpDownCounter
	BlockStoreActive  instrument.Int64UpDownCounter

	// Asynchronous
	PeerQueueActive  instrument.Int64ObservableGauge
	PeerQueuePending instrument.Int64ObservableGauge
}

func NewDecisionMetrics(meterProvider metric.MeterProvider) (*DecisionMetrics, error) {
	m := meterProvider.Meter(DecisionScope.Name, metric.WithInstrumentationVersion(DecisionScope.Version), metric.WithSchemaURL(DecisionScope.SchemaURL))

	blockStorePending, err := m.Int64UpDownCounter(
		"bitswap.decision.blockstore_pending",
		instrument.WithUnit(unit.Dimensionless),
		instrument.WithDescription("Number of pending blockstore tasks"),
	)
	if err != nil {
		return nil, err
	}

	blockStoreActive, err := m.Int64UpDownCounter(
		"bitswap.decision.blockstore_active",
		instrument.WithUnit(unit.Dimensionless),
		instrument.WithDescription("Number of active blockstore tasks"),
	)
	if err != nil {
		return nil, err
	}

	peerQueueActive, err := m.Int64ObservableGauge(
		"bitswap.decision.peer_queue_active",
		instrument.WithUnit(unit.Dimensionless),
		instrument.WithDescription("Number of active peer queue tasks"),
	)
	if err != nil {
		return nil, err
	}

	peerQueuePending, err := m.Int64ObservableGauge(
		"bitswap.decision.peer_queue_pending",
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

func (m *DecisionMetrics) RegisterCallback(cb func(context.Context, metric.Observer) error) error {
	insts := []instrument.Asynchronous{
		m.PeerQueueActive,
		m.PeerQueuePending,
	}
	_, err := m.m.RegisterCallback(cb, insts...)
	return err
}

// func PendingBlocksGauge(ctx context.Context) metrics.Gauge {
// 	return metrics.NewCtx(ctx, "pending_block_tasks", "Total number of pending blockstore tasks").Gauge()
// }

// func ActiveBlocksGauge(ctx context.Context) metrics.Gauge {
// 	return metrics.NewCtx(ctx, "active_block_tasks", "Total number of active blockstore tasks").Gauge()
// }
