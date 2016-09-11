package main

import (
	"github.com/echocat/caretakerd/errors"
	"sort"
)

// PickedDefinitions represents every Definition that is picked for rendering.
type PickedDefinitions struct {
	// RootId is the ID of the root.
	RootID IDType
	// Source contains the original Definitions that are contained in this instance.
	Source *Definitions
	// IDToDefinition is a map of every ID of a Definition to its instance.
	IDToDefinition map[string]Definition
	// IDToInlinedMarker is a map of every ID of a Definition to its instance as WithInlinedMarker type.
	IDToInlinedMarker map[string]WithInlinedMarker
	// TopLevelDefinitions contains every Definitions that are at the top level and have no children.
	TopLevelDefinitions []Definition
}

// PickDefinitionsFrom picks every Definition to be displayed from the given Definitions.
func PickDefinitionsFrom(source *Definitions, rootElementID IDType) (*PickedDefinitions, error) {
	pd := &PickedDefinitions{
		Source:              source,
		IDToDefinition:      map[string]Definition{},
		IDToInlinedMarker:   map[string]WithInlinedMarker{},
		TopLevelDefinitions: []Definition{},
	}
	err := enrichWithElementAndItsChildren(pd, rootElementID)
	if err != nil {
		return nil, err
	}
	defs := definitions{}
	for _, definition := range pd.IDToDefinition {
		if inlinedMarker, ok := definition.(WithInlinedMarker); ok && inlinedMarker.Inlined() {
			pd.IDToInlinedMarker[definition.ID().String()] = inlinedMarker
		} else if definition.IsTopLevel() {
			defs = append(defs, definition)
		}
	}
	sort.Sort(defs)
	pd.TopLevelDefinitions = defs
	pd.RootID = rootElementID
	return pd, nil
}

func enrichWithElementAndItsChildren(pd *PickedDefinitions, elementID IDType) error {
	if pd.IDToDefinition[elementID.String()] != nil {
		return nil
	}
	element, err := pd.GetSourceElementBy(elementID)
	if err != nil {
		return err
	}
	if element == nil {
		return nil
	}
	pd.IDToDefinition[elementID.String()] = element
	if valueType, ok := element.(WithValueType); ok {
		for _, idType := range ExtractAllIDTypesFrom(valueType.ValueType()) {
			err := enrichWithElementAndItsChildren(pd, idType)
			if err != nil {
				return errors.New("Could not extract valueType '%v' of type '%s'.", idType, elementID).CausedBy(err)
			}
		}
	}
	if children, ok := element.(WithChildren); ok {
		for _, child := range children.Children() {
			err := enrichWithElementAndItsChildren(pd, child.ID())
			if err != nil {
				return errors.New("Could not extract child '%v' of type '%s'.", child.ID(), elementID).CausedBy(err)
			}
		}
	}

	return nil
}

// GetSourceElementBy returns the original Definition for the given ID.
// Returns nil if this Definition is a primitive one or an error if the Definition does not exist.
func (instance *PickedDefinitions) GetSourceElementBy(id IDType) (Definition, error) {
	if id.Primitive {
		return nil, nil
	}
	result := instance.Source.GetBy(id)
	if result == nil {
		return nil, errors.New("Could not find expected element '%s'.", id)
	}
	return result, nil
}

// FindInlinedFor returns the original inlined Definition for the given ID.
// nil is returned if this Definition does not exist.
func (instance *PickedDefinitions) FindInlinedFor(id IDType) WithInlinedMarker {
	result, ok := instance.IDToInlinedMarker[id.String()]
	if ok {
		return result
	}
	return nil
}

type definitions []Definition

func (instance definitions) Len() int {
	return len(instance)
}

func (instance definitions) Swap(i, j int) {
	instance[i], instance[j] = instance[j], instance[i]
}

func (instance definitions) Less(i, j int) bool {
	return instance[i].ID().String() < instance[j].ID().String()
}
