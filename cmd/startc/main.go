package main

import (
    "github.com/abdulbahajaj/startc/pkg/namespaces"
)


func main() {
    desc := namespaces.Desc{
        Mount: true,
        User: true,
        Pid: true,
        Uts: true,
        Ipc: false,
        Network: true,
        Cgroup: false,
        MountProc: true,
        Cmd: "/bin/sh",
        MountPath: "/home/ubuntu/projects/mount-points/newroot",
    }


    namespaces.Create(desc)
}
