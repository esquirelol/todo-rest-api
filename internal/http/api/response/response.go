package response

type Response struct {
	Status  string
	Message string
}

const (
	statusOk  = "OK"
	statusErr = "Error"
)

func OK(msg string) Response {
	return Response{statusOk, msg}
}

func Error(msg string) Response {
	return Response{statusErr, msg}
}
