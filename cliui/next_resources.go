package cliui

import (
	awsu "r53/aws_utils"

	"github.com/rivo/tview"
)

type ResourceInput struct {
	Type   awsu.ResourceType
	UIName string
}

type NextResourcesInput struct {
	Next                   []ResourceInput
	OnNextResourceSelected func(i *NextResourcesInput, selected *ResourceInput, isBack bool)
}

func GetResourcesList(existingList *tview.List, i *NextResourcesInput) *tview.List {

	list := existingList

	if existingList == nil {
		list = tview.NewList()
	}

	list.Clear()

	for idx, r := range i.Next {
		resource := r
		list.AddItem(resource.UIName, "", rune(97+idx), func() {
			i.OnNextResourceSelected(i, &resource, false)
		})
	}

	list.AddItem("Previous", "", 'q', func() {
		i.OnNextResourceSelected(i, nil, true)
	})

	return list
}

type NextResourcesView interface {
	Get() tview.Primitive
	SetNextAvailableResource(i *NextResourcesInput)
}

type NextResourcesIndicator struct {
	nextAvailable *tview.List
}

func NewNextResourceIndicator() NextResourcesView {
	return &NextResourcesIndicator{
		nextAvailable: tview.NewList(),
	}

}

func (ni *NextResourcesIndicator) SetNextAvailableResource(i *NextResourcesInput) {
	GetResourcesList(ni.nextAvailable, i)
}

func (ni *NextResourcesIndicator) Get() tview.Primitive {
	return ni.nextAvailable
}

func GetNextResourcesView() NextResourcesView {
	return NewNextResourceIndicator()
}
