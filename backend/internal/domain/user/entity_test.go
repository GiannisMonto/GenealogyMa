package user

import (
	"testing"
)

func TestUserValidate(t *testing.T) {
	tests := []struct {
		name    string
		user    *User
		wantErr bool
	}{
		{
			name: "valid user",
			user: &User{
				Username: "testuser",
				Email:    "test@example.com",
				Status:   StatusActive,
			},
			wantErr: false,
		},
		{
			name: "empty username",
			user: &User{
				Username: "",
				Email:    "test@example.com",
				Status:   StatusActive,
			},
			wantErr: true,
		},
		{
			name: "username too short",
			user: &User{
				Username: "ab",
				Email:    "test@example.com",
				Status:   StatusActive,
			},
			wantErr: true,
		},
		{
			name: "username too long",
			user: &User{
				Username: string(make([]byte, 51)),
				Email:    "test@example.com",
				Status:   StatusActive,
			},
			wantErr: true,
		},
		{
			name: "empty email",
			user: &User{
				Username: "testuser",
				Email:    "",
				Status:   StatusActive,
			},
			wantErr: true,
		},
		{
			name: "invalid status",
			user: &User{
				Username: "testuser",
				Email:    "test@example.com",
				Status:   "invalid",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.user.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("User.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestUserIsActive(t *testing.T) {
	tests := []struct {
		name   string
		status UserStatus
		want   bool
	}{
		{"active", StatusActive, true},
		{"inactive", StatusInactive, false},
		{"banned", StatusBanned, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &User{Status: tt.status}
			if got := u.IsActive(); got != tt.want {
				t.Errorf("User.IsActive() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUserCanLogin(t *testing.T) {
	tests := []struct {
		name   string
		status UserStatus
		wantErr bool
	}{
		{"active can login", StatusActive, false},
		{"inactive cannot login", StatusInactive, true},
		{"banned cannot login", StatusBanned, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &User{Status: tt.status}
			err := u.CanLogin()
			if (err != nil) != tt.wantErr {
				t.Errorf("User.CanLogin() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestUserHasRole(t *testing.T) {
	u := &User{
		Roles: []*Role{
			{Code: "super_admin"},
			{Code: "verified_user"},
		},
	}

	if !u.HasRole("super_admin") {
		t.Error("expected user to have super_admin role")
	}
	if u.HasRole("guest") {
		t.Error("expected user to not have guest role")
	}
}

func TestUserGetRoleCodes(t *testing.T) {
	u := &User{
		Roles: []*Role{
			{Code: "super_admin"},
			{Code: "verified_user"},
		},
	}

	codes := u.GetRoleCodes()
	if len(codes) != 2 {
		t.Fatalf("expected 2 role codes, got %d", len(codes))
	}
	if codes[0] != "super_admin" || codes[1] != "verified_user" {
		t.Errorf("unexpected role codes: %v", codes)
	}
}

func TestUserGetPermissions(t *testing.T) {
	u := &User{
		Roles: []*Role{
			{
				Code: "admin",
				Permissions: []*Permission{
					{Code: "person:read"},
					{Code: "person:write"},
				},
			},
			{
				Code: "verified_user",
				Permissions: []*Permission{
					{Code: "person:read"},
					{Code: "culture:read"},
				},
			},
		},
	}

	perms := u.GetPermissions()
	if len(perms) != 3 {
		t.Errorf("expected 3 unique permissions, got %d: %v", len(perms), perms)
	}
}

func TestUserHasPermission(t *testing.T) {
	u := &User{
		Roles: []*Role{
			{Code: "verified_user", Permissions: []*Permission{{Code: "person:read"}}},
		},
	}

	if !u.HasPermission("person:read") {
		t.Error("expected user to have person:read permission")
	}
	if u.HasPermission("person:write") {
		t.Error("expected user to not have person:write permission")
	}
}

func TestUserHasPermission_SuperAdmin(t *testing.T) {
	u := &User{
		Roles: []*Role{
			{Code: "super_admin", Permissions: []*Permission{}},
		},
	}

	if !u.HasPermission("any:permission") {
		t.Error("super_admin should have all permissions")
	}
}

func TestRoleValidate(t *testing.T) {
	tests := []struct {
		name    string
		role    *Role
		wantErr bool
	}{
		{
			name:    "valid role",
			role:    &Role{Code: "editor", Name: "Editor"},
			wantErr: false,
		},
		{
			name:    "empty code",
			role:    &Role{Code: "", Name: "Editor"},
			wantErr: true,
		},
		{
			name:    "code too long",
			role:    &Role{Code: string(make([]byte, 51)), Name: "Editor"},
			wantErr: true,
		},
		{
			name:    "empty name",
			role:    &Role{Code: "editor", Name: ""},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.role.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Role.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestPermissionValidate(t *testing.T) {
	tests := []struct {
		name    string
		perm    *Permission
		wantErr bool
	}{
		{
			name:    "valid permission",
			perm:    &Permission{Code: "test:read", Name: "Test Read", Module: "test"},
			wantErr: false,
		},
		{
			name:    "empty code",
			perm:    &Permission{Code: "", Name: "Test Read", Module: "test"},
			wantErr: true,
		},
		{
			name:    "code too long",
			perm:    &Permission{Code: string(make([]byte, 101)), Name: "Test Read", Module: "test"},
			wantErr: true,
		},
		{
			name:    "empty name",
			perm:    &Permission{Code: "test:read", Name: "", Module: "test"},
			wantErr: true,
		},
		{
			name:    "empty module",
			perm:    &Permission{Code: "test:read", Name: "Test Read", Module: ""},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.perm.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Permission.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
