package state

import "github.com/AlexGustafsson/grapevine/internal/webpush"

type ConfigFile struct {
	Topics map[string]Topic `json:"topics"`
}

type Topic struct {
	Name      string `json:"name"`
	ShortName string `json:"shortName"`
}

type SecretsFile struct {
	Clients map[string]ClientSecrets `json:"clients"`
}

type ClientSecrets struct {
	PrivateKey string `json:"privateKey"`
}

type SubscriptionsFile struct {
	Topics map[string]map[string]webpush.Subscription `json:"topics"`
}
