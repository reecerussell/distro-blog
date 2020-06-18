package model

import (
	"github.com/reecerussell/distro-blog/domain/datamodel"
	"github.com/reecerussell/distro-blog/domain/dto"
)

// Scope is domain model for the scope entity.
type Scope struct {
	id string
	name string
}

// Name returns the scope's name.
func (s *Scope) Name() string {
	return s.name
}

// DTO returns a new *dto.Scope populated with the scopes's values.
func (s *Scope) DTO() *dto.Scope {
	return &dto.Scope{
		ID: s.id,
		Name: s.name,
	}
}

// ScopeFromDataModel returns a new Scope domain model populated with
// values from the given data model.
func ScopeFromDataModel(dm *datamodel.UserScope) *Scope {
	return &Scope{
		id: dm.ScopeID,
		name: dm.ScopeName,
	}
}

