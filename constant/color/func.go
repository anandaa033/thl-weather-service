package color

import "fmt"

// const (
// 	RED    = "\033[31m"
// 	GREEN  = "\033[32m"
// 	YELLOW = "\033[33m"
// 	BLUE   = "\033[34m"
// 	PURPLE = "\033[35m"
// 	WHITE  = "\033[37m"
// )

// const (
// 	printColor = "\033[38;5;%dm%s\033[39;49m\n"
// )

// func all() {
// 	for j := 0; j < 256; j++ {
// 		fmt.Printf(printColor, j, "Hello!")
// 	}
// }

func Red(text string) string {
	return fmt.Sprintf("\033[38;5;%dm%s\033[39;49m", 9, text)
}

func Green(text string) string {
	return fmt.Sprintf("\033[38;5;%dm%s\033[39;49m", 10, text)
}

func Yellow(text string) string {
	return fmt.Sprintf("\033[38;5;%dm%s\033[39;49m", 11, text)
}

func Blue(text string) string {
	return fmt.Sprintf("\033[38;5;%dm%s\033[39;49m", 12, text)
}

func Purple(text string) string {
	return fmt.Sprintf("\033[38;5;%dm%s\033[39;49m", 13, text)
}

func CadetBlue(text string) string {
	return fmt.Sprintf("\033[38;5;%dm%s\033[39;49m", 14, text)
}

func White(text string) string {
	return fmt.Sprintf("\033[38;5;%dm%s\033[39;49m", 15, text)
}
