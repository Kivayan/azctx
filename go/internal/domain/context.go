package domain

import (
	"fmt"
	"regexp"
	"strings"
	"time"

	"gopkg.in/yaml.v3"
)

const ContextIDPattern = `^[a-zA-Z0-9_-]{1,20}$`

const timestampLayout = "2006-01-02T15:04:05.999999999"

var contextIDRegex = regexp.MustCompile(ContextIDPattern)

type Context struct {
	ContextID        string
	ContextName      string
	SubscriptionID   string
	SubscriptionName string
	TenantID         string
	TenantName       string
	Username         string
	CreatedAt        time.Time
}

func ValidateContextID(contextID string) bool {
	if contextID == "" {
		return false
	}
	return contextIDRegex.MatchString(contextID)
}

func (c Context) MarshalYAML() (interface{}, error) {
	type alias struct {
		ContextID        string `yaml:"context_id"`
		ContextName      string `yaml:"context_name"`
		SubscriptionID   string `yaml:"subscription_id"`
		SubscriptionName string `yaml:"subscription_name"`
		TenantID         string `yaml:"tenant_id"`
		TenantName       string `yaml:"tenant_name"`
		Username         string `yaml:"username"`
		CreatedAt        string `yaml:"created_at"`
	}

	return alias{
		ContextID:        c.ContextID,
		ContextName:      c.ContextName,
		SubscriptionID:   c.SubscriptionID,
		SubscriptionName: c.SubscriptionName,
		TenantID:         c.TenantID,
		TenantName:       c.TenantName,
		Username:         c.Username,
		CreatedAt:        c.CreatedAt.Format(timestampLayout),
	}, nil
}

func (c *Context) UnmarshalYAML(node *yaml.Node) error {
	type alias struct {
		ContextID        string `yaml:"context_id"`
		ContextName      string `yaml:"context_name"`
		SubscriptionID   string `yaml:"subscription_id"`
		SubscriptionName string `yaml:"subscription_name"`
		TenantID         string `yaml:"tenant_id"`
		TenantName       string `yaml:"tenant_name"`
		Username         string `yaml:"username"`
		CreatedAt        string `yaml:"created_at"`
	}

	var payload alias
	if err := node.Decode(&payload); err != nil {
		return err
	}

	createdAt := strings.TrimSpace(payload.CreatedAt)
	if createdAt == "" {
		return fmt.Errorf("missing created_at")
	}

	parsedCreatedAt, err := ParseTimestamp(createdAt)
	if err != nil {
		return fmt.Errorf("invalid created_at: %w", err)
	}

	*c = Context{
		ContextID:        payload.ContextID,
		ContextName:      payload.ContextName,
		SubscriptionID:   payload.SubscriptionID,
		SubscriptionName: payload.SubscriptionName,
		TenantID:         payload.TenantID,
		TenantName:       payload.TenantName,
		Username:         payload.Username,
		CreatedAt:        parsedCreatedAt,
	}

	return nil
}

func ParseTimestamp(value string) (time.Time, error) {
	layouts := []string{
		time.RFC3339Nano,
		timestampLayout,
		"2006-01-02T15:04:05",
	}

	var lastErr error
	for _, layout := range layouts {
		parsed, err := time.Parse(layout, value)
		if err == nil {
			return parsed, nil
		}
		lastErr = err
	}

	return time.Time{}, lastErr
}