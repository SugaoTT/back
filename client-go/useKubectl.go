package main

import (
	"fmt"
	"os/exec"
)

func main() {

	//scanner := bufio.NewScanner(os.Stdin)
	fmt.Print("input > ")
	/*
		for scanner.Scan() {
			if scanner.Text() == "stop" {
				break
			}
			input := strings.Split(scanner.Text(), " ")

			a := ""
			for _, v := range input {
				a += v + " "
			}
			a = strings.TrimSpace(a)
			fmt.Println(a)
	*/
	//result, _ := exec.Command(" ", input...).Output()
	result, _ := exec.Command("ping", "-c", "3", "127.0.0.1").Output()
	fmt.Println(string(result))
	fmt.Print("input > ")
	//}

}
