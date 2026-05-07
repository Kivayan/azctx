package service

import (
	"errors"
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/AlecAivazis/survey/v2"
	"github.com/AlecAivazis/survey/v2/terminal"
	"github.com/Kivayan/azctx/go/internal/azure"
	"github.com/Kivayan/azctx/go/internal/domain"
	"github.com/Kivayan/azctx/go/internal/storage"
)

type ContextSummary struct {
	ContextID   string
	ContextName string
}

type SwitchResult struct {
	Success      bool
	Message      string
	Context      *domain.Context
	Error        string
	AvailableIDs []string
}

type AddResult struct {
	Success        bool
	Message        string
	Context        *domain.Context
	CurrentAccount *azure.Account
	Error          string
}

type StatusResult struct {
	Success        bool
	Message        string
	Context        *domain.Context
	CurrentAccount *azure.Account
	IsManaged      bool
	Verbose        bool
	Error          string
}

type ListResult struct {
	Success   bool
	Message   string
	Contexts  []domain.Context
	Summaries []ContextSummary
	Verbose   bool
	Error     string
}

type DeleteResult struct {
	Success     bool
	Message     string
	ContextID   string
	ContextName string
	Error       string
}

type Manager struct {
	store *storage.Store
	azure *azure.Client
}

func NewManager(store *storage.Store, azureClient *azure.Client) *Manager {
	return &Manager{store: store, azure: azureClient}
}

func (m *Manager) SwitchContextInteractive() SwitchResult {
	if !m.azure.CheckInstalled() {
		return SwitchResult{Success: false, Message: azure.ErrAzureCliNotFound.Error(), Error: "azure_cli_not_found"}
	}

	contexts, err := m.store.LoadContexts()
	if err != nil {
		return SwitchResult{Success: false, Message: fmt.Sprintf("Failed to load contexts: %v", err), Error: "storage_error"}
	}

	if len(contexts) == 0 {
		return SwitchResult{
			Success: false,
			Message: "No saved contexts found. Use 'azctx add' to save your current context first.",
			Error:   "empty_list",
		}
	}

	if len(contexts) == 1 {
		context := contexts[0]
		return SwitchResult{
			Success: false,
			Message: fmt.Sprintf(
				"Only one context saved: %s (%s). Use 'azctx add' to save more contexts before switching.",
				context.ContextName,
				context.ContextID,
			),
			Context: &context,
			Error:   "single_context",
		}
	}

	choices := make([]string, 0, len(contexts))
	for _, context := range contexts {
		choices = append(choices, formatChoice(context))
	}

	var selected string
	err = survey.AskOne(&survey.Select{
		Message: "Select a context to switch to:",
		Options: choices,
	}, &selected)
	if isCancelled(err) {
		return SwitchResult{Success: false, Message: "Context switch cancelled", Error: "cancelled"}
	}
	if err != nil {
		return SwitchResult{Success: false, Message: fmt.Sprintf("Failed to switch context: %v", err), Error: "unknown"}
	}

	contextID := parseChoiceID(selected)
	selectedContext, err := m.store.GetContextByID(contextID)
	if err != nil {
		return SwitchResult{Success: false, Message: fmt.Sprintf("Failed to load context '%s': %v", contextID, err), Error: "storage_error"}
	}
	if selectedContext == nil {
		return SwitchResult{Success: false, Message: fmt.Sprintf("Context '%s' not found", contextID), Error: "not_found"}
	}

	switched, err := m.azure.SetAccount(selectedContext.SubscriptionID)
	if err != nil {
		if errors.Is(err, azure.ErrNoActiveSession) {
			return SwitchResult{Success: false, Message: err.Error(), Error: "no_session"}
		}
		if errors.Is(err, azure.ErrAzureCliNotFound) {
			return SwitchResult{Success: false, Message: err.Error(), Error: "azure_cli_not_found"}
		}
		return SwitchResult{Success: false, Message: fmt.Sprintf("Failed to switch context: %v", err), Error: "unknown"}
	}
	if !switched {
		return SwitchResult{Success: false, Message: fmt.Sprintf("Failed to switch to context '%s'.", contextID), Context: selectedContext, Error: "switch_failed"}
	}

	currentAccount, err := m.azure.GetCurrentAccount()
	if err != nil {
		if errors.Is(err, azure.ErrNoActiveSession) {
			return SwitchResult{Success: false, Message: err.Error(), Error: "no_session"}
		}
		return SwitchResult{Success: false, Message: fmt.Sprintf("Failed to switch context: %v", err), Error: "unknown"}
	}

	if currentAccount != nil && currentAccount.ID != selectedContext.SubscriptionID {
		return SwitchResult{
			Success: false,
			Message: fmt.Sprintf(
				"Failed to switch to context '%s'. Azure CLI returned different subscription.",
				contextID,
			),
			Context: selectedContext,
			Error:   "verification_failed",
		}
	}

	return SwitchResult{
		Success: true,
		Message: fmt.Sprintf("Successfully switched to context: %s", selectedContext.ContextID),
		Context: selectedContext,
	}
}

