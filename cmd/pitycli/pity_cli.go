package main

import (
	"github.com/chzyer/readline"
	"io"
	"log"
	"strings"
)

func usage(w io.Writer) {
	io.WriteString(w, "commands:\n")
	io.WriteString(w, completer.Tree("    "))
}

var completer = readline.NewPrefixCompleter(
	readline.PcItem("exit"),
	readline.PcItem("help"),
)

var defaultPrompt = "pitydb-cli\033[32m»\033[0m "

func main() {
	l, err := readline.NewEx(&readline.Config{
		Prompt:          defaultPrompt,
		HistoryFile:     "/tmp/readline.tmp",
		AutoComplete:    completer,
		InterruptPrompt: "^C",
		EOFPrompt:       "exit",

		HistorySearchFold: true,
	})
	if err != nil {
		panic(err)
	}
	defer l.Close()

	log.SetOutput(l.Stderr())
	var cmds []string
	for {
		line, err := l.Readline()
		if err != nil {
			break
		}
		line = strings.TrimSpace(line)
		if len(line) == 0 {
			continue
		}
		cmds = append(cmds, line)
		if !strings.HasSuffix(line, ";") {
			l.SetPrompt(">>> ")
			continue
		}
		cmd := strings.Join(cmds, " ")
		cmds = cmds[:0]
		l.SetPrompt(defaultPrompt)
		l.SaveHistory(cmd)
		println(cmd)

		//line, err := l.Readline()
		//if err == readline.ErrInterrupt {
		//	if len(line) == 0 {
		//		break
		//	} else {
		//		continue
		//	}
		//} else if err == io.EOF {
		//	break
		//}
		//
		//line = strings.TrimSpace(line)
		//switch {
		//case line == "help":
		//	usage(l.Stderr())
		//case line == "exit":
		//	goto exit
		//case line == "":
		//default:
		//	log.Println("Unknown:", strconv.Quote(line))
		//}
	}
}
