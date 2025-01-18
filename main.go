package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/khatibomar/virtual-consensus/virtuallog"
)

func main() {
	if err := realMain(); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
}

func realMain() error {
	menuOptions := []string{"append", "check tail", "read", "seal", "reconfigure", "print", "exit"}

	virtualLog := virtuallog.NewVirtualLog[string]()

	var option int
menu:
	for {
		fmt.Println("Choose an option:")
		for i, option := range menuOptions {
			fmt.Printf("%d. %s\n", i+1, option)
		}
		if _, err := fmt.Scanln(&option); err != nil {
			return err
		}

		switch option {
		case 1:
			fmt.Println("Enter the value to append:")
			reader := bufio.NewReader(os.Stdin)
			var value string
			value, err := reader.ReadString('\n')
			if err != nil {
				fmt.Println(err)
			}
			value = strings.TrimSuffix(value, "\n")
			pos, err := virtualLog.Append(value)
			if err != nil {
				fmt.Println(err)
				continue
			}
			fmt.Printf("Appended at position %d\n", pos)
		case 2:
			fmt.Printf("Tail: %v\n", virtualLog.CheckTail())
		case 3:
			fmt.Println("Enter the start and end positions:")
			var start, end int64
			if _, err := fmt.Scanln(&start, &end); err != nil {
				fmt.Println(err)
				break menu
			}
			values, err := virtualLog.ReadNext(start, end)
			if err != nil {
				fmt.Println(err)
				continue
			}
			fmt.Println("Values:")
			for _, value := range values {
				fmt.Printf("\t%s\n", value)
			}
		case 4:
			virtualLog.Seal()
			fmt.Println("Sealed")
		case 5:
			err := virtualLog.Reconfigure()
			if err != nil {
				fmt.Println(err)
			} else {
				fmt.Println("Reconfigured")
			}
		case 6:
			fmt.Println(virtualLog)
		case 7:
			break menu
		default:
			fmt.Println("Invalid option")
		}
	}
	return nil
}