func (m *Manager) SwitchContextByID(contextID string) SwitchResult {
	contextID = strings.TrimSpace(contextID)
	if contextID == "" {
		return SwitchResult{Success: false, Message: "Context ID cannot be empty", Error: "empty_id"}
	}

	if !m.azure.CheckInstalled() {
		return SwitchResult{Success: false, Message: azure.ErrAzureCliNotFound.Error(), Error: "azure_cli_not_found"}
	}

	contexts, err := m.store.LoadContexts()
	if err != nil {
		return SwitchResult{Success: false, Message: fmt.Sprintf("Failed to load contexts: %v", err), Error: "storage_error"}
	}

	if len(contexts) == 0 {
		return SwitchResult{
			Success: false,
			Message: "No saved contexts found. Use 'azctx add' to save your current context first.",
			Error:   "empty_list",
		}
	}

	selectedContext, err := m.store.GetContextByID(contextID)
	if err != nil {
		return SwitchResult{Success: false, Message: fmt.Sprintf("Failed to load context '%s': %v", contextID, err), Error: "storage_error"}
	}
	if selectedContext == nil {
		availableIDs := make([]string, 0, len(contexts))
		for _, context := range contexts {
			availableIDs = append(availableIDs, context.ContextID)
		}
		sort.Strings(availableIDs)
		return SwitchResult{
			Success:      false,
			Message:      fmt.Sprintf("Context '%s' not found.", contextID),
			Error:        "not_found",
			AvailableIDs: availableIDs,
		}
	}

	currentAccount, err := m.azure.GetCurrentAccount()
	if err != nil {
		if errors.Is(err, azure.ErrNoActiveSession) {
			return SwitchResult{Success: false, Message: err.Error(), Error: "no_session"}
		}
		if errors.Is(err, azure.ErrAzureCliNotFound) {
			return SwitchResult{Success: false, Message: err.Error(), Error: "azure_cli_not_found"}
		}
	} else if currentAccount != nil && currentAccount.ID == selectedContext.SubscriptionID {
		return SwitchResult{
			Success: false,
			Message: fmt.Sprintf("Context '%s' is already active.", contextID),
			Context: selectedContext,
			Error:   "already_active",
		}
	}

	switched, err := m.azure.SetAccount(selectedContext.SubscriptionID)
	if err != nil {
		if errors.Is(err, azure.ErrNoActiveSession) {
			return SwitchResult{Success: false, Message: err.Error(), Context: selectedContext, Error: "no_session"}
		}
		if errors.Is(err, azure.ErrAzureCliNotFound) {
			return SwitchResult{Success: false, Message: err.Error(), Error: "azure_cli_not_found"}
		}
		return SwitchResult{Success: false, Message: fmt.Sprintf("Failed to switch context: %v", err), Context: selectedContext, Error: "unknown"}
	}
	if !switched {
		return SwitchResult{Success: false, Message: fmt.Sprintf("Failed to switch to context '%s'.", contextID), Context: selectedContext, Error: "switch_failed"}
	}

	verifiedAccount, err := m.azure.GetCurrentAccount()
	if err == nil && verifiedAccount != nil && verifiedAccount.ID != selectedContext.SubscriptionID {
		return SwitchResult{
			Success: false,
			Message: fmt.Sprintf("Failed to verify switch to context '%s'. Azure CLI returned different subscription.", contextID),
			Context: selectedContext,
			Error:   "verification_failed",
		}
	}

	return SwitchResult{
		Success: true,
		Message: fmt.Sprintf("Successfully switched to context: %s", contextID),
		Context: selectedContext,
	}
}

