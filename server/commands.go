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

	RespNoSuchKey = &resp.Resp{
		Type:  resp.ErrorResp,
		Error: "no such key",
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

func RespErr(err error) *resp.Resp {
	return &resp.Resp{
		Type:  resp.ErrorResp,
		Error: err.Error(),
	}
}

func cmdJdocSet(r *resp.Resp, client *session) *resp.Resp {
	if len(r.Multi) != 3 {
		return RespInvalidParam
	}

	k, err := r.Key()
	if err != nil {
		log.Warning(err)
		return RespErr(err)
	}

	var val interface{}
	err = json.Unmarshal(r.Multi[2].Bulk, &val)
	if err != nil {
		log.Warning(err)
		return RespErr(err)
	}

	err = client.srv.db.PutDoc(string(k), val)
	if err != nil {
		log.Warning(err)
		return RespErr(err)
	}

	return RespOk
}

func cmdJdocGet(r *resp.Resp, client *session) *resp.Resp {
	k, err := r.Key()
	if err != nil {
		log.Warning(err)
		return RespErr(err)
	}

	val, _ := client.srv.db.GetDoc(string(k))
	if val == nil {
		return RespNoSuchKey
	}
	b, err := json.Marshal(val)
	if err != nil {
		log.Warning(err)
		return RespErr(err)
	}

	return &resp.Resp{
		Type: resp.BulkResp,
		Bulk: b,
	}
}

func cmdJSet(r *resp.Resp, client *session) *resp.Resp {
	if len(r.Multi) != 4 {
		return RespInvalidParam
	}

	k, err := r.Key()
	if err != nil {
		log.Warning(err)
		return RespErr(err)
	}

	path := string(r.Multi[2].Bulk)

	var val interface{}
	err = json.Unmarshal(r.Multi[3].Bulk, &val)
	if err != nil {
		log.Warning(err)
		return RespErr(err)
	}

	err = client.srv.db.PutPath(string(k), path, val)
	if err != nil {
		log.Warning(err)
		return RespErr(err)
	}

	return RespOk
}

func cmdJGet(r *resp.Resp, client *session) *resp.Resp {
	if len(r.Multi) != 3 {
		return RespInvalidParam
	}

	k, err := r.Key()
	if err != nil {
		log.Warning(err)
		return RespErr(err)
	}

	path := string(r.Multi[2].Bulk)
	ret, err := client.srv.db.GetPath(string(k), path)
	if err != nil {
		log.Warning(err)
		return RespErr(err)
	}
	if ret == nil {
		return &resp.Resp{
			Type:  resp.ErrorResp,
			Error: "no such field",
		}
	}

	b, err := json.Marshal(ret)
	if err != nil {
		log.Warning(err)
		return RespErr(err)
	}

	return &resp.Resp{
		Type: resp.BulkResp,
		Bulk: b,
	}
}
