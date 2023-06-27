package main

import (
	"fmt"
	"regexp"
	"strings"
)

func main() {
	str := `vtysh -c "show ip route"`

	// Split by space but respect quoted strings as single element
	re := regexp.MustCompile(`[^\s"']+|"([^"]*)"|'([^']*)`)
	parts := re.FindAllString(str, -1)

	// Remove quotes from quoted strings
	for i, part := range parts {
		if strings.HasPrefix(part, `"`) && strings.HasSuffix(part, `"`) {
			parts[i] = strings.Trim(part, `"`)
		}
	}

	fmt.Println(parts[2]) // Output: [vtysh -c show ip route]
}
