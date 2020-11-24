package main

// ErrorHandlerFunc ... Function to process an error
type ErrorHandlerFunc func(error)

// InfoHandlerFunc ... Function to process information logging
type InfoHandlerFunc func(string, ...interface{})

// Invoker ... Call handler
type Invoker struct {
	errorHandler ErrorHandlerFunc
	infoHanlder  InfoHandlerFunc
}

// Invocable ... Abstract generic function that can be invoked
type Invocable func(...interface{}) error

// Call ... Invoke a function with a panic condition on erroring
//
// Params:
// - Invocable inv: Function to invoke
// - ...InvParams params: Parameters to pass to given function
func (i *Invoker) call(inv Invocable, params ...interface{}) {
	i.handleErrors(inv(params))
}

// HandleErrors ... Invoke a panic call if an error is thrown
func (i *Invoker) handleErrors(err error) {
	if err != nil {
		i.errorHandler(err)
	}
}

// HandledInvocationGroup ... Handle errors for an arbitrary size invocable group
func (i *Invoker) handledInvocationGroup(throwables ...error) {
	for _, err := range throwables {
		i.handleErrors(err)
	}
}
