package main

// InvParams ... Generic parameters for an invocable function
type InvParams interface{}

// Invocable ... Abstract generic function that can be invoked
type Invocable func(...InvParams) error

// Call ... Invoke a function with a panic condition on erroring
//
// Params:
// - Invocable inv: Function to invoke
// - ...InvParams params: Parameters to pass to given function
func call(inv Invocable, params ...InvParams) {
	if err := inv(params); err != nil {
		panic(err)
	}
}
