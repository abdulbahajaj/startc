package main

import (
    // "fmt"
    // "log"
    "github.com/abdulbahajaj/startc/pkg/namespaces"
)

func main() {
    desc := namespaces.Desc{
        Mount: true,
        Uts: false,
        Ipc: false,
        Network: false,
        Pid: true,
        Cgroup: false,
        User: false,
        Fork: true,
        MountProc: true,
        Cmd: "bash",
    }

    namespaces.Create(desc)
}
