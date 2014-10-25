package server

import (
	"encoding/json"
	"jj/resp"

	log "github.com/ngaut/logging"
)

var (
	RespNoSuchCmd = &resp.Resp{
		Type:  resp.ErrorResp,
		Error: "unknown command",
	}

	RespInvalidParam = &resp.Resp{
		Type:  resp.ErrorResp,
		Error: "invalid parameter",
	}

	RespInnerError = &resp.Resp{
		Type:  resp.ErrorResp,
		Error: "inner error",
	}

	RespOk = &resp.Resp{
		Type:   resp.SimpleString,
		Status: "OK",
	}
)

func cmdJdocSet(r *resp.Resp, client *session) (*resp.Resp, error) {
	if len(r.Multi) != 3 {
		return RespInvalidParam, nil
	}

	k, err := r.Key()
	if err != nil {
		log.Warning(err)
		return nil, err
	}

	var val interface{}
	err = json.Unmarshal(r.Multi[2].Bulk, &val)
	if err != nil {
		return &resp.Resp{
			Type:  resp.ErrorResp,
			Error: err.Error(),
		}, nil
	}

	client.srv.db.PutDoc(string(k), val)

	return RespOk, nil
}

func cmdJdocGet(r *resp.Resp, client *session) (*resp.Resp, error) {
	k, err := r.Key()
	if err != nil {
		log.Warning(err)
		return nil, err
	}

	val, _ := client.srv.db.GetDoc(string(k))
	b, err := json.Marshal(val)
	if err != nil {
		return &resp.Resp{
			Type:  resp.ErrorResp,
			Error: err.Error(),
		}, nil
	}

	ret := &resp.Resp{
		Type: resp.BulkResp,
		Bulk: b,
	}

	return ret, nil
}

func cmdJSet(r *resp.Resp, client *session) (*resp.Resp, error) {
	return RespOk, nil
}

func cmdJGet(r *resp.Resp, client *session) (*resp.Resp, error) {
	return RespOk, nil
}
