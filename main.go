package main

import "wcurl/app"

func main() {
	app.Run()
}

// args := os.Args[1:]
// command := flag.String("command", "", "Pass curl command")
// if len(args) > 0 && args[0] == "curl" {
// 	command := strings.Join(args[1:], " ")
//
// 	re := regexp.MustCompile(`^(?:https?:\/\/)?[^\/]+(\/.*)`)
//
// 	matches := re.FindAllStringSubmatch(command, -1)
// 	fmt.Println(matches)
// 	fmt.Println(matches[0][1])
// }
