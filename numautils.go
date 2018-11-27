package numautils

import (
       "errors"
       "io/ioutil"
       "os"
       "path"
       "strconv"
       "strings"

       "code.cloudfoundry.org/bytefmt"
)

// Checks directory is available or not
func ExistsDir(path string) (bool) {
    if _, err := os.Stat(path); os.IsNotExist(err) {
        return false
    }
    return true
}

// Lists the directories in particular path
func ListDir(dir string) ([]string, error) {
    var dirs []string
    fileInfo, err := ioutil.ReadDir(dir)
    if err != nil {
        return nil, err
    }

    for _, file := range fileInfo {
        fileName := file.Name()
        if file.IsDir() || (file.Mode()&os.ModeSymlink == os.ModeSymlink) {
            updatedDirPath := path.Join(dir, fileName)
            if ExistsDir(updatedDirPath) {
                dirs = append(dirs, updatedDirPath)
            }
        }
    }
    return dirs, nil
}

// Lists only the NUMA node directories
func GetNumaNodeDirs() ([]string, error) {
    var numaNodeDirs []string
    numaNodePath := "/sys/devices/system/node/"
    if _, err := os.Stat(numaNodePath); os.IsNotExist(err) {
        return nil, err
    } else {
        dirs, err := ListDir(numaNodePath)
        if err != nil {
            return nil, err
        }
        for _, dir := range dirs {
            baseName := path.Base(dir)
            if strings.HasPrefix(baseName, "node") {
                numaNodeDirs = append(numaNodeDirs, dir)
            }
         }
         return numaNodeDirs, nil
    }
}

// Gets the total memory for each NUMA nodes
func GetNodesMemoryInfo() (map[int]string, error){
     var ram = map[int]string{}
     dirs, err := GetNumaNodeDirs()
     if err != nil {
         return nil, err
     }
     for _, numaNodeDir := range dirs {
          baseNumaNodeDir := path.Base(numaNodeDir)
          if !strings.HasPrefix(baseNumaNodeDir, "node") {
              continue
          }
          NumaNodeID, err := strconv.Atoi(baseNumaNodeDir[4:])
          if err != nil {
              return nil, err
          }
          memInfoFileName := path.Join(numaNodeDir, "meminfo")
          memInfo, err := ioutil.ReadFile(memInfoFileName)
          if err!=nil {
              return nil, err
          }
          lines := strings.Split(string(memInfo), "\n")
          value := ""
          for _, line := range lines {
              if strings.Contains(line, "MemTotal") {
                  value = strings.Trim(strings.Split(line, ":")[1], " ")
              }
          }
          bytesVal, err := bytefmt.ToBytes(strings.Replace(value, " ", "", 1))
          if err!=nil {
              return nil, err
          }
          ram[NumaNodeID] = bytefmt.ByteSize(bytesVal)
     }
     return ram, nil
}

// Core defines Core ID and  thread siblingsinformation
type CPUInfo struct {
        CPU            int
        ThreadSiblings []int
}

func GetNodesCoresInfo() (map[int][]*CPUInfo, error){
     var cpus = map[int][]*CPUInfo{}
     dirs, err := GetNumaNodeDirs()
     if err != nil {
         return nil, err
     }
     for _, numaNodeDir := range dirs {
          var cpusInfo = []*CPUInfo{}
          var cores = map[int][]int{}
          baseNumaNodeDir := path.Base(numaNodeDir)
          if !strings.HasPrefix(baseNumaNodeDir, "node") {
              continue
          }
          NumaNodeID, err := strconv.Atoi(baseNumaNodeDir[4:])
          if err != nil {
              return nil, err
          }
          threadDirs, err := ListDir(numaNodeDir)
          if err != nil {
              return nil, err
          }
          for _, threadDir := range threadDirs {
               baseThreadDir := path.Base(threadDir)
               if !strings.HasPrefix(baseThreadDir, "cpu") {
                   continue
               }
               threadID, err := strconv.Atoi(baseThreadDir[3:])
               if err!=nil {
                   return nil, err
               }
               cpuFileName := path.Join(threadDir, "topology", "core_id")
               cpuData, err := ioutil.ReadFile(cpuFileName)
               if err!=nil {
                   return nil, err
               }
               cpuID, _ := strconv.Atoi(strings.TrimSpace(string(cpuData[:])))
               cores[cpuID] = append(cores[cpuID], threadID)
          }
          for cpuID, threads := range cores {
                c := &CPUInfo{
                        CPU:            cpuID,
                        ThreadSiblings: threads,
                }
                cpusInfo = append(cpusInfo, c)
          }
          cpus[NumaNodeID] = cpusInfo
     }
     return cpus, nil
}

// Gets NICs info for each NUMA nodes
func GetNodesNicsInfo() (map[int][]string, error){
     var nics = map[int][]string{}
     nicDevicePath := "/sys/class/net/"
     if ExistsDir(nicDevicePath) {
         nicDirs, err := ListDir(nicDevicePath)
         if err != nil {
             return nil, nil
             // return nil, err
         }
         for _, dir := range nicDirs {
             if !ExistsDir(path.Join(dir, "device")) {
                 continue
             }
             nicInfoFileName := path.Join(dir, "device", "numa_node")
             nicInfo, err := ioutil.ReadFile(nicInfoFileName)
             if err != nil {
                 continue
                 // return nil, err
             }
             baseNicDir := path.Base(dir)
             numaNode, err := strconv.Atoi(strings.TrimSpace(string(nicInfo)))
             if err != nil && numaNode < 0{
                 continue
                 // return nil, err
             }
             nics[numaNode] =  append(nics[numaNode], baseNicDir)
         }
     }

     return nics, nil
}

// NUMATopology defines NUMA topology information
type NUMATopology struct {
        NUMA       int64
        RAM        string
        NICs       []string
        CPUs       []*CPUInfo
}

func GetNumaTopology() ([]*NUMATopology, error) {
     var numaTopology []*NUMATopology
     ramList, err := GetNodesMemoryInfo()
     if err != nil {
          return nil, errors.New("Unable to determine memory details.")
     }
     // fmt.Println(ramList)

     cpus, err := GetNodesCoresInfo()
     if err != nil {
          return nil, errors.New("Unable to determine CPUs details.")
     }

     nics, err := GetNodesNicsInfo()
     if err != nil {
          return nil, errors.New("Unable to determine NICs details.")
     }

     for node, ram := range ramList {
         r := &NUMATopology{
              NUMA:       int64(node),
              RAM:        ram,
              NICs:       nics[node],
              CPUs:       cpus[node],
         }
         numaTopology = append(numaTopology, r)
     }
     return numaTopology, nil
}
