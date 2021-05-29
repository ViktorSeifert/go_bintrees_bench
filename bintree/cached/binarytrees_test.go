// Tests for program #8.
//
// Mainly to ensure that derived programs produce the same result.

package main

import (
	"regexp"
	"strings"
	"testing"
)

func TestResultsAreAsBenchmark(t *testing.T) {
	var data = []struct {
		depth uint32
		iter uint32
		expect uint32
	}{
		{4, 2097152, 65011712},
		{6, 524288, 66584576},
		{8, 131072, 66977792},
	}

	for _, d := range(data) {
		if r := innerImpl(d.depth, d.iter); r != d.expect {
			t.Errorf("innerImpl(%d, %d) = %d", d.depth, d.iter, r)
		}
	}
}

const bench21Output = `stretch tree of depth 22         check: 8388607
2097152  trees of depth 4        check: 65011712
524288   trees of depth 6        check: 66584576
131072   trees of depth 8        check: 66977792
32768    trees of depth 10       check: 67076096
8192     trees of depth 12       check: 67100672
2048     trees of depth 14       check: 67106816
512      trees of depth 16       check: 67108352
128      trees of depth 18       check: 67108736
32       trees of depth 20       check: 67108832
long lived tree of depth 21      check: 4194303`

func TestOutputSameAsOriginalBenchmark(t *testing.T) {
	r := run(21)
	if normalize(r) != normalize(bench21Output) {
		t.Errorf("Expected:\n%s\nbut got:\n%s\n", bench21Output, r)
	}
}

var re = regexp.MustCompile(`\s+`)
func normalize(in string) string {
	return strings.TrimSpace(re.ReplaceAllLiteralString(in, " "))
}
