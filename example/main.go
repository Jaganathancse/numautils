package main

import(
       "fmt"
       "github.com/Jaganathancse/numautils"
)

// Gets NUMA topology details
func main() {
    numa, err:= numautils.GetNumaTopology()
    if err == nil {
        for node, info := range numa { 
             fmt.Println("NUMA: ", node)
             fmt.Println("RAM: ", info.RAM)
             fmt.Println("NICs: ", info.NICs)
             fmt.Println("CPUs: ")
             for _, cpu := range info.CPUs {
                fmt.Println(cpu.CPU)
                fmt.Println(cpu.ThreadSiblings)
             }
        }
    } else {
          fmt.Println(err)
    }
}

