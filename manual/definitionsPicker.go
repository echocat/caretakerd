package main

import (
	"github.com/echocat/caretakerd/errors"
	"sort"
)

type PickedDefinitions struct {
	RootId              IdType
	Source              *Definitions
	NameToDefinition    map[string]Definition
	NameToInlinedMarker map[string]WithInlinedMarker
	TopLevelDefinitions []Definition
}

func PickDefinitionsFrom(source *Definitions, rootElementId IdType) (*PickedDefinitions, error) {
	pd := &PickedDefinitions{
		Source:              source,
		NameToDefinition:    map[string]Definition{},
		NameToInlinedMarker: map[string]WithInlinedMarker{},
		TopLevelDefinitions: []Definition{},
	}
	err := enrichWithElementAndItsChildren(pd, rootElementId)
	if err != nil {
		return nil, err
	}
	defs := definitions{}
	for _, definition := range pd.NameToDefinition {
		if inlinedMarker, ok := definition.(WithInlinedMarker); ok && inlinedMarker.Inlined() {
			pd.NameToInlinedMarker[definition.Id().String()] = inlinedMarker
		} else if definition.IsTopLevel() {
			defs = append(defs, definition)
		}
	}
	sort.Sort(defs)
	pd.TopLevelDefinitions = defs
	pd.RootId = rootElementId
	return pd, nil
}

func enrichWithElementAndItsChildren(pd *PickedDefinitions, elementId IdType) error {
	if pd.NameToDefinition[elementId.String()] != nil {
		return nil
	}
	element, err := pd.GetSourceElementBy(elementId)
	if err != nil {
		return err
	}
	if element == nil {
		return nil
	}
	pd.NameToDefinition[elementId.String()] = element
	if valueType, ok := element.(WithValueType); ok {
		for _, idType := range ExtractAllIdTypesFrom(valueType.ValueType()) {
			err := enrichWithElementAndItsChildren(pd, idType)
			if err != nil {
				return errors.New("Could not extract valueType '%v' of type '%s'.", idType, elementId).CausedBy(err)
			}
		}
	}
	if children, ok := element.(WithChildren); ok {
		for _, child := range children.Children() {
			err := enrichWithElementAndItsChildren(pd, child.Id())
			if err != nil {
				return errors.New("Could not extract child '%v' of type '%s'.", child.Id(), elementId).CausedBy(err)
			}
		}
	}

	return nil
}

func (instance *PickedDefinitions) GetSourceElementBy(id IdType) (Definition, error) {
	if id.Primitive {
		return nil, nil
	}
	result := instance.Source.GetBy(id)
	if result == nil {
		return nil, nil // TODO! Should we fail here or just ignore it??
		//return nil, errors.New("Could not find expected element '%s'.", id)
	}
	return result, nil
}

func (instance *PickedDefinitions) FindInlinedFor(id IdType) WithInlinedMarker {
	result, ok := instance.NameToInlinedMarker[id.String()]
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
	return instance[i].Id().String() < instance[j].Id().String()
}
