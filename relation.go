package main

// Relations creates a map where the source id is linked to the target
// Fact Sheets.
func Relations(source string, targets []interface{}) map[string][]Relation {
	relations := make([]Relation, 0)
	for _, t := range targets {
		relations = append(relations, NewRelation(t))
	}
	return map[string][]Relation{
		source: relations,
	}
}

// NewRelation creates a new relation from the source Fact Sheet
// to the target Fact Sheet.
func NewRelation(target interface{}) Relation {
	return map[string]interface{}{
		"uid":     target,
		"relName": "relToRequiredBy",
	}
}
