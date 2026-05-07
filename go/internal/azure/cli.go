package azure

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

var (
	ErrAzureCliNotFound = errors.New("Azure CLI is not installed or not in PATH")
	ErrNoActiveSession  = errors.New("No active Azure session. Run 'az login' first.")
)

var windowsFallbackPaths = []string{
	`C:\Program Files\Microsoft SDKs\Azure\CLI2\wbin\az.cmd`,
	`C:\Program Files (x86)\Microsoft SDKs\Azure\CLI2\wbin\az.cmd`,
}

type User struct {
	Name string `json:"name"`
}

type Account struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	TenantID string `json:"tenantId"`
	User     User   `json:"user"`
}

type Client struct{}

type commandResult struct {
	stdout   string
	stderr   string
	exitCode int
}

func NewClient() *Client {
	return &Client{}
}

func (c *Client) CheckInstalled() bool {
	result, err := c.runCommand([]string{"--version"}, 30*time.Second)
	if err != nil {
		return false
	}

	return result.exitCode == 0
}

func (c *Client) GetCurrentAccount() (*Account, error) {
	if !c.CheckInstalled() {
		return nil, ErrAzureCliNotFound
	}

	result, err := c.runCommand([]string{"account", "show", "--output", "json"}, 30*time.Second)
	if err != nil {
		return nil, err
	}

	if result.exitCode != 0 {
		if strings.Contains(strings.ToLower(result.stderr), "az login") {
			return nil, ErrNoActiveSession
		}
		return nil, fmt.Errorf("%w: %s", ErrNoActiveSession, strings.TrimSpace(result.stderr))
	}

	var account Account
	if err := json.Unmarshal([]byte(result.stdout), &account); err != nil {
		return nil, fmt.Errorf("failed to parse Azure CLI output: %w", err)
	}

	return &account, nil
}

func (c *Client) SetAccount(subscriptionID string) (bool, error) {
	if !c.CheckInstalled() {
		return false, ErrAzureCliNotFound
	}

	result, err := c.runCommand([]string{"account", "set", "--subscription", subscriptionID}, 5*time.Second)
	if err != nil {
		if errors.Is(err, ErrNoActiveSession) || errors.Is(err, ErrAzureCliNotFound) {
			return false, err
		}
		return false, nil
	}

	return result.exitCode == 0, nil
}

func (c *Client) runCommand(args []string, timeout time.Duration) (*commandResult, error) {
	azExecutable, err := findExecutable()
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	command, err := buildCommand(ctx, azExecutable, args)
	if err != nil {
		return nil, err
	}

	command.Env = os.Environ()

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	command.Stdout = &stdout
	command.Stderr = &stderr

	runErr := command.Run()
	if ctx.Err() == context.DeadlineExceeded {
		return nil, fmt.Errorf("Azure CLI command timed out")
	}

	result := &commandResult{
		stdout: strings.TrimSpace(stdout.String()),
		stderr: strings.TrimSpace(stderr.String()),
	}

	var exitErr *exec.ExitError
	if errors.As(runErr, &exitErr) {
		result.exitCode = exitErr.ExitCode()
		return result, nil
	}

	if runErr != nil {
		return nil, runErr
	}

	result.exitCode = 0
	return result, nil
}

func findExecutable() (string, error) {
	if azPath, err := exec.LookPath("az"); err == nil {
		return azPath, nil
	}

	if runtime.GOOS == "windows" {
		for _, candidate := range windowsFallbackPaths {
			if _, err := os.Stat(candidate); err == nil {
				return candidate, nil
			}
		}
	}

	return "", ErrAzureCliNotFound
}

func buildCommand(ctx context.Context, azExecutable string, args []string) (*exec.Cmd, error) {
	if runtime.GOOS == "windows" && strings.EqualFold(filepath.Ext(azExecutable), ".cmd") {
		commandArgs := append([]string{"/c", azExecutable}, args...)
		return exec.CommandContext(ctx, "cmd.exe", commandArgs...), nil
	}

	return exec.CommandContext(ctx, azExecutable, args...), nil
}