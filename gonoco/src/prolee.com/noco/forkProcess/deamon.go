package daemon

/**
 * Golang的syscall.ForkExec()函数需要指定要执行的程序. 不能像C语言一样分叉执行代码.
 * 这时我们可以通过一个小技巧来实现父子进程执行不同的代码, 这个技巧就是通过参数来实现.
 * 我们可以在执行子进程程序时传递一个特有的参数来区分当前进程是否子进程, 例如我们可以传递”--daemon”参数.
 * 因为父进程没有接收到”--daemon”参数, 所以被认为是父进程, 而子进程收到”--daemon”参数, 所以知道是子进程.
 *
 * 在Daemon()函数中, 首先判断是否有”--daemon”参数, 如果有这个参数说明是子进程, 那么就初始化子进程的运行环境, 然后返回.
 * 如果没有”--daemon”参数, 说明是父进程, 那么就调用syscall.ForkExec()函数来执行当前程序.
 * 执行程序之前记得要添加”--daemon”参数给子进程.
 */
import (
    "errors"
    "os"
    "runtime"
    "syscall"
)

const daemonFlagName = "--daemon"

func initDaemonRuntime() {
    // 创建新回话
    _, err := syscall.Setsid()
    if err != nil {
        return
    }
    // 把标准输入输出指向null
    fd, err := os.OpenFile("/dev/null", os.O_RDWR, 0)
    if err != nil {
        return
    }
    _ = syscall.Dup2(int(fd.Fd()), int(os.Stdin.Fd()))
    _ = syscall.Dup2(int(fd.Fd()), int(os.Stdout.Fd()))
    _ = syscall.Dup2(int(fd.Fd()), int(os.Stderr.Fd()))
    if fd.Fd() > os.Stderr.Fd() {
        _ = fd.Close()
    }
}
func Daemon() (int, error) {
    if runtime.GOOS == "windows" {
        return -1, errors.New("unsupported windows operating system")
    }
    isDaemon := false
    for i := 1; i < len(os.Args); i++ {
        if os.Args[i] == daemonFlagName {
            isDaemon = true
        }
    }
    if isDaemon { // daemon process
        initDaemonRuntime()
        return 0, nil
    }
    procPath := os.Args[0]
    // 添加"--daemon"参数
    args := make([]string, 0, len(os.Args)+1)
    args = append(args, os.Args...)
    args = append(args, daemonFlagName)
    attr := &syscall.ProcAttr{
        Env:   os.Environ(),
        Files: []uintptr{os.Stdin.Fd(), os.Stdout.Fd(), os.Stderr.Fd()},
    }
    pid, err := syscall.ForkExec(procPath, args, attr)
    if err != nil {
        return -1, err
    }
    return pid, nil
}
