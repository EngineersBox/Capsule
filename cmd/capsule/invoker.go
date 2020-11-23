package main

// Invocable ... Abstract generic function that can be invoked
type Invocable func(...interface{}) error

// Call ... Invoke a function with a panic condition on erroring
//
// Params:
// - Invocable inv: Function to invoke
// - ...InvParams params: Parameters to pass to given function
func call(inv Invocable, params ...interface{}) {
	handleErrors(inv(params))
}

func handleErrors(err error) {
	if err != nil {
		panic(err)
	}
}
