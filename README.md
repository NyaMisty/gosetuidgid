# gosetuidgid

This is a simple tool that based on [gosetuidgid](https://github.com/tianon/gosetuidgid), which supports setting uid(setuid), gid(setgid), and supplementary gid (i.e. sgid, setgroup)

This tool does not handle user to uid map, and you must specify uid/gid in number.

The first GID specified will be the new gid of process 

```console
$ gosetuidgid
Usage: ./gosetuidgid user-spec command [args]
   eg: ./gosetuidgid tianon bash
       ./gosetuidgid nobody:root bash -c 'whoami && id'
       ./gosetuidgid 1000:1 id

./gosetuidgid version: 1.1 (go1.3.1 on linux/amd64; gc)
```

Once the user/group is processed, we switch to that user, then we `exec` the specified process and `gosetuidgid` itself is no longer resident or involved in the process lifecycle at all.  This avoids all the issues of signal passing and TTY, and punts them to the process invoking `gosetuidgid` and the process being invoked by `gosetuidgid`, where they belong.

# Why

The main point here is to provide a better setuidgid implementation on Android.
- The setuidgid in busybox cannot handle supplementary gid, but in Android various permission are actually expressed by supplementary groups, such as 3003(inet).
- Unlike Linux, you don't various utils in android, so a static-compiled go setuidgid with such feature will be extremely useful.