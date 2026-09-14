package scripts

import (
	"app/data"
)

type EntrySingleGroup struct{}

func (EntrySingleGroup) GetName() string {
	return "entry/group"
}

func (EntrySingleGroup) Run(_ string, entry *data.Entry) error {
	if len(entry.Groups) > 0 {
		entry.Group = entry.Groups[0]
		entry.Groups = nil
	}

	return nil
}
