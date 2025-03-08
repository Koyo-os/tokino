package main

import (
	"fmt"
	"sync"

	"github.com/koyo-os/tokino/src/models"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/network"
	"github.com/libp2p/go-libp2p/core/peer"
)

func trackPeers(host host.Host, peersReached chan struct{}) {
	var (
		mu    sync.Mutex
		peers = make(map[peer.ID]struct{})
	)

	host.Network().Notify(&network.NotifyBundle{
		ConnectedF: func(net network.Network, conn network.Conn) {
			mu.Lock()
			defer mu.Unlock()

			peers[conn.RemotePeer()] = struct{}{}
			fmt.Printf("Подключен новый пир: %s (всего: %d)\n", conn.RemotePeer(), len(peers))
			if len(peers) >= 3 {
				select {
				case peersReached <- struct{}{}:
				default:
				}
			}
		},
		DisconnectedF: func(net network.Network, conn network.Conn) {
			mu.Lock()
			defer mu.Unlock()

			delete(peers, conn.RemotePeer())
			fmt.Printf("Пир отключен: %s (всего: %d)\n", conn.RemotePeer(), len(peers))
		},
	})
}

func GetBalance(address string, chain []models.Block) int {
	balance := 0
	for _, block := range chain {
			for _, out := range block.Transaction.Outputs {
				if out.PubKey == address {
					balance += out.Value
				}
			}
	}
	return balance
}