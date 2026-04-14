package rules

import (
	"context"
	"strings"

	"buf.build/go/bufplugin/check"
	"google.golang.org/protobuf/reflect/protoreflect"
)

const patchRequestFieldMaskRuleID = "PATCH_REQUEST_FIELD_MASK"

func PatchRequestFieldMaskRule() *check.RuleSpec {
	return &check.RuleSpec{
		ID:      patchRequestFieldMaskRuleID,
		Default: true,
		Purpose: "Ensures PatchXxxRequest messages have a google.protobuf.FieldMask update_mask field to support partial updates.",
		Type:    check.RuleTypeLint,
		Handler: check.RuleHandlerFunc(handlePatchRequestFieldMask),
	}
}

func handlePatchRequestFieldMask(
	_ context.Context,
	responseWriter check.ResponseWriter,
	request check.Request,
) error {
	for _, fileDescriptor := range request.FileDescriptors() {
		if fileDescriptor.IsImport() {
			continue
		}

		messages := fileDescriptor.ProtoreflectFileDescriptor().Messages()
		for i := range messages.Len() {
			message := messages.Get(i)

			name := string(message.Name())
			if !strings.HasPrefix(name, "Patch") || !strings.HasSuffix(name, "Request") {
				continue
			}

			if !hasFieldMaskField(message) {
				responseWriter.AddAnnotation(
					check.WithDescriptor(message),
					check.WithMessagef(
						"patch request message %q must have a google.protobuf.FieldMask update_mask field to support partial updates",
						message.Name(),
					),
				)
			}
		}
	}

	return nil
}

func hasFieldMaskField(message protoreflect.MessageDescriptor) bool {
	fields := message.Fields()
	for i := range fields.Len() {
		field := fields.Get(i)
		if field.Name() != "patch_mask" {
			continue
		}

		return field.Kind() == protoreflect.MessageKind &&
			field.Message().FullName() == "google.protobuf.FieldMask"
	}

	return false
}
