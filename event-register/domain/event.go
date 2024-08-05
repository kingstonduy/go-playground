package domain

// Define your event struct
type Event struct {
	AggregateID   string `json:"AGGREGATE_ID"`
	AggregateType string `json:"AGGREGATE_TYPE"`
	CommandID     string `json:"COMMAND_ID"`
	CommandType   string `json:"COMMAND_TYPE"`
	Payload       string `json:"PAYLOAD"`
	Processed     int8   `json:"PROCESSED"`
	ProcessAt     int64  `json:"PROCESS_AT"`
	TraceParentID string `json:"TRACE_PARENT_ID"`
}

// Define your payload structs
type Foo struct {
	Message string `json:"message"`
}
type Bar struct {
	Value int `json:"value"`
}
type Nop struct {
	Result string `json:"result"`
}
