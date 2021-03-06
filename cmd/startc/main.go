package main

import (
    "github.com/abdulbahajaj/startc/pkg/namespaces"
)


func main() {
    desc := namespaces.Desc{
        Cgroup: true,
        Ipc: true,
        Mount: true,
        User: true,
        Pid: true,
        Uts: true,
        Network: true,
        MountProc: true,

        CPUShare: 100,
        MemoryLimit: 2000,

        Cmd: "/bin/sh",
        MountPath: "/home/ubuntu/projects/mount-points/newroot",
    }
    namespaces.Create(desc)
}
