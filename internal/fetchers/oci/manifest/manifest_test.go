// Copyright The Enterprise Contract Contributors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//
// SPDX-License-Identifier: Apache-2.0

//go:build unit

package manifest

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/google/go-containerregistry/pkg/crane"
	"github.com/google/go-containerregistry/pkg/name"
	"github.com/google/go-containerregistry/pkg/v1/remote"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/enterprise-contract/ec-cli/internal/fetchers/oci"
	"github.com/enterprise-contract/ec-cli/internal/fetchers/oci/fake"
)

func TestImageManifests(t *testing.T) {
	ref := name.MustParseReference("registry.io/repository/image:tag")

	image, err := crane.Image(map[string][]byte{
		"autoexec.bat":              []byte(`@ECHO OFF`),
		"manifests/a.json":          []byte(`{"a":1}`),
		"manifests/b.yaml":          []byte(`b: 2`),
		"manifests/c.xml":           []byte(`<?xml version="1.0" encoding="UTF-8"?>`),
		"manifests/unreadable.yaml": []byte(`***`),
		"manifests/unreadable.json": []byte(`***`),
	})
	require.NoError(t, err)

	var opts []remote.Option = nil

	client := fake.FakeClient{}
	client.On("Image", ref, opts).Return(image, nil)

	ctx := oci.WithClient(context.Background(), &client)

	manifests, err := ImageManifests(ctx, ref)

	assert.NoError(t, err)

	assert.Equal(t, map[string]json.RawMessage{
		"manifests/a.json": []byte(`{"a":1}`),
		"manifests/b.yaml": []byte(`{"b":2}`),
	}, manifests)
}
