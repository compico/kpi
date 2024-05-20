package fact

import (
	"bytes"
	"io"
	"mime/multipart"
	"strconv"
)

// Fact стандартная структура из документации, добавил id для работы в mysql
type Fact struct {
	Id                  int    `form:"-" json:"-" gorm:"primaryKey"`
	IndicatorToMoFactId int    `form:"indicator_to_mo_fact_id" json:"indicator_to_mo_fact_id" gorm:"index"` //indicator_to_mo_fact_id:0
	PeriodStart         string `form:"period_start" json:"period_start"`                                    //period_start:2024-05-01
	PeriodEnd           string `form:"period_end" json:"period_end"`                                        //period_end:2024-05-31
	PeriodKey           string `form:"period_key" json:"period_key"`                                        //period_key:month
	IndicatorToMoId     int    `form:"indicator_to_mo_id" json:"indicator_to_mo_id"`                        //indicator_to_mo_id:227373
	Value               int    `form:"value" json:"value"`                                                  //value:1
	FactTime            string `form:"fact_time" json:"fact_time"`                                          //fact_time:2024-05-31
	IsPlan              bool   `form:"is_plan" json:"is_plan"`                                              //is_plan:0
	AuthUserId          int    `form:"auth_user_id" json:"auth_user_id"`                                    //auth_user_id:40
	Comment             string `form:"comment" json:"comment"`                                              //comment: buffer Last_name
}

type Collection []*Fact

func (f *Fact) ToFromData() map[string]string {
	return map[string]string{
		"period_start":            f.PeriodStart,
		"period_end":              f.PeriodEnd,
		"period_key":              f.PeriodKey,
		"indicator_to_mo_id":      strconv.Itoa(f.IndicatorToMoId),
		"indicator_to_mo_fact_id": strconv.Itoa(f.IndicatorToMoFactId),
		"value":                   strconv.Itoa(f.Value),
		"fact_time":               f.FactTime,
		"is_plan":                 f.IsPlanString(),
		"auth_user_id":            strconv.Itoa(f.AuthUserId),
		"comment":                 f.Comment,
	}
}

func (f *Fact) Payload() (io.Reader, string, error) {
	payload := &bytes.Buffer{}
	writer := multipart.NewWriter(payload)
	_ = writer.WriteField("period_start", f.PeriodStart)
	_ = writer.WriteField("period_end", f.PeriodEnd)
	_ = writer.WriteField("period_key", f.PeriodKey)
	_ = writer.WriteField("indicator_to_mo_id", strconv.Itoa(f.IndicatorToMoId))
	_ = writer.WriteField("indicator_to_mo_fact_id", strconv.Itoa(f.IndicatorToMoFactId))
	_ = writer.WriteField("value", strconv.Itoa(f.Value))
	_ = writer.WriteField("fact_time", f.FactTime)
	_ = writer.WriteField("is_plan", f.IsPlanString())
	_ = writer.WriteField("auth_user_id", strconv.Itoa(f.AuthUserId))
	_ = writer.WriteField("comment", f.Comment)

	return payload, writer.FormDataContentType(), writer.Close()
}

func (f *Fact) IsPlanString() string {
	if f.IsPlan {
		return "1"
	}
	return "0"
}
