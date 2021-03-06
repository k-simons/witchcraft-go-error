package werror

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

const stackTraceString = ".*github.com/palantir/witchcraft-go-error.TestErrorFormatting\n" +
	".*github.com/palantir/witchcraft-go-error/werror_printer_test.*\n" +
	"testing.tRunner\n" +
	".*src/testing/testing.go.*\n" +
	"runtime.goexit\n" +
	".*src/runtime.*"

func TestErrorFormatting(t *testing.T) {
	for _, currCase := range []struct {
		name                    string
		err                     error
		expectedRegex           string
		outputEveryCallingStack bool
	}{
		{
			name: "simple error",
			err:  Error("simple_error"),
			expectedRegex: "" +
				"simple_error\n\n" +
				stackTraceString,
		},
		{
			name: "simple error with param",
			err:  Error("simple_error", SafeParam("safeParamKey", "safeParamValue")),
			expectedRegex: "" +
				"simple_error safeParamKey:safeParamValue\n\n" +
				stackTraceString,
		},
		{
			name: "simple error with many params",
			err:  Error("simple_error", SafeParam("safeParamKey", "safeParamValue"), SafeParam("safeParamKey2", "safeParamValue2")),
			expectedRegex: "" +
				"simple_error safeParamKey:safeParamValue, safeParamKey2:safeParamValue2\n\n" +
				stackTraceString,
		},
		{
			name: "simple wrapped error",
			err:  Wrap(Error("simple_error"), "simple_error_2"),
			expectedRegex: "" +
				"simple_error_2\n" +
				"simple_error\n\n" +
				stackTraceString,
		},
		{
			name: "simple wrapped error with forced stacks",
			err:  Wrap(Error("simple_error"), "simple_error_2"),
			expectedRegex: "" +
				"simple_error_2\n" +
				"simple_error\n\n" +
				stackTraceString +
				"\n" +
				stackTraceString,
			outputEveryCallingStack: true,
		},
		{
			name: "double wrapped error with params",
			err: Wrap(Wrap(
				Error("inner0Message", SafeParam("inner0ParamKey", "inner0VParamValue"), SafeParam("inner0ParamKey1", "inner0VParamValue1"), SafeParam("inner0ParamKey2", "inner0VParamValue2")),
				"inner1Message", SafeParam("inner1ParamKey", "inner1ValueKey")), "inner2Message"),
			expectedRegex: "" +
				"inner2Message\n" +
				"inner1Message inner1ParamKey:inner1ValueKey\n" +
				"inner0Message inner0ParamKey:inner0VParamValue, inner0ParamKey1:inner0VParamValue1, inner0ParamKey2:inner0VParamValue2\n\n" +
				stackTraceString,
		},
	} {
		t.Run(currCase.name, func(t *testing.T) {
			assert.Regexp(t, currCase.expectedRegex, GenerateErrorString(currCase.err, currCase.outputEveryCallingStack))
		})
	}
}
