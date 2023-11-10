//go:build js && wasm

package main

import (
	conv "diagram-converter/internal/conversion"
	"fmt"
	"syscall/js"
)

func convertExcalidrawToGliffy(input string) string {
	result, err := conv.ConvertExcalidrawToGliffy(input)
	if err != nil {
		return err.Error()
	}
	return result
}

func convertExcalidrawToGliffyWrapper() js.Func {
	return js.FuncOf(func(this js.Value, args []js.Value) any {
		if len(args) != 1 {
			return "Invalid no of arguments passed"
		}
		x := args[0].String()
		// fmt.Printf("Input: %s\n", x)
		return convertExcalidrawToGliffy(x)
	})
}

func main() {
	fmt.Println("WASM loaded")

	js.Global().Set("convertExcalidrawToGliffy", convertExcalidrawToGliffyWrapper())

	<-make(chan bool)
}
