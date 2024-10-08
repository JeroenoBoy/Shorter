package models

type Permissions int

const (
	PermissionsAdmin        Permissions = ^0
	PermissionsManageUsers  Permissions = 1 << 0
	PermissionsManageShorts Permissions = 1 << 1
	PermissionsManageServer Permissions = 1 << 2
    PermissionsAnyDashboardAccess = PermissionsManageUsers + PermissionsManageShorts + PermissionsManageServer
)

func (p Permissions) HasAll(permissions Permissions) bool {
	return (p & permissions) == permissions
}

func (p Permissions) HasAny(permissions Permissions) bool {
    return (p & permissions) > 0
}
