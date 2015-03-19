package framework

import (
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/vault/logical"
)

// Secret is a type of secret that can be returned from a backend.
type Secret struct {
	// Type is the name of this secret type. This is used to setup the
	// vault ID and to look up the proper secret structure when revocation/
	// renewal happens. Once this is set this should not be changed.
	//
	// The format of this must match (case insensitive): ^a-Z0-9_$
	Type string

	// Fields is the mapping of data fields and schema that comprise
	// the structure of this secret.
	Fields map[string]*FieldSchema

	// Renewable is whether or not this secret type can be renewed.
	Renewable bool

	// DefaultDuration and DefaultGracePeriod are the default values for
	// the duration of the lease for this secret and its grace period. These
	// can be manually overwritten with the result of Response().
	DefaultDuration    time.Duration
	DefaultGracePeriod time.Duration
}

// SecretType is the type of the secret with the given ID.
func SecretType(id string) string {
	idx := strings.Index(id, "-")
	if idx < 0 {
		return ""
	}

	return id[:idx]
}

func (s *Secret) Response(
	data map[string]interface{}) (*logical.Response, error) {
	uuid, err := logical.UUID()
	if err != nil {
		return nil, err
	}

	id := fmt.Sprintf("%s-%s", s.Type, uuid)
	return &logical.Response{
		IsSecret: true,
		Lease: &logical.Lease{
			VaultID:     id,
			Renewable:   s.Renewable,
			Duration:    s.DefaultDuration,
			GracePeriod: s.DefaultGracePeriod,
		},
		Data: data,
	}, nil
}