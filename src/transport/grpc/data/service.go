package Case

import "golang.org/x/net/context"

type Service interface {
	Case(ctx context.Context, _ string, _ int64) (context.Context, string, error)
}

type CaseRequest struct {
	A string
	B int64
}

type CaseResponse struct {
	Ctx context.Context
	V   string
}