func (m *Manager) AddContextInteractive() AddResult {
	if !m.azure.CheckInstalled() {
		return AddResult{Success: false, Message: azure.ErrAzureCliNotFound.Error(), Error: "azure_cli_not_found"}
	}

	currentAccount, err := m.azure.GetCurrentAccount()
	if err != nil {
		if errors.Is(err, azure.ErrNoActiveSession) {
			return AddResult{Success: false, Message: err.Error(), Error: "no_session"}
		}
		if errors.Is(err, azure.ErrAzureCliNotFound) {
			return AddResult{Success: false, Message: err.Error(), Error: "azure_cli_not_found"}
		}
		return AddResult{Success: false, Message: fmt.Sprintf("Failed to load Azure account: %v", err), Error: "unknown"}
	}

	contexts, err := m.store.LoadContexts()
	if err != nil {
		return AddResult{Success: false, Message: fmt.Sprintf("Failed to load contexts: %v", err), CurrentAccount: currentAccount, Error: "storage_error"}
	}

	for _, context := range contexts {
		if context.SubscriptionID == currentAccount.ID && context.TenantID == currentAccount.TenantID && context.Username == currentAccount.User.Name {
			matched := context
			return AddResult{
				Success:        false,
				CurrentAccount: currentAccount,
				Context:        &matched,
				Message: fmt.Sprintf(
					"This Azure context is already managed as '%s' (%s).\nSubscription: %s\nTenant: %s\nAccount: %s",
					context.ContextName,
					context.ContextID,
					context.SubscriptionName,
					context.TenantID,
					context.Username,
				),
				Error: "already_exists",
			}
		}
	}

	contextName, cancelled, err := askRequiredInput(
		"Enter a friendly name for this context (1-100 characters):",
		func(answer string) error {
			trimmed := strings.TrimSpace(answer)
			if len(trimmed) < 1 || len(trimmed) > 100 {
				return fmt.Errorf("name must be between 1 and 100 characters")
			}
			return nil
		},
	)
	if cancelled {
		return AddResult{Success: false, Message: "Context addition cancelled", CurrentAccount: currentAccount, Error: "cancelled"}
	}
	if err != nil {
		return AddResult{Success: false, Message: fmt.Sprintf("Failed to read context name: %v", err), CurrentAccount: currentAccount, Error: "unknown"}
	}

	var contextID string
	for {
		value, cancelled, err := askRequiredInput(
			"Enter a short ID for this context (1-20 chars, alphanumeric/hyphens/underscores):",
			func(answer string) error {
				if !domain.ValidateContextID(strings.TrimSpace(answer)) {
					return fmt.Errorf("ID must be 1-20 characters (alphanumeric, hyphens, underscores only)")
				}
				return nil
			},
		)
		if cancelled {
			return AddResult{Success: false, Message: "Context addition cancelled", CurrentAccount: currentAccount, Error: "cancelled"}
		}
		if err != nil {
			return AddResult{Success: false, Message: fmt.Sprintf("Failed to read context ID: %v", err), CurrentAccount: currentAccount, Error: "unknown"}
		}

		exists, err := m.store.ContextIDExists(value)
		if err != nil {
			return AddResult{Success: false, Message: fmt.Sprintf("Failed to validate context ID: %v", err), CurrentAccount: currentAccount, Error: "storage_error"}
		}
		if exists {
			fmt.Printf("Error: Context ID '%s' already exists. Please choose a different ID.\n", value)
			continue
		}

		contextID = value
		break
	}

	newContext := domain.Context{
		ContextID:        contextID,
		ContextName:      contextName,
		SubscriptionID:   currentAccount.ID,
		SubscriptionName: currentAccount.Name,
		TenantID:         currentAccount.TenantID,
		TenantName:       currentAccount.TenantID,
		Username:         currentAccount.User.Name,
		CreatedAt:        time.Now(),
	}

	if err := m.store.AddContext(newContext); err != nil {
		return AddResult{Success: false, Message: fmt.Sprintf("Failed to save context: %v", err), CurrentAccount: currentAccount, Error: "storage_error"}
	}

	return AddResult{
		Success:        true,
		Message:        fmt.Sprintf("Successfully added context: %s", contextName),
		Context:        &newContext,
		CurrentAccount: currentAccount,
	}
}

