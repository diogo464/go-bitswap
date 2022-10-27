package telemetry

import (
	"context"

	"github.com/diogo464/telemetry"
	"github.com/ipfs/go-bitswap"
	logging "github.com/ipfs/go-log"
	"go.opentelemetry.io/otel/metric/instrument"
	"go.opentelemetry.io/otel/metric/instrument/asyncint64"
	"go.opentelemetry.io/otel/metric/unit"
)

var log = logging.Logger("bitswap/telemetry")

func Start(bs *bitswap.Bitswap) error {
	t := telemetry.GetGlobalTelemetry()
	return registerMetrics(t, bs)
}

func registerMetrics(t telemetry.Telemetry, bs *bitswap.Bitswap) error {
	var (
		err error

		blocksReceived    asyncint64.Counter
		dataReceived      asyncint64.Counter
		dupBlocksReceived asyncint64.Counter
		dupDataReceived   asyncint64.Counter
		messagesReceived  asyncint64.Counter
		discoverySuccess  asyncint64.Counter
		discoveryFailure  asyncint64.Counter
		blocksSent        asyncint64.Counter
		dataSent          asyncint64.Counter
	)

	m := t.Meter("libp2p.io/bitswap")

	if blocksReceived, err = m.AsyncInt64().Counter(
		"blocks_received",
		instrument.WithUnit(unit.Dimensionless),
		instrument.WithDescription("Total number of received blocks"),
	); err != nil {
		return err
	}

	if dataReceived, err = m.AsyncInt64().Counter(
		"data_received",
		instrument.WithUnit(unit.Bytes),
		instrument.WithDescription("Total number of data bytes received"),
	); err != nil {
		return err
	}

	if dupBlocksReceived, err = m.AsyncInt64().Counter(
		"dup_blocks_received",
		instrument.WithUnit(unit.Dimensionless),
		instrument.WithDescription("Total number of duplicate blocks received"),
	); err != nil {
		return err
	}

	if dupDataReceived, err = m.AsyncInt64().Counter(
		"dup_data_received",
		instrument.WithUnit(unit.Bytes),
		instrument.WithDescription("Total number of duplicate data bytes received"),
	); err != nil {
		return err
	}

	if messagesReceived, err = m.AsyncInt64().Counter(
		"messages_received",
		instrument.WithUnit(unit.Dimensionless),
		instrument.WithDescription("Total number of messages received"),
	); err != nil {
		return err
	}

	if discoverySuccess, err = m.AsyncInt64().Counter(
		"discovery_success",
		instrument.WithUnit(unit.Dimensionless),
		instrument.WithDescription("Total number times discovery succeeded"),
	); err != nil {
		return err
	}

	if discoveryFailure, err = m.AsyncInt64().Counter(
		"discovery_failure",
		instrument.WithUnit(unit.Dimensionless),
		instrument.WithDescription("Total number times discovery failed"),
	); err != nil {
		return err
	}

	if blocksSent, err = m.AsyncInt64().Counter(
		"blocks_sent",
		instrument.WithUnit(unit.Dimensionless),
		instrument.WithDescription("Total number of blocks sent"),
	); err != nil {
		return err
	}

	if dataSent, err = m.AsyncInt64().Counter(
		"data_sent",
		instrument.WithUnit(unit.Bytes),
		instrument.WithDescription("Total number of data bytes sent"),
	); err != nil {
		return err
	}

	m.RegisterCallback([]instrument.Asynchronous{
		blocksReceived,
		dataReceived,
		dupBlocksReceived,
		dupDataReceived,
		messagesReceived,
		discoverySuccess,
		discoveryFailure,
		blocksSent,
		dataSent,
	}, func(ctx context.Context) {
		stat, err := bs.Stat()
		if err != nil {
			log.Errorf("bitswap.Stat() failed", "error", err)
			return
		}

		blocksReceived.Observe(ctx, int64(stat.BlocksReceived))
		dataReceived.Observe(ctx, int64(stat.DataReceived))
		dupBlocksReceived.Observe(ctx, int64(stat.DupBlksReceived))
		dupDataReceived.Observe(ctx, int64(stat.DupDataReceived))
		messagesReceived.Observe(ctx, int64(stat.MessagesReceived))
		discoverySuccess.Observe(ctx, int64(stat.DiscoverySuccess))
		discoveryFailure.Observe(ctx, int64(stat.DiscoveryFailure))
		blocksSent.Observe(ctx, int64(stat.BlocksSent))
		dataSent.Observe(ctx, int64(stat.DataSent))

	})

	return nil
}
