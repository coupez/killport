<p align="center">
<img src="https://github.com/coupez/killport/blob/main/logo.png?raw=true" alt="killport" />
</p>

# killport

A cross-platform utility to kill processes running on specific ports.

## Features

- ✅ Terminates processes using specified TCP ports
- 🖥️ Works on multiple platforms:
  - **macOS** (using lsof)
  - **Linux** (using lsof)
  - **Windows** (using netstat and taskkill)
- 🔢 Supports multiple ports in a single command
- 🛡️ Safely handles errors and invalid inputs
- 💬 Provides clear success/failure feedback

## Installation

```bash
go install github.com/coupez/killport@latest
```

## Usage

### Basic Usage

Kill the process running on a specific port:

```bash
killport 8080
```

### Multiple Ports

You can specify multiple ports in a single command:

```bash
killport 3000 443 8000
```

### Exit Codes

- `0`: All processes successfully terminated
- `1`: Failed to terminate at least one process (with details printed)

## Examples

```bash
# Kill process on default HTTP port
killport 80

# Kill processes for a typical web development setup
killport 3000 8080 5432

# Kill process for HTTP and HTTPS
killport 80 443
```

## How It Works

- On **macOS** and **Linux**: Uses `lsof` to find processes and `kill` to terminate them
- On **Windows**: Uses `netstat` to find processes and `taskkill` to terminate them

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.
