
package models

import (

	"github.com/hexya-erp/hexya/src/models/fieldtype"
	"github.com/satori/go.uuid"
)

// A UUIDField is a field for storing UUID.
type UUIDField struct {
	JSON          string
	String        string
	Help          string
	Stored        bool
	Required      bool
	ReadOnly      bool
	RequiredFunc  func(Environment) (bool, Conditioner)
	ReadOnlyFunc  func(Environment) (bool, Conditioner)
	InvisibleFunc func(Environment) (bool, Conditioner)
	Unique        bool
	Index         bool
	Compute       Methoder
	Depends       []string
	Related       string
	NoCopy        bool
	GoType        interface{}
	Translate     bool
	OnChange      Methoder
	Constraint    Methoder
	Inverse       Methoder
	Contexts      FieldContexts
	Default       func(Environment) interface{}
}

// DeclareField creates a html field for the given FieldsCollection with the given name.
func (uuidf UUIDField) DeclareField(fc *FieldsCollection, name string) *Field {
	fInfo := genericDeclareField(fc, &uuidf, name, fieldtype.UUID, new(uuid.UUID))
	return fInfo
}


