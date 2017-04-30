package dtrans

import (
  "os/exec"
  "fmt"
  "strings"
)

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
  for i := 2; i < 10; i++ {
		name := fmt.Sprintf(host, i)
		if name == hname {
			// return string(rune('A' + (i - 2)))
      return string(i)
		}
	}
  return ""
}
