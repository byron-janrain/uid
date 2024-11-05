# UID

It's always bugged me that the libraries in widespread use for generating UUIDs can return errors. Ultimately this is
because they are using "Too Much Crypto" https://eprint.iacr.org/2019/1492.pdf.

While it's idiomatic, wrapping every `New` in a `Must`, if you're optimistic, or `Log(err)`, if you are responsible. Is
slow and, at worst, dangerous. Draining `/dev/random` can have nasty side effects to every randomness or introduce
races or latencies. Using "strong" crypto is wasteful for generating non-secrets that you don't verify.

I present you yet another UUID library: `uid`.

`uid` constructs RFC compliant v4 and v7 UUIDs an order of magnitude faster than Google/Gofrs, with zero allocations,
and, most importantly, no errors.

** Strictly speaking, the parser indicates a parsing error, but instead of making you juggle an esoteric log line, it
returns bool.

## How To
