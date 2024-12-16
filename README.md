# Yet another UUID library!?

The wonderful libraries by Google and Gofrs have served us quite well, however, they have two fatal flaws. First, they
use "Too Much Crypto" https://eprint.iacr.org/2019/1492.pdf. Second, ironically given the first, they can return errors.

The idiom to wrap every `New` in a  `Log(err)` (responsible), or `Must` (optimistic), is verbose, inefficient, and
possibly dangerous.

This library is opinionated about what UUIDs are worthwhile (v4 and v7), how you should handle errors when parsing or
unmarshalling (sentinel), and even which compact serializations are useful (NCName).

## But the crypto!

This library follows Go's `math/rand/v2` and Linux's `/dev/random` changes to use ChaCha20-based cryptographic
pseudorandom number generators to ensure error-free generation and speed. UUIDs are not cryptographic keys or secrets.

## But the errors!

Errors returned from unmarshalling functions are anonymous, message-free sentinels. With no text to translate or
sanitize they are functionally boolean: `nil` or not.

Boolean success and sentinel error returns free (require) you to handle parsing/unmarshalling failures your way.

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

# How To

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

## Short Serializations

The "hex-and-dash" encoding of a canonical UUID is already URL-safe and contains no ambiguous characters. Omitting the
dashes (which are positional anyway) gives you a short (32-runes), case-insensitive, URL-safe identifier string.

Sometimes an even shorter (but still non-binary) string is helpful. `uid` supports Compact UUIDs and ShortUUIDs.

### Compact UUIDs for Constrained Grammars (NCName)

`Parse` supports automatic detection and decoding of `UUID-NCName-32` and `UUID-NCName-64` compact encodings for
constrained grammars.

`UUID.Compact64()` and `UUID.Compact32()` return the Base64 and Base32 NCName encoded values, respectively.

These formats achieve or preserve the goals of compaction, URL-safety, and CSS/DOM identifier safety.

More info: https://datatracker.ietf.org/doc/draft-taylor-uuid-ncname/

### ShortUUID support

Python ShortUUID is problematic in multiple ways.

1. The common implementation accepts ANY alphabet (and padding) endangering transferability.
2. The encoding algorithm does not encode standard alphabets using standard mappings. If you "ShortUUID" encode using
the Base64 alphabet, you cannot Base64 decode the result back into the original bytes.
3. Base57 (default alphabet) has no other usage.
4. Optimizing for "manual human entry" is problematic in itself but Base57 still includes the `o` rune. The more
commonly used Base56 omits `o`.
5. Base57 alphabet ShortUUIDs may contain leading digits (often due to left-padding with `2`) making them unsuitable for
DOM and CSS identifiers without escaping.

Despite all this, it's a popular library and you may be interacting with a system that already uses them so the
following helpers

`FromPythonShort` enables decoding of Python ShortUUID encoded UUIDs using the default alphabet (Base57) and padding
(22).

`ToPythonShort` encodes a given `UUID` into a Python ShortUUID using the default alphabet (Base57) and padding (22).
