package error

// CodeToUserMessages is for example
var CodeToUserMessages = map[Code]string{
	RequestNotValid: "Invalid Request. Please try again.",
	SqlInsertError:  "Something went wrong. Please try again.",
	SqlUpdateError:  "Something went wrong. Please try again.",
	SqlFetchError:   "Something went wrong. Please try again.",
}
