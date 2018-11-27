# NUMA Utils

NUMA utils is used to get the NUMA topology details like RAM, CPU Cores,
thread siblings and NICs data for each NUMA node.

Currently supports for RHEL, fedora and centos environments.

### Installation

```
cd $GOPATH
mkdir -p ./src/github.com/Jaganathancse/numautils
git clone https://github.com/Jaganathancse/numautils.git ./src/github.com/Jaganathancse/numautils
govendor fetch github.com/cloudfoundry/bytefmt
```

### Usage

```
import "github.com/Jaganathancse/numautils"
numaInfo, err := numautils.GetNumaTopology()
```

### Example

NUMATopology:
```
  - CPUs:
      - CPU: 0
        ThreadSiblings:
          - 0
          - 4
      - CPU: 2
        ThreadSiblings:
          - 2
          - 6
    NICs:
      - ens1
      - ens2
    NUMA: 0
    RAM: 3.9G

  - CPUs:
      - CPU: 1
        ThreadSiblings:
          - 1
          - 5
      - CPU: 3
        ThreadSiblings:
          - 3
          - 7
    NICs:
      - ens3
      - ens4
    NUMA: 1
    RAM: 2.5G
```
