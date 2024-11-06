# Yet another UUID library!?

It's always bugged me that the libraries in widespread use for generating UUIDs can return errors. Ultimately this is
because they are using "Too Much Crypto" https://eprint.iacr.org/2019/1492.pdf.

While it's idiomatic to wrap every `New` in a `Must` (if you're optimistic), or `Log(err)` (if you are responsible),
it's wasteful, slow, and possibly dangerous.

`uid` constructs RFC compliant v4 and v7 UUIDs an order of magnitude faster than Google/Gofrs, with zero allocations,
and, most importantly, no errors.

Essentially, `uid` is just following Go's `math/rand/v2` and Linux's `/dev/random` changes to use ChaCha20-based
cryptographic pseudorandom number generators (CPRNG). If it's good enough for them, it's good enough for you.

As the benchmarks show, it's way faster, does not allocate, and cannot return errors.

** Strictly speaking, `UnmarshalXXX` methods return errors because their interfaces require it.

** Technically speaking, Parse returns errors if you squint hard enough at that `bool`. The nice thing about returning
`bool` from Parse is that validation is just  `if _, ok := Parse(input); ok {...}`.

# How To

New Random (v4) UUID...
```go
id := uid.NewV4()
```

New Sortable (v7) UUID (with "method 3", extended precision monotonicity)
```go
id := uid.NewV7()
```

New Sortable (v7) UUID with strict local monotonicity. You don't want this.
```go
id := uid.NewV7Strict()
```

# TODO
Optimize parser.
