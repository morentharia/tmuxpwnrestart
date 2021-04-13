package main

import (
	"encoding/base64"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/pkg/errors"
)

const GDBWindowname = "GDB"

func command(cmd string) (string, error) {
	fmt.Printf("$ %s\n", cmd)
	out, err := exec.Command("bash", "-c", cmd).Output()
	if err != nil {
		log.Printf("command('%s')", err)
		return "", errors.WithMessagef(err, "command('%s')", cmd)
	}
	log.Printf("%s", string(out))
	return string(out), nil
}

func E(s string) string {
	res := fmt.Sprintf("%#v", s)
	res = strings.ReplaceAll(res, "$", "\\$")
	// res = strings.ReplaceAll(res, "\\\\", "\\")
	return res
}

func main() {
	// log.Printf(" = %+v\n", E("azzzzz $esp  "))
	log.Printf("cmd = %+v\n", strings.Join(os.Args[1:], " "))
	log.Printf("cmd = %+v\n", strings.Join(os.Args[1:], " "))
	cmd := strings.Join(os.Args[1:], " ")
	// command("tmux splitw -h " + E(strings.Join(os.Args[1:], " ")))
	// out, _ := command("tmux splitw -h " + E(strings.Join(os.Args[1:], " ")))
	// log.Printf("out = %#v\n", out)

	SessionName, err := command("tmux display-message -p '#S'")
	if err != nil {
		log.Print(err)
	}
	SessionName = strings.TrimSuffix(SessionName, "\n")

	command(fmt.Sprintf("tmux switch -t %s ", E(fmt.Sprintf("%s:$", SessionName))))

	WindowName, err := command("tmux display-message -p '#W'")
	if err != nil {
		log.Fatal("error:", err)
	}
	WindowName = strings.TrimSuffix(WindowName, "\n")
	log.Printf("WindowName = %#v\n", WindowName)

	if strings.HasPrefix(cmd, "SEND ") {
		if WindowName == GDBWindowname {
			cmd := cmd[len("SEND "):]

			data, err := base64.StdEncoding.DecodeString(cmd)
			if err != nil {
				log.Fatal("error:", err)
			}

			fmt.Printf("%q\n", data)

			log.Printf("cmd = %#v\n", string(data))
			command(fmt.Sprintf("tmux send-keys -t %s %s C-m", E(fmt.Sprintf("%s:$", SessionName)), E(string(data))))

		}
	} else {
		if WindowName == GDBWindowname {
			// command(fmt.Sprintf("tmux kill-window -t %s", E(fmt.Sprintf("%s:$", SessionName))))
			command(fmt.Sprintf("tmux send-keys -t %s %s C-m", E(fmt.Sprintf("%s:$", SessionName)), E("quit")))
			command(fmt.Sprintf("tmux send-keys -t %s %s C-m", E(fmt.Sprintf("%s:$", SessionName)), E("quit")))
		} else {
			command(fmt.Sprintf("tmux new-window -d -n %s", GDBWindowname))
			command(fmt.Sprintf("tmux switch -t %s ", E(fmt.Sprintf("%s:$", SessionName))))
		}

		log.Printf("cmd = %#v\n", cmd)

		// command(fmt.Sprintf("tmux send-keys -t %s %s C-m", E(fmt.Sprintf("%s:$", SessionName)), E(fmt.Sprintf("sudo %s", cmd))))
		command(fmt.Sprintf("tmux send-keys -t %s %s C-m", E(fmt.Sprintf("%s:$", SessionName)), E(cmd)))
	}

}
