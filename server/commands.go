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

	RespNoImplement = &resp.Resp{
		Type:  resp.ErrorResp,
		Error: "no implement",
	}

	RespOk = &resp.Resp{
		Type:   resp.SimpleString,
		Status: "OK",
	}

	RespNil = &resp.Resp{
		Type: resp.BulkResp,
		Bulk: nil,
	}
)

func RespErr(err error) *resp.Resp {
	return &resp.Resp{
		Type:  resp.ErrorResp,
		Error: err.Error(),
	}
}

func generalSetPathVal(r *resp.Resp, client *session, fn func(string, string, interface{}) error) *resp.Resp {
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

	err = fn(string(k), path, val)
	if err != nil {
		if err == ErrNoSuchKey {
			return RespNil
		}
		log.Warning(err)
		return RespErr(err)
	}

	return RespOk
}

func generalGetPathVal(r *resp.Resp, client *session, fn func(string, string) (interface{}, error)) *resp.Resp {
	if len(r.Multi) != 3 {
		return RespInvalidParam
	}

	k, err := r.Key()
	if err != nil {
		log.Warning(err)
		return RespErr(err)
	}

	path := string(r.Multi[2].Bulk)
	ret, err := fn(string(k), path)
	if err != nil {
		if err == ErrNoSuchKey {
			return RespNil
		}
		log.Warning(err)
		return RespErr(err)
	}
	if ret == nil {
		return RespNil
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
		return RespNil
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
	return generalSetPathVal(r, client, client.srv.db.PutPath)
}

func cmdJGet(r *resp.Resp, client *session) *resp.Resp {
	return generalGetPathVal(r, client, client.srv.db.GetPath)
}

func cmdJPush(r *resp.Resp, client *session) *resp.Resp {
	return generalSetPathVal(r, client, client.srv.db.PushPath)
}

func cmdJPop(r *resp.Resp, client *session) *resp.Resp {
	return generalGetPathVal(r, client, client.srv.db.PopPath)
}

func cmdJIncr(r *resp.Resp, client *session) *resp.Resp {
	return generalSetPathVal(r, client, client.srv.db.IncrPath)
}

func cmdSave(r *resp.Resp, client *session) *resp.Resp {
	return RespNoImplement
}

func cmdBgSave(r *resp.Resp, client *session) *resp.Resp {
	return RespNoImplement
}
