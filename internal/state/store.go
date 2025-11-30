package state

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"
	"log/slog"
	"maps"
	"os"
	"path/filepath"
	"slices"
	"sync"

	"github.com/AlexGustafsson/grapevine/internal/webpush"
)

var (
	ErrTopicNotFound        = errors.New("topic not found")
	ErrSubscriptionNotFound = errors.New("subscription not found")
)

type Client struct {
	topic      string
	name       string
	shortName  string
	privateKey *ecdsa.PrivateKey
}

func (c *Client) Topic() string {
	return c.topic
}

func (c *Client) Name() string {
	return c.name
}

func (c *Client) ShortName() string {
	return c.shortName
}

func (c *Client) Subject() string {
	return "https://example.com/" + c.topic // TODO
}

func (c *Client) WebPushClient() webpush.Client {
	keyExchangeKey, err := c.privateKey.ECDH()
	if err != nil {
		panic(err)
	}

	return webpush.NewClient(c.Subject(), c.privateKey, keyExchangeKey)
}

type Store struct {
	mutex         sync.RWMutex
	basePath      string
	clients       map[string]Client
	subscriptions map[string]map[string]webpush.Subscription
}

func (s *Store) BasePath() string {
	return s.basePath
}

func Migrate(basePath string) error {
	var config ConfigFile
	err := readJSON(filepath.Join(basePath, "config.json"), &config)
	if err != nil {
		return err
	}

	var secrets SecretsFile
	err = readJSON(filepath.Join(basePath, "secrets.json"), &secrets)
	if errors.Is(err, os.ErrNotExist) {
		secrets.Clients = make(map[string]ClientSecrets)
	} else if err != nil {
		return err
	}

	var subscriptions SubscriptionsFile
	err = readJSON(filepath.Join(basePath, "subscriptions.json"), &subscriptions)
	if errors.Is(err, os.ErrNotExist) {
		subscriptions.Topics = make(map[string]map[string]webpush.Subscription)
	} else if err != nil {
		return err
	}

	// Remove secrets for topics that don't exist
	for clientTopicName := range secrets.Clients {
		found := false
		for topicName := range config.Topics {
			if clientTopicName == topicName {
				found = true
				break
			}
		}

		if !found {
			slog.Info("Identified removed topic, removing secrets", slog.String("topic", clientTopicName))
			delete(secrets.Clients, clientTopicName)
		}
	}

	// Generate secrets for new topics
	for topicName := range config.Topics {
		_, ok := secrets.Clients[topicName]
		if !ok {
			slog.Info("Identified new topic, generating secrets", slog.String("topic", topicName))

			privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
			if err != nil {
				return err
			}

			privateKeyBytes, err := privateKey.Bytes()
			if err != nil {
				return err
			}

			privateKeyPEM := pem.EncodeToMemory(&pem.Block{
				Type:  "PRIVATE KEY",
				Bytes: privateKeyBytes,
			})

			secrets.Clients[topicName] = ClientSecrets{
				PrivateKey: string(privateKeyPEM),
			}
		}
	}

	// Add subscriptions for new topics
	for topicName := range config.Topics {
		_, ok := subscriptions.Topics[topicName]
		if !ok {
			subscriptions.Topics[topicName] = make(map[string]webpush.Subscription)
		}
	}

	err = writeJSON(filepath.Join(basePath, "secrets.json"), &secrets)
	if err != nil {
		return err
	}

	err = writeJSON(filepath.Join(basePath, "subscriptions.json"), &subscriptions)
	if err != nil {
		return err
	}

	return nil
}

func Load(basePath string) (*Store, error) {
	var config ConfigFile
	err := readJSON(filepath.Join(basePath, "config.json"), &config)
	if err != nil {
		return nil, err
	}

	var secrets SecretsFile
	err = readJSON(filepath.Join(basePath, "secrets.json"), &secrets)
	if err != nil {
		return nil, err
	}

	var subscriptions SubscriptionsFile
	err = readJSON(filepath.Join(basePath, "subscriptions.json"), &subscriptions)
	if err != nil {
		return nil, err
	}

	clients := make(map[string]Client)
	for topicName, topic := range config.Topics {
		secrets, ok := secrets.Clients[topicName]
		if !ok {
			return nil, fmt.Errorf("invalid secrets file")
		}

		block, rest := pem.Decode([]byte(secrets.PrivateKey))
		if len(rest) > 0 {
			return nil, fmt.Errorf("invalid secret in secrets file")
		}

		privateKey, err := ecdsa.ParseRawPrivateKey(elliptic.P256(), block.Bytes)
		if err != nil {
			return nil, fmt.Errorf("invalid secret in secrets file: %w", err)
		}

		clients[topicName] = Client{
			topic:      topicName,
			name:       topic.Name,
			shortName:  topic.ShortName,
			privateKey: privateKey,
		}
	}

	return &Store{
		basePath:      basePath,
		clients:       clients,
		subscriptions: subscriptions.Topics,
	}, nil
}

func (s *Store) Client(topic string) (Client, bool) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	client, ok := s.clients[topic]
	return client, ok
}

func (s *Store) AddSubscription(topic string, id string, subscription webpush.Subscription) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	subscriptions, ok := s.subscriptions[topic]
	if !ok {
		return ErrTopicNotFound
	}

	subscriptions[id] = subscription
	return nil
}

func (s *Store) GetSubscription(topic string, id string) (webpush.Subscription, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	subscriptions, ok := s.subscriptions[topic]
	if !ok {
		return webpush.Subscription{}, ErrTopicNotFound
	}

	subscription, ok := subscriptions[id]
	if !ok {
		return webpush.Subscription{}, ErrSubscriptionNotFound
	}

	return subscription, nil
}

func (s *Store) DeleteSubscription(topic string, id string) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	subscriptions, ok := s.subscriptions[topic]
	if !ok {
		return ErrTopicNotFound
	}

	_, ok = subscriptions[id]
	if !ok {
		return ErrSubscriptionNotFound
	}

	delete(subscriptions, id)
	return nil
}

func (s *Store) GetSubscriptions(topic string) ([]webpush.Subscription, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	subscriptions, ok := s.subscriptions[topic]
	if !ok {
		return nil, ErrTopicNotFound
	}

	return slices.Collect(maps.Values(subscriptions)), nil
}

func (s *Store) Save(basePath string) error {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	secrets := SecretsFile{
		Clients: make(map[string]ClientSecrets),
	}

	for topic, client := range s.clients {
		privateKeyBytes, err := client.privateKey.Bytes()
		if err != nil {
			return err
		}

		privateKey := pem.EncodeToMemory(&pem.Block{
			Type:  "PRIVATE KEY",
			Bytes: privateKeyBytes,
		})

		secrets.Clients[topic] = ClientSecrets{
			PrivateKey: string(privateKey),
		}
	}

	subscriptions := SubscriptionsFile{
		Topics: s.subscriptions,
	}

	err := writeJSON(filepath.Join(basePath, "secrets.json"), &secrets)
	if err != nil {
		return err
	}

	err = writeJSON(filepath.Join(basePath, "subscriptions.json"), &subscriptions)
	if err != nil {
		return err
	}

	return nil
}

func readJSON(path string, v any) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	return json.NewDecoder(file).Decode(v)
}

func writeJSON(path string, v any) error {
	file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return err
	}

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(v); err != nil {
		file.Close()
		return err
	}

	return file.Close()
}
