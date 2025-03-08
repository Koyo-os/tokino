package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/koyo-os/tokino/src/db"
	"github.com/koyo-os/tokino/src/models"
	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/libp2p/go-libp2p/p2p/discovery/mdns"
)

type BlockChain struct{
	blocks []*models.Block
	data *db.Data
	logger *Logger
	host host.Host
}

type discoveryNotifee struct {
	host host.Host
	logger *Logger
}


func (n *discoveryNotifee) HandlePeerFound(pi peer.AddrInfo) {
	n.logger.Info().Any("id", pi.ID).Msg("new peer")
	if err := n.host.Connect(context.Background(), pi); err != nil {
		n.logger.Error().Any("id", pi.ID).Msg("cant conntect")
	} else {
		n.logger.Info().Any("id", pi.ID).Msg("success connect!")
	}
}

const DiscoveryServiceTag = "p2p-network"

func NewChain() (*BlockChain, error) {
	logger := NewLogger()
	logger.Info().Msg("starting tokino")

	db, err := db.New()
	if err != nil{
		logger.Error().Err(err)
		return nil,err
	}

	host,err := libp2p.New(
		libp2p.ListenAddrStrings("/ip4/0.0.0.0/tcp/0"),
	)

	if err != nil{
		logger.Error().Err(err)
		return nil, err
	}

	logger.Info().Any("addr", host.Addrs()).Any("perr id", host.ID()).Msg("new peer!")
	return &BlockChain{
		host: host,
		logger: &logger,	
		data: db,
	}, nil
}

func (b *BlockChain) Start() {
	discovery := mdns.NewMdnsService(b.host, DiscoveryServiceTag, &discoveryNotifee{
		host: b.host,
		logger: b.logger,
	})

	if err := discovery.Start();err != nil{
		b.logger.Error().Err(err)
		return
	}
	defer discovery.Close()

	peersReached := make(chan struct{})
	go trackPeers(b.host, peersReached)

	<-peersReached
	fmt.Println("i get 3 peers")

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	<-ch
	b.logger.Info().Msg("stop...")
}

func (b *BlockChain) RouteCommand(cmd string) error {
	comand := strings.Split(cmd, " ")
	if comand[0] == "create" {
		if len(comand) < 4 {
			fmt.Println("Usage: send to from price")
		}
	}
}


func (b *BlockChain) AddBlock(data string, prev string, index int) (bool, error) {
	b.logger.Info().Msg("starting check new block")

	last,err := b.data.GetLast()
	if err != nil{
		return false, err
	}

	if last.SelfHash == prev {
		block := models.NewBlock(index, data, prev)
		b.data.Add(block)
		b.blocks = append(b.blocks, block)
		b.logger.Info().Msg("check successfully!")
		return true, nil
	}

	b.logger.Info().Msg("check unsuccessfully")

	return false, nil
}
