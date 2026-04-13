package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"

	schema "github.com/open-feature/flagd-schemas/json"
	"github.com/open-feature/open-feature-operator/apis/core/v1beta1"
	"github.com/xeipuuv/gojsonschema"
	"gopkg.in/yaml.v3"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run main.go <path-to-featureflags>")
		os.Exit(1)
	}

	path := os.Args[1]

	featureFlags, err := loadFeatureFlagsMultiDoc(path)
	if err != nil {
		fmt.Printf("❌ Error loading feature flags: %v\n", err)
		os.Exit(1)
	}

	for _, ff := range featureFlags {
		err := validateFeatureFlagFlags(ff.Spec.FlagSpec.Flags)
		if err != nil {
			fmt.Printf("❌ FeatureFlag %s validation failed: %v\n", ff.Name, err)
			os.Exit(1)
		} else {
			fmt.Printf("✅ FeatureFlag %s is valid.\n", ff.Name)
		}
	}
}

func loadFeatureFlagsMultiDoc(path string) ([]v1beta1.FeatureFlag, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	decoder := yaml.NewDecoder(f)
	var featureFlags []v1beta1.FeatureFlag

	for {
		var doc map[string]interface{}
		if err := decoder.Decode(&doc); err != nil {
			if err == io.EOF {
				break
			}
			return nil, fmt.Errorf("failed to decode YAML: %w", err)
		}

		u := &unstructured.Unstructured{Object: doc}

		if u.GetKind() != "FeatureFlag" {
			continue
		}

		bytes, err := json.Marshal(u.Object)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal unstructured to JSON: %w", err)
		}

		var ff v1beta1.FeatureFlag
		if err := json.Unmarshal(bytes, &ff); err != nil {
			return nil, fmt.Errorf("failed to unmarshal JSON to FeatureFlag: %w", err)
		}

		featureFlags = append(featureFlags, ff)
	}

	return featureFlags, nil
}

func initSchemas() (*gojsonschema.Schema, error) {
	schemaLoader := gojsonschema.NewSchemaLoader()
	if err := schemaLoader.AddSchemas(gojsonschema.NewStringLoader(schema.TargetingSchema)); err != nil {
		return nil, err
	}
	return schemaLoader.Compile(gojsonschema.NewStringLoader(schema.FlagSchema))
}

func validateFeatureFlagFlags(flags v1beta1.Flags) error {
	b, err := json.Marshal(flags)
	if err != nil {
		return err
	}

	documentLoader := gojsonschema.NewStringLoader(string(b))
	compiledSchema, err := initSchemas()
	if err != nil {
		return fmt.Errorf("unable to initialize Schema: %w", err)
	}

	result, err := compiledSchema.Validate(documentLoader)
	if err != nil {
		return err
	}

	if !result.Valid() {
		errStr := ""
		for _, desc := range result.Errors() {
			errStr += fmt.Sprintf("- %s\n", desc.String())
		}

		return fmt.Errorf(errStr)
	}
	return nil
}
