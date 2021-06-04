package model

import (
	"encoding/json"
)

func (x *Define) UnmarshalJSON(data []byte) error {
	type tmp struct {
		DefineName string
		RegexSteps json.RawMessage
	}
	var y tmp
	err := json.Unmarshal(data, &y)
	if err != nil {
		return err
	}
	x.DefineName = y.DefineName
	var stepInstances []map[string]interface{}
	err = json.Unmarshal(y.RegexSteps, &stepInstances)
	if err != nil {
		return err
	}
	for _, stepInstance := range stepInstances {
		if stepInstance["Type"] == "PositionSaveStep" {
			var z PositionSaveStep
			data, err := json.Marshal(stepInstance)
			if err != nil {
				return err
			}
			err = json.Unmarshal(data, &z)
			if err != nil {
				return err
			}
			x.RegexSteps = append(x.RegexSteps, z)
		}
		if stepInstance["Type"] == "CallStep" {
			var z CallStep
			data, err := json.Marshal(stepInstance)
			if err != nil {
				return err
			}
			err = json.Unmarshal(data, &z)
			if err != nil {
				return err
			}
			x.RegexSteps = append(x.RegexSteps, z)
		}
		if stepInstance["Type"] == "MatchCombineStep" {
			var z MatchCombineStep
			data, err := json.Marshal(stepInstance)
			if err != nil {
				return err
			}
			err = json.Unmarshal(data, &z)
			if err != nil {
				return err
			}
			x.RegexSteps = append(x.RegexSteps, z)
		}
		if stepInstance["Type"] == "MatchSaveStep" {
			var z MatchSaveStep
			data, err := json.Marshal(stepInstance)
			if err != nil {
				return err
			}
			err = json.Unmarshal(data, &z)
			if err != nil {
				return err
			}
			x.RegexSteps = append(x.RegexSteps, z)
		}
		if stepInstance["Type"] == "MatchStep" {
			var z MatchStep
			data, err := json.Marshal(stepInstance)
			if err != nil {
				return err
			}
			err = json.Unmarshal(data, &z)
			if err != nil {
				return err
			}
			x.RegexSteps = append(x.RegexSteps, z)
		}
	}
	return nil
}

func (x PositionSaveStep) MarshalJSON() ([]byte, error) {
	type marshal PositionSaveStep
	y := marshal(x)
	y.Type = "PositionSaveStep"
	return json.Marshal(y)
}

func (x CallStep) MarshalJSON() ([]byte, error) {
	type marshal CallStep
	y := marshal(x)
	y.Type = "CallStep"
	return json.Marshal(y)
}

func (x MatchCombineStep) MarshalJSON() ([]byte, error) {
	type marshal MatchCombineStep
	y := marshal(x)
	y.Type = "MatchCombineStep"
	return json.Marshal(y)
}

func (x MatchSaveStep) MarshalJSON() ([]byte, error) {
	type marshal MatchSaveStep
	y := marshal(x)
	y.Type = "MatchSaveStep"
	return json.Marshal(y)
}

func (x MatchStep) MarshalJSON() ([]byte, error) {
	type marshal MatchStep
	y := marshal(x)
	y.Type = "MatchStep"
	return json.Marshal(y)
}
