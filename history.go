package music

import (
	"encoding/base64"
	"strconv"
	"time"
)

type History struct {
	ID          int64     `json:"id"`
	Name        string    `json:"name"`
	Type        string    `json:"type"`
	Reason      string    `json:"reason"`
	DeletedDate time.Time `json:"deletedDate"`
}

type HistoryJSON struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Type        string `json:"type"`
	Reason      string `json:"reason"`
	DeletedDate string `json:"deleted_date"`
}

func (h *History) SetFromJSON(history HistoryJSON) error {
	decoded, err := base64.StdEncoding.DecodeString(history.ID)
	if err != nil {
		return err
	}
	h.ID, err = strconv.ParseInt(string(decoded), 10, 64)
	if err != nil {
		return err
	}
	decoded, err = base64.StdEncoding.DecodeString(history.Name)
	if err != nil {
		return err
	}
	h.Name = string(decoded)
	decoded, err = base64.StdEncoding.DecodeString(history.Type)
	if err != nil {
		return err
	}
	h.Type = string(decoded)
	decoded, err = base64.StdEncoding.DecodeString(history.Reason)
	if err != nil {
		return err
	}
	h.Reason = string(decoded)

	decoded, err = base64.StdEncoding.DecodeString(history.DeletedDate)
	if err != nil {
		return err
	}
	h.DeletedDate, err = time.Parse("2006-01-02 15:04:05", string(decoded))
	if err != nil {
		return err
	}

	return err
}
