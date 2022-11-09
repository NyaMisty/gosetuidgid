package main // import "github.com/tianon/gosetuidgid"

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"syscall"
	"text/template"
)

func init() {
	// make sure we only have one process and that it runs on the main thread (so that ideally, when we Exec, we keep our user switches and stuff)
	runtime.GOMAXPROCS(1)
	runtime.LockOSThread()
}

func version() string {
	return fmt.Sprintf(`%s (%s on %s/%s; %s)`, Version, runtime.Version(), runtime.GOOS, runtime.GOARCH, runtime.Compiler)
}

func usage() string {
	t := template.Must(template.New("usage").Parse(`
Usage: {{ .Self }} uid gids command [args]
   eg: {{ .Self }} 1000 1000 sh
       {{ .Self }} 1000 1000,1001,1002 sh

{{ .Self }} version: {{ .Version }}
{{ .Self }} license: Apache-2.0 (full text at https://github.com/tianon/gosetuidgid)
`))
	var b bytes.Buffer
	template.Must(t, t.Execute(&b, struct {
		Self    string
		Version string
	}{
		Self:    filepath.Base(os.Args[0]),
		Version: version(),
	}))
	return strings.TrimSpace(b.String()) + "\n"
}

func main() {
	log.SetFlags(0) // no timestamps on our logs

	if len(os.Args) <= 3 {
		log.Println(usage())
		os.Exit(1)
	}

	uidStr := os.Args[1]
	uid, err := strconv.Atoi(uidStr)
	if err != nil {
		fmt.Println(usage())
		os.Exit(1)
	}
	gidStr := os.Args[2]
	commands := os.Args[3:]

	gidsStrList := strings.Split(gidStr, ",")
	gids := make([]int, 0)
	for _, gid := range gidsStrList {
		gidNum, err := strconv.Atoi(gid)
		if err != nil {
			fmt.Printf("Unknown gid: %s\n", gidNum)
			os.Exit(1)
		}
		gids = append(gids, gidNum)
	}

	if err := syscall.Setgroups(gids); err != nil {
		fmt.Printf("Setgroups failed: %v\n", err)
	}
	if err := syscall.Setgid(gids[0]); err != nil {
		fmt.Printf("Setgid failed: %v\n", err)
	}
	if err := syscall.Setuid(uid); err != nil {
		fmt.Printf("Setuid failed: %v\n", err)
	}

	name, err := exec.LookPath(commands[0])
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	if err = syscall.Exec(name, commands, os.Environ()); err != nil {
		log.Fatalf("error: exec failed: %v", err)
	}
}
