package quartz

type QuartzResponseType int

const (
	QUARTZ_RESP_TYPE_ACK      QuartzResponseType = iota
	QUARTZ_RESP_TYPE_ERR      QuartzResponseType = iota
	QUARTZ_RESP_TYPE_PWRON    QuartzResponseType = iota
	QUARTZ_RESP_TYPE_UPDATE   QuartzResponseType = iota
	QUARTZ_RESP_TYPE_READ_DST QuartzResponseType = iota
	QUARTZ_RESP_TYPE_READ_SRC QuartzResponseType = iota
	QUARTZ_RESP_TYPE_READ_LVL QuartzResponseType = iota
	QUARTZ_RESP_TYPE_LOCK_STS QuartzResponseType = iota
)

type QuartzResponse interface {
	GetType() QuartzResponseType
	GetRaw() string
}

type ResponseAcknowledge struct {
	RawData string
}

func (r *ResponseAcknowledge) GetType() QuartzResponseType {
	return QUARTZ_RESP_TYPE_ACK
}

func (r *ResponseAcknowledge) GetRaw() string {
	return r.RawData
}

type ResponseError struct {
	RawData string
}

func (r *ResponseError) GetType() QuartzResponseType {
	return QUARTZ_RESP_TYPE_ERR
}

func (r *ResponseError) GetRaw() string {
	return r.RawData
}

type ResponsePowerOn struct {
	RawData string
}

func (r *ResponsePowerOn) GetType() QuartzResponseType {
	return QUARTZ_RESP_TYPE_PWRON
}

func (r *ResponsePowerOn) GetRaw() string {
	return r.RawData
}

type ResponseUpdate struct {
	RawData     string
	Levels      []QuartzLevel
	Destination uint
	Source      uint
}

func (r *ResponseUpdate) GetType() QuartzResponseType {
	return QUARTZ_RESP_TYPE_UPDATE
}

func (r *ResponseUpdate) GetRaw() string {
	return r.RawData
}

type ResponseReadDestination struct {
	RawData     string
	Destination uint
	Name        string
}

func (r *ResponseReadDestination) GetType() QuartzResponseType {
	return QUARTZ_RESP_TYPE_READ_DST
}

func (r *ResponseReadDestination) GetRaw() string {
	return r.RawData
}

type ResponseReadSource struct {
	RawData string
	Source  uint
	Name    string
}

func (r *ResponseReadSource) GetType() QuartzResponseType {
	return QUARTZ_RESP_TYPE_READ_SRC
}

func (r *ResponseReadSource) GetRaw() string {
	return r.RawData
}

type ResponseReadLevel struct {
	RawData string
	Level   QuartzLevel
	Name    string
}

func (r *ResponseReadLevel) GetType() QuartzResponseType {
	return QUARTZ_RESP_TYPE_READ_LVL
}

func (r *ResponseReadLevel) GetRaw() string {
	return r.RawData
}

type ResponseLockStatus struct {
	RawData     string
	Destination uint
	Locked      bool
}

func (r *ResponseLockStatus) GetType() QuartzResponseType {
	return QUARTZ_RESP_TYPE_LOCK_STS
}

func (r *ResponseLockStatus) GetRaw() string {
	return r.RawData
}
