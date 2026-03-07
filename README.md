# transitext

`transitext` is a Go module for translation pipelines with a single public
contract and multiple backend providers.

The core package gives you stable request/response models, validation,
chunking, provider registry, and reusable wrappers for retry, fallback,
rate-limiting, and caching. Provider packages focus on transport and protocol
specific logic, so integrations can switch backends without changing pipeline
code.

The module is designed for build tools, localization automation, batch
translation jobs, and CLI utilities where you need predictable behavior,
explicit errors, and easy runtime configuration.

## Install

```bash
go get github.com/woozymasta/transitext
```

## Configuration Reference

If you use `transitext` from external config files, use these artifacts:

* [CONFIG.md](CONFIG.md) for human-readable parameter reference.
* [schema.json](schema.json) for JSON Schema validation and editor tooling.
* [config package](config/types.go) (`github.com/woozymasta/transitext/config`)
  for reusable typed config contracts in your application code.

## Providers

`transitext` supports both production-oriented API providers and experimental
web/free backends.

Paid or API-key providers (recommended for production):

* `google`
* `azure`
* `deepl`
* `yandex`
* `libre` (self-hosted or managed instance)
* `openai` (OpenAI-compatible chat/completions translation flow)

Free or unofficial providers (experimental):

* `googlefree`
* `bingfree`
* `deeplfree`
* `microsoft`
* `yandexfree`

> [!WARNING]  
> Free/unofficial providers rely on undocumented web endpoints and may stop
> working at any time without notice. Behavior, limits, and response formats
> can change unexpectedly. Use them at your own risk, with no guarantees.
> For production workloads, use official API providers with keys.

## OpenAI-Compatible Backends

The `openai` provider is not limited to OpenAI itself. It works with many
OpenAI-compatible servers and gateways, for example Ollama, LM Studio,
vLLM-based gateways, DeepSeek-compatible endpoints, Together AI-compatible
endpoints, and similar API layers that expose OpenAI-style chat interfaces.

## Request Tuning for Free Providers

Free providers expose a shared `request` block for low-level HTTP shaping.
You can override `user_agent`, add custom `headers`, and provide `cookies`
without patching or rebuilding the module. Each provider still keeps its own
safe defaults.

```yaml
provider: deeplfree
options:
  request:
    user_agent: "Mozilla/5.0 (Windows NT 10.0; Win64; x64)"
    headers:
      Accept-Language: "en-US,en;q=0.9"
    cookies:
      dl_session: "..."
```

## Partial Results and Item Retry

If you need to stop wasting API budget on full-batch retries, wrap provider
with `PartialTranslator`. It translates items one-by-one, retries one failed
item (`item_retries`) instead of replaying the whole batch, and keeps
per-item errors in `result.items[i].error`.

Default behavior is strict:

* stop immediately on temporary/system errors;
* return already translated items plus error for failed and not-processed rows.

```go
pipeline := transitext.Wrap(baseTranslator).
    Partial(transitext.PartialOptions{
        ItemRetries: 1,
    }).
    Build()
```

## Context Passing

> [!NOTE]  
> This is an experimental hack. It can improve disambiguation in some cases,
> but it can also degrade output or be ignored by provider models. Treat it
> as best-effort behavior only.

For ambiguous words you can wrap any translator with `ContextTranslator`.
The wrapper injects context and source text into symbolic marker blocks, then
strips marker envelope from translated output.

Practical note: context text works better when it is written in the same
language as `request.SourceLang`. Mixed-language context can be ignored or
misinterpreted by some providers.

Observed behavior from live free-provider probes in this workspace:

* `deeplfree` and `yandexfree` reacted to context more often.
* `googlefree` reacted in some cases, but consistency was lower.
* `microsoft` and `bingfree` were mostly insensitive to context
  in tested pairs.

So this feature is useful for experimentation and batch prefill, but not a
reliable semantic control mechanism for production.

```text
source text: замок
target: en

context: "in the context of a door"     -> lock
context: "in the context of a building" -> castle
```

```go
pipeline := transitext.Wrap(baseTranslator).
    Context(transitext.ContextOptions{
        ContextByID: map[string]string{
            "item-door":     "in the context of a door",
            "item-building": "in the context of a building",
        },
    }).
    Build()
```

## Minimal Example

```go
package main

import (
    "context"
    "log"

    "github.com/woozymasta/transitext"
    "github.com/woozymasta/transitext/providers"
)

func main() {
    registry, err := providers.NewDefaultRegistry()
    if err != nil {
        log.Fatal(err)
    }

    translator, err := registry.Build("googlefree", nil)
    if err != nil {
        log.Fatal(err)
    }

    result, err := translator.Translate(context.Background(), transitext.Request{
        SourceLang: "en",
        TargetLang: "ru",
        Items: []transitext.Item{
            {ID: "greeting", Text: "Hello world"},
        },
    })
    if err != nil {
        log.Fatal(err)
    }

    log.Printf("%s => %s", result.Items[0].ID, result.Items[0].Text)
}
```

## Language Mapping in Translation Flow

`langmap` helps when input language names come from config, UI, or legacy
tools and are not guaranteed to be strict provider codes.

```go
package main

import (
    "context"
    "log"

    "github.com/woozymasta/transitext"
    "github.com/woozymasta/transitext/langmap"
    "github.com/woozymasta/transitext/providers"
)

func main() {
    registry, _ := providers.NewDefaultRegistry()
    translator, _ := registry.Build("deepl", nil)

    source, ok := langmap.ResolveForProvider("deepl", "english")
    if !ok {
        log.Fatal("unknown source language")
    }
    target, ok := langmap.ResolveForProvider("deepl", "chinesesimp")
    if !ok {
        log.Fatal("unknown target language")
    }

    result, err := translator.Translate(context.Background(), transitext.Request{
        SourceLang: source,
        TargetLang: target,
        Items:      []transitext.Item{{ID: "1", Text: "About"}},
    })
    if err != nil {
        log.Fatal(err)
    }

    log.Println(result.Items[0].Text)
}
```

You can also validate support before translation and select only compatible
providers for current language:

```go
code, ok := langmap.SupportedByProvider("deepl", "alapmunte")
if !ok {
    log.Fatal("deepl does not support requested language")
}

providers := langmap.SupportingProviders("ukrainian")
log.Printf("providers supporting ukrainian: %v", providers)
_ = code
```

## Testing Notes

Network integration checks for free providers are part of the provider test
suite and are skipped in short mode (`go test -short ./...`).
