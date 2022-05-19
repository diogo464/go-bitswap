module github.com/ipfs/go-bitswap

replace github.com/libp2p/go-libp2p => ../go-libp2p/

replace github.com/diogo464/telemetry => ../../

require (
	github.com/benbjohnson/clock v1.3.0
	github.com/cskr/pubsub v1.0.2
	github.com/diogo464/telemetry v0.0.0-00010101000000-000000000000
	github.com/gogo/protobuf v1.3.2
	github.com/google/uuid v1.3.0
	github.com/ipfs/go-block-format v0.0.3
	github.com/ipfs/go-cid v0.1.0
	github.com/ipfs/go-datastore v0.5.1
	github.com/ipfs/go-detect-race v0.0.1
	github.com/ipfs/go-ipfs-blockstore v1.2.0
	github.com/ipfs/go-ipfs-blocksutil v0.0.1
	github.com/ipfs/go-ipfs-delay v0.0.1
	github.com/ipfs/go-ipfs-exchange-interface v0.1.0
	github.com/ipfs/go-ipfs-routing v0.2.1
	github.com/ipfs/go-ipfs-util v0.0.2
	github.com/ipfs/go-ipld-format v0.4.0
	github.com/ipfs/go-log v1.0.5
	github.com/ipfs/go-metrics-interface v0.0.1
	github.com/ipfs/go-peertaskqueue v0.7.1
	github.com/jbenet/goprocess v0.1.4
	github.com/libp2p/go-buffer-pool v0.0.2
	github.com/libp2p/go-libp2p v0.19.1
	github.com/libp2p/go-libp2p-core v0.15.1
	github.com/libp2p/go-libp2p-loggables v0.1.0
	github.com/libp2p/go-libp2p-netutil v0.1.0
	github.com/libp2p/go-libp2p-testing v0.9.2
	github.com/libp2p/go-msgio v0.2.0
	github.com/multiformats/go-multiaddr v0.5.0
	github.com/multiformats/go-multistream v0.3.0
	github.com/stretchr/testify v1.7.1
	go.uber.org/zap v1.21.0
)

go 1.16
