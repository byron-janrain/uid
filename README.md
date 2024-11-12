# Yet another UUID library!?

The wonderful libraries by Google and Gofrs have served us quite well, however, they have two fatal flaws. First, they
use "Too Much Crypto" https://eprint.iacr.org/2019/1492.pdf. Second, ironically given the first, they can return errors.

It's idiomatic to wrap every `New` in a  `Log(err)` (if you're responsible), or `Must` (if you're optimistic), but it's
verbose, inefficient, and possibly dangerous. All errors should be informative, safe, and actionable. I rarely see
UUID errors translated so they are effectively sentinels outside of English. Is it always okay to send the raw error to
the client considering the source of errors for a UUID constructor is a randomness unavailability or underflow? Finally,
what action would the caller take other than trying again until it works? Should they have backoff?

IMO a UUID library should be incapable of constructor failure and only return sentinel errors for parse failures.

`uid` constructs RFC compliant v4 and v7 UUIDs with no errors. Moreover, it parses without errors and, as a bonus, is an
order of magnitude faster on construction than Google/Gofrs, without adding any allocations.

This library follows Go's `math/rand/v2` and Linux's `/dev/random` changes to use ChaCha20-based cryptographic
pseudorandom number generators (CPRNG) to ensure no errors during random fills.

## But Unmarshal and Parse return errors!

`UnmarshalXXX` methods return errors because their interfaces require it. `uid.ParseError`, however, is purely a
sentinel error. Nothing to translate or sanitize. At worst you could check for it by type in an `errors.Is|As` chain,
which you're already doing as necessary.

Technically speaking, Parse's "ok" idiom is also an "error pattern". However, the expected usage of `uid.Parse` is
likely an API context (users should never be asked for a UUID) where input validation strictness needs only check for
validity before rejecting the entire request see the example below.

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

Because `uid` does not include error messages, you are free to observe and translate the failure (or not) as necessary.

```go
id, ok := uid.Parse(r.PathValue("id"))
if !ok {
    // observe it your way
    slog.Log("bad ID: %s", input)
    badIDCounter.Inc()
    // translate error messages your way
    http.Error(w, messagePrinter.Sprint("invalid ID"), http.StatusBadRequest)
    return
}
```
