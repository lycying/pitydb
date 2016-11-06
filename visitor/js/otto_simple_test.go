package js_test

import (
	"errors"
	"fmt"
	"github.com/robertkrimen/otto"
	"os"
	"testing"
	"time"
)

var halt = errors.New("Stahp")

func runUnsafe(unsafe string) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		if caught := recover(); caught != nil {
			if caught == halt {
				fmt.Fprintf(os.Stderr, "Some code took to long! Stopping after: %v\n", duration)
				return
			}
			panic(caught) // Something else happened, repanic!
		}
		fmt.Fprintf(os.Stderr, "Ran code successfully: %v\n", duration)
	}()

	vm := otto.New()
	vm.Interrupt = make(chan func(), 1) // The buffer prevents blocking

	go func() {
		time.Sleep(2 * time.Second) // Stop after two seconds
		vm.Interrupt <- func() {
			panic(halt)
		}
	}()

	vm.Run(unsafe) // Here be dragons (risky code)
}

func TestSimpleOtto(t *testing.T) {
	vm := otto.New()
	vm.Run(`
    abc = 2 + 2;
    console.log("The value of abc is " + abc); // 4
`)
	if value, err := vm.Get("abc"); err == nil {
		if value_int, err := value.ToInteger(); err == nil {
			fmt.Printf("%v,%v", value_int, err)
		}
	}

	vm.Set("xyzzy", "Nothing happens.")
	vm.Run(`
    console.log(xyzzy.length); // 16
`)

	value, _ := vm.Run("xyzzy.length")
	{
		// value is an int64 with a value of 16
		value, _ := value.ToInteger()
		fmt.Println(value)
	}

	value, err := vm.Run("abcdefghijlmnopqrstuvwxyz.length")
	if err != nil {
		// err = ReferenceError: abcdefghijlmnopqrstuvwxyz is not defined
		// If there is an error, then value.IsUndefined() is true
		println(err.Error())
	}

	vm.Set("sayHello", func(call otto.FunctionCall) otto.Value {
		fmt.Printf("Hello, %s.\n", call.Argument(0).String())
		return otto.Value{}
	})

	vm.Set("twoPlus", func(call otto.FunctionCall) otto.Value {
		right, _ := call.Argument(0).ToInteger()
		result, _ := vm.ToValue(2 + right)
		return result
	})

	result, _ := vm.Run(`
    sayHello("Xyzzy".split(""));      // Hello, Xyzzy.
    sayHello();             // Hello, undefined

    result = twoPlus(2.0); // 4
`)
	fmt.Printf("%v", result)

	runUnsafe(`var abc = [];`)
	runUnsafe(`
    while (true) {
        // Loop forever
    }`)

}
