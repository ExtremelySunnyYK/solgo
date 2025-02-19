package utils

import (
	"fmt"
	"math/big"
)

type Network string
type NetworkID uint64

func (n NetworkID) String() string {
	return n.ToBig().String()
}

func (n NetworkID) ToNetwork() Network {
	switch n {
	case EthereumNetworkID:
		return Ethereum
	case BscNetworkID:
		return Bsc
	case PolygonNetworkID:
		return Polygon
	case AvalancheNetworkID:
		return Avalanche
	case FantomNetworkID:
		return Fantom
	case ArbitrumNetworkID:
		return Arbitrum
	case OptimismNetworkID:
		return Optimism
	default:
		return Ethereum
	}
}

func (n NetworkID) Uint64() uint64 {
	return uint64(n)
}

const (
	AnvilNetwork Network = "anvil"
	Ethereum     Network = "ethereum"
	Bsc          Network = "bsc"
	Polygon      Network = "polygon"
	Avalanche    Network = "avalanche"
	Fantom       Network = "fantom"
	Arbitrum     Network = "arbitrum"
	Optimism     Network = "optimism"

	// Mainnets
	EthereumNetworkID  NetworkID = 1
	BscNetworkID       NetworkID = 56
	PolygonNetworkID   NetworkID = 137
	AvalancheNetworkID NetworkID = 43114
	FantomNetworkID    NetworkID = 250
	ArbitrumNetworkID  NetworkID = 42161
	OptimismNetworkID  NetworkID = 10

	// Special Mainnets
	AnvilMainnetNetworkID NetworkID = 1

	// Testnets
	RopstenNetworkID    NetworkID = 3
	RinkebyNetworkID    NetworkID = 4
	GoerliNetworkID     NetworkID = 5
	KovanNetworkID      NetworkID = 42
	BscTestnetNetworkID NetworkID = 97
	MumbaiNetworkID     NetworkID = 80001
	FujiNetworkID       NetworkID = 43113
	FantomTestNetworkID NetworkID = 4002
	ArbitrumRinkebyID   NetworkID = 421611
	OptimismKovanID     NetworkID = 69

	// Localnets
	AnvilNetworkID NetworkID = 1337
)

func (n Network) String() string {
	return string(n)
}

func (n NetworkID) ToBig() *big.Int {
	return new(big.Int).SetUint64(uint64(n))
}

func GetNetworkID(network Network) NetworkID {
	switch network {
	case Ethereum:
		return EthereumNetworkID
	case Bsc:
		return BscNetworkID
	case Polygon:
		return PolygonNetworkID
	case Avalanche:
		return AvalancheNetworkID
	case Fantom:
		return FantomNetworkID
	case Arbitrum:
		return ArbitrumNetworkID
	case Optimism:
		return OptimismNetworkID
	case AnvilNetwork:
		return AnvilMainnetNetworkID
	default:
		return 0
	}
}

func GetNetworkFromID(id NetworkID) (Network, error) {
	switch id {
	case EthereumNetworkID:
		return Ethereum, nil
	case BscNetworkID:
		return Bsc, nil
	case PolygonNetworkID:
		return Polygon, nil
	case AvalancheNetworkID:
		return Avalanche, nil
	case FantomNetworkID:
		return Fantom, nil
	case ArbitrumNetworkID:
		return Arbitrum, nil
	case OptimismNetworkID:
		return Optimism, nil
	default:
		return "", fmt.Errorf("unknown network ID '%d' provided", id)
	}
}

func GetNetworkFromString(network string) (Network, error) {
	switch network {
	case "ethereum":
		return Ethereum, nil
	case "bsc":
		return Bsc, nil
	case "polygon":
		return Polygon, nil
	case "avalanche":
		return Avalanche, nil
	case "fantom":
		return Fantom, nil
	case "arbitrum":
		return Arbitrum, nil
	case "optimism":
		return Optimism, nil
	case "anvil":
		return AnvilNetwork, nil
	default:
		return "", fmt.Errorf("unknown network '%s' provided", network)
	}
}
