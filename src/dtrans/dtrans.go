package dtrans

import (
  "bufio"
  "fmt"
  "os"
  "regexp"
  "node"
)

var currentId int32

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
          runCommand(res, i)
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
      os.Exit(1)
    }
    currentId = tid
    fmt.Printf("Beginning transaction %v\n", tid)

  } else if cmds[i] == "SET" {
    err := Set(cmds[i + 1], cmds[i + 2], cmds[i + 3], currentId)
    if err != nil {
      fmt.Println("Could not set: ", err)
      os.Exit(1)
    }
    fmt.Printf("SETTING %v.%v = %v\n", cmds[i + 1], cmds[i + 2], cmds[i + 3])

  } else if cmds[i] == "GET" {
    // get cmds[i + 1].cmds[i + 2]
    fmt.Printf("GETTING %v.%v\n", cmds[i + 1], cmds[i + 2])
  } else if cmds[i] == "COMMIT" {
    // commit_transaction
    fmt.Println("COMMITTING")
  } else if cmds[i] == "ABORT" {
    // abort_transaction
    fmt.Println("ABORTING")
  } else {
    fmt.Println("Error: Invalid command")
  }
}
