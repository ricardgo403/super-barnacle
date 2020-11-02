package main

import (
	"fmt"
	"time"
)

var mostrarProcesos bool
var channel chan uint64

func proceso(id uint64) {
	i := uint64(0)
	isActive := true
	for isActive {
		select {
		case myID := <-channel:
			if id == myID {
				isActive = false
				fmt.Println("El proceso con el ID ", id, " ha sido terminado")
				break
			} else {
				channel <- myID
			}
		default:
			if mostrarProcesos {
				fmt.Printf("id %d: %d\n", id, i)
			}
			i = i + 1
		}
		time.Sleep(time.Millisecond * 500)
	}
}

func generadorID() func() uint64 {
	i := uint64(0)
	return func() uint64 {
		var par = i
		i++
		return par
	}
}

func sliceContains(slice []uint64, id uint64) (bool, int) {
	isContained := false
	index := -1
	for i, thisID := range slice {
		if thisID == id {
			isContained = true
			index = i
		}
	}
	return isContained, index
}

func main() {
	mostrarProcesos = false
	channel = make(chan uint64)
	ids := make([]uint64, 0)
	var input string
	nextID := generadorID()
	for {
		fmt.Println("Menu")
		fmt.Println("1) Agregar proceso")
		fmt.Println("2) Mostrar procesos")
		fmt.Println("3) Eliminar proceso")
		fmt.Println("4) Salir")
		fmt.Scanln(&input)
		if input == "1" {
			id := nextID()
			ids = append(ids, id)
			go proceso(id)
		} else if input == "2" {
			mostrarProcesos = !mostrarProcesos
			input = ""
		} else if input == "3" {
			var thisID uint64
			fmt.Print("Ingrese el Id del proceso a matar: ")
			fmt.Scanln(&thisID)
			isContained, index := sliceContains(ids, thisID)
			if isContained {
				channel <- thisID
				ids = append(ids[:index], ids[index+1:]...)
			} else {
				fmt.Println("El proceso con ese ID no existe...")
			}
		} else if input == "4" {
			break
		}
	}
}
