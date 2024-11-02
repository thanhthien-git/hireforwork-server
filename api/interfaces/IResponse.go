package interfaces

type IResponse[T any] struct {
	Doc T `json:"doc"`
}
