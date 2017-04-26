package node

import (
  "fmt"
  "os/exec"
  "strings"
  "node/coordinator"
  "node/participant"
)

var host string = "sp17-cs425-g26-0%d.cs.illinois.edu"
var nodeId int

func Start() {
  hostname := getHostName()
  for i := 1; i < 10; i++ {
    name := fmt.Sprintf(host, i)
    if name == hostname {
      nodeId = i
      break
    }
  }
  if nodeId == 1 {
    // if id is 1, is Coordinator
    go coordinator.Start()
  } else {
    // otherwise participant
    go participant.Start(hostname, nodeId)
  }
  // handle everything else there
}

func getHostName() string {
	out, err := exec.Command("hostname").Output()
	if err != nil {
		fmt.Println("Failed to obtain hostname")
    return ""
	}
	return strings.TrimSpace(string(out))
}
