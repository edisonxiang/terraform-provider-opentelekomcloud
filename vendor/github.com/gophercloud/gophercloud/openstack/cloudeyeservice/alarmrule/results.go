package alarmrule

import (
	"github.com/gophercloud/gophercloud"
)

type CreateResponse struct {
	AlarmID string `json:"alarm_id"`
}

type CreateResult struct {
	gophercloud.Result
}

// Extract is a function that accepts a result and extracts a router.
func (c CreateResult) Extract() (*CreateResponse, error) {
	r := &CreateResponse{}
	return r, c.ExtractInto(r)
}

type DimensionInfo struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type MetricInfo struct {
	Namespace  string          `json:"namespace"`
	MetricName string          `json:"metric_name"`
	Dimensions []DimensionInfo `json:"dimensions"`
}

type ConditionInfo struct {
	Period             int    `json:"period"`
	Filter             string `json:"filter"`
	ComparisonOperator string `json:"comparison_operator"`
	Value              int    `json:"value"`
	Unit               string `json:"unit"`
	Count              int    `json:"count"`
}

type ActionInfo struct {
	Type             string   `json:"type"`
	NotificationList []string `json:"notificationList"`
}

type AlarmRule struct {
	AlarmName               string        `json:"alarm_name"`
	AlarmDescription        string        `json:"alarm_description"`
	Metric                  MetricInfo    `json:"metric"`
	Condition               ConditionInfo `json:"condition"`
	AlarmActions            []ActionInfo  `json:"alarm_actions"`
	InsufficientdataActions []ActionInfo  `json:"insufficientdata_actions"`
	OkActions               []ActionInfo  `json:"ok_actions"`
	AlarmEnabled            bool          `json:"alarm_enabled"`
	AlarmActionEnabled      bool          `json:"alarm_action_enabled"`
	ID                      string        `json:"alarm_id"`
	UpdateTime              int64         `json:"update_time"`
	AlarmState              string        `json:"alarm_state"`
}

type GetResult struct {
	gophercloud.Result
}

// Extract is a function that accepts a result and extracts a router.
func (g GetResult) Extract() (*AlarmRule, error) {
	r := &AlarmRule{}
	return r, g.ExtractInto(r)
}
