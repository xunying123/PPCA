# 按程序分流

代理服务器在接受连接时, 可以获取到连接对方的 TCP 端口号。在 Linux 系统下, 通过读取 procfs 中 /proc/[pid]/fd/* 和 /proc/[pid]/fdinfo/* 的信息, 可以获取到每个进程打开的所有文件描述符 (file descriptor) 的信息, 进而可以通过遍历进程的方式来确定哪一进程发起了这个连接。`lsof` 命令的 `-i` 选项就实现了这一功能。

在读取到连接对方的 PID 之后, 可以通过 /proc/[pid]/exe 读取到连接对方所执行的命令, 通过 /proc/[pid]/cmdline 读取到连接对方的命令行参数, 进而进行分流。

非常不建议在除 Linux 外的任何操作系统上完成此功能。