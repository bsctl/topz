package main

import (
        "log"
        "net/http"
        "sync"
        "text/tabwriter"
        "time"
        "fmt"

        "github.com/shirou/gopsutil/process"
)

type ProcInfo struct {
        PID           int32
        MemoryPercent float32
        MemoryInfo    *process.MemoryInfoStat
        CPUPercent    float64
        Command       string
}

func handleError(res http.ResponseWriter, err error) {
        res.WriteHeader(http.StatusInternalServerError)
        res.Write([]byte(err.Error()))
}

func handleTopz(res http.ResponseWriter, req *http.Request) {
        pids, err := process.Pids()
        if err != nil {
                handleError(res, err)
                return
        }
        procs := make([]*ProcInfo, len(pids))
        wg := sync.WaitGroup{}
        wg.Add(len(pids))
        for ix := range pids {
                proc, err := process.NewProcess(pids[ix])
                if err != nil {
                        handleError(res, err)
                        return
                }

                go func(i int) {
                        var err error
                        p := &ProcInfo{}
                        p.PID = pids[i]
                        if p.Command, err = proc.Cmdline(); err != nil {
                                log.Printf("Error getting Command Line: %v", err)
                        }
                        if p.MemoryInfo, err = proc.MemoryInfo(); err != nil {
                                log.Printf("Error getting memory info: %v", err)
                        }
                        if p.MemoryPercent, err = proc.MemoryPercent(); err != nil {
                                log.Printf("Error getting Memory: %v", err)
                        }
                        if p.CPUPercent, err = proc.Percent(100 * time.Millisecond); err != nil {
                                log.Printf("Error getting CPU: %v", err)
                        }
                        if len(p.Command) > 0 {
                                procs[i] = p
                        }
                        wg.Done()
                }(ix)
        }
        wg.Wait()

        res.WriteHeader(http.StatusOK)
        w := tabwriter.NewWriter(res, 0, 0, 1, ' ', 0)
        fmt.Fprintf(w, "%s\t%s\t%s\t%s\n", "PID", "CPU %", "Memory %", "Command")
        for ix := range procs {
                proc := procs[ix]
                if proc == nil {
                        continue
                }
                fmt.Fprintf(w, "%d\t%g\t%g\t%s\n", proc.PID, proc.CPUPercent, proc.MemoryPercent, proc.Command)
        }
        w.Flush()
}


func main() {
        http.HandleFunc("/topz", handleTopz)
        log.Fatal(http.ListenAndServe(":1080", nil))
}