func (m *Manager) GetStatus(verbose bool) StatusResult {
	if !m.azure.CheckInstalled() {
		return StatusResult{Success: false, Message: azure.ErrAzureCliNotFound.Error(), Verbose: verbose, Error: "azure_cli_not_found"}
	}

	currentAccount, err := m.azure.GetCurrentAccount()
	if err != nil {
		if errors.Is(err, azure.ErrNoActiveSession) {
			return StatusResult{Success: false, Message: err.Error(), Verbose: verbose, Error: "no_session"}
		}
		return StatusResult{Success: false, Message: fmt.Sprintf("Failed to load Azure account: %v", err), Verbose: verbose, Error: "unknown"}
	}

	contexts, err := m.store.LoadContexts()
	if err != nil {
		return StatusResult{Success: false, Message: fmt.Sprintf("Failed to load contexts: %v", err), CurrentAccount: currentAccount, Verbose: verbose, Error: "storage_error"}
	}

	for _, context := range contexts {
		if context.SubscriptionID == currentAccount.ID {
			matched := context
			return StatusResult{
				Success:        true,
				Message:        "Current context is managed by azctx",
				Context:        &matched,
				CurrentAccount: currentAccount,
				IsManaged:      true,
				Verbose:        verbose,
			}
		}
	}

	return StatusResult{
		Success:        true,
		Message:        "Current context is not managed. Use 'azctx add' to save it.",
		CurrentAccount: currentAccount,
		IsManaged:      false,
		Verbose:        verbose,
	}
}

func (m *Manager) GetCurrentAccount() (*azure.Account, error) {
	if !m.azure.CheckInstalled() {
		return nil, azure.ErrAzureCliNotFound
	}

	return m.azure.GetCurrentAccount()
}

func (m *Manager) ListContexts(verbose bool) ListResult {
	contexts, err := m.store.LoadContexts()
	if err != nil {
		return ListResult{Success: false, Message: fmt.Sprintf("Failed to load contexts: %v", err), Verbose: verbose, Error: "storage_error"}
	}

	if len(contexts) == 0 {
		return ListResult{
			Success: true,
			Message: "No saved contexts found. Use 'azctx add' to save your current context.",
			Verbose: verbose,
			Error:   "empty_list",
		}
	}

	if verbose {
		return ListResult{Success: true, Message: fmt.Sprintf("Found %d saved context(s)", len(contexts)), Contexts: contexts, Verbose: true}
	}

	summaries := make([]ContextSummary, 0, len(contexts))
	for _, context := range contexts {
		summaries = append(summaries, ContextSummary{ContextID: context.ContextID, ContextName: context.ContextName})
	}

	return ListResult{Success: true, Message: fmt.Sprintf("Found %d saved context(s)", len(contexts)), Summaries: summaries, Verbose: false}
}

