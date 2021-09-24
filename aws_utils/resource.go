package aws_utils

type ResourceKey string

const (
	ARNAttr     ResourceKey = "arn"
	Ec2IDAttr   ResourceKey = "ec2id"
	CertArnAttr ResourceKey = "certarn"
	RegionAttr  ResourceKey = "region"
	DNSAttr     ResourceKey = "dns"
)

type ResourceType string

type Resource interface {
	ID() string
	Type() ResourceType
	// describer input
	GetQueryInput() interface{}
	// describer output
	GetQueryOutput() interface{}
	// optional: if input was constructed dynamically from some source
	GetInputSource() interface{}
	// Optional: preset bag of values
	TryGetValue(key ResourceKey) ([]string, bool)
}

type ResourceDescriptionResult struct {
	rType ResourceType
	// some unique id
	id string
	//if available its the source struct that was used to create  the query input
	// i.e r53RecordSet and via that queryInput to get ELB was created
	inputSource interface{}
	// the input used to get the resource description initially
	queryInput interface{}
	// the actuall result from the description output such as some elb description
	queryOutput interface{}
	// bad of values
	optionalResourceValues map[ResourceKey][]string
}

func NewResourceDescriptionResult(inputSource, input, output interface{}, rType ResourceType, optionalValues map[ResourceKey][]string, id string) Resource {
	return &ResourceDescriptionResult{
		inputSource:            inputSource,
		queryInput:             input,
		queryOutput:            output,
		id:                     id,
		rType:                  rType,
		optionalResourceValues: optionalValues,
	}
}

func (r *ResourceDescriptionResult) TryGetValue(key ResourceKey) ([]string, bool) {
	if r.optionalResourceValues != nil {
		vals, exist := r.optionalResourceValues[key]
		return vals, exist
	}
	return nil, false
}

func (r *ResourceDescriptionResult) ID() string {
	return r.id
}

func (r *ResourceDescriptionResult) Type() ResourceType {
	return r.rType
}

func (r *ResourceDescriptionResult) GetQueryInput() interface{} {
	return r.queryInput
}
func (r *ResourceDescriptionResult) GetQueryOutput() interface{} {
	return r.queryOutput
}

func (r *ResourceDescriptionResult) GetInputSource() interface{} {
	return r.inputSource
}
