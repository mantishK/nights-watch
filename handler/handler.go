package handler

import (
	"database/sql"
	"net/http"
	"plivo/nights-watch/config"

	ae "plivo/nights-watch/apperror"
	"plivo/nights-watch/model"

	"plivo/nights-watch/reply"
)

type smsRequestBody struct {
	AuthId string `json:"auth_id"`
	Src    string `json:"src"`
	Dest   string `json:"dst"`
	Text   string `json:"text"`
}

//validation
func (body *smsRequestBody) OK() *ae.Error {
	if len(body.Src) == 0 {
		return ae.Required("", "src")
	} else if len(body.Dest) == 0 {
		return ae.Required("", "dst")
	} else if len(body.Text) == 0 {
		return ae.Required("", "text")
	}
	return nil
}

//Handler that sends SMS
func SendSMS(w http.ResponseWriter, r *http.Request) {
	reqBody := smsRequestBody{}
	appErr := decode(r, &reqBody)
	if appErr != nil {
		reply.Err(w, appErr)
		return
	}

	//get max count from config
	max, ok := config.GetInt("max_text_count")
	if !ok {
		reply.Err(w, ae.Config())
		return
	}

	//get minimum used no from the auth_id
	if accountID, err := model.GetAccountID(reqBody.AuthId); err != nil {
		reply.Err(w, ae.DB("", err))
		return
	} else if tx, newFrmNo, err := model.GetMinUsedPhoneNo(reqBody.Dest, reqBody.AuthId, accountID, max); err != nil {
		tx.Rollback()
		if err == sql.ErrNoRows {
			reply.Err(w, ae.Forbidden("Numbers exhausted"))
		} else {
			reply.Err(w, ae.DB("", err))
		}
		return
	} else if err := model.IncPhoneUsedCount(newFrmNo, reqBody.AuthId, tx); err != nil {
		tx.Rollback()
		reply.Err(w, ae.DB("", err))
		return
	} else if err := model.FireSMS(newFrmNo, reqBody.Dest, reqBody.AuthId, reqBody.Text); err != nil {
		_ = model.DecPhoneUsedCount(newFrmNo, reqBody.AuthId)
		reply.Err(w, ae.Internal("", err))
		return
	}

	reply.OK(w, "")
}

type stopRequestBody struct {
	From string `json:"from"`
	To   string `json:"to"`
	Text string `json:"text"`
}

//validation
func (body *stopRequestBody) OK() *ae.Error {
	if len(body.From) == 0 {
		return ae.Required("", "from")
	} else if len(body.To) == 0 {
		return ae.Required("", "to")
	} else if len(body.Text) == 0 {
		return ae.Required("", "text")
	} else if body.Text != "STOP" {
		return ae.InvalidInput("Unrecognized Text", "text")
	}
	return nil
}

//handler to block the phone number
func StopNo(w http.ResponseWriter, r *http.Request) {
	reqBody := stopRequestBody{}
	appErr := decode(r, &reqBody)
	if appErr != nil {
		reply.Err(w, appErr)
		return
	}

	//Fetch the auth id
	authID, err := model.GetAuthID(reqBody.To)
	if err != nil {
		reply.Err(w, ae.DB("", err))
		return
	}

	//Assign from and to inversly
	s := model.Stopped{authID, reqBody.To, reqBody.From}

	//Insert if the row doesn't exist
	if exists, err := s.Exists(); err != nil {
		reply.Err(w, ae.DB("", err))
		return
	} else if !exists {
		if err := s.Add(); err != nil {
			reply.Err(w, ae.DB("", err))
			return
		}
	}
	reply.OK(w, "")
}
