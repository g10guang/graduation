package tools

// use NoCopy as sub-struct
// `go vet` can detect and fail
type NoCopy struct {

}

func (*NoCopy) Lock() {}

func (*NoCopy) Unlock() {}

