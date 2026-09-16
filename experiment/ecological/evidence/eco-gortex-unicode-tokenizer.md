# Evidence key: Gortex Unicode tokenizer panic

> Rater/controller material. Never mount this directory, `candidates.json`, snapshot `controller/`, or upstream evidence in a respondent process. The merged patch is strong evidence, not the only acceptable remedy.

## Frozen task

- Task: `eco-gortex-unicode-tokenizer`
- Parent: `a1c26e19aa9dd8d840b0128bd470b83a98fdb91f`
- Evidence: `7ff32d7b7be63f2aa05be74ff4dcf3ce63714207`
- Primary parent location: `internal/search/rerank/tokens.go`, lookahead in `tokenize`
- Controller-focused target: `go test -count=1 ./internal/search/rerank -run '^TestTokenize$'`
- Observed parent result: exit 1 with `index out of range [1] with length 1`
- Observed evidence result: exit 0

## Reference resolution shown to raters

- Original repository: `zzet/gortex`
- Merged pull request: [#569](https://github.com/zzet/gortex/pull/569)
- Title: *fix(search): stop the rerank tokenizer panicking on trailing multi-byte uppercase runes*

Succinct original-repository account: the SCREAMING-to-Camel lookahead compared a byte offset with a byte length before indexing a rune slice. A trailing multi-byte uppercase rune could pass that guard even though the suffix contained only one rune. Lowercase query text could still trigger the panic because reranking tokenises retrieved candidate text. The applied fix decodes the next rune from the correct UTF-8 boundary and avoids materialising the whole suffix. The merged tests cover the panic inputs, a multi-byte acronym-to-Camel split, and established ASCII behavior.

The blinded rating bundle includes the exact parent-to-evidence patch and changed test diff, generated from the two frozen commits above. Raters are told that this is an accepted baseline, not a gold standard. Equivalent idiomatic remedies and supported findings beyond PR #569 remain eligible.

## Primary mechanism

A Go `range` over a string yields `i` as a **byte offset** and `r` as a rune. In the SCREAMING-to-Camel branch, the parent checks `i+1 < len(s)`, where `len(s)` is also bytes, and then evaluates:

```go
next := []rune(s[i:])[1]
```

The check proves only that at least one byte follows the first byte position. It does not prove that a second rune follows the current rune. If the current uppercase rune is multibyte and is the final rune, `i+1 < len(s)` can be true because `i+1` points inside that rune's UTF-8 encoding. `[]rune(s[i:])` then has length one, and index 1 panics. The whole-suffix conversion also allocates on a tokenisation hot path.

The previous rune must be uppercase to enter this lookahead branch. A representative minimal class is a token ending in consecutive uppercase runes with a non-ASCII uppercase final rune.

## Why query text is not sufficient

The rerank package tokenises more than the query. Signal calculation tokenises candidate signatures, names, and qualified names while scoring repository-derived candidates. A lowercase Cyrillic query does not itself satisfy the uppercase branch, but it can retrieve or score corpus text that does. Whether a request panics can therefore depend on indexed repository content and the candidates reached by the query. Claims that the index must be corrupt, or that lowercase query text alone directly enters the faulty branch, are unsupported.

## Correction evidence

The merged correction imports `unicode/utf8` and replaces suffix-wide rune conversion with boundary-aware decoding:

```go
if unicode.IsUpper(r) && unicode.IsUpper(prev) {
    next, size := utf8.DecodeRuneInString(s[i+utf8.RuneLen(r):])
    if size > 0 && unicode.IsLower(next) {
        flush()
    }
}
```

This advances from the current byte offset by the current rune's UTF-8 width, then decodes only the next rune. At end of string, `DecodeRuneInString("")` returns size zero, so no lookahead is used. It preserves the acronym-to-Camel split and avoids allocating `[]rune` for every uppercase lookahead.

Equivalent remedies are acceptable if they:

- never treat a byte offset or byte count as a rune index/count;
- safely represent the absence of a next rune;
- preserve ASCII and Unicode acronym/camel boundaries;
- avoid repeated whole-suffix or whole-input rune allocation in the hot path;
- remain scoped to tokenisation unless broader changes are independently justified.

A one-time `[]rune(s)` conversion can be functionally correct but does not fully satisfy the brief's allocation constraint. A special case for Cyrillic or for the exact panic string is not Unicode-correct.

## Discriminating regression evidence

The merged evidence test includes:

| Input | Expected tokens | Purpose |
|---|---|---|
| `""` | `nil` | Existing empty behavior. |
| `"ParseHTTPHeader"` | `parse`, `http`, `header` | Existing ASCII camel/acronym behavior. |
| `"validate_user_token"` | `validate`, `user`, `token` | Existing separator behavior. |
| `"HTTPHeader"` | `http`, `header` | Existing SCREAMING-to-Camel split. |
| `"простой текст"` | `простой`, `текст` | Lowercase Cyrillic and separator behavior. |
| `"ТЕКСТ"` | `текст` | Direct pre-fix panic at a multibyte uppercase ending. |
| `"ПРОСТОЙ ТЕКСТ"` | `простой`, `текст` | Panic class across a separator. |
| `"CAFÉ"` | `café` | Latin multibyte uppercase ending. |
| `"ЖКХStatus"` | `жкх`, `status` | Unicode acronym-to-ASCII Camel lookahead remains correct. |

A strong alternative test matrix need not copy these literals. It must include a pre-fix-failing multibyte uppercase end, a real next-rune case that must still split, and controls for established behavior.

## Supported novel findings

Rate additional findings on their own evidence. Do not reject one merely because PR 569 did not change it. De-duplicate restatements of the byte/rune mismatch as one primary finding. Suggestions about malformed UTF-8, alternative state machines, or broader tokenizer performance receive credit only when the output explains current Go behavior and proposes a scoped, verifiable test.
