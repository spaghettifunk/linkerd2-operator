package certs

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenerateCA(t *testing.T) {
	it, err := GenerateTrustAnchorsCertificates("identity.linkerd.cluster.local")
	assert.NotNil(t, it)
	assert.NotNil(t, err)
}
