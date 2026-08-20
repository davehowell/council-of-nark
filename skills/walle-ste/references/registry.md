# Registry

The authoritative index. One row per active check: id, tier, citation, home, symptom. If any
other file disagrees with this table, this table wins and the other file is the bug. Anything
added to the catalog must merge with an existing check or replace one — and must land here.

## Active checks

| Id | Tier | Cite | Home | Symptom |
|---|---|---|---|---|
| A1 | M+J | 3.6 | kernel | Passive voice, agentful or agentless |
| A3 | M+J | 3.4 | checks-grammar | Modal passive, requirement with no actor |
| A5 | M+J | 3.7, 1.13 | kernel | Nominalization: "perform a migration of" |
| A6 | M+J | 1.2, 1.4, 9.2 | checks-words | Verbed noun, coined inflection |
| A7 | M+J | 1.12 | kernel | Umbrella verb, missing object or direction |
| A9 | M+J | 4.2 | checks-grammar | Dropped noun, verb, or subject |
| A10 | J | 3.6 M4 | checks-grammar | Document avoids `you` and `we` entirely |
| B1 | M+J | Part 2 2-0-19 | kernel | Modal that hides a fact |
| B2 | M+J | Part 2 `should` | checks-grammar | `should` stating a requirement |
| B3 | J | 4.1 | kernel | Abstract claim, no mechanism, no magnitude |
| B4 | M+J | 4.1, Part 2 | checks-grammar | Unquantified degree word |
| B5 | M+J | 4.1, 1.5 | checks-grammar | Comparative defects |
| B8 | M+J | 1.5, EXT | kernel | Unverifiable abstraction |
| B9 | M | Part 2 2-0-13 | kernel | Filler word |
| B10 | M+J | EXT | kernel | Reference point not stated |
| B12 | M+J | 4.1 | checks-grammar | Prohibition instead of an action |
| B13 | M+J | EXT | kernel | Mixed or missing unit and time convention |
| B14 | M+J | EXT | kernel | Rhetorical padding, closing summary |
| C1 | M+J | 5.1, 6.3, 4.4, 6.2 | kernel | Over-length sentence, bad split |
| C4 | M+J | 3.2, Part 2 2-0-7 | checks-grammar | Verb form outside the six permitted |
| C7 | M+J | 3.5 | checks-grammar | `-ing` verbal, chained, or dangling |
| C8 | M+J | 2.1, 2.2, 8.7 | kernel | Noun pile-up, hyphen chain |
| C9 | M+J | 5.4, GR-2 | kernel | Buried conditional in a step |
| C10 | M | 8.1 | checks-punctuation | Semicolon in prose |
| C11 | J | 6.1, 4.1 | checks-structure | Two topics in one sentence |
| C12 | M+J | 4.3, 6.4 | checks-grammar | Inline enumeration |
| C13 | M+J | 4.3 | checks-grammar | Broken vertical list |
| C14 | M+J | 4.3, 5#, 6# | checks-grammar | Facts and steps mixed in one list |
| C15 | M+J | 6.4–6.6, 9.4 | kernel | Paragraph and section structure |
| C18 | J | 8.3 | checks-punctuation | Parenthetical carrying a sentence |
| C21 | M+J | 6.2, 6.5 | checks-structure | Key word not carried forward |
| C24 | J | 9.1 | checks-structure | Prose duplicating a table |
| C25 | M | 4.2 | checks-grammar | Contraction |
| C27 | J | 4.3, 1.5 | checks-grammar | List items do not fit the lead-in |
| C28 | M+J | 8#, 5.4 | checks-punctuation | Comma permitting two readings |
| D1 | M+J | GR-3, GR-4, 4.4 | kernel | Unresolvable `this` or `it` |
| D3 | M+J | GR-1 | checks-punctuation | Dropped `that` |
| D4 | M+J | GR-2, 3.5 | checks-punctuation | Buried instrument: `with`, `using` |
| D6 | M+J | 9.2, Part 2 | checks-punctuation | Multi-sense ordinary word |
| D8 | M+J | 4.5, 4.2 | checks-grammar | Articles wrong in either direction |
| D10 | J | 2.1, 2.2, 8.2 | checks-words | Modifier attachment ambiguity |
| D11 | M | GR-6 | checks-punctuation | Latin abbreviation |
| D12 | M+J | GR-3, GR-7 | checks-punctuation | Gendered or exclusionary term |
| D13 | M+J | GR-8 | checks-punctuation | Possessive apostrophe |
| D14 | J | GR-5 | checks-punctuation | False friend |
| E1 | M+J | 1.11, 6.2 | kernel | Synonym rotation |
| E2 | M+J | 9.4 | checks-words | One action, many wordings |
| E3 | M+J | 1.3, 9.2 | checks-words | One word, two meanings |
| E4 | M+J | 1.6, 1.8, 2.2 | checks-words | Paraphrased canonical name |
| E5 | M+J | 2.2 M1, 8.3, 6.6 | kernel | Abbreviation failures |
| E6 | M+J | 9.3 | checks-words | Phrasal verb |
| E7 | M+J | 1.10 | checks-words | Slang, in-group jargon |
| E8 | M+J | 1.12, Part 2 | kernel | Inflated word |
| E13 | M+J | 1.5 note | checks-words | Capitalized common noun |
| E15 | M+J | 8.2 | checks-punctuation | Hyphenation |
| F1 | M+J | 5.3, 3.6 | kernel | Step not a clean imperative |
| F2 | M+J | 5.2 | checks-structure | Two actions in one step |
| F3 | M+J | 5.2 | checks-structure | Unnumbered steps |
| F4 | M+J | 5.5 | kernel | Load-bearing note |
| F5 | M+J | 5.5 | checks-structure | Requirement or danger in a note |
| F7 | M+J | 5.5, 6# | checks-structure | Imperative in descriptive text |
| G1 | M+J | 7# | kernel | Destructive operation, no warning |
| G2 | J | 7.1–7.3 | kernel | Safety block missing a slot (never skipped) |
| G3 | M+J | 7.1, 7.3 | kernel | Vague or softened warning |
| G5 | M+J | 7.1 | checks-safety | Wrong severity label |
| G6 | M+J | 7#, 7.1 | checks-safety | Label vocabulary and rendering |
| G7 | M+J | 7.1, 7.2 | checks-safety | Warning placement, split severity |
| G9 | M+J | 5.1, EXT | checks-safety | Over-long warning |
| G10 | J | 7# | checks-safety | Warning in a descriptive document |
| G11 | M | 4.3 | checks-safety | Prohibition only in the lead-in |

70 active checks: 23 kernel, 47 catalog.

## Retired

| Id | Fate | Why |
|---|---|---|
| C2 | Merged into C1 | Guard rails on splits belong to the split check |
| C16 | Merged into C1 | Reconnecting after a split belongs to the split check |
| E9 | Deleted | Six of twelve alternations were dead regex. A spell-checker locale owns spelling variants. Convention: US spelling, because the identifiers are US-spelled |

Never reuse a retired id.
