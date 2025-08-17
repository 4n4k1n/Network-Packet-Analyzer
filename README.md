# Network Packet Analyzer

A command-line network packet analyzer built in Go that captures and analyzes network traffic in real-time. This tool provides detailed statistics about network packets including protocol breakdown, IP address activity, and connection patterns.

## Features

- **Real-time packet capture** from network interfaces
- **Protocol analysis** with support for TCP, UDP, and other protocols
- **Traffic statistics** including packet counts, rates, and protocol distribution
- **IP address tracking** to identify most active sources and destinations
- **Flexible filtering** by IP address, protocol, and port
- **Configurable capture duration** and network interface selection

## Installation

### Prerequisites

- Go 1.24.5 or later
- libpcap development libraries
- Root/administrator privileges (required for packet capture)

### System Dependencies

**Ubuntu/Debian:**
```bash
sudo apt-get install libpcap-dev
```

**macOS (with Homebrew):**
```bash
brew install libpcap
```

**CentOS/RHEL/Fedora:**
```bash
sudo yum install libpcap-devel    # CentOS/RHEL
sudo dnf install libpcap-devel    # Fedora
```

### Build from Source

```bash
git clone <repository-url>
cd Network-Packet-Analyzer
go mod download
go build -o analyzer ./src
```

## Usage

### Basic Usage

```bash
# Capture packets for 30 seconds on default interface (wlan0)
sudo ./analyzer

# Capture on specific interface for 60 seconds
sudo ./analyzer -device eth0 -time 60
```

### Command Line Options

| Flag | Description | Default | Example |
|------|-------------|---------|---------|
| `-device` | Network interface to capture from | `wlan0` | `-device eth0` |
| `-time` | Capture duration in seconds | `30` | `-time 120` |
| `-ip` | Filter by specific IP address | None | `-ip 192.168.1.1` |
| `-protocol` | Filter by protocol (tcp/udp) | None | `-protocol tcp` |
| `-port` | Filter by port number | None | `-port 80` |

### Examples

```bash
# Capture HTTP traffic only
sudo ./analyzer -protocol tcp -port 80

# Monitor specific IP for 5 minutes
sudo ./analyzer -ip 8.8.8.8 -time 300

# Capture UDP DNS traffic
sudo ./analyzer -protocol udp -port 53

# Combine filters: HTTPS traffic to/from specific IP
sudo ./analyzer -ip 192.168.1.100 -protocol tcp -port 443
```

## Sample Output

### Real-time Packet Display
```
src IP           dest IP          protocol  src port    dest port
192.168.1.100    8.8.8.8         UDP       54234       53
192.168.1.100    172.217.16.142  TCP       45678       443
10.0.0.1         192.168.1.100   TCP       80          54321
```

### Capture Summary
```
=== Capture Summary ===
Total packets captured  : 1247
Capture duration        : 30
Average rate            : 41.6
Unique source IPs       : 12
Unique destination IPs  : 8
Protocol breakdown
   TCP   : 891 packets (71.5%)
   UDP   : 298 packets (23.9%)
   other : 58 packets (4.6%)
Most active source      : 192.168.1.100 (423 packets)
Most active destination : 8.8.8.8 (156 packets)
```

## Project Structure

```
Network-Packet-Analyzer/
├── README.md           # This file
├── ROADMAP.md         # Development roadmap and learning guide
├── go.mod             # Go module dependencies
├── go.sum             # Dependency checksums
├── analyzer           # Compiled binary
└── src/               # Source code
    ├── main.go        # Main program entry point
    ├── types.go       # Data structure definitions
    ├── parse.go       # Command-line argument parsing
    ├── getters.go     # Packet data extraction
    ├── prints.go      # Output formatting and display
    └── utils.go       # Utility functions
```

## Code Overview

### Core Components

- **main.go**: Program entry point, packet capture loop, and statistics collection
- **types.go**: Defines data structures for parsed arguments, packet data, and statistics
- **parse.go**: Handles command-line argument parsing and BPF filter creation
- **getters.go**: Extracts IP addresses, protocols, and port information from packets
- **prints.go**: Formats and displays packet information and capture statistics
- **utils.go**: Helper functions for data processing

### Key Data Structures

```go
// Packet information
type Data struct {
    src_ip   string
    dst_ip   string
    protocol layers.IPProtocol
    tcp      *layers.TCP
    udp      *layers.UDP
    src_port uint16
    dst_port uint16
}

// Capture statistics
type Stats_data struct {
    total_packets     int
    captured_duration int
    average_rate      float32
    unique_src_ids    int
    unique_dst_ids    int
    tcp_packets       int
    udp_packets       int
    other_packets     int
    most_active_src   string
    most_active_dst   string
    src_ip_counts     map[string]int
    dst_ip_counts     map[string]int
}
```

## Dependencies

- **github.com/google/gopacket**: Packet capture and analysis library
- **github.com/dariubs/percent**: Percentage calculation utilities

## Security Considerations

- Requires root privileges for raw socket access
- Only captures packet headers, not payload data
- Filters can be applied to limit captured traffic
- No data is stored persistently by default

## Troubleshooting

### Common Issues

**"Permission denied" error:**
```bash
# Run with sudo
sudo ./analyzer
```

**"No such device" error:**
```bash
# List available interfaces
ip link show
# or
ifconfig -a

# Use correct interface name
sudo ./analyzer -device eth0
```

**"No packets captured":**
- Verify network interface is active and has traffic
- Check if other packet capture tools work: `sudo tcpdump -i <interface>`
- Try capturing on a different interface
- Remove filters to capture all traffic

## Development

### Building for Development
```bash
# Run without building binary
sudo go run ./src -device eth0 -time 10

# Build with debug information
go build -gcflags="-N -l" -o analyzer-debug ./src
```

### Testing
```bash
# Generate test traffic
ping google.com &
curl -s http://example.com > /dev/null &

# Capture the test traffic
sudo ./analyzer -time 5
```

## Contributing

1. Fork the repository
2. Create a feature branch: `git checkout -b feature-name`
3. Make your changes
4. Test thoroughly with different network conditions
5. Submit a pull request

## License

[Add your license information here]

## Related Projects

- [Wireshark](https://www.wireshark.org/) - Full-featured network protocol analyzer
- [tcpdump](https://www.tcpdump.org/) - Command-line packet analyzer
- [gopacket](https://github.com/google/gopacket) - Go library for packet processing

## Roadmap

See [ROADMAP.md](ROADMAP.md) for detailed development plans and learning milestones.