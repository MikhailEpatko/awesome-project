package main

import (
	"fmt"
	"os"
)

func main() {
	cwd, err := os.Getwd()
	if err != nil {
		fmt.Println("Ошибка получения текущей рабочей директории:", err)
		return
	}

	fmt.Println("Текущая рабочая директория:", cwd)
}
