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

// HandleErrors ... Invoke a panic call if an error is thrown
func handleErrors(err error) {
	if err != nil {
		panic(err)
	}
}

// HandledInvocationGroup ... Handle errors for an arbitrary size invocable group
func handledInvocationGroup(throwables ...error) {
	for i, err := range throwables {
		handleErrors(err)
	}
}
