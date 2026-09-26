package service

import (
	"testing"

	"github.com/assethub/assethub/internal/constants"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestCanModerate(t *testing.T) {
	t.Parallel()
	tests := []struct {
		role string
		want bool
	}{
		{role: string(constants.RoleAdmin), want: true},
		{role: string(constants.RoleModerator), want: true},
		{role: string(constants.RoleUploader), want: false},
		{role: string(constants.RoleViewer), want: false},
		{role: "", want: false},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.role, func(t *testing.T) {
			t.Parallel()
			if got := canModerate(tt.role); got != tt.want {
				t.Fatalf("canModerate(%q) = %v, want %v", tt.role, got, tt.want)
			}
		})
	}
}

func TestContainsObjectID(t *testing.T) {
	t.Parallel()
	a := primitive.NewObjectID()
	b := primitive.NewObjectID()
	tests := []struct {
		name   string
		list   []primitive.ObjectID
		target primitive.ObjectID
		want   bool
	}{
		{name: "present first", list: []primitive.ObjectID{a, b}, target: a, want: true},
		{name: "present last", list: []primitive.ObjectID{a, b}, target: b, want: true},
		{name: "missing", list: []primitive.ObjectID{a}, target: b, want: false},
		{name: "empty", list: nil, target: a, want: false},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if got := containsObjectID(tt.list, tt.target); got != tt.want {
				t.Fatalf("containsObjectID() = %v, want %v", got, tt.want)
			}
		})
	}
}
