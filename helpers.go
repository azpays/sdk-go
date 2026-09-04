package azpays

// String returns a pointer to the given string value.
// Useful for setting optional string fields in request structs.
//
//	req := &azpays.CreatePaymentRequest{
//	    Description: azpays.String("Order #1234"),
//	}
func String(v string) *string {
	return &v
}

// Float64 returns a pointer to the given float64 value.
func Float64(v float64) *float64 {
	return &v
}

// Int returns a pointer to the given int value.
func Int(v int) *int {
	return &v
}

// Bool returns a pointer to the given bool value.
func Bool(v bool) *bool {
	return &v
}
