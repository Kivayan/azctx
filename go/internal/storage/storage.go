package storage

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/Kivayan/azctx/go/internal/domain"
	"gopkg.in/yaml.v3"
)

type Store struct {
	storagePath string
}

type filePayload struct {
	Contexts []domain.Context `yaml:"contexts"`
}

type rawContext struct {
	ContextID        string `yaml:"context_id"`
	ContextName      string `yaml:"context_name"`
	SubscriptionID   string `yaml:"subscription_id"`
	SubscriptionName string `yaml:"subscription_name"`
	TenantID         string `yaml:"tenant_id"`
	TenantName       string `yaml:"tenant_name"`
	Username         string `yaml:"username"`
	CreatedAt        string `yaml:"created_at"`
 }

func NewStore(storagePath string) (*Store, error) {
	if storagePath == "" {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return nil, fmt.Errorf("failed to determine home directory: %w", err)
		}
		storagePath = filepath.Join(homeDir, ".azctx", "contexts.yaml")
	}

	return &Store{storagePath: storagePath}, nil
}

func (s *Store) Path() string {
	return s.storagePath
}

func (s *Store) LoadContexts() ([]domain.Context, error) {
	content, err := os.ReadFile(s.storagePath)
	if err != nil {
		if os.IsNotExist(err) {
			return []domain.Context{}, nil
		}
		return nil, fmt.Errorf("failed to read contexts file: %w", err)
	}

	if len(strings.TrimSpace(string(content))) == 0 {
		return []domain.Context{}, nil
	}

	var root yaml.Node
	if err := yaml.Unmarshal(content, &root); err != nil {
		return nil, fmt.Errorf("failed to parse contexts file: %w", err)
	}

	sequenceNode := findContextsSequence(&root)
	if sequenceNode == nil || len(sequenceNode.Content) == 0 {
		return []domain.Context{}, nil
	}

	contexts := make([]domain.Context, 0, len(sequenceNode.Content))
	for _, entry := range sequenceNode.Content {
		if entry == nil {
			fmt.Fprintln(os.Stderr, "Warning: skipping invalid context entry: empty node")
			continue
		}

		var raw rawContext
		if err := entry.Decode(&raw); err != nil {
			fmt.Fprintf(os.Stderr, "Warning: skipping invalid context entry: %v\n", err)
			continue
		}

		createdAt, err := domain.ParseTimestamp(strings.TrimSpace(raw.CreatedAt))
		if err != nil {
			fmt.Fprintf(os.Stderr, "Warning: skipping invalid context entry: %v\n", err)
			continue
		}

		context := domain.Context{
			ContextID:        raw.ContextID,
			ContextName:      raw.ContextName,
			SubscriptionID:   raw.SubscriptionID,
			SubscriptionName: raw.SubscriptionName,
			TenantID:         raw.TenantID,
			TenantName:       raw.TenantName,
			Username:         raw.Username,
			CreatedAt:        createdAt,
		}

		contexts = append(contexts, context)
	}

	return contexts, nil
}

func findContextsSequence(root *yaml.Node) *yaml.Node {
	if root == nil {
		return nil
	}

	node := root
	if node.Kind == yaml.DocumentNode && len(node.Content) > 0 {
		node = node.Content[0]
	}

	if node.Kind != yaml.MappingNode {
		return nil
	}

	for index := 0; index+1 < len(node.Content); index += 2 {
		keyNode := node.Content[index]
		valueNode := node.Content[index+1]
		if keyNode != nil && keyNode.Value == "contexts" && valueNode != nil && valueNode.Kind == yaml.SequenceNode {
			return valueNode
		}
	}

	return nil
}

func (s *Store) SaveContexts(contexts []domain.Context) error {
	storageDir := filepath.Dir(s.storagePath)
	if err := os.MkdirAll(storageDir, 0o755); err != nil {
		return fmt.Errorf("failed to create storage directory: %w", err)
	}

	encoded, err := yaml.Marshal(filePayload{Contexts: contexts})
	if err != nil {
		return fmt.Errorf("failed to encode contexts: %w", err)
	}

	tempFile, err := os.CreateTemp(storageDir, "contexts-*.yaml")
	if err != nil {
		return fmt.Errorf("failed to create temp file: %w", err)
	}

	tempPath := tempFile.Name()
	cleanup := func() {
		_ = tempFile.Close()
		_ = os.Remove(tempPath)
	}

	if _, err := tempFile.Write(encoded); err != nil {
		cleanup()
		return fmt.Errorf("failed to write temp file: %w", err)
	}

	if err := tempFile.Close(); err != nil {
		cleanup()
		return fmt.Errorf("failed to close temp file: %w", err)
	}

	if _, err := os.Stat(s.storagePath); err == nil {
		if err := os.Remove(s.storagePath); err != nil {
			_ = os.Remove(tempPath)
			return fmt.Errorf("failed to replace contexts file: %w", err)
		}
	}

	if err := os.Rename(tempPath, s.storagePath); err != nil {
		_ = os.Remove(tempPath)
		return fmt.Errorf("failed to move temp file into place: %w", err)
	}

	return nil
	}

func (s *Store) AddContext(context domain.Context) error {
	contexts, err := s.LoadContexts()
	if err != nil {
		return err
	}

	contexts = append(contexts, context)
	return s.SaveContexts(contexts)
}

func (s *Store) DeleteContext(contextID string) error {
	contexts, err := s.LoadContexts()
	if err != nil {
		return err
	}

	filtered := make([]domain.Context, 0, len(contexts))
	for _, context := range contexts {
		if context.ContextID != contextID {
			filtered = append(filtered, context)
		}
	}

	return s.SaveContexts(filtered)
}

func (s *Store) GetContextByID(contextID string) (*domain.Context, error) {
	contexts, err := s.LoadContexts()
	if err != nil {
		return nil, err
	}

	for _, context := range contexts {
		if context.ContextID == contextID {
			matched := context
			return &matched, nil
		}
	}

	return nil, nil
}

func (s *Store) ContextIDExists(contextID string) (bool, error) {
	context, err := s.GetContextByID(contextID)
	if err != nil {
		return false, err
	}

	return context != nil, nil
}