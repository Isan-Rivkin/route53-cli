package expander

import (
	"fmt"
	u "r53/aws_utils"
)

type Cache interface {
	GetExpansion(from u.Resource, to u.ResourceType) ([]u.Resource, bool)
	SetExpansion(from u.Resource, to u.ResourceType, expansions []u.Resource)
}

type ExpanderCache struct {
	expansionsCache map[string][]u.Resource
}

func NewExpanderCache() Cache {
	return &ExpanderCache{
		expansionsCache: map[string][]u.Resource{},
	}
}

func (c *ExpanderCache) genKey(from u.Resource, to u.ResourceType) string {
	return fmt.Sprintf("%s-%s", to, from.ID())
}

func (c *ExpanderCache) GetExpansion(from u.Resource, to u.ResourceType) ([]u.Resource, bool) {
	key := c.genKey(from, to)

	expansions, found := c.expansionsCache[key]
	return expansions, found
}

func (c *ExpanderCache) SetExpansion(from u.Resource, to u.ResourceType, expansions []u.Resource) {
	key := c.genKey(from, to)

	_, found := c.GetExpansion(from, to)

	if found {
		c.expansionsCache[key] = append(c.expansionsCache[key], expansions...)
	} else {
		c.expansionsCache[key] = expansions
	}
}
