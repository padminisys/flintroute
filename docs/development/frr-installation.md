# FRR Installation & Testing Guide

## Table of Contents
- [Overview](#overview)
- [FRR Installation on Debian 12](#frr-installation-on-debian-12)
- [FRR Configuration](#frr-configuration)
- [gRPC Northbound API Setup](#grpc-northbound-api-setup)
- [Testing with Containerlab](#testing-with-containerlab)
- [Testing with GNS3](#testing-with-gns3)
- [Testing with FRR in Docker](#testing-with-frr-in-docker)
- [Testing with Mininet](#testing-with-mininet)
- [Sample Network Topologies](#sample-network-topologies)
- [Multi-Node BGP Scenarios](#multi-node-bgp-scenarios)
- [Troubleshooting](#troubleshooting)

---

## Overview

This guide covers installing FRR stable on Debian 12 and setting up various testing environments for FlintRoute development. We focus on **Containerlab** as the primary testing tool, with alternatives for different use cases.

### Testing Tools Comparison

| Tool | Best For | Complexity | Resource Usage | Setup Time |
|------|----------|------------|----------------|------------|
| **Containerlab** | Modern, container-based labs | Low | Low | 5 min |
| **GNS3** | Full network simulation | High | High | 30 min |
| **FRR in Docker** | Simple, quick tests | Very Low | Very Low | 2 min |
| **Mininet** | SDN/OpenFlow testing | Medium | Medium | 15 min |

---

## FRR Installation on Debian 12

### Add FRR Repository

```bash
# Install prerequisites
sudo apt-get update
sudo apt-get install -y \
    curl \
    gnupg \
    lsb-release \
    apt-transport-https

# Add FRR GPG key
curl -s https://deb.frrouting.org/frr/keys.asc | sudo apt-key add -

# Add FRR stable repository
FRRVER="frr-stable"
echo "deb https://deb.frrouting.org/frr $(lsb_release -s -c) $FRRVER" | \
    sudo tee /etc/apt/sources.list.d/frr.list

# Update package list
sudo apt-get update
```

### Install FRR Packages

```bash
# Install FRR and Python tools
sudo apt-get install -y \
    frr \
    frr-pythontools \
    frr-doc

# Verify installation
frr --version
# Expected output: FRRouting 8.4.x or higher

# Check installed daemons
dpkg -l | grep frr
```

### Enable FRR Daemons

Edit `/etc/frr/daemons`:

```bash
sudo nano /etc/frr/daemons
```

Enable required daemons:

```bash
# FRR daemons configuration
bgpd=yes
ospfd=no
ospf6d=no
ripd=no
ripngd=no
isisd=no
pimd=no
ldpd=no
nhrpd=no
eigrpd=no
babeld=no
sharpd=no
pbrd=no
bfdd=no
fabricd=no
vrrpd=no

# Enable zebra (required)
zebra=yes

# Enable integrated config
vtysh_enable=yes
zebra_options="  -A 127.0.0.1 -s 90000000"
bgpd_options="   -A 127.0.0.1"
ospfd_options="  -A 127.0.0.1"
```

### Start FRR Service

```bash
# Enable FRR service
sudo systemctl enable frr

# Start FRR service
sudo systemctl start frr

# Check status
sudo systemctl status frr

# Verify daemons are running
sudo systemctl status frr@zebra
sudo systemctl status frr@bgpd
```

### Basic FRR Configuration

```bash
# Enter FRR shell
sudo vtysh

# Configure hostname
configure terminal
hostname frr-dev
exit

# Save configuration
write memory
exit
```

---

## FRR Configuration

### Configure Basic Settings

Create `/etc/frr/frr.conf`:

```bash
sudo nano /etc/frr/frr.conf
```

Add basic configuration:

```
! FRR Configuration for FlintRoute Development
!
frr version 8.4
frr defaults traditional
hostname frr-dev
log syslog informational
service integrated-vtysh-config
!
! Enable gRPC
grpc
 address 0.0.0.0
 port 50051
!
! BGP Configuration (example)
router bgp 65000
 bgp router-id 10.0.0.1
 no bgp ebgp-requires-policy
 no bgp network import-check
!
line vty
!
end
```

### Apply Configuration

```bash
# Reload FRR configuration
sudo systemctl reload frr

# Or restart FRR
sudo systemctl restart frr

# Verify configuration
sudo vtysh -c "show running-config"
```

---

## gRPC Northbound API Setup

### Enable gRPC in FRR

Edit `/etc/frr/frr.conf` and add:

```
grpc
 address 0.0.0.0
 port 50051
 no tls
!
```

### Verify gRPC is Running

```bash
# Check if gRPC port is listening
sudo netstat -tlnp | grep 50051

# Or use ss
sudo ss -tlnp | grep 50051

# Expected output:
# tcp   LISTEN 0   128   0.0.0.0:50051   0.0.0.0:*   users:(("bgpd",pid=1234,fd=10))
```

### Install gRPC Tools

```bash
# Install grpcurl for testing
go install github.com/fullstorydev/grpcurl/cmd/grpcurl@latest

# Test gRPC connection
grpcurl -plaintext localhost:50051 list

# Expected output: List of available gRPC services
# frr.Northbound
```

### Test gRPC API

```bash
# List available methods
grpcurl -plaintext localhost:50051 list frr.Northbound

# Get capabilities
grpcurl -plaintext localhost:50051 frr.Northbound.GetCapabilities

# Get configuration
grpcurl -plaintext localhost:50051 \
    -d '{"type": "STATE", "encoding": "JSON"}' \
    frr.Northbound.Get
```

---

## Testing with Containerlab

**Containerlab** is the recommended testing tool for FlintRoute development. It provides fast, container-based network labs.

### Install Containerlab

```bash
# Install Containerlab
bash -c "$(curl -sL https://get.containerlab.dev)"

# Verify installation
containerlab version

# Or install via package manager
sudo apt-get install -y containerlab
```

### Basic Containerlab Topology

Create `lab-basic.yml`:

```yaml
name: flintroute-basic

topology:
  nodes:
    # FRR Router 1
    frr1:
      kind: linux
      image: frrouting/frr:v8.4.0
      ports:
        - "50051:50051"  # gRPC
        - "2601:2605"    # vtysh
      binds:
        - ./configs/frr1:/etc/frr
      env:
        DAEMONS: "bgpd zebra"
    
    # FRR Router 2
    frr2:
      kind: linux
      image: frrouting/frr:v8.4.0
      ports:
        - "50052:50051"
        - "2611:2605"
      binds:
        - ./configs/frr2:/etc/frr
      env:
        DAEMONS: "bgpd zebra"
    
    # FRR Router 3
    frr3:
      kind: linux
      image: frrouting/frr:v8.4.0
      ports:
        - "50053:50051"
        - "2621:2605"
      binds:
        - ./configs/frr3:/etc/frr
      env:
        DAEMONS: "bgpd zebra"

  links:
    # Connect routers
    - endpoints: ["frr1:eth1", "frr2:eth1"]
    - endpoints: ["frr2:eth2", "frr3:eth1"]
    - endpoints: ["frr3:eth2", "frr1:eth2"]
```

### Deploy Containerlab Topology

```bash
# Create config directories
mkdir -p configs/{frr1,frr2,frr3}

# Deploy lab
sudo containerlab deploy -t lab-basic.yml

# Check status
sudo containerlab inspect -t lab-basic.yml

# Access router console
sudo docker exec -it clab-flintroute-basic-frr1 vtysh

# Destroy lab when done
sudo containerlab destroy -t lab-basic.yml
```

### BGP Peering Topology

Create `lab-bgp-peering.yml`:

```yaml
name: flintroute-bgp

topology:
  nodes:
    # Route Server (AS 65000)
    rs1:
      kind: linux
      image: frrouting/frr:v8.4.0
      ports:
        - "50051:50051"
      binds:
        - ./configs/rs1:/etc/frr
      env:
        DAEMONS: "bgpd zebra"
    
    # Peer 1 (AS 65001)
    peer1:
      kind: linux
      image: frrouting/frr:v8.4.0
      ports:
        - "50061:50051"
      binds:
        - ./configs/peer1:/etc/frr
      env:
        DAEMONS: "bgpd zebra"
    
    # Peer 2 (AS 65002)
    peer2:
      kind: linux
      image: frrouting/frr:v8.4.0
      ports:
        - "50062:50051"
      binds:
        - ./configs/peer2:/etc/frr
      env:
        DAEMONS: "bgpd zebra"
    
    # Peer 3 (AS 65003)
    peer3:
      kind: linux
      image: frrouting/frr:v8.4.0
      ports:
        - "50063:50051"
      binds:
        - ./configs/peer3:/etc/frr
      env:
        DAEMONS: "bgpd zebra"

  links:
    - endpoints: ["rs1:eth1", "peer1:eth1"]
    - endpoints: ["rs1:eth2", "peer2:eth1"]
    - endpoints: ["rs1:eth3", "peer3:eth1"]
```

### Configure BGP Peers

**Route Server (configs/rs1/frr.conf):**

```
frr version 8.4
frr defaults traditional
hostname rs1
!
grpc
 address 0.0.0.0
 port 50051
!
router bgp 65000
 bgp router-id 10.0.0.1
 no bgp ebgp-requires-policy
 no bgp network import-check
 !
 neighbor 10.1.1.2 remote-as 65001
 neighbor 10.1.1.2 description Peer1
 !
 neighbor 10.1.2.2 remote-as 65002
 neighbor 10.1.2.2 description Peer2
 !
 neighbor 10.1.3.2 remote-as 65003
 neighbor 10.1.3.2 description Peer3
!
end
```

**Peer 1 (configs/peer1/frr.conf):**

```
frr version 8.4
frr defaults traditional
hostname peer1
!
grpc
 address 0.0.0.0
 port 50051
!
router bgp 65001
 bgp router-id 10.1.1.2
 no bgp ebgp-requires-policy
 !
 neighbor 10.1.1.1 remote-as 65000
 neighbor 10.1.1.1 description RouteServer
 !
 address-family ipv4 unicast
  network 192.168.1.0/24
 exit-address-family
!
end
```

---

## Testing with GNS3

GNS3 provides full network simulation with support for various network devices.

### Install GNS3

```bash
# Add GNS3 repository
sudo add-apt-repository ppa:gns3/ppa
sudo apt-get update

# Install GNS3
sudo apt-get install -y gns3-gui gns3-server

# Install dependencies
sudo apt-get install -y \
    qemu-kvm \
    libvirt-daemon-system \
    libvirt-clients \
    bridge-utils \
    virtinst

# Add user to required groups
sudo usermod -aG kvm,libvirt,docker $USER

# Logout and login for group changes to take effect
```

### Import FRR Appliance

1. Download FRR appliance from GNS3 marketplace
2. Import into GNS3: `File > Import Appliance`
3. Select FRR appliance file

### Create GNS3 Topology

1. Open GNS3
2. Create new project: `flintroute-test`
3. Drag FRR routers onto canvas
4. Connect routers with links
5. Configure each router

### Configure FRR in GNS3

Right-click router > Console:

```bash
# Enter configuration mode
vtysh
configure terminal

# Configure interfaces
interface eth0
 ip address 10.1.1.1/24
 no shutdown
exit

# Configure BGP
router bgp 65000
 bgp router-id 10.0.0.1
 neighbor 10.1.1.2 remote-as 65001
exit

# Enable gRPC
grpc
 address 0.0.0.0
 port 50051
exit

# Save configuration
write memory
exit
```

---

## Testing with FRR in Docker

Simplest method for quick testing.

### Run FRR Container

```bash
# Run single FRR instance
docker run -d \
  --name frr-test \
  --privileged \
  -p 50051:50051 \
  -p 2605:2605 \
  -v $(pwd)/frr.conf:/etc/frr/frr.conf \
  frrouting/frr:v8.4.0

# Access FRR shell
docker exec -it frr-test vtysh

# View logs
docker logs frr-test

# Stop and remove
docker stop frr-test
docker rm frr-test
```

### Docker Compose Setup

Create `docker-compose.yml`:

```yaml
version: '3.8'

services:
  frr1:
    image: frrouting/frr:v8.4.0
    container_name: frr1
    privileged: true
    ports:
      - "50051:50051"
      - "2605:2605"
    volumes:
      - ./configs/frr1.conf:/etc/frr/frr.conf
    networks:
      testnet:
        ipv4_address: 172.20.0.10

  frr2:
    image: frrouting/frr:v8.4.0
    container_name: frr2
    privileged: true
    ports:
      - "50052:50051"
      - "2606:2605"
    volumes:
      - ./configs/frr2.conf:/etc/frr/frr.conf
    networks:
      testnet:
        ipv4_address: 172.20.0.11

networks:
  testnet:
    driver: bridge
    ipam:
      config:
        - subnet: 172.20.0.0/16
```

Start the lab:

```bash
# Start containers
docker-compose up -d

# Check status
docker-compose ps

# Access router
docker-compose exec frr1 vtysh

# Stop containers
docker-compose down
```

---

## Testing with Mininet

Mininet is useful for SDN and OpenFlow testing.

### Install Mininet

```bash
# Install Mininet
sudo apt-get install -y mininet

# Verify installation
sudo mn --version
```

### Create Mininet Topology with FRR

Create `frr_topology.py`:

```python
#!/usr/bin/env python3

from mininet.net import Mininet
from mininet.node import Controller, Host
from mininet.cli import CLI
from mininet.log import setLogLevel, info
from mininet.link import TCLink

def frrTopology():
    """Create a network with FRR routers"""
    
    net = Mininet(controller=Controller, link=TCLink)
    
    info('*** Adding controller\n')
    net.addController('c0')
    
    info('*** Adding hosts\n')
    h1 = net.addHost('h1', ip='10.0.1.10/24')
    h2 = net.addHost('h2', ip='10.0.2.10/24')
    h3 = net.addHost('h3', ip='10.0.3.10/24')
    
    info('*** Adding routers\n')
    r1 = net.addHost('r1', ip='10.0.1.1/24')
    r2 = net.addHost('r2', ip='10.0.2.1/24')
    r3 = net.addHost('r3', ip='10.0.3.1/24')
    
    info('*** Creating links\n')
    net.addLink(h1, r1)
    net.addLink(h2, r2)
    net.addLink(h3, r3)
    net.addLink(r1, r2)
    net.addLink(r2, r3)
    net.addLink(r3, r1)
    
    info('*** Starting network\n')
    net.start()
    
    info('*** Configuring FRR on routers\n')
    # Start FRR daemons on routers
    for router in [r1, r2, r3]:
        router.cmd('sysctl -w net.ipv4.ip_forward=1')
        router.cmd('/usr/lib/frr/zebra -d')
        router.cmd('/usr/lib/frr/bgpd -d')
    
    info('*** Running CLI\n')
    CLI(net)
    
    info('*** Stopping network\n')
    net.stop()

if __name__ == '__main__':
    setLogLevel('info')
    frrTopology()
```

Run the topology:

```bash
# Make executable
chmod +x frr_topology.py

# Run with sudo
sudo ./frr_topology.py
```

---

## Sample Network Topologies

### Topology 1: Simple BGP Peering

```
┌─────────┐         ┌─────────┐
│  FRR1   │─────────│  FRR2   │
│ AS65001 │         │ AS65002 │
└─────────┘         └─────────┘
```

**Use Case**: Basic BGP peer testing

### Topology 2: Route Server

```
        ┌─────────┐
        │   RS    │
        │ AS65000 │
        └────┬────┘
             │
    ┌────────┼────────┐
    │        │        │
┌───▼───┐ ┌──▼───┐ ┌─▼────┐
│ Peer1 │ │Peer2 │ │Peer3 │
│AS65001│ │AS65002│ │AS65003│
└───────┘ └──────┘ └──────┘
```

**Use Case**: IXP route server simulation

### Topology 3: Transit Provider

```
┌─────────┐
│Upstream │
│ AS174   │
└────┬────┘
     │
┌────▼────┐
│  Edge   │
│ AS65000 │
└────┬────┘
     │
┌────▼────┐
│Customer │
│ AS65001 │
└─────────┘
```

**Use Case**: Transit and customer peering

### Topology 4: Full Mesh

```
┌─────────┐         ┌─────────┐
│  FRR1   │─────────│  FRR2   │
│ AS65001 │         │ AS65002 │
└────┬────┘         └────┬────┘
     │                   │
     │    ┌─────────┐    │
     └────│  FRR3   │────┘
          │ AS65003 │
          └─────────┘
```

**Use Case**: Full mesh BGP testing

---

## Multi-Node BGP Scenarios

### Scenario 1: BGP Session Establishment

**Objective**: Test basic BGP session establishment

```bash
# Deploy topology
sudo containerlab deploy -t lab-bgp-peering.yml

# Wait for BGP sessions to establish
sleep 30

# Check BGP summary on route server
sudo docker exec clab-flintroute-bgp-rs1 vtysh -c "show bgp summary"

# Expected output: All peers in Established state
```

### Scenario 2: Route Advertisement

**Objective**: Test route advertisement and reception

```bash
# On peer1, advertise a network
sudo docker exec clab-flintroute-bgp-peer1 vtysh << EOF
configure terminal
router bgp 65001
address-family ipv4 unicast
network 192.168.1.0/24
exit-address-family
exit
exit
EOF

# On route server, verify route is received
sudo docker exec clab-flintroute-bgp-rs1 vtysh -c "show bgp ipv4 unicast"

# Should see 192.168.1.0/24 from peer1
```

### Scenario 3: Peer Down Detection

**Objective**: Test peer down detection and alerting

```bash
# Stop peer1 container
sudo docker stop clab-flintroute-bgp-peer1

# On route server, check BGP summary
sudo docker exec clab-flintroute-bgp-rs1 vtysh -c "show bgp summary"

# Peer1 should show as "Active" or "Connect"

# Restart peer1
sudo docker start clab-flintroute-bgp-peer1

# Wait for session to re-establish
sleep 30

# Verify session is back up
sudo docker exec clab-flintroute-bgp-rs1 vtysh -c "show bgp summary"
```

### Scenario 4: Configuration Changes

**Objective**: Test dynamic configuration changes

```bash
# Add a new peer via gRPC
grpcurl -plaintext localhost:50051 \
    -d '{
      "type": "CANDIDATE",
      "encoding": "JSON",
      "config": {
        "frr-routing:routing": {
          "control-plane-protocols": {
            "control-plane-protocol": [{
              "type": "frr-bgp:bgp",
              "name": "main",
              "frr-bgp:bgp": {
                "neighbors": {
                  "neighbor": [{
                    "remote-address": "10.1.4.2",
                    "remote-as": 65004
                  }]
                }
              }
            }]
          }
        }
      }
    }' \
    frr.Northbound.Commit

# Verify new peer is configured
sudo docker exec clab-flintroute-bgp-rs1 vtysh -c "show running-config"
```

### Scenario 5: Route Filtering

**Objective**: Test route maps and prefix lists

```bash
# Configure route map on route server
sudo docker exec clab-flintroute-bgp-rs1 vtysh << EOF
configure terminal
!
ip prefix-list ALLOWED-PREFIXES seq 5 permit 192.168.0.0/16 le 24
!
route-map FILTER-IN permit 10
 match ip address prefix-list ALLOWED-PREFIXES
!
route-map FILTER-IN deny 20
!
router bgp 65000
 neighbor 10.1.1.2 route-map FILTER-IN in
!
exit
exit
EOF

# Clear BGP session to apply filter
sudo docker exec clab-flintroute-bgp-rs1 vtysh -c "clear bgp 10.1.1.2"

# Verify filtering is working
sudo docker exec clab-flintroute-bgp-rs1 vtysh -c "show bgp ipv4 unicast"
```

---

## Troubleshooting

### FRR Service Issues

```bash
# Check FRR service status
sudo systemctl status frr

# Check individual daemon status
sudo systemctl status frr@zebra
sudo systemctl status frr@bgpd

# View FRR logs
sudo journalctl -u frr -f

# Check FRR configuration syntax
sudo vtysh -C -f /etc/frr/frr.conf
```

### gRPC Connection Issues

```bash
# Verify gRPC is listening
sudo netstat -tlnp | grep 50051

# Test gRPC connection
grpcurl -plaintext localhost:50051 list

# Check firewall rules
sudo iptables -L -n | grep 50051

# Allow gRPC port if blocked
sudo iptables -A INPUT -p tcp --dport 50051 -j ACCEPT
```

### Containerlab Issues

```bash
# Check Containerlab version
containerlab version

# List running labs
sudo containerlab inspect --all

# View container logs
sudo docker logs clab-<lab-name>-<node-name>

# Clean up failed deployments
sudo containerlab destroy --cleanup

# Remove all Containerlab networks
sudo docker network prune
```

### BGP Session Issues

```bash
# Check BGP summary
sudo vtysh -c "show bgp summary"

# Check specific neighbor
sudo vtysh -c "show bgp neighbor 10.1.1.2"

# Debug BGP
sudo vtysh -c "debug bgp neighbor-events"
sudo vtysh -c "debug bgp updates"

# View BGP logs
sudo tail -f /var/log/frr/bgpd.log
```

### Docker Network Issues

```bash
# List Docker networks
docker network ls

# Inspect network
docker network inspect <network-name>

# Check container connectivity
docker exec <container> ping <ip-address>

# Restart Docker networking
sudo systemctl restart docker
```

---

## Performance Testing

### Load Testing BGP Sessions

```bash
# Create 100 BGP peers using Containerlab
# This requires a topology file with 100 peer definitions

# Monitor resource usage
watch -n 1 'docker stats --no-stream'

# Monitor FRR memory usage
ps aux | grep bgpd

# Check BGP performance
sudo vtysh -c "show bgp performance-statistics"
```

### Stress Testing

```bash
# Advertise many routes
for i in {1..1000}; do
    sudo docker exec clab-test-peer1 vtysh << EOF
configure terminal
router bgp 65001
address-family ipv4 unicast
network 10.$((i/256)).$((i%256)).0/24
exit-address-family
exit
exit
EOF
done

# Monitor route processing
sudo docker exec clab-test-rs1 vtysh -c "show bgp memory"
```

---

## Next Steps

1. **Start Testing**: Deploy a basic topology and test FlintRoute integration
2. **Develop Features**: Use test environment for feature development
3. **Run Tests**: Follow [Testing Guide](testing.md) for automated testing
4. **Implement Phase 1**: Review [Implementation Plan](../phase1/implementation-plan.md)

---

## Quick Reference

### Essential Commands

```bash
# Containerlab
sudo containerlab deploy -t <topology.yml>
sudo containerlab inspect -t <topology.yml>
sudo containerlab destroy -t <topology.yml>

# FRR
sudo vtysh
sudo systemctl restart frr
sudo vtysh -c "show running-config"

# gRPC
grpcurl -plaintext localhost:50051 list
grpcurl -plaintext localhost:50051 frr.Northbound.GetCapabilities

# Docker
docker exec -it <container> vtysh
docker logs <container>
docker-compose up -d
```

---

**Last Updated**: 2024-01-15  
**Version**: 0.1.0-alpha