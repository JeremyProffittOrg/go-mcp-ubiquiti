package managers

import (
	"context"
	"encoding/json"
	"fmt"
)

// UserGroup represents a user group (bandwidth profile).
type UserGroup struct {
	ID           string `json:"_id,omitempty"`
	SiteID       string `json:"site_id,omitempty"`
	Name         string `json:"name"`
	QOSRateMaxUp   int  `json:"qos_rate_max_up,omitempty"`   // kbps, -1 = unlimited
	QOSRateMaxDown int  `json:"qos_rate_max_down,omitempty"` // kbps, -1 = unlimited
}

// UserGroupManager handles user group operations.
type UserGroupManager struct {
	conn *ConnectionManager
}

// NewUserGroupManager creates a new user group manager.
func NewUserGroupManager(conn *ConnectionManager) *UserGroupManager {
	return &UserGroupManager{conn: conn}
}

// ListUserGroups returns all user groups.
func (um *UserGroupManager) ListUserGroups(ctx context.Context) ([]UserGroup, error) {
	endpoint := um.conn.GetSitePath("/rest/usergroup")

	data, err := um.conn.Request(ctx, "GET", endpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get user groups: %w", err)
	}

	var groups []UserGroup
	if err := json.Unmarshal(data, &groups); err != nil {
		return nil, fmt.Errorf("failed to parse user groups: %w", err)
	}

	return groups, nil
}

// GetUserGroupDetails returns details for a specific user group.
func (um *UserGroupManager) GetUserGroupDetails(ctx context.Context, groupID string) (*UserGroup, error) {
	groups, err := um.ListUserGroups(ctx)
	if err != nil {
		return nil, err
	}

	for _, g := range groups {
		if g.ID == groupID {
			return &g, nil
		}
	}

	return nil, fmt.Errorf("user group not found: %s", groupID)
}

// CreateUserGroup creates a new user group.
func (um *UserGroupManager) CreateUserGroup(ctx context.Context, group *UserGroup) (*UserGroup, error) {
	endpoint := um.conn.GetSitePath("/rest/usergroup")

	data, err := um.conn.Request(ctx, "POST", endpoint, group)
	if err != nil {
		return nil, fmt.Errorf("failed to create user group: %w", err)
	}

	var groups []UserGroup
	if err := json.Unmarshal(data, &groups); err != nil {
		return nil, fmt.Errorf("failed to parse created group: %w", err)
	}

	if len(groups) == 0 {
		return nil, fmt.Errorf("no group returned after creation")
	}

	return &groups[0], nil
}

// UpdateUserGroup updates an existing user group.
func (um *UserGroupManager) UpdateUserGroup(ctx context.Context, groupID string, updates map[string]interface{}) error {
	endpoint := um.conn.GetSitePath("/rest/usergroup/" + groupID)

	_, err := um.conn.Request(ctx, "PUT", endpoint, updates)
	if err != nil {
		return fmt.Errorf("failed to update user group: %w", err)
	}

	return nil
}

// DeleteUserGroup deletes a user group.
func (um *UserGroupManager) DeleteUserGroup(ctx context.Context, groupID string) error {
	endpoint := um.conn.GetSitePath("/rest/usergroup/" + groupID)

	_, err := um.conn.Request(ctx, "DELETE", endpoint, nil)
	if err != nil {
		return fmt.Errorf("failed to delete user group: %w", err)
	}

	return nil
}
