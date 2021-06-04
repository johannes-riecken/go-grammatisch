package model

import (
	"encoding/json"
	"fmt"
)

func (x *Alternative) UnmarshalJSON(data []byte) error {
	type tmp struct {
		Elements json.RawMessage
	}
	var y tmp
	err := json.Unmarshal(data, &y)
	if err != nil {
		return err
	}
	var elementInstances []map[string]interface{}
	err = json.Unmarshal(y.Elements, &elementInstances)
	if err != nil {
		return err
	}
	for _, elementInstance := range elementInstances {
		switch elementInstance["Type"] {
		case "RuleRef":
			var z RuleRef
			data, err := json.Marshal(elementInstance)
			if err != nil {
				return err
			}
			err = json.Unmarshal(data, &z)
			if err != nil {
				return err
			}
			x.Elements = append(x.Elements, z)
		case "Quoted":
			var z Quoted
			data, err := json.Marshal(elementInstance)
			if err != nil {
				return err
			}
			err = json.Unmarshal(data, &z)
			if err != nil {
				return err
			}
			x.Elements = append(x.Elements, z)
		default:
			return fmt.Errorf("unknown type: %v", elementInstance["Type"])
		}
	}
	return nil
}

func (x RuleRef) MarshalJSON() ([]byte, error) {
	type marshal RuleRef
	y := marshal(x)
	y.Type = "RuleRef"
	return json.Marshal(y)
}

func (x Quoted) MarshalJSON() ([]byte, error) {
	type marshal Quoted
	y := marshal(x)
	y.Type = "Quoted"
	return json.Marshal(y)
}
