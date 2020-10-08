// Package jsonplan provides methods for parsing a JSON plan file into a Plan struct for testing.
package jsonplan

// While this module closely follows Terraform's internal plan structure,
// Hashicorp do not guarantee that the structure of a plan will not change
// between versions.

// Therefore, we maintain our own internal Plan structure which may differ in
// some ways from Terraform's internal representation.

// The types in here are a slightly modified version of the types in:
// https://github.com/hashicorp/terraform/tree/master/command/jsonplan
//
// Changes:
// - resource.go: "Index" keys are now type json.RawMessage instead of
//                Terraform's internal addrs.InstanceKey
