
# FW - Simple Firewall Tool

FW is a simple firewall command-line tool written in Go, designed to control the Windows firewall using the `netsh` command. This tool allows you to enable or disable the firewall, allow or deny traffic on specific ports, and list firewall rules in a simplified manner (Program Name, Port, Direction).

## Features

- **Enable/Disable Firewall**: Quickly enable or disable the firewall.
- **Allow/Deny Ports**: Add rules to allow or deny traffic on specific ports.
- **List Rules**: Simplified listing of firewall rules showing only Program Name, Port, and Direction.
- **Cross-Platform (Windows)**: Works specifically with Windows using `netsh`.

## Usage

### Running the FW Tool

You can run FW in different modes:
1. **Enable/Disable Firewall**: Control the state of the Windows firewall.
2. **Allow/Deny Ports**: Add rules to allow or block traffic on specific ports.
3. **List Rules**: List the current firewall rules in a simplified format.
4. **Delete Rules**: Remove firewall rules for a specific port.

### Flags

- `enable`: Enable the firewall for all profiles.
- `disable`: Disable the firewall for all profiles.
- `allow <port> [in|out]`: Allow traffic on the specified port (default: inbound).
- `deny <port> [in|out]`: Block traffic on the specified port (default: inbound).
- `delete <port>`: Delete the firewall rule for the specified port.
- `list`: List current firewall rules (Program Name, Port, Direction).
- `status`: Show the current status of the firewall.

### Examples

#### 1. Enabling the Firewall

To enable the firewall on all profiles:

```bash
go run main.go enable
```

#### 2. Disabling the Firewall

To disable the firewall on all profiles:

```bash
go run main.go disable
```

#### 3. Allowing Traffic on a Port

To allow inbound traffic on port `8080`:

```bash
go run main.go allow 8080
```

To allow outbound traffic on port `8080`:

```bash
go run main.go allow 8080 out
```

#### 4. Denying Traffic on a Port

To block inbound traffic on port `22`:

```bash
go run main.go deny 22
```

#### 5. Deleting a Port Rule

To delete a firewall rule for port `8080`:

```bash
go run main.go delete 8080
```

#### 6. Listing Firewall Rules

To list the current firewall rules (showing Program Name, Port, and Direction):

```bash
go run main.go list
```

#### 7. Firewall Status

To check the current firewall status:

```bash
go run main.go status
```

### Usage Example

```bash
$ go run main.go list
Port        Direction        Program Name
8080        Inbound          AllowInboundPort8080
80          Inbound          AllowInboundPort80
22          Inbound          DenyInboundPort22
```

### How It Works

- **Enabling/Disabling the Firewall**: The tool uses the `netsh advfirewall` command to enable or disable the firewall.
- **Adding/Denying Rules**: It creates rules using the `netsh advfirewall firewall add rule` command for allowing or blocking traffic on specific ports.
- **Listing Rules**: The tool parses the output of the `netsh advfirewall firewall show rule` command to display the relevant firewall rules.

## Supported Platforms

- **Windows**: This tool works on Windows only, as it uses the `netsh` command-line tool to manage the firewall.

## Security Warning

This tool is intended for educational purposes or system administration tasks. Use it responsibly and ensure you have appropriate permissions to modify firewall rules on the system.

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.
