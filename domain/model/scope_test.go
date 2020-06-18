package model

import (
	"github.com/reecerussell/distro-blog/domain/datamodel"
	"testing"
)

func TestScope_Name(t *testing.T) {
	s := &Scope{
		name: "my scope",
	}

	if v := s.Name(); v != s.name {
		t.Errorf("expected '%s' but got '%s'", s.name, v)
	}
}

func TestScope_DTO(t *testing.T) {
	s := &Scope{
		id: "34693423",
		name: "scope:test",
	}

	d := s.DTO()

	if d.ID != s.id {
		t.Errorf("id was expected to be '%s' but was '%s'", s.id, d.ID)
	}

	if d.Name != s.name {
		t.Errorf("name was expected to be '%s' but was '%s'", s.name, d.Name)
	}
}

func TestScopeFromDataModel(t *testing.T) {
	dm := &datamodel.UserScope{
		ScopeID: "364923",
		ScopeName: "scope:test",
	}

	s := ScopeFromDataModel(dm)

	if dm.ScopeID != s.id {
		t.Errorf("id was expected to be '%s' but was '%s'", dm.ScopeID, s.id)
	}

	if dm.ScopeName != s.name {
		t.Errorf("name was expected to be '%s' but was '%s'", dm.ScopeName, s.name)
	}
}
