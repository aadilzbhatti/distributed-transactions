package dtrans

import (
	"bufio"
	"fmt"
	"node"
	"os"
	"regexp"
)

var currentId int32 = 0

func Start() {
	go node.Start()

	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Starting transaction interface")

	const usage = `Operations:
  BEGIN
  SET <server>.<key> <value>
  GET <server>.<key>
  COMMIT
  ABORT`

	fmt.Println(usage)
	r, _ := regexp.Compile(`(BEGIN)|(SET) (.*)\.(.+) (.*)|(GET) (.*)\.(.*)|(COMMIT)|(ABORT)`)
	for {
		text, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Error reading from stdin:", err)
			os.Exit(1)
		}
		if r.MatchString(text) {
			res := r.FindStringSubmatch(text)
			for i := range res {
				if i > 0 && res[i] != "" {
					go runCommand(res, i)
					break
				}
			}
		} else {
			fmt.Println("Error: Could not interpret input")
			fmt.Println(usage)
		}
	}
}

func runCommand(cmds []string, i int) {
	if cmds[i] == "BEGIN" {
		err, tid := Begin()
		if err != nil {
			fmt.Println("Cannot begin a transaction")
			return
		}
		currentId = tid
		fmt.Printf("Beginning transaction %v\n", tid)

	} else if cmds[i] == "SET" {
		if currentId == 0 {
			fmt.Println("Error: Must begin transaction before calling SET")
			return
		}
		err := Set(cmds[i+1], cmds[i+2], cmds[i+3], currentId)
		if err != nil {
			fmt.Println("Could not set:", err)
			currentId = 0
			return
		}
		fmt.Printf("SETTING %v.%v = %v\n", cmds[i+1], cmds[i+2], cmds[i+3])

	} else if cmds[i] == "GET" {
		if currentId == 0 {
			fmt.Println("Error: Must begin transaction before calling GET")
			return
		}
		res, err := Get(cmds[i+1], cmds[i+2], currentId)
		if err != nil {
			fmt.Println("Could not get:", err)
			currentId = 0
			return
		}
		if res != "" {
			fmt.Printf("%v.%v = %v\n", cmds[i+1], cmds[i+2], res)
		}

	} else if cmds[i] == "COMMIT" {
		if currentId == 0 {
			fmt.Println("Error: Must begin transaction before calling COMMIT")
			return
		}
		err := Commit()
		if err != nil {
			fmt.Println("Error in commit:", err)
			return
		}
		fmt.Println("OK")
		currentId = 0

	} else if cmds[i] == "ABORT" {
		if currentId == 0 {
			fmt.Println("Error: Must begin transaction before calling ABORT")
			return
		}
		err := Abort()
		if err != nil {
			fmt.Println("Error in abort:", err)
			return
		}

		fmt.Println("OK")
		currentId = 0

	} else {
		fmt.Println("Error: Invalid command")
	}
}
