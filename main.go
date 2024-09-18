package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"text/tabwriter"
)

// RunNetshCommand executes netsh command with the given arguments
func RunNetshCommand(args ...string) (string, error) {
	cmd := exec.Command("netsh", args...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("command error: %w - %s", err, string(output))
	}
	return string(output), nil
}

// EnableFirewall enables the firewall
func EnableFirewall() {
	fmt.Println("Enabling firewall...")
	_, err := RunNetshCommand("advfirewall", "set", "allprofiles", "state", "on")
	if err != nil {
		log.Fatalf("Error enabling firewall: %v", err)
	}
	fmt.Println("Firewall enabled.")
}

// DisableFirewall disables the firewall
func DisableFirewall() {
	fmt.Println("Disabling firewall...")
	_, err := RunNetshCommand("advfirewall", "set", "allprofiles", "state", "off")
	if err != nil {
		log.Fatalf("Error disabling firewall: %v", err)
	}
	fmt.Println("Firewall disabled.")
}

// AllowPort allows traffic on the specified port and direction (default is in)
func AllowPort(port string, direction string) {
	if direction == "" {
		direction = "in" // Default to inbound if no direction is specified
	}
	fmt.Printf("Allowing %sbound traffic on port %s...\n", direction, port)
	_, err := RunNetshCommand("advfirewall", "firewall", "add", "rule",
		"name=AllowInboundPort"+port, "protocol=TCP", "dir="+direction, "action=allow", "localport="+port)
	if err != nil {
		log.Fatalf("Error allowing port %s: %v", port, err)
	}
	fmt.Printf("%sbound traffic allowed on port %s.\n", strings.Title(direction), port)
}

// DenyPort denies traffic on the specified port and direction (default is in)
func DenyPort(port string, direction string) {
	if direction == "" {
		direction = "in" // Default to inbound if no direction is specified
	}
	fmt.Printf("Denying %sbound traffic on port %s...\n", direction, port)
	_, err := RunNetshCommand("advfirewall", "firewall", "add", "rule",
		"name=DenyInboundPort"+port, "protocol=TCP", "dir="+direction, "action=block", "localport="+port)
	if err != nil {
		log.Fatalf("Error denying port %s: %v", port, err)
	}
	fmt.Printf("%sbound traffic denied on port %s.\n", strings.Title(direction), port)
}


// DeletePortRule removes a rule for the specified port
func DeletePortRule(port string) {
	fmt.Printf("Deleting rule for port %s...\n", port)
	_, err := RunNetshCommand("advfirewall", "firewall", "delete", "rule",
		"name=AllowInboundPort"+port)
	// Also try to delete deny rule
	_, _ = RunNetshCommand("advfirewall", "firewall", "delete", "rule",
		"name=DenyInboundPort"+port)
	if err != nil {
		log.Fatalf("Error deleting rule for port %s: %v", port, err)
	}
	fmt.Printf("Rule for port %s deleted.\n", port)
}

// SimpleListRules lists only the program name, port, and direction of firewall rules


func SimpleListRules() {
	// Set up a tabwriter for column alignment (removing AlignRight to allow left alignment for Program Name)
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)

	// Header
	fmt.Fprintln(w, "Port\tDirection\\Rule Name\t")
	fmt.Fprintln(w, strings.Repeat("-", 50)) // Separator line

	output, err := RunNetshCommand("advfirewall", "firewall", "show", "rule", "name=all")
	if err != nil {
		log.Fatalf("Error listing firewall rules: %v", err)
	}

	lines := strings.Split(output, "\n")
	var name, port, direction string

	// Regular expressions to extract relevant fields
	reName := regexp.MustCompile(`(?i)^Rule Name:\s*(.*)`)
	rePort := regexp.MustCompile(`(?i)^LocalPort:\s*(\d+)`)
	reDirection := regexp.MustCompile(`(?i)^Direction:\s*(In|Out)`)

	for _, line := range lines {
		line = strings.TrimSpace(line)

		if matches := reName.FindStringSubmatch(line); len(matches) == 2 {
			name = matches[1]
		}
		if matches := rePort.FindStringSubmatch(line); len(matches) == 2 {
			port = matches[1]
		}
		if matches := reDirection.FindStringSubmatch(line); len(matches) == 2 {
			direction = matches[1]

			// Print the formatted output using tabs for alignment: Port, Direction, Program Name (left-aligned)
			fmt.Fprintf(w, "%s\t%s\t%s\t\n", port, direction, name)

			// Reset for next rule
			name, port, direction = "", "", ""
		}
	}

	// Flush the tabwriter to ensure output is printed
	w.Flush()
}



// FirewallStatus shows the current firewall status
func FirewallStatus() {
	fmt.Println("Firewall status:")
	output, err := RunNetshCommand("advfirewall", "show", "allprofiles")
	if err != nil {
		log.Fatalf("Error showing firewall status: %v", err)
	}
	fmt.Println(output)
}

// Usage provides help for using the CLI
func Usage() {
	fmt.Println("Usage:")
	fmt.Println("  fw enable                 - Enable the firewall")
	fmt.Println("  fw disable                - Disable the firewall")
	fmt.Println("  fw allow <port> [in|out]  - Allow traffic on a port (default: in)")
	fmt.Println("  fw deny <port> [in|out]   - Deny traffic on a port (default: in)")
	fmt.Println("  fw delete <port>          - Delete the rule for a port")
	fmt.Println("  fw status                 - Show the current firewall status and rules")
	fmt.Println("  fw list                   - List rules (Program Name | Port | Direction)")
	os.Exit(1)
}

func main() {
	if len(os.Args) < 2 {
		Usage()
	}

	// Command line arguments
	switch os.Args[1] {
	case "enable":
		EnableFirewall()
	case "disable":
		DisableFirewall()
	case "allow":
		if len(os.Args) < 3 || len(os.Args) > 4 {
			Usage()
		}
		// If no direction is specified, it defaults to inbound
		direction := "in"
		if len(os.Args) == 4 {
			direction = os.Args[3] // Optional direction argument
		}
		AllowPort(os.Args[2], direction)
	case "deny":
		if len(os.Args) < 3 || len(os.Args) > 4 {
			Usage()
		}
		// If no direction is specified, it defaults to inbound
		direction := "in"
		if len(os.Args) == 4 {
			direction = os.Args[3] // Optional direction argument
		}
		DenyPort(os.Args[2], direction)
	case "delete":
		if len(os.Args) != 3 {
			Usage()
		}
		DeletePortRule(os.Args[2])
	case "status":
		FirewallStatus()
	case "list":
		SimpleListRules()
	default:
		Usage()
	}
}
