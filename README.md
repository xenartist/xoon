# xoon - master x1 like a pro

Simple GUI for X1/Solana/SVM Hardcore Operations

![GitHub Repo stars](https://img.shields.io/github/stars/xenartist/xoon?style=flat)
 ![GitHub forks](https://img.shields.io/github/forks/xenartist/xoon?style=flat)
 ![GitHub Downloads (all assets, all releases)](https://img.shields.io/github/downloads/xenartist/xoon/total) ![GitHub License](https://img.shields.io/github/license/xenartist/xoon)

## Features

- Terminal-based GUI interface
- X1/Solana/SVM validator management
- Support for both mainnet and testnet
- Easy configuration and monitoring
- Mouse support for all operations

## Binary Installation

The easiest way to get started is to download the pre-built binary.

### Linux (x64)

```bash
# Download the latest release
wget https://github.com/xenartist/xoon/releases/download/vX.Y.Z/xoon-linux-x64-vX.Y.Z.tar.gz

# Extract
tar zxvf xoon-linux-x64-vX.Y.Z.tar.gz

# Enter directory
cd xoon-x.y.z

# Run
./xoon
```

## Building from Source

If you prefer to build from source, follow these steps:

### Prerequisites

1. Install Rust:
```bash
curl --proto '=https' --tlsv1.2 -sSf https://sh.rustup.rs | sh
```

2. Install development dependencies (Ubuntu/Debian):
```bash
sudo apt update
sudo apt install -y build-essential pkg-config libssl-dev
sudo apt install -y libncurses5-dev libncursesw5-dev
```

For other Linux distributions, install the equivalent packages using your package manager.

### Build and Run

1. Clone the repository:
```bash
git clone https://github.com/xenartist/xoon.git
cd xoon
```

2. Build:
```bash
cargo build --debug
```
or
```bash
cargo build --release
```

3. Run:
```bash
cargo run --debug
```
or
```bash
cargo run --release
```

## License

[GPL-3.0](https://github.com/xenartist/xoon/blob/main/LICENSE)

## Follow me on X

[xen_artist](https://x.com/xen_artist)

