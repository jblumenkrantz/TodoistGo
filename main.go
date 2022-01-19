package main

import "fmt"

func main() {
	tc := NewTodoistClient("0da0034f8290c454a4416502d1baadeeb871e8e9")
	tasks, err := tc.getTasks()
	if err != nil {
		fmt.Printf("Error %s", err)
	}
	for _, task := range tasks {
		fmt.Println(task.ID)
		cerr := tc.closeTask(task)
		if cerr != nil {
			fmt.Printf("Error: %s", err)
		}
	}
}
