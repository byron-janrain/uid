package uid_test

import (
	"strings"
	"testing"

	"github.com/hoodie-ninja/uid"
)

//nolint:cyclop // invariant checks
func FuzzParse(f *testing.F) {
	for _, seed := range []string{
		uid.NilCanonical, uid.MaxCanonical, ref4, ref7,
		`"` + ref4 + `"`, "{" + ref4 + "}",
		ref4b32, ref7b32, `"` + ref4b32 + `"`,
		ref4b64, ref7b64, `"` + ref4b64 + `"`,
		string(ref4Bytes), string(ref7Bytes),
	} {
		f.Add(seed)
	}
	f.Fuzz(func(t *testing.T, s string) {
		id, ok := uid.Parse(s) // must never panic
		if !ok {
			if !id.IsNil() {
				t.Fatalf("failed parse of %q returned non-nil %s", s, id)
			}
			return
		}
		// accepted ids satisfy the type invariant
		switch id.Version() {
		case uid.Version4, uid.Version7:
			if id.Variant() != uid.Variant9562 {
				t.Fatalf("parse of %q produced bad variant: %s", s, id)
			}
		case uid.VersionNil, uid.VersionMax:
		default:
			t.Fatalf("parse of %q produced bad version: %s", s, id)
		}
		// every encoding round-trips
		for _, enc := range []string{id.String(), id.Compact32(), id.Compact64(), string(id.Bytes())} {
			if rt, rok := uid.Parse(enc); !rok || rt != id {
				t.Fatalf("round trip of %q via %q failed", s, enc)
			}
		}
	})
}

func FuzzFromPythonShort(f *testing.F) {
	for _, seed := range []string{
		uid.NilPythonShort, uid.MaxPythonShort,
		"CXc85b4rqinB7s5J52TRYb", "zzzzzzzzzzzzzzzzzzzzzz", " CXc85b4rqinB7s5J52TRYb\t",
	} {
		f.Add(seed)
	}
	f.Fuzz(func(t *testing.T, s string) {
		id, ok := uid.FromPythonShort(s) // must never panic
		if !ok {
			return
		}
		if got, want := uid.ToPythonShort(id), strings.TrimSpace(s); got != want {
			t.Fatalf("round trip of %q: got %q", want, got)
		}
	})
}
