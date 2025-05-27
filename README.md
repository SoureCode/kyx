
# kyx

This tools provides a simple way to start, stop, reset and deploy symfony based applications.

## Installation

For a linux amd64 system, you can use the following commands to install `kyx`:

```bash
sudo mkdir -p /opt/kyx
curl -L https://github.com/SoureCode/kyx/releases/latest/download/kyx_linux_amd64.tar.gz | sudo tar xz --no-same-owner -C /opt/kyx
sudo ln -s /opt/kyx/kyx /usr/local/bin/kyx
kyx --help
```