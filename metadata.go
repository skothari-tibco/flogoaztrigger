package flogoaztrigger

import (
	"github.com/project-flogo/core/data/coerce"
)

type Output struct {
	Body string `md:"body"`
}

func (o *Output) FromMap(values map[string]interface{}) error {

	var err error
	o.Body, err = coerce.ToString(values["body"])
	if err != nil {
		return err
	}

	return nil
}

func (o *Output) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"body": o.Body,
	}
}

type Reply struct {
	Code int         `md:"code"`
	Data interface{} `md:"data"`
}

func (r *Reply) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"code": r.Code,
		"data": r.Data,
	}
}

func (r *Reply) FromMap(values map[string]interface{}) error {

	var err error
	r.Code, err = coerce.ToInt(values["code"])
	if err != nil {
		return err
	}
	r.Data, _ = values["data"]

	return nil
}
