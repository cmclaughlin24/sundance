package spec

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/getkin/kin-openapi/openapi2"
	"github.com/getkin/kin-openapi/openapi2conv"
	"github.com/getkin/kin-openapi/openapi3"
)

func parseV2(ctx context.Context, data []byte) (*openapi3.T, error) {
	var doc2 openapi2.T
	if err := json.Unmarshal(data, &doc2); err != nil {
		return nil, fmt.Errorf("failed to unmarshal openapi v2 spec: %w", err)
	}

	doc3, err := openapi2conv.ToV3(&doc2)
	if err != nil {
		return nil, fmt.Errorf("failed to convert openapi v2 spec to v3: %w", err)
	}

	return doc3, nil
}

func parseV3x(ctx context.Context, data []byte) (*openapi3.T, error) {
	loader := openapi3.NewLoader()
	loader.IsExternalRefsAllowed = true

	doc, err := loader.LoadFromData(data)
	if err != nil {
		return nil, fmt.Errorf("failed to parse openapi v3 spec: %w", err)
	}

	if err := doc.Validate(ctx); err != nil {
		// FIXME: Log a warning that the document failed validation and may cause issues with registered tools.
	}

	return doc, nil
}
