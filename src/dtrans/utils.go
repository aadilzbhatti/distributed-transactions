package dtrans

import (
  "os/exec"
  "os"
  "fmt"
  "strings"
  "strconv"
)

var host string = "sp17-cs425-g26-0%d.cs.illinois.edu"

func getHostName() string {
	out, err := exec.Command("hostname").Output()
	if err != nil {
		fmt.Println("Failed to obtain hostname")
		return ""
	}
	return strings.TrimSpace(string(out))
}

func getNodeId() string {
  hname := getHostName()
  for i := 1; i < 10; i++ {
		name := fmt.Sprintf(host, i)
		if name == hname {
      return strconv.Itoa(i)
		}
	}
  os.Exit(1)
  return ""
}
