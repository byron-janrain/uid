# Yet another UUID library!?

The wonderful libraries by Google and Gofrs have served us quite well, however, they have two fatal flaws. First, they
use "Too Much Crypto" https://eprint.iacr.org/2019/1492.pdf. Second, ironically given the first, they can return errors.

It's idiomatic to wrap every `New` in a  `Log(err)` (if you're responsible), or `Must` (if you're optimistic), but it's
verbose, inefficient, and possibly dangerous. All errors should be informative, safe, and actionable.

Is it always okay to send the raw error to the client? Should they know there is a randomness underflow? Untranslated
(non-internationalized) errors are only sentinels to unfluent readers anyway. What action (besides logging or
panicking) can the caller take with a "bad" UUID? Retry until it works? Will they properly back off?

This library constructs UUIDs without errors. Parse/validation failure is explicitly the sentinel it's always been.

`uid` constructs RFC compliant v4 and v7 UUIDs with no errors. Moreover, it parses without errors and, as a bonus, is an
order of magnitude faster on construction than Google/Gofrs.

This library follows Go's `math/rand/v2` and Linux's `/dev/random` changes to use ChaCha20-based cryptographic
pseudorandom number generators (CPRNG) to ensure no errors during random fills.

## What about Short UUIDs?

`uid.UUID` implements compact UUIDs encoded per https://datatracker.ietf.org/doc/draft-taylor-uuid-ncname/. Specifically
Base32 (shorter and case-insensitive) and Base64url (shortest but case-sensitive). Arbitrary shorteners and any
library based on the Python `shortuuid` algorithm still can produce encodings with leading digits which are prohibited in
DOM IDs and require escaping in CSS classes. NCName-encoded UUIDs are safe for use in CSS classes, DOM IDs, and URIs
without escaping.

## But Unmarshal and Parse return errors!

`UnmarshalXXX` methods return errors because the implemented interfaces require it. `uid.ParseError`, however, is purely
a sentinel error. No text to translate or sanitize. Check it for `nil` and move on like you would have before, "bad UUID"
is probably not your service's primary concern.

Technically speaking, `Parse`'s "ok" idiom is an error "pattern". However, the expected usage of `uid.Parse` is
in an API context (users should never be asked to type a UUID) where input validation strictness needs only check for
validity before rejecting the entire request see the example below.

# How To

## Generate

New Random UUID (v4)...
```go
id := uid.NewV4()
```

New Sortable UUID (v7 with "method 3", extended precision monotonicity)
```go
id := uid.NewV7()
```

New Sortable UUID (v7 with "method 3" monotonicity and strict process-local uniqueness... You don't want this, but it's
here if you need it.)
```go
id := uid.NewV7Strict()
```

## Parse

Because `uid` does not include error messages, you are free to handle (count, log, translate, etc...) the failure (or
not) your way.

```go
id, ok := uid.Parse(r.PathValue("id"))
if !ok {
    // observe it your way
    slog.Log("bad ID: %s", sanitizeForLog(input))
    badIDCounter.Inc()
    // translate responses your way
    http.Error(w, messagePrinter.Sprint("invalid ID"), http.StatusBadRequest)
    return
}
```
