package expander

import (
	"fmt"
	u "r53/aws_utils"

	log "github.com/sirupsen/logrus"
)

type Expander interface {
	Expand(from u.Resource, to u.ResourceType) ([]u.Resource, error)
	IsExpanded(from u.Resource, to u.ResourceType) bool
	IsExpandable(from u.Resource, to u.ResourceType) bool
	GetPotentialExpansions(from u.ResourceType) []u.ResourceType
}

// responsible for expanding from resource to another one
type SingleThreadedExpander struct {
	describer u.ResourceDescriber
	cache     Cache
}

func NewExpander(describer u.ResourceDescriber, cache Cache) Expander {
	return &SingleThreadedExpander{
		cache:     cache,
		describer: describer,
	}
}

// buildInput for describe request
func (e *SingleThreadedExpander) buildInput(from u.Resource, to u.ResourceType) ([]interface{}, error) {
	input, err := u.BuildResourceInputFromInput(from, to)

	if err != nil {
		return nil, err
	}

	return input.Inputs, nil
}

func (e *SingleThreadedExpander) Expand(from u.Resource, to u.ResourceType) ([]u.Resource, error) {

	var expansionErrs string
	var descriptions []u.Resource
	var err error

	if !e.IsExpandable(from, to) {
		return nil, fmt.Errorf("NoExpansionPath")
	}

	// check if cached already

	if expansions, found := e.cache.GetExpansion(from, to); found {
		return expansions, nil
	}

	// get the correct input

	inputs, err := e.buildInput(from, to)

	if err != nil {
		return nil, err
	}

	for _, in := range inputs {
		// describer resource

		description, err := e.describer.DescribeResource(to, in)

		if err != nil {
			log.WithError(err).Debug("failed expanding resource input")
			expansionErrs += fmt.Sprintf(" | %s", err.Error())
			continue
		}

		descriptions = append(descriptions, description)
	}

	if expansionErrs != "" {
		err = fmt.Errorf("%s", expansionErrs)
	}

	// set cache
	e.cache.SetExpansion(from, to, descriptions)

	return descriptions, err
}

func (e *SingleThreadedExpander) IsExpanded(from u.Resource, to u.ResourceType) bool {
	_, found := e.cache.GetExpansion(from, to)
	return found
}

// check if there potentially a path from resource to another one
func (e *SingleThreadedExpander) IsExpandable(from u.Resource, to u.ResourceType) bool {

	if expandableTypes, exist := u.AdjacentResources[from.Type()]; exist {

		for _, rType := range expandableTypes {
			if rType == to {
				return true
			}
		}

	}

	return false
}

func (e *SingleThreadedExpander) GetPotentialExpansions(from u.ResourceType) []u.ResourceType {
	return u.AdjacentResources[from]
}
