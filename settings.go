package main

import (
	"encoding/json"
	"fmt"
	"os"
)

type Setting struct {
	ServerHost, ServerPort, PgHost, PgPort, PgUser, PgPass, PgBase, Data, Assets, HTML string
}

var cfg Setting

func init() {
	// Открыть файл конфигурации
	file, e := os.Open("settings.cfg")
	if e != nil {
		fmt.Println(e.Error())
		panic("Не удалось открыть файл конфигурации")
	}
	defer file.Close()

	//	Прочитать статистику о файле
	stat, e := file.Stat()
	if e != nil {
		fmt.Println(e.Error())
		panic("Не удалось прочитать статистику о файле конфигурации")
	}

	//	Читаем из файла, пишем в массив байт
	readByte := make([]byte, stat.Size())
	_, e = file.Read(readByte)
	if e != nil {
		fmt.Println(e.Error())
		panic("Не удалось прочитать из файла конфигурации ")
	}

	//	Делаем из строки объект
	e = json.Unmarshal(readByte, &cfg)
	if e != nil {
		fmt.Println(e.Error())
		panic("Не удалось сделать Unmarshal ")
	}
}
