package main

import(
       "fmt"
        "github.com/Jaganathancse/numautils"
)

// Gets NUMA topology details
func main() {
    cpus, _:= numautils.GetNodesCoresInfo()
    fmt.Println(cpus)
    ram, _ := numautils.GetNodesMemoryInfo()
    fmt.Println(ram)
    nics,_ := numautils.GetNodesNicsInfo()
    fmt.Println(nics)
}

