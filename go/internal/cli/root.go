package cli

import (
	"errors"
	"fmt"
	"time"

	"github.com/spf13/cobra"

	"github.com/Kivayan/azctx/go/internal/azure"
	"github.com/Kivayan/azctx/go/internal/service"
	"github.com/Kivayan/azctx/go/internal/storage"
	"github.com/Kivayan/azctx/go/internal/ui"
)

var ErrVersionRequested = errors.New("version requested")

type ExitError struct {
	Code int
	Err  error
}

func (e ExitError) Error() string {
	if e.Err != nil {
		return e.Err.Error()
	}
	return fmt.Sprintf("process exited with code %d", e.Code)
}

func NewRootCommand(version string) *cobra.Command {
	store, err := storage.NewStore("")
	if err != nil {
		return &cobra.Command{RunE: func(cmd *cobra.Command, args []string) error { return ExitError{Code: 1, Err: err} }}
	}

	manager := service.NewManager(store, azure.NewClient())
	var showVersion bool

	root := &cobra.Command{
		Use:           "azctx",
		Short:         "Azure CLI context switcher - manage and switch between Azure accounts quickly.",
		SilenceUsage:  true,
		SilenceErrors: true,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			if showVersion {
				fmt.Printf("azctx version %s\n", version)
				return ErrVersionRequested
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
	}

	root.PersistentFlags().BoolVarP(&showVersion, "version", "V", false, "Show version information.")
	root.CompletionOptions.DisableDefaultCmd = true

	root.AddCommand(newSwitchCommand(manager))
	root.AddCommand(newStatusCommand(manager))
	root.AddCommand(newAddCommand(manager))
	root.AddCommand(newListCommand(manager))
	root.AddCommand(newDeleteCommand(manager))

	return root
}

func newSwitchCommand(manager *service.Manager) *cobra.Command {
	var contextID string

	command := &cobra.Command{
		Use:   "switch",
		Short: "Switch to a different saved context.",
		RunE: func(cmd *cobra.Command, args []string) error {
			var result service.SwitchResult
			if contextID != "" {
				result = manager.SwitchContextByID(contextID)
			} else {
				result = manager.SwitchContextInteractive()
			}

			if result.Success {
				context := result.Context
				ui.PrintFieldPanel(
					fmt.Sprintf("Successfully Switched to %s", context.ContextName),
					[]ui.Field{{Name: "Name", Value: context.ContextName}, {Name: "ID", Value: context.ContextID}, {Name: "Subscription", Value: context.SubscriptionName}, {Name: "Tenant", Value: context.TenantName}, {Name: "Account", Value: context.Username}},
				)
				return ExitError{Code: 0}
			}

			switch result.Error {
			case "cancelled":
				ui.PrintMessagePanel("Context Switch Cancelled", result.Message)
				return ExitError{Code: 130}
			case "already_active":
				ui.PrintMessagePanel("Already Active", result.Message)
				return ExitError{Code: 0}
			case "not_found":
				message := result.Message
				if len(result.AvailableIDs) > 0 {
					message = fmt.Sprintf("%s\n\nAvailable contexts: %s", message, joinWithComma(result.AvailableIDs))
				}
				ui.PrintMessagePanel("Context Not Found", message)
				return ExitError{Code: 1}
			case "empty_list", "single_context":
				ui.PrintMessagePanel("Cannot Switch Context", result.Message)
				return ExitError{Code: 1}
			case "azure_cli_not_found":
				ui.PrintMessagePanel("Azure CLI Not Found", result.Message+"\n\nPlease install Azure CLI: https://aka.ms/azure-cli")
				return ExitError{Code: 1}
			default:
				ui.PrintMessagePanel("Failed to Switch Context", result.Message)
				return ExitError{Code: 1}
			}
		},
	}

	command.Flags().StringVarP(&contextID, "id", "i", "", "Context ID to switch to directly (case-sensitive). Omit for interactive selection.")
	return command
}

func newStatusCommand(manager *service.Manager) *cobra.Command {
	var verbose bool

	command := &cobra.Command{
		Use:   "status",
		Short: "Show currently active Azure context.",
		RunE: func(cmd *cobra.Command, args []string) error {
			result := manager.GetStatus(verbose)
			if !result.Success {
				title := "Error"
				message := result.Message
				if result.Error == "no_session" {
					title = "No Active Azure Session"
					message += "\n\nPlease run 'az login' to authenticate with Azure first."
				}
				if result.Error == "azure_cli_not_found" {
					title = "Azure CLI Not Found"
					message += "\n\nPlease install Azure CLI: https://aka.ms/azure-cli"
				}
				ui.PrintMessagePanel(title, message)
				return ExitError{Code: 1}
			}

			if result.IsManaged {
				context := result.Context
				fields := []ui.Field{{Name: "Name", Value: context.ContextName}, {Name: "ID", Value: context.ContextID}, {Name: "Subscription", Value: context.SubscriptionName}, {Name: "Tenant", Value: context.TenantName}, {Name: "Account", Value: context.Username}}
				extra := []string{}
				if verbose {
					extra = append(extra, "", "Detailed Information:", fmt.Sprintf("Subscription ID: %s", context.SubscriptionID), fmt.Sprintf("Tenant ID: %s", context.TenantID), fmt.Sprintf("Created: %s", formatTimestamp(context.CreatedAt)))
				}
				ui.PrintFieldPanel(fmt.Sprintf("Managed Context: %s", context.ContextName), fields, extra...)
				return ExitError{Code: 0}
			}

			current := result.CurrentAccount
			extra := []string{"", "This context is not managed by azctx.", "Run azctx add to save it with a friendly name."}
			if verbose {
				extra = append(extra, "", "Detailed Information:", fmt.Sprintf("Subscription ID: %s", current.ID), fmt.Sprintf("Tenant ID: %s", current.TenantID))
			}
			ui.PrintFieldPanel("Unmanaged Context", []ui.Field{{Name: "Subscription", Value: current.Name}, {Name: "Tenant", Value: current.TenantID}, {Name: "Account", Value: current.User.Name}}, extra...)
			return ExitError{Code: 0}
		},
	}

	command.Flags().BoolVarP(&verbose, "verbose", "v", false, "Show detailed information")
	return command
}

func newAddCommand(manager *service.Manager) *cobra.Command {
	return &cobra.Command{
		Use:   "add",
		Short: "Add current Azure CLI context with friendly identifiers.",
		RunE: func(cmd *cobra.Command, args []string) error {
			currentAccount, err := managerCurrentAccount(manager)
			if err == nil && currentAccount != nil {
				ui.PrintFieldPanel("Current Azure Context", []ui.Field{{Name: "Subscription", Value: currentAccount.Name}, {Name: "Tenant", Value: currentAccount.TenantID}, {Name: "Account", Value: currentAccount.User.Name}})
				fmt.Println()
			}

			result := manager.AddContextInteractive()
			if result.Success {
				context := result.Context
				ui.PrintFieldPanel(fmt.Sprintf("Successfully Added Context: %s", context.ContextName), []ui.Field{{Name: "Name", Value: context.ContextName}, {Name: "ID", Value: context.ContextID}, {Name: "Subscription", Value: context.SubscriptionName}, {Name: "Tenant", Value: context.TenantName}, {Name: "Account", Value: context.Username}})
				return ExitError{Code: 0}
			}

			switch result.Error {
			case "cancelled":
				ui.PrintMessagePanel("Context Addition Cancelled", result.Message)
				return ExitError{Code: 130}
			case "already_exists":
				ui.PrintMessagePanel("Context Already Managed", result.Message)
				return ExitError{Code: 1}
			case "no_session":
				ui.PrintMessagePanel("No Active Azure Session", result.Message+"\n\nPlease run 'az login' to authenticate with Azure first.")
				return ExitError{Code: 1}
			case "azure_cli_not_found":
				ui.PrintMessagePanel("Azure CLI Not Found", result.Message+"\n\nPlease install Azure CLI: https://aka.ms/azure-cli")
				return ExitError{Code: 1}
			default:
				ui.PrintMessagePanel("Failed to Add Context", result.Message)
				return ExitError{Code: 1}
			}
		},
	}
}

func newListCommand(manager *service.Manager) *cobra.Command {
	var verbose bool

	command := &cobra.Command{
		Use:   "list",
		Short: "List all saved contexts.",
		RunE: func(cmd *cobra.Command, args []string) error {
			result := manager.ListContexts(verbose)
			if result.Error == "empty_list" {
				ui.PrintMessagePanel("No Saved Contexts", result.Message)
				return ExitError{Code: 0}
			}
			if !result.Success {
				ui.PrintMessagePanel("Error", result.Message)
				return ExitError{Code: 1}
			}

			if verbose {
				fmt.Printf("Saved Contexts (%d):\n\n", len(result.Contexts))
				for _, context := range result.Contexts {
					ui.PrintFieldPanel(context.ContextID, []ui.Field{{Name: "Name", Value: context.ContextName}, {Name: "Subscription", Value: fmt.Sprintf("%s (%s)", context.SubscriptionName, context.SubscriptionID)}, {Name: "Tenant ID", Value: context.TenantID}, {Name: "Username", Value: context.Username}, {Name: "Created", Value: formatTimestamp(context.CreatedAt)}})
				}
				return ExitError{Code: 0}
			}

			ui.PrintContextTable(fmt.Sprintf("Saved Contexts (%d)", len(result.Summaries)), result.Summaries)
			return ExitError{Code: 0}
		},
	}

	command.Flags().BoolVarP(&verbose, "verbose", "v", false, "Show detailed information")
	return command
}

func newDeleteCommand(manager *service.Manager) *cobra.Command {
	return &cobra.Command{
		Use:   "delete",
		Short: "Delete a saved context.",
		RunE: func(cmd *cobra.Command, args []string) error {
			result := manager.DeleteContextInteractive()
			if result.Success {
				ui.PrintFieldPanel(fmt.Sprintf("Successfully Deleted Context: %s", result.ContextName), []ui.Field{{Name: "Context ID", Value: result.ContextID}, {Name: "Context Name", Value: result.ContextName}}, "", "Note: Your active Azure CLI session remains unchanged.")
				return ExitError{Code: 0}
			}

			switch result.Error {
			case "cancelled":
				ui.PrintMessagePanel("Deletion Cancelled", result.Message)
				return ExitError{Code: 130}
			case "empty_list":
				ui.PrintMessagePanel("No Contexts to Delete", result.Message)
				return ExitError{Code: 0}
			case "not_found":
				ui.PrintMessagePanel("Context Not Found", result.Message)
				return ExitError{Code: 1}
			default:
				ui.PrintMessagePanel("Failed to Delete Context", result.Message)
				return ExitError{Code: 1}
			}
		},
	}
}

func managerCurrentAccount(manager *service.Manager) (*azure.Account, error) {
	return manager.GetCurrentAccount()
}

func joinWithComma(values []string) string {
	if len(values) == 0 {
		return ""
	}

	result := values[0]
	for _, value := range values[1:] {
		result += ", " + value
	}
	return result
}

func formatTimestamp(value time.Time) string {
	return value.Format("2006-01-02 15:04:05.999999999")
}