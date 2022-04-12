package util

type Handler struct {
	IO     PackageIO
	Status int
}

const (
	STATUS_ESTABLISHED = iota
	STATUS_LOGINED
	STATUS_DISCONNECTED
)
