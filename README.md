# Yet another UUID library!?

It's always bugged me that the libraries in widespread use for generating UUIDs can return errors. Ultimately this is
because they are using "Too Much Crypto" https://eprint.iacr.org/2019/1492.pdf.

While it's idiomatic to wrap every `New` in a  `Log(err)` (if you're responsible), or `Must` (if you're optimistic),
it's bloaty, slow, and possibly dangerous.

`uid` constructs RFC compliant v4 and v7 UUIDs an order of magnitude faster than Google/Gofrs, with zero allocations,
and, most importantly, no errors.

Essentially, `uid` is just following Go's `math/rand/v2` and Linux's `/dev/random` changes to use ChaCha20-based
cryptographic pseudorandom number generators (CPRNG). If it's good enough for them, it's good enough for you.

As the benchmarks show, generating IDs is way faster, but still does not allocate, and cannot return errors.
Parsing IDs is comparably fast, still does not allocate, and returns blind success indicators. It's up to your input
validation to take the failure and make it observable or craft a safe and appropriate response to the caller.

** Strictly speaking, `UnmarshalXXX` methods return errors because their interfaces require it. But `uid.ParseError` is
a purely sentinel error. No string to translate, sanitize, or unwrap.

** Technically speaking, Parse also returns "errors" if you squint hard enough at that `bool`. However, in addition to
being inherently safe, simple validations are simply `if !ok` instead of the `if errors.Is`, `if errors.As`, or
`if err != nil` patterns.

# How To

## Generate

New Random (v4) UUID...
```go
id := uid.NewV4()
```

New Sortable (v7) UUID (with "method 3", extended precision monotonicity)
```go
id := uid.NewV7()
```

New Sortable (v7) UUID with strict local (still method 3) monotonicity. You don't want this.
```go
id := uid.NewV7Strict()
```

## Parse

```go
id, ok := uid.Parse(input)
if !ok {
    slog.Log("bad input: %s", input) // observe it your way
    // you can't accidentaly leak an error if you don't have one.
    http.Error(w, "bad input", http.StatusBadRequest)
    return
}
```
