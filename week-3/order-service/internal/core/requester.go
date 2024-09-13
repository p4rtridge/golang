package core

import "context"

type ReqType string

const KeyRequester ReqType = "requester"

type Requester interface {
	GetSubject() string
	GetTokenID() string
}

type requesterData struct {
	Sub string `json:"user_id"`
	Tid string `json:"tid"`
}

func NewRequester(sub, tid string) Requester {
	return &requesterData{
		Sub: sub,
		Tid: tid,
	}
}

func (r *requesterData) GetSubject() string {
	return r.Sub
}

func (r *requesterData) GetTokenID() string {
	return r.Tid
}

func GetRequester(ctx context.Context) Requester {
	if requester, ok := ctx.Value(KeyRequester).(Requester); ok {
		return requester
	}

	return nil
}

func ContextWithRequester(ctx context.Context, requester Requester) context.Context {
	return context.WithValue(ctx, KeyRequester, requester)
}
