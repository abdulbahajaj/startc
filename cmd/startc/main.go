package main

import (
    // "fmt"
    // "log"
    "github.com/abdulbahajaj/startc/pkg/namespaces"
    // "github.com/moby/moby/pkg/reexec"
    // "os/exec"
    // "path/filepath"
    // "os"
    // "fmt"
)


func main() {
    desc := namespaces.Desc{
        Mount: true,
        User: true,
        Pid: true,
        Uts: false,
        Ipc: false,
        Network: true,
        Cgroup: false,
        MountProc: true,
        Cmd: "/bin/sh",
    }
    namespaces.Create(desc)
}