func (m *Manager) DeleteContextInteractive() DeleteResult {
	contexts, err := m.store.LoadContexts()
	if err != nil {
		return DeleteResult{Success: false, Message: fmt.Sprintf("Failed to load contexts: %v", err), Error: "storage_error"}
	}

	if len(contexts) == 0 {
		return DeleteResult{Success: false, Message: "No saved contexts found. Nothing to delete.", Error: "empty_list"}
	}

	choices := make([]string, 0, len(contexts))
	for _, context := range contexts {
		choices = append(choices, formatChoice(context))
	}

	var selected string
	err = survey.AskOne(&survey.Select{Message: "Select a context to delete:", Options: choices}, &selected)
	if isCancelled(err) {
		return DeleteResult{Success: false, Message: "Deletion cancelled", Error: "cancelled"}
	}
	if err != nil {
		return DeleteResult{Success: false, Message: fmt.Sprintf("Failed to delete context: %v", err), Error: "unknown"}
	}

	contextID := parseChoiceID(selected)
	selectedContext, err := m.store.GetContextByID(contextID)
	if err != nil {
		return DeleteResult{Success: false, Message: fmt.Sprintf("Failed to load context '%s': %v", contextID, err), Error: "storage_error"}
	}
	if selectedContext == nil {
		return DeleteResult{Success: false, Message: fmt.Sprintf("Context '%s' not found", contextID), Error: "not_found"}
	}

	confirmed, cancelled, err := askConfirmation(
		fmt.Sprintf("Are you sure you want to delete context '%s' (%s)?", selectedContext.ContextName, selectedContext.ContextID),
		false,
	)
	if cancelled {
		return DeleteResult{Success: false, Message: "Deletion cancelled", Error: "cancelled"}
	}
	if err != nil {
		return DeleteResult{Success: false, Message: fmt.Sprintf("Failed to delete context: %v", err), Error: "unknown"}
	}
	if !confirmed {
		return DeleteResult{Success: false, Message: "Deletion cancelled", Error: "cancelled"}
	}

	if err := m.store.DeleteContext(contextID); err != nil {
		return DeleteResult{Success: false, Message: fmt.Sprintf("Failed to delete context: %v", err), Error: "storage_error"}
	}

	return DeleteResult{
		Success:     true,
		Message:     fmt.Sprintf("Successfully deleted context: %s (%s)", selectedContext.ContextName, selectedContext.ContextID),
		ContextID:   selectedContext.ContextID,
		ContextName: selectedContext.ContextName,
	}
}

func askRequiredInput(message string, validate func(string) error) (string, bool, error) {
	for {
		var value string
		validator := func(answer interface{}) error {
			text, ok := answer.(string)
			if !ok {
				return fmt.Errorf("invalid input type")
			}
			return validate(text)
		}

		err := survey.AskOne(&survey.Input{Message: message}, &value, survey.WithValidator(validator))
		if isCancelled(err) {
			return "", true, nil
		}
		if err != nil {
			return "", false, err
		}

		trimmed := strings.TrimSpace(value)
		if trimmed == "" {
			continue
		}

		return trimmed, false, nil
	}
}

func askConfirmation(message string, defaultValue bool) (bool, bool, error) {
	var confirmed bool
	err := survey.AskOne(&survey.Confirm{Message: message, Default: defaultValue}, &confirmed)
	if isCancelled(err) {
		return false, true, nil
	}
	if err != nil {
		return false, false, err
	}
	return confirmed, false, nil
}

func isCancelled(err error) bool {
	return errors.Is(err, terminal.InterruptErr)
}

func formatChoice(context domain.Context) string {
	return fmt.Sprintf("%s (%s)", context.ContextName, context.ContextID)
}

func parseChoiceID(choice string) string {
	lastOpenParen := strings.LastIndex(choice, "(")
	if lastOpenParen == -1 {
		return strings.TrimSpace(choice)
	}
	return strings.TrimSuffix(strings.TrimSpace(choice[lastOpenParen+1:]), ")")
}