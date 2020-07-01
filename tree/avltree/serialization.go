package avltree

import (
	"encoding/json"

	"github.com/morganxf/algorithm/util"
)

func (t *Tree) ToJSON() ([]byte, error) {
	elements := make(map[string]interface{})
	it := t.Iterator()
	for it.Next() {
		elements[util.ToString(it.Key())] = it.Value()
	}
	return json.Marshal(elements)
}

func (t *Tree) FromJSON(data []byte) error {
	elements := make(map[string]interface{})
	err := json.Unmarshal(data, &elements)
	if err != nil {
		return err
	}
	t.Clear()
	for k, v := range elements {
		t.Put(k, v)
	}
	return nil
}
