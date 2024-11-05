//nolint:errcheck // benchmarks
package uid_test

import (
	"testing"

	"github.com/byron-janrain/uid"
	gofrsuuid "github.com/gofrs/uuid"
	googleuuid "github.com/google/uuid"
)

func BenchmarkV4(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = uid.NewV4()
	}
}

func BenchmarkGoogleV4(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, _ = googleuuid.NewRandom() // ignoring error is best-case for performance comparison but don't
	}
}

func BenchmarkGofrsV4(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, _ = gofrsuuid.NewV4() // ignoring error is best-case for performance comparison but don't
	}
}

func BenchmarkV7(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = uid.NewV7()
	}
}

func BenchmarkV7Batch(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = uid.NewV7Batch()
	}
}

func BenchmarkGoogleV7(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, _ = googleuuid.NewV7() // ignoring error is unrealistic but best-case for performance comparison
	}
}

func BenchmarkGofrsV7(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, _ = gofrsuuid.NewV7() // ignoring error is unrealistic but best-case for performance comparison
	}
}
