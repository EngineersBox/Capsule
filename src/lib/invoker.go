package capsule

// ErrorHandlerFunc ... Function to process an error
type ErrorHandlerFunc func(error)

// InfoHandlerFunc ... Function to process information logging
type InfoHandlerFunc func(string, ...interface{})

// Handler ... Call handler
type Handler struct {
	errorHandler ErrorHandlerFunc
	infoHanlder  InfoHandlerFunc
}

// Handlable ... Abstract generic function that can be invoked
type Handlable func(...interface{}) error

// Call ... Invoke a function with a panic condition on erroring
//
// Params:
// - Invocable inv: Function to invoke
// - ...interface{} params: Parameters to pass to given function
func (h *Handler) Call(inv Handlable, params ...interface{}) {
	h.HandleErrors(inv(params))
}

// HandleErrors ... Invoke a panic call if an error is thrown
func (h *Handler) HandleErrors(err error) {
	if err != nil {
		h.errorHandler(err)
	}
}

// HandledInvocationGroup ... Handle errors for an arbitrary size handled group
func (h *Handler) HandledInvocationGroup(throwables ...error) {
	for _, err := range throwables {
		h.HandleErrors(err)
	}
}

// AsyncHandledInvocationGroup ... Asynchronously handle errors for an arbitrary size invocable group
func (h *Handler) AsyncHandledInvocationGroup(throwables ...error) {
	for _, err := range throwables {
		go func(err error) {
			h.HandleErrors(err)
		}(err)
	}
}
