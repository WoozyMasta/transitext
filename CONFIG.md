<!-- Automatically generated file, do not modify! -->

# transitext config reference

* Source file: [`schema.json`](https://github.com/woozymasta/transitext/blob/HEAD/schema.json)
* Source URL: [Raw schema URL](https://raw.githubusercontent.com/woozymasta/transitext/HEAD/schema.json)
* Schema identifier: `https://github.com/woozymasta/transitext/config/config`
* JSON Schema version: `https://json-schema.org/draft/2020-12/schema`
* Version support: `supported (2020-12)`
* Root reference: `#/$defs/Config`

## Contents

* [Config](#config)
  * [Providers](#providers)
    * [azure_Options](#azure_options)
    * [bingfree_Options](#bingfree_options)
      * [transitext_HTTPRequestOptions](#transitext_httprequestoptions)
    * [deepl_Options](#deepl_options)
    * [deeplfree_Options](#deeplfree_options)
    * [google_Options](#google_options)
    * [googlefree_Options](#googlefree_options)
    * [libre_Options](#libre_options)
    * [microsoft_Options](#microsoft_options)
    * [openai_Options](#openai_options)
    * [yandex_Options](#yandex_options)
    * [yandexfree_Options](#yandexfree_options)
  * [Wrappers](#wrappers)
    * [transitext_CacheOptions](#transitext_cacheoptions)
    * [transitext_ContextOptions](#transitext_contextoptions)
    * [transitext_LongTextOptions](#transitext_longtextoptions)
    * [transitext_PartialOptions](#transitext_partialoptions)
    * [transitext_RateLimitOptions](#transitext_ratelimitoptions)
    * [transitext_RetryOptions](#transitext_retryoptions)
* [Example yaml document](#example-yaml-document)

## Config

Config is a full transitext configuration document.

| Attribute | Value |
| --- | --- |
| Type | `object` |
| Properties | 2 |
| Additional properties | boolean schema=false |

### Config.Providers

Key: `providers`

Providers contains per-provider transport and API options.

| Attribute | Value |
| --- | --- |
| Required | yes |
| Reference | [`Providers`](#providers) (`#/$defs/Providers`) |

### Config.Wrappers

Key: `wrappers`

Wrappers contains pipeline wrapper options applied around provider.

| Attribute | Value |
| --- | --- |
| Required | yes |
| Reference | [`Wrappers`](#wrappers) (`#/$defs/Wrappers`) |

## Providers

Providers groups built-in provider option blocks.

| Attribute | Value |
| --- | --- |
| Type | `object` |
| Properties | 11 |
| Additional properties | boolean schema=false |

### Providers.microsoft_Options

Key: `microsoft`

Path: [`providers`](#configproviders).`microsoft`

Microsoft configures unofficial Microsoft Edge backend.

| Attribute | Value |
| --- | --- |
| Required | yes |
| Reference | [`microsoft_Options`](#microsoft_options) (`#/$defs/microsoft_Options`) |

### Providers.google_Options

Key: `google`

Path: [`providers`](#configproviders).`google`

Google configures official Google Translate API backend.

| Attribute | Value |
| --- | --- |
| Required | yes |
| Reference | [`google_Options`](#google_options) (`#/$defs/google_Options`) |

### Providers.yandex_Options

Key: `yandex`

Path: [`providers`](#configproviders).`yandex`

Yandex configures official Yandex Cloud Translate backend.

| Attribute | Value |
| --- | --- |
| Required | yes |
| Reference | [`yandex_Options`](#yandex_options) (`#/$defs/yandex_Options`) |

### Providers.azure_Options

Key: `azure`

Path: [`providers`](#configproviders).`azure`

Azure configures official Azure Translator backend.

| Attribute | Value |
| --- | --- |
| Required | yes |
| Reference | [`azure_Options`](#azure_options) (`#/$defs/azure_Options`) |

### Providers.libre_Options

Key: `libre`

Path: [`providers`](#configproviders).`libre`

Libre configures LibreTranslate backend.

| Attribute | Value |
| --- | --- |
| Required | yes |
| Reference | [`libre_Options`](#libre_options) (`#/$defs/libre_Options`) |

### Providers.deepl_Options

Key: `deepl`

Path: [`providers`](#configproviders).`deepl`

DeepL configures official DeepL API backend.

| Attribute | Value |
| --- | --- |
| Required | yes |
| Reference | [`deepl_Options`](#deepl_options) (`#/$defs/deepl_Options`) |

### Providers.bingfree_Options

Key: `bingfree`

Path: [`providers`](#configproviders).`bingfree`

BingFree configures unofficial Bing web backend.

| Attribute | Value |
| --- | --- |
| Required | yes |
| Reference | [`bingfree_Options`](#bingfree_options) (`#/$defs/bingfree_Options`) |

### Providers.yandexfree_Options

Key: `yandexfree`

Path: [`providers`](#configproviders).`yandexfree`

YandexFree configures unofficial Yandex web backend.

| Attribute | Value |
| --- | --- |
| Required | yes |
| Reference | [`yandexfree_Options`](#yandexfree_options) (`#/$defs/yandexfree_Options`) |

### Providers.deeplfree_Options

Key: `deeplfree`

Path: [`providers`](#configproviders).`deeplfree`

DeepLFree configures unofficial DeepL web backend.

| Attribute | Value |
| --- | --- |
| Required | yes |
| Reference | [`deeplfree_Options`](#deeplfree_options) (`#/$defs/deeplfree_Options`) |

### Providers.googlefree_Options

Key: `googlefree`

Path: [`providers`](#configproviders).`googlefree`

GoogleFree configures unofficial Google web backend.

| Attribute | Value |
| --- | --- |
| Required | yes |
| Reference | [`googlefree_Options`](#googlefree_options) (`#/$defs/googlefree_Options`) |

### Providers.openai_Options

Key: `openai`

Path: [`providers`](#configproviders).`openai`

OpenAI configures OpenAI-compatible chat translation backend.

| Attribute | Value |
| --- | --- |
| Required | yes |
| Reference | [`openai_Options`](#openai_options) (`#/$defs/openai_Options`) |

## Wrappers

Wrappers groups wrapper options used by transitext.Wrap(...).

| Attribute | Value |
| --- | --- |
| Type | `object` |
| Properties | 6 |
| Additional properties | boolean schema=false |

### Wrappers.transitext_ContextOptions

Key: `context`

Path: [`wrappers`](#configwrappers).`context`

Context configures context marker injection mode.

| Attribute | Value |
| --- | --- |
| Required | yes |
| Reference | [`transitext_ContextOptions`](#transitext_contextoptions) (`#/$defs/transitext_ContextOptions`) |

### Wrappers.transitext_RetryOptions

Key: `retry`

Path: [`wrappers`](#configwrappers).`retry`

Retry configures automatic retries for retryable errors.

| Attribute | Value |
| --- | --- |
| Required | yes |
| Reference | [`transitext_RetryOptions`](#transitext_retryoptions) (`#/$defs/transitext_RetryOptions`) |

### Wrappers.transitext_PartialOptions

Key: `partial`

Path: [`wrappers`](#configwrappers).`partial`

Partial configures per-item fallback behavior on batch failures.

| Attribute | Value |
| --- | --- |
| Required | yes |
| Reference | [`transitext_PartialOptions`](#transitext_partialoptions) (`#/$defs/transitext_PartialOptions`) |

### Wrappers.transitext_LongTextOptions

Key: `long_text`

Path: [`wrappers`](#configwrappers).`long_text`

LongText configures auto-splitting for oversized text items.

| Attribute | Value |
| --- | --- |
| Required | yes |
| Reference | [`transitext_LongTextOptions`](#transitext_longtextoptions) (`#/$defs/transitext_LongTextOptions`) |

### Wrappers.transitext_RateLimitOptions

Key: `rate_limit`

Path: [`wrappers`](#configwrappers).`rate_limit`

RateLimit configures minimum interval between outbound requests.

| Attribute | Value |
| --- | --- |
| Required | yes |
| Reference | [`transitext_RateLimitOptions`](#transitext_ratelimitoptions) (`#/$defs/transitext_RateLimitOptions`) |

### Wrappers.transitext_CacheOptions

Key: `cache`

Path: [`wrappers`](#configwrappers).`cache`

Cache configures in-memory response caching.

| Attribute | Value |
| --- | --- |
| Required | yes |
| Reference | [`transitext_CacheOptions`](#transitext_cacheoptions) (`#/$defs/transitext_CacheOptions`) |

## azure_Options

Options controls Azure provider behavior.

| Attribute | Value |
| --- | --- |
| Type | `object` |
| Properties | 6 |
| Additional properties | boolean schema=false |

### azure_Options.base_url

Key: `base_url`

Path: [`providers`](#configproviders).[`azure`](#providersazure_options).`base_url`

BaseURL overrides Azure translate endpoint URL.

| Attribute | Value |
| --- | --- |
| Type | `string` |
| Required | no |
| Examples | `https://api.cognitive.microsofttranslator.com/translate` |
| Format | `uri` |

### azure_Options.batch_max_chars

Key: `batch_max_chars`

Path: [`providers`](#configproviders).[`azure`](#providersazure_options).`batch_max_chars`

BatchMaxChars limits request batch size by total chars.

| Attribute | Value |
| --- | --- |
| Type | `integer` |
| Required | no |
| Default | `50000` |
| Constraints | minimum=1; maximum=50000 |

### azure_Options.batch_max_items

Key: `batch_max_items`

Path: [`providers`](#configproviders).[`azure`](#providersazure_options).`batch_max_items`

BatchMaxItems limits request batch size by item count.

| Attribute | Value |
| --- | --- |
| Type | `integer` |
| Required | no |
| Default | `1000` |
| Constraints | minimum=1; maximum=1000 |

### azure_Options.key

Key: `key`

Path: [`providers`](#configproviders).[`azure`](#providersazure_options).`key`

Key is Azure Translator subscription key.

| Attribute | Value |
| --- | --- |
| Type | `string` |
| Required | no |
| Constraints | minLength=1 |

### azure_Options.region

Key: `region`

Path: [`providers`](#configproviders).[`azure`](#providersazure_options).`region`

Region is Azure Translator resource region.

| Attribute | Value |
| --- | --- |
| Type | `string` |
| Required | no |
| Examples | `westeurope` |
| Constraints | maxLength=64 |

### azure_Options.timeout

Key: `timeout`

Path: [`providers`](#configproviders).[`azure`](#providersazure_options).`timeout`

Timeout is request timeout when HTTPClient is not provided.

| Attribute | Value |
| --- | --- |
| Type | `integer` |
| Required | no |
| Default | `20000000000` |
| Constraints | minimum=0 |

## bingfree_Options

Options controls bingfree provider behavior.

| Attribute | Value |
| --- | --- |
| Type | `object` |
| Properties | 7 |
| Additional properties | boolean schema=false |

### bingfree_Options.host_url

Key: `host_url`

Path: [`providers`](#configproviders).[`bingfree`](#providersbingfree_options).`host_url`

HostURL overrides Bing host URL.

| Attribute | Value |
| --- | --- |
| Type | `string` |
| Required | no |
| Examples | `https://www.bing.com` |
| Format | `uri` |

### bingfree_Options.max_chars

Key: `max_chars`

Path: [`providers`](#configproviders).[`bingfree`](#providersbingfree_options).`max_chars`

MaxChars limits total chars per one transitext batch.

| Attribute | Value |
| --- | --- |
| Type | `integer` |
| Required | no |
| Default | `1000` |
| Constraints | minimum=1; maximum=20000 |

### bingfree_Options.max_items

Key: `max_items`

Path: [`providers`](#configproviders).[`bingfree`](#providersbingfree_options).`max_items`

MaxItems limits items per one transitext batch.

| Attribute | Value |
| --- | --- |
| Type | `integer` |
| Required | no |
| Default | `10` |
| Constraints | minimum=1; maximum=100 |

### bingfree_Options.max_text_chars

Key: `max_text_chars`

Path: [`providers`](#configproviders).[`bingfree`](#providersbingfree_options).`max_text_chars`

MaxTextChars limits one input text length.

| Attribute | Value |
| --- | --- |
| Type | `integer` |
| Required | no |
| Default | `1000` |
| Constraints | minimum=1; maximum=3000 |

### bingfree_Options.transitext_HTTPRequestOptions

Key: `request`

Path: [`providers`](#configproviders).[`bingfree`](#providersbingfree_options).`request`

Request controls low-level HTTP header/cookie/user-agent shaping.

| Attribute | Value |
| --- | --- |
| Required | no |
| Reference | [`transitext_HTTPRequestOptions`](#transitext_httprequestoptions) (`#/$defs/transitext_HTTPRequestOptions`) |

### bingfree_Options.timeout

Key: `timeout`

Path: [`providers`](#configproviders).[`bingfree`](#providersbingfree_options).`timeout`

Timeout is request timeout when HTTPClient is not provided.

| Attribute | Value |
| --- | --- |
| Type | `integer` |
| Required | no |
| Default | `20000000000` |
| Constraints | minimum=0 |

### bingfree_Options.user_agent

Key: `user_agent`

Path: [`providers`](#configproviders).[`bingfree`](#providersbingfree_options).`user_agent`

UserAgent overrides default request user agent.

| Attribute | Value |
| --- | --- |
| Type | `string` |
| Required | no |
| Constraints | maxLength=512 |

## deepl_Options

Options controls DeepL provider behavior.

| Attribute | Value |
| --- | --- |
| Type | `object` |
| Properties | 10 |
| Additional properties | boolean schema=false |

### deepl_Options.auth_key

Key: `auth_key`

Path: [`providers`](#configproviders).[`deepl`](#providersdeepl_options).`auth_key`

Auth key for DeepL API.

| Attribute | Value |
| --- | --- |
| Type | `string` |
| Required | no |
| Constraints | minLength=1 |

### deepl_Options.batch_max_chars

Key: `batch_max_chars`

Path: [`providers`](#configproviders).[`deepl`](#providersdeepl_options).`batch_max_chars`

BatchMaxChars limits request batch size by total chars.

| Attribute | Value |
| --- | --- |
| Type | `integer` |
| Required | no |
| Default | `50000` |
| Constraints | minimum=1; maximum=131072 |

### deepl_Options.batch_max_items

Key: `batch_max_items`

Path: [`providers`](#configproviders).[`deepl`](#providersdeepl_options).`batch_max_items`

BatchMaxItems limits request batch size by item count.

| Attribute | Value |
| --- | --- |
| Type | `integer` |
| Required | no |
| Default | `50` |
| Constraints | minimum=1; maximum=50 |

### deepl_Options.formality

Key: `formality`

Path: [`providers`](#configproviders).[`deepl`](#providersdeepl_options).`formality`

Formality controls formality mode where supported.

| Attribute | Value |
| --- | --- |
| Type | `string` |
| Required | no |
| Enum | `default`, `more`, `less`, `prefer_more`, `prefer_less` |

### deepl_Options.preserve_formatting

Key: `preserve_formatting`

Path: [`providers`](#configproviders).[`deepl`](#providersdeepl_options).`preserve_formatting`

PreserveFormatting preserves source formatting when true.

| Attribute | Value |
| --- | --- |
| Type | `boolean` |
| Required | no |

### deepl_Options.source_lang

Key: `source_lang`

Path: [`providers`](#configproviders).[`deepl`](#providersdeepl_options).`source_lang`

SourceLang sets default source language when request.SourceLang is empty.

| Attribute | Value |
| --- | --- |
| Type | `string` |
| Required | no |
| Examples | `EN` |
| Constraints | minLength=2; maxLength=16 |

### deepl_Options.split_sentences

Key: `split_sentences`

Path: [`providers`](#configproviders).[`deepl`](#providersdeepl_options).`split_sentences`

SplitSentences controls DeepL sentence splitting behavior.

| Attribute | Value |
| --- | --- |
| Type | `string` |
| Required | no |
| Enum | `0`, `1`, `nonewlines` |

### deepl_Options.timeout

Key: `timeout`

Path: [`providers`](#configproviders).[`deepl`](#providersdeepl_options).`timeout`

Timeout is request timeout when HTTPClient is not provided.

| Attribute | Value |
| --- | --- |
| Type | `integer` |
| Required | no |
| Default | `20000000000` |
| Constraints | minimum=0 |

### deepl_Options.url

Key: `url`

Path: [`providers`](#configproviders).[`deepl`](#providersdeepl_options).`url`

URL overrides DeepL endpoint URL.

| Attribute | Value |
| --- | --- |
| Type | `string` |
| Required | no |
| Examples | `https://api.deepl.com/v2/translate` |
| Format | `uri` |

### deepl_Options.use_free_api

Key: `use_free_api`

Path: [`providers`](#configproviders).[`deepl`](#providersdeepl_options).`use_free_api`

UseFreeAPI selects free endpoint when URL is empty.

| Attribute | Value |
| --- | --- |
| Type | `boolean` |
| Required | no |

## deeplfree_Options

Options controls deeplfree provider behavior.

| Attribute | Value |
| --- | --- |
| Type | `object` |
| Properties | 11 |
| Additional properties | boolean schema=false |

### deeplfree_Options.accept_language

Key: `accept_language`

Path: [`providers`](#configproviders).[`deeplfree`](#providersdeeplfree_options).`accept_language`

AcceptLanguage overrides Accept-Language header value.

| Attribute | Value |
| --- | --- |
| Type | `string` |
| Required | no |
| Default | `en-US` |
| Constraints | maxLength=128 |

### deeplfree_Options.dl_session

Key: `dl_session`

Path: [`providers`](#configproviders).[`deeplfree`](#providersdeeplfree_options).`dl_session`

DLSession optionally sends dl_session cookie for authenticated web mode.

| Attribute | Value |
| --- | --- |
| Type | `string` |
| Required | no |
| Constraints | maxLength=512 |

### deeplfree_Options.max_chars

Key: `max_chars`

Path: [`providers`](#configproviders).[`deeplfree`](#providersdeeplfree_options).`max_chars`

MaxChars limits total chars per one transitext batch.

| Attribute | Value |
| --- | --- |
| Type | `integer` |
| Required | no |
| Default | `5000` |
| Constraints | minimum=1; maximum=50000 |

### deeplfree_Options.max_items

Key: `max_items`

Path: [`providers`](#configproviders).[`deeplfree`](#providersdeeplfree_options).`max_items`

MaxItems limits items per one transitext batch.

| Attribute | Value |
| --- | --- |
| Type | `integer` |
| Required | no |
| Default | `20` |
| Constraints | minimum=1; maximum=100 |

### deeplfree_Options.max_text_chars

Key: `max_text_chars`

Path: [`providers`](#configproviders).[`deeplfree`](#providersdeeplfree_options).`max_text_chars`

MaxTextChars limits one input text length.

| Attribute | Value |
| --- | --- |
| Type | `integer` |
| Required | no |
| Default | `5000` |
| Constraints | minimum=1; maximum=5000 |

### deeplfree_Options.transitext_HTTPRequestOptions

Key: `request`

Path: [`providers`](#configproviders).[`deeplfree`](#providersdeeplfree_options).`request`

Request controls low-level HTTP header/cookie/user-agent shaping.

| Attribute | Value |
| --- | --- |
| Required | no |
| Reference | [`transitext_HTTPRequestOptions`](#transitext_httprequestoptions) (`#/$defs/transitext_HTTPRequestOptions`) |

### deeplfree_Options.request_alternatives

Key: `request_alternatives`

Path: [`providers`](#configproviders).[`deeplfree`](#providersdeeplfree_options).`request_alternatives`

RequestAlternatives sets requestAlternatives for each text.

| Attribute | Value |
| --- | --- |
| Type | `integer` |
| Required | no |
| Default | `3` |
| Constraints | minimum=1; maximum=10 |

### deeplfree_Options.split_mode

Key: `split_mode`

Path: [`providers`](#configproviders).[`deeplfree`](#providersdeeplfree_options).`split_mode`

SplitMode controls DeepL splitting mode (for example "newlines").

| Attribute | Value |
| --- | --- |
| Type | `string` |
| Required | no |
| Default | `newlines` |
| Constraints | maxLength=32 |

### deeplfree_Options.timeout

Key: `timeout`

Path: [`providers`](#configproviders).[`deeplfree`](#providersdeeplfree_options).`timeout`

Timeout is request timeout when HTTPClient is not provided.

| Attribute | Value |
| --- | --- |
| Type | `integer` |
| Required | no |
| Default | `20000000000` |
| Constraints | minimum=0 |

### deeplfree_Options.url

Key: `url`

Path: [`providers`](#configproviders).[`deeplfree`](#providersdeeplfree_options).`url`

URL overrides DeepL web endpoint URL.

| Attribute | Value |
| --- | --- |
| Type | `string` |
| Required | no |
| Examples | `https://www2.deepl.com/jsonrpc` |
| Format | `uri` |

### deeplfree_Options.user_agent

Key: `user_agent`

Path: [`providers`](#configproviders).[`deeplfree`](#providersdeeplfree_options).`user_agent`

UserAgent overrides default request user agent.

| Attribute | Value |
| --- | --- |
| Type | `string` |
| Required | no |
| Constraints | maxLength=512 |

## google_Options

Options controls Google provider behavior.

| Attribute | Value |
| --- | --- |
| Type | `object` |
| Properties | 7 |
| Additional properties | boolean schema=false |

### google_Options.base_url

Key: `base_url`

Path: [`providers`](#configproviders).[`google`](#providersgoogle_options).`base_url`

BaseURL overrides official Google endpoint URL.

| Attribute | Value |
| --- | --- |
| Type | `string` |
| Required | no |
| Examples | `https://translation.googleapis.com/language/translate/v2` |
| Format | `uri` |

### google_Options.batch_max_chars

Key: `batch_max_chars`

Path: [`providers`](#configproviders).[`google`](#providersgoogle_options).`batch_max_chars`

BatchMaxChars limits request batch size by total chars.

| Attribute | Value |
| --- | --- |
| Type | `integer` |
| Required | no |
| Default | `30000` |
| Constraints | minimum=1; maximum=30000 |

### google_Options.batch_max_items

Key: `batch_max_items`

Path: [`providers`](#configproviders).[`google`](#providersgoogle_options).`batch_max_items`

BatchMaxItems limits request batch size by item count.

| Attribute | Value |
| --- | --- |
| Type | `integer` |
| Required | no |
| Constraints | minimum=1; maximum=1000 |

### google_Options.format

Key: `format`

Path: [`providers`](#configproviders).[`google`](#providersgoogle_options).`format`

Format controls source text format: "text" or "html".

| Attribute | Value |
| --- | --- |
| Type | `string` |
| Required | no |
| Default | `text` |
| Enum | `text`, `html` |

### google_Options.key

Key: `key`

Path: [`providers`](#configproviders).[`google`](#providersgoogle_options).`key`

API key for Google Translate API.

| Attribute | Value |
| --- | --- |
| Type | `string` |
| Required | no |
| Constraints | minLength=1 |

### google_Options.model

Key: `model`

Path: [`providers`](#configproviders).[`google`](#providersgoogle_options).`model`

Model selects Google translation model when supported.

| Attribute | Value |
| --- | --- |
| Type | `string` |
| Required | no |
| Constraints | maxLength=64 |

### google_Options.timeout

Key: `timeout`

Path: [`providers`](#configproviders).[`google`](#providersgoogle_options).`timeout`

Timeout is request timeout when HTTPClient is not provided.

| Attribute | Value |
| --- | --- |
| Type | `integer` |
| Required | no |
| Default | `20000000000` |
| Constraints | minimum=0 |

## googlefree_Options

Options controls googlefree translator behavior.

| Attribute | Value |
| --- | --- |
| Type | `object` |
| Properties | 9 |
| Additional properties | boolean schema=false |

### googlefree_Options.client_value

Key: `client_value`

Path: [`providers`](#configproviders).[`googlefree`](#providersgooglefree_options).`client_value`

ClientValue overrides "client" query parameter.

| Attribute | Value |
| --- | --- |
| Type | `string` |
| Required | no |
| Default | `gtx` |
| Constraints | maxLength=64 |

### googlefree_Options.concurrency

Key: `concurrency`

Path: [`providers`](#configproviders).[`googlefree`](#providersgooglefree_options).`concurrency`

Concurrency limits parallel per-item HTTP calls.

| Attribute | Value |
| --- | --- |
| Type | `integer` |
| Required | no |
| Default | `2` |
| Constraints | minimum=1; maximum=64 |

### googlefree_Options.max_chars

Key: `max_chars`

Path: [`providers`](#configproviders).[`googlefree`](#providersgooglefree_options).`max_chars`

MaxChars limits total chars per one transitext batch.

| Attribute | Value |
| --- | --- |
| Type | `integer` |
| Required | no |
| Default | `5000` |
| Constraints | minimum=1; maximum=30000 |

### googlefree_Options.max_items

Key: `max_items`

Path: [`providers`](#configproviders).[`googlefree`](#providersgooglefree_options).`max_items`

MaxItems limits items per one transitext batch.

| Attribute | Value |
| --- | --- |
| Type | `integer` |
| Required | no |
| Default | `10` |
| Constraints | minimum=1; maximum=100 |

### googlefree_Options.max_text_chars

Key: `max_text_chars`

Path: [`providers`](#configproviders).[`googlefree`](#providersgooglefree_options).`max_text_chars`

MaxTextChars limits one input text length.

| Attribute | Value |
| --- | --- |
| Type | `integer` |
| Required | no |
| Default | `5000` |
| Constraints | minimum=1; maximum=5000 |

### googlefree_Options.transitext_HTTPRequestOptions

Key: `request`

Path: [`providers`](#configproviders).[`googlefree`](#providersgooglefree_options).`request`

Request controls low-level HTTP header/cookie/user-agent shaping.

| Attribute | Value |
| --- | --- |
| Required | no |
| Reference | [`transitext_HTTPRequestOptions`](#transitext_httprequestoptions) (`#/$defs/transitext_HTTPRequestOptions`) |

### googlefree_Options.service_hosts

Key: `service_hosts`

Path: [`providers`](#configproviders).[`googlefree`](#providersgooglefree_options).`service_hosts`

ServiceHosts contains hosts for /translate_a/single endpoint.

| Attribute | Value |
| --- | --- |
| Type | `array` |
| Required | no |
| Items type | `string` |
| Items examples | `translate.googleapis.com` |
| Constraints | minItems=1; maxItems=16 |

### googlefree_Options.timeout

Key: `timeout`

Path: [`providers`](#configproviders).[`googlefree`](#providersgooglefree_options).`timeout`

Timeout is request timeout when HTTPClient is not provided.

| Attribute | Value |
| --- | --- |
| Type | `integer` |
| Required | no |
| Default | `20000000000` |
| Constraints | minimum=0 |

### googlefree_Options.user_agent

Key: `user_agent`

Path: [`providers`](#configproviders).[`googlefree`](#providersgooglefree_options).`user_agent`

UserAgent overrides default request user agent.

| Attribute | Value |
| --- | --- |
| Type | `string` |
| Required | no |
| Constraints | maxLength=512 |

## libre_Options

Options controls LibreTranslate provider behavior.

| Attribute | Value |
| --- | --- |
| Type | `object` |
| Properties | 6 |
| Additional properties | boolean schema=false |

### libre_Options.api_key

Key: `api_key`

Path: [`providers`](#configproviders).[`libre`](#providerslibre_options).`api_key`

API key for LibreTranslate instance (optional for some deployments).

| Attribute | Value |
| --- | --- |
| Type | `string` |
| Required | no |
| Constraints | minLength=1 |

### libre_Options.base_url

Key: `base_url`

Path: [`providers`](#configproviders).[`libre`](#providerslibre_options).`base_url`

BaseURL overrides LibreTranslate endpoint URL.

| Attribute | Value |
| --- | --- |
| Type | `string` |
| Required | no |
| Examples | `https://libretranslate.com/translate` |
| Format | `uri` |

### libre_Options.batch_max_chars

Key: `batch_max_chars`

Path: [`providers`](#configproviders).[`libre`](#providerslibre_options).`batch_max_chars`

BatchMaxChars limits request batch size by total chars.

| Attribute | Value |
| --- | --- |
| Type | `integer` |
| Required | no |
| Constraints | minimum=1; maximum=1000000 |

### libre_Options.batch_max_items

Key: `batch_max_items`

Path: [`providers`](#configproviders).[`libre`](#providerslibre_options).`batch_max_items`

BatchMaxItems limits request batch size by item count.

| Attribute | Value |
| --- | --- |
| Type | `integer` |
| Required | no |
| Constraints | minimum=1; maximum=1000 |

### libre_Options.format

Key: `format`

Path: [`providers`](#configproviders).[`libre`](#providerslibre_options).`format`

Format controls source text format: "text" or "html".

| Attribute | Value |
| --- | --- |
| Type | `string` |
| Required | no |
| Default | `text` |
| Enum | `text`, `html` |

### libre_Options.timeout

Key: `timeout`

Path: [`providers`](#configproviders).[`libre`](#providerslibre_options).`timeout`

Timeout is request timeout when HTTPClient is not provided.

| Attribute | Value |
| --- | --- |
| Type | `integer` |
| Required | no |
| Default | `20000000000` |
| Constraints | minimum=0 |

## microsoft_Options

Options controls microsoft translator behavior.

| Attribute | Value |
| --- | --- |
| Type | `object` |
| Properties | 10 |
| Additional properties | boolean schema=false |

### microsoft_Options.auth_url

Key: `auth_url`

Path: [`providers`](#configproviders).[`microsoft`](#providersmicrosoft_options).`auth_url`

AuthURL overrides auth endpoint URL.

| Attribute | Value |
| --- | --- |
| Type | `string` |
| Required | no |
| Examples | `https://edge.microsoft.com/translate/auth` |
| Format | `uri` |

### microsoft_Options.authentication_headers

Key: `authentication_headers`

Path: [`providers`](#configproviders).[`microsoft`](#providersmicrosoft_options).`authentication_headers`

AuthenticationHeaders are applied directly in custom_headers mode.

| Attribute | Value |
| --- | --- |
| Type | `object` |
| Required | no |
| Additional properties type | `string` |

### microsoft_Options.max_chars

Key: `max_chars`

Path: [`providers`](#configproviders).[`microsoft`](#providersmicrosoft_options).`max_chars`

MaxChars limits total chars per one provider HTTP request.

| Attribute | Value |
| --- | --- |
| Type | `integer` |
| Required | no |
| Default | `4000` |
| Constraints | minimum=1; maximum=50000 |

### microsoft_Options.max_items

Key: `max_items`

Path: [`providers`](#configproviders).[`microsoft`](#providersmicrosoft_options).`max_items`

MaxItems limits items per one provider HTTP request.

| Attribute | Value |
| --- | --- |
| Type | `integer` |
| Required | no |
| Default | `20` |
| Constraints | minimum=1; maximum=1000 |

### microsoft_Options.mode

Key: `mode`

Path: [`providers`](#configproviders).[`microsoft`](#providersmicrosoft_options).`mode`

Mode selects auth mode: "edge_free" or "custom_headers".

| Attribute | Value |
| --- | --- |
| Type | `string` |
| Required | no |
| Default | `edge_free` |
| Enum | `edge_free`, `custom_headers` |

### microsoft_Options.transitext_HTTPRequestOptions

Key: `request`

Path: [`providers`](#configproviders).[`microsoft`](#providersmicrosoft_options).`request`

Request controls low-level HTTP header/cookie/user-agent shaping.

| Attribute | Value |
| --- | --- |
| Required | no |
| Reference | [`transitext_HTTPRequestOptions`](#transitext_httprequestoptions) (`#/$defs/transitext_HTTPRequestOptions`) |

### microsoft_Options.timeout

Key: `timeout`

Path: [`providers`](#configproviders).[`microsoft`](#providersmicrosoft_options).`timeout`

Timeout is request timeout when HTTPClient is not provided.

| Attribute | Value |
| --- | --- |
| Type | `integer` |
| Required | no |
| Default | `20000000000` |
| Constraints | minimum=0 |

### microsoft_Options.translate_options

Key: `translate_options`

Path: [`providers`](#configproviders).[`microsoft`](#providersmicrosoft_options).`translate_options`

TranslateOptions adds optional query params to translate endpoint.

| Attribute | Value |
| --- | --- |
| Type | `object` |
| Required | no |
| Additional properties type | `string` |

### microsoft_Options.translate_url

Key: `translate_url`

Path: [`providers`](#configproviders).[`microsoft`](#providersmicrosoft_options).`translate_url`

TranslateURL overrides translate endpoint URL.

| Attribute | Value |
| --- | --- |
| Type | `string` |
| Required | no |
| Examples | `https://api-edge.cognitive.microsofttranslator.com/translate` |
| Format | `uri` |

### microsoft_Options.user_agent

Key: `user_agent`

Path: [`providers`](#configproviders).[`microsoft`](#providersmicrosoft_options).`user_agent`

UserAgent overrides default request user agent.

| Attribute | Value |
| --- | --- |
| Type | `string` |
| Required | no |
| Constraints | maxLength=512 |

## openai_Options

Options controls OpenAI-compatible provider behavior.

| Attribute | Value |
| --- | --- |
| Type | `object` |
| Properties | 13 |
| Additional properties | boolean schema=false |

### openai_Options.auth_token

Key: `auth_token`

Path: [`providers`](#configproviders).[`openai`](#providersopenai_options).`auth_token`

AuthToken is bearer token for API authentication.

| Attribute | Value |
| --- | --- |
| Type | `string` |
| Required | no |
| Constraints | minLength=1 |

### openai_Options.base_url

Key: `base_url`

Path: [`providers`](#configproviders).[`openai`](#providersopenai_options).`base_url`

BaseURL is OpenAI-compatible API base URL.

| Attribute | Value |
| --- | --- |
| Type | `string` |
| Required | no |
| Examples | `https://api.openai.com/v1` |
| Format | `uri` |

### openai_Options.batch_max_chars

Key: `batch_max_chars`

Path: [`providers`](#configproviders).[`openai`](#providersopenai_options).`batch_max_chars`

BatchMaxChars limits request batch size by chars.

| Attribute | Value |
| --- | --- |
| Type | `integer` |
| Required | no |
| Constraints | minimum=1; maximum=1000000 |

### openai_Options.batch_max_items

Key: `batch_max_items`

Path: [`providers`](#configproviders).[`openai`](#providersopenai_options).`batch_max_items`

BatchMaxItems limits request batch size by items.

| Attribute | Value |
| --- | --- |
| Type | `integer` |
| Required | no |
| Constraints | minimum=1; maximum=1000 |

### openai_Options.instruction_prefix

Key: `instruction_prefix`

Path: [`providers`](#configproviders).[`openai`](#providersopenai_options).`instruction_prefix`

InstructionPrefix is prepended to request instructions.

| Attribute | Value |
| --- | --- |
| Type | `string` |
| Required | no |
| Constraints | maxLength=4000 |

### openai_Options.instruction_suffix

Key: `instruction_suffix`

Path: [`providers`](#configproviders).[`openai`](#providersopenai_options).`instruction_suffix`

InstructionSuffix is appended to request instructions.

| Attribute | Value |
| --- | --- |
| Type | `string` |
| Required | no |
| Constraints | maxLength=4000 |

### openai_Options.max_tokens

Key: `max_tokens`

Path: [`providers`](#configproviders).[`openai`](#providersopenai_options).`max_tokens`

MaxTokens sets response token cap when supported.

| Attribute | Value |
| --- | --- |
| Type | `integer` |
| Required | no |
| Constraints | minimum=1; maximum=200000 |

### openai_Options.model

Key: `model`

Path: [`providers`](#configproviders).[`openai`](#providersopenai_options).`model`

Model is model identifier.

| Attribute | Value |
| --- | --- |
| Type | `string` |
| Required | no |
| Default | `gpt-4o-mini` |
| Constraints | maxLength=128 |

### openai_Options.strict_json_array

Key: `strict_json_array`

Path: [`providers`](#configproviders).[`openai`](#providersopenai_options).`strict_json_array`

StrictJSONArray requires strict JSON array response parsing.

| Attribute | Value |
| --- | --- |
| Type | `boolean` |
| Required | no |

### openai_Options.system_prompt

Key: `system_prompt`

Path: [`providers`](#configproviders).[`openai`](#providersopenai_options).`system_prompt`

SystemPrompt overrides default system prompt.

| Attribute | Value |
| --- | --- |
| Type | `string` |
| Required | no |
| Constraints | maxLength=8000 |

### openai_Options.temperature

Key: `temperature`

Path: [`providers`](#configproviders).[`openai`](#providersopenai_options).`temperature`

Temperature sets sampling temperature.

| Attribute | Value |
| --- | --- |
| Type | `number` |
| Required | no |
| Constraints | minimum=0; maximum=2 |

### openai_Options.timeout

Key: `timeout`

Path: [`providers`](#configproviders).[`openai`](#providersopenai_options).`timeout`

Timeout is HTTP timeout when HTTPClient is not provided.

| Attribute | Value |
| --- | --- |
| Type | `integer` |
| Required | no |
| Default | `60000000000` |
| Constraints | minimum=0 |

### openai_Options.top_p

Key: `top_p`

Path: [`providers`](#configproviders).[`openai`](#providersopenai_options).`top_p`

TopP sets nucleus sampling parameter.

| Attribute | Value |
| --- | --- |
| Type | `number` |
| Required | no |
| Constraints | minimum=0; maximum=1 |

## transitext_CacheOptions

CacheOptions controls cache wrapper behavior.

| Attribute | Value |
| --- | --- |
| Type | `object` |
| Properties | 2 |
| Additional properties | boolean schema=false |

### transitext_CacheOptions.include_hints

Key: `include_hints`

Path: [`wrappers`](#configwrappers).[`cache`](#wrapperstransitext_cacheoptions).`include_hints`

IncludeHints adds `hints` fields into cache key. Enable this when hints can
change translation output.

| Attribute | Value |
| --- | --- |
| Type | `boolean` |
| Required | no |

### transitext_CacheOptions.include_metadata

Key: `include_metadata`

Path: [`wrappers`](#configwrappers).[`cache`](#wrapperstransitext_cacheoptions).`include_metadata`

IncludeMetadata adds request metadata into cache key. Enable this only if
metadata affects provider output.

| Attribute | Value |
| --- | --- |
| Type | `boolean` |
| Required | no |

## transitext_ContextOptions

ContextOptions controls experimental context-injection wrapper behavior.

| Attribute | Value |
| --- | --- |
| Type | `object` |
| Properties | 5 |
| Additional properties | boolean schema=false |

### transitext_ContextOptions.context

Key: `context`

Path: [`wrappers`](#configwrappers).[`context`](#wrapperstransitext_contextoptions).`context`

Context is default translation context injected for each item.

| Attribute | Value |
| --- | --- |
| Type | `string` |
| Required | no |
| Constraints | maxLength=2000 |

### transitext_ContextOptions.context_by_id

Key: `context_by_id`

Path: [`wrappers`](#configwrappers).[`context`](#wrapperstransitext_contextoptions).`context_by_id`

ContextByID overrides context for specific item IDs.

| Attribute | Value |
| --- | --- |
| Type | `object` |
| Required | no |
| Additional properties type | `string` |

### transitext_ContextOptions.context_token

Key: `context_token`

Path: [`wrappers`](#configwrappers).[`context`](#wrapperstransitext_contextoptions).`context_token`

ContextToken is marker inserted before context payload.

| Attribute | Value |
| --- | --- |
| Type | `string` |
| Required | no |
| Constraints | minLength=1; maxLength=64 |

### transitext_ContextOptions.end_token

Key: `end_token`

Path: [`wrappers`](#configwrappers).[`context`](#wrapperstransitext_contextoptions).`end_token`

EndToken is marker appended after source text payload.

| Attribute | Value |
| --- | --- |
| Type | `string` |
| Required | no |
| Constraints | minLength=1; maxLength=64 |

### transitext_ContextOptions.text_token

Key: `text_token`

Path: [`wrappers`](#configwrappers).[`context`](#wrapperstransitext_contextoptions).`text_token`

TextToken is marker inserted before source text payload.

| Attribute | Value |
| --- | --- |
| Type | `string` |
| Required | no |
| Constraints | minLength=1; maxLength=64 |

## transitext_HTTPRequestOptions

HTTPRequestOptions controls generic outbound request shaping for providers.

| Attribute | Value |
| --- | --- |
| Type | `object` |
| Properties | 3 |
| Additional properties | boolean schema=false |

### transitext_HTTPRequestOptions.cookies

Key: `cookies`

Paths:

* [`providers`](#configproviders).[`bingfree`](#providersbingfree_options).[`request`](#bingfree_optionstransitext_httprequestoptions).`cookies`
* [`providers`](#configproviders).[`deeplfree`](#providersdeeplfree_options).[`request`](#deeplfree_optionstransitext_httprequestoptions).`cookies`
* [`providers`](#configproviders).[`googlefree`](#providersgooglefree_options).[`request`](#googlefree_optionstransitext_httprequestoptions).`cookies`
* [`providers`](#configproviders).[`microsoft`](#providersmicrosoft_options).[`request`](#microsoft_optionstransitext_httprequestoptions).`cookies`
* [`providers`](#configproviders).[`yandexfree`](#providersyandexfree_options).[`request`](#yandexfree_optionstransitext_httprequestoptions).`cookies`

Cookies sets default cookies applied when request cookie header is empty.

| Attribute | Value |
| --- | --- |
| Type | `object` |
| Required | no |
| Additional properties type | `string` |

### transitext_HTTPRequestOptions.headers

Key: `headers`

Paths:

* [`providers`](#configproviders).[`bingfree`](#providersbingfree_options).[`request`](#bingfree_optionstransitext_httprequestoptions).`headers`
* [`providers`](#configproviders).[`deeplfree`](#providersdeeplfree_options).[`request`](#deeplfree_optionstransitext_httprequestoptions).`headers`
* [`providers`](#configproviders).[`googlefree`](#providersgooglefree_options).[`request`](#googlefree_optionstransitext_httprequestoptions).`headers`
* [`providers`](#configproviders).[`microsoft`](#providersmicrosoft_options).[`request`](#microsoft_optionstransitext_httprequestoptions).`headers`
* [`providers`](#configproviders).[`yandexfree`](#providersyandexfree_options).[`request`](#yandexfree_optionstransitext_httprequestoptions).`headers`

Headers sets default headers applied when missing in request.

| Attribute | Value |
| --- | --- |
| Type | `object` |
| Required | no |
| Additional properties type | `string` |

### transitext_HTTPRequestOptions.user_agent

Key: `user_agent`

Paths:

* [`providers`](#configproviders).[`bingfree`](#providersbingfree_options).[`request`](#bingfree_optionstransitext_httprequestoptions).`user_agent`
* [`providers`](#configproviders).[`deeplfree`](#providersdeeplfree_options).[`request`](#deeplfree_optionstransitext_httprequestoptions).`user_agent`
* [`providers`](#configproviders).[`googlefree`](#providersgooglefree_options).[`request`](#googlefree_optionstransitext_httprequestoptions).`user_agent`
* [`providers`](#configproviders).[`microsoft`](#providersmicrosoft_options).[`request`](#microsoft_optionstransitext_httprequestoptions).`user_agent`
* [`providers`](#configproviders).[`yandexfree`](#providersyandexfree_options).[`request`](#yandexfree_optionstransitext_httprequestoptions).`user_agent`

UserAgent sets default User-Agent when request does not set it explicitly.

| Attribute | Value |
| --- | --- |
| Type | `string` |
| Required | no |
| Constraints | maxLength=512 |

## transitext_LongTextOptions

LongTextOptions controls single-item overflow handling.

| Attribute | Value |
| --- | --- |
| Type | `object` |
| Properties | 2 |
| Additional properties | boolean schema=false |

### transitext_LongTextOptions.error_on_overflow

Key: `error_on_overflow`

Path: [`wrappers`](#configwrappers).[`long_text`](#wrapperstransitext_longtextoptions).`error_on_overflow`

ErrorOnOverflow fails when item exceeds MaxTextChars. By default long items are
split and merged automatically.

| Attribute | Value |
| --- | --- |
| Type | `boolean` |
| Required | no |

### transitext_LongTextOptions.max_text_chars

Key: `max_text_chars`

Path: [`wrappers`](#configwrappers).[`long_text`](#wrapperstransitext_longtextoptions).`max_text_chars`

MaxTextChars overrides provider single-item limit.

| Attribute | Value |
| --- | --- |
| Type | `integer` |
| Required | no |
| Constraints | minimum=1; maximum=1000000 |

## transitext_PartialOptions

PartialOptions controls partial-result wrapper behavior.

| Attribute | Value |
| --- | --- |
| Type | `object` |
| Properties | 3 |
| Additional properties | boolean schema=false |

### transitext_PartialOptions.continue_on_context

Key: `continue_on_context`

Path: [`wrappers`](#configwrappers).[`partial`](#wrapperstransitext_partialoptions).`continue_on_context`

ContinueOnContext keeps processing after context cancellation/deadline.

| Attribute | Value |
| --- | --- |
| Type | `boolean` |
| Required | no |

### transitext_PartialOptions.continue_on_temporary

Key: `continue_on_temporary`

Path: [`wrappers`](#configwrappers).[`partial`](#wrapperstransitext_partialoptions).`continue_on_temporary`

ContinueOnTemporary keeps processing after temporary provider errors.

| Attribute | Value |
| --- | --- |
| Type | `boolean` |
| Required | no |

### transitext_PartialOptions.item_retries

Key: `item_retries`

Path: [`wrappers`](#configwrappers).[`partial`](#wrapperstransitext_partialoptions).`item_retries`

ItemRetries is extra retry count for one failed item.

| Attribute | Value |
| --- | --- |
| Type | `integer` |
| Required | no |
| Constraints | minimum=0; maximum=16 |

## transitext_RateLimitOptions

RateLimitOptions controls rate limit wrapper behavior.

| Attribute | Value |
| --- | --- |
| Type | `object` |
| Properties | 1 |
| Additional properties | boolean schema=false |

### transitext_RateLimitOptions.min_interval

Key: `min_interval`

Path: [`wrappers`](#configwrappers).[`rate_limit`](#wrapperstransitext_ratelimitoptions).`min_interval`

MinInterval is minimal delay between outbound requests. Set `0` to disable
throttling.

| Attribute | Value |
| --- | --- |
| Type | `integer` |
| Required | no |
| Constraints | minimum=0 |

## transitext_RetryOptions

RetryOptions controls retry wrapper behavior.

| Attribute | Value |
| --- | --- |
| Type | `object` |
| Properties | 5 |
| Additional properties | boolean schema=false |

### transitext_RetryOptions.attempts

Key: `attempts`

Path: [`wrappers`](#configwrappers).[`retry`](#wrapperstransitext_retryoptions).`attempts`

Attempts is total number of attempts for one request, including first try.

| Attribute | Value |
| --- | --- |
| Type | `integer` |
| Required | no |
| Default | `3` |
| Constraints | minimum=1; maximum=32 |

### transitext_RetryOptions.backoff

Key: `backoff`

Path: [`wrappers`](#configwrappers).[`retry`](#wrapperstransitext_retryoptions).`backoff`

Backoff multiplies delay after each failed attempt.

| Attribute | Value |
| --- | --- |
| Type | `number` |
| Required | no |
| Default | `2` |
| Constraints | minimum=1; maximum=10 |

### transitext_RetryOptions.delay

Key: `delay`

Path: [`wrappers`](#configwrappers).[`retry`](#wrapperstransitext_retryoptions).`delay`

Delay defines delay before first retry.

| Attribute | Value |
| --- | --- |
| Type | `integer` |
| Required | no |
| Default | `300000000` |
| Constraints | minimum=0 |

### transitext_RetryOptions.max_delay

Key: `max_delay`

Path: [`wrappers`](#configwrappers).[`retry`](#wrapperstransitext_retryoptions).`max_delay`

MaxDelay defines maximum delay cap between retries.

| Attribute | Value |
| --- | --- |
| Type | `integer` |
| Required | no |
| Constraints | minimum=0 |

### transitext_RetryOptions.retry_permanent

Key: `retry_permanent`

Path: [`wrappers`](#configwrappers).[`retry`](#wrapperstransitext_retryoptions).`retry_permanent`

RetryPermanent allows retry for permanent provider errors.

| Attribute | Value |
| --- | --- |
| Type | `boolean` |
| Required | no |

## yandex_Options

Options controls Yandex provider behavior.

| Attribute | Value |
| --- | --- |
| Type | `object` |
| Properties | 7 |
| Additional properties | boolean schema=false |

### yandex_Options.api_key

Key: `api_key`

Path: [`providers`](#configproviders).[`yandex`](#providersyandex_options).`api_key`

API key for Yandex Cloud Translate.

| Attribute | Value |
| --- | --- |
| Type | `string` |
| Required | no |
| Constraints | minLength=1 |

### yandex_Options.base_url

Key: `base_url`

Path: [`providers`](#configproviders).[`yandex`](#providersyandex_options).`base_url`

BaseURL overrides Yandex translate endpoint URL.

| Attribute | Value |
| --- | --- |
| Type | `string` |
| Required | no |
| Examples | `https://translate.api.cloud.yandex.net/translate/v2/translate` |
| Format | `uri` |

### yandex_Options.batch_max_chars

Key: `batch_max_chars`

Path: [`providers`](#configproviders).[`yandex`](#providersyandex_options).`batch_max_chars`

BatchMaxChars limits request batch size by total chars.

| Attribute | Value |
| --- | --- |
| Type | `integer` |
| Required | no |
| Default | `10000` |
| Constraints | minimum=1; maximum=10000 |

### yandex_Options.batch_max_items

Key: `batch_max_items`

Path: [`providers`](#configproviders).[`yandex`](#providersyandex_options).`batch_max_items`

BatchMaxItems limits request batch size by item count.

| Attribute | Value |
| --- | --- |
| Type | `integer` |
| Required | no |
| Constraints | minimum=1; maximum=1000 |

### yandex_Options.folder_id

Key: `folder_id`

Path: [`providers`](#configproviders).[`yandex`](#providersyandex_options).`folder_id`

FolderID is Yandex Cloud folder id for API key auth mode.

| Attribute | Value |
| --- | --- |
| Type | `string` |
| Required | no |
| Constraints | maxLength=128 |

### yandex_Options.iam_token

Key: `iam_token`

Path: [`providers`](#configproviders).[`yandex`](#providersyandex_options).`iam_token`

IAMToken can be used instead of APIKey.

| Attribute | Value |
| --- | --- |
| Type | `string` |
| Required | no |
| Constraints | minLength=1 |

### yandex_Options.timeout

Key: `timeout`

Path: [`providers`](#configproviders).[`yandex`](#providersyandex_options).`timeout`

Timeout is request timeout when HTTPClient is not provided.

| Attribute | Value |
| --- | --- |
| Type | `integer` |
| Required | no |
| Default | `20000000000` |
| Constraints | minimum=0 |

## yandexfree_Options

Options controls yandexfree provider behavior.

| Attribute | Value |
| --- | --- |
| Type | `object` |
| Properties | 7 |
| Additional properties | boolean schema=false |

### yandexfree_Options.base_url

Key: `base_url`

Path: [`providers`](#configproviders).[`yandexfree`](#providersyandexfree_options).`base_url`

BaseURL overrides Yandex free API base URL.

| Attribute | Value |
| --- | --- |
| Type | `string` |
| Required | no |
| Examples | `https://translate.yandex.net/api/v1/tr.json` |
| Format | `uri` |

### yandexfree_Options.max_chars

Key: `max_chars`

Path: [`providers`](#configproviders).[`yandexfree`](#providersyandexfree_options).`max_chars`

MaxChars limits total chars per one transitext batch.

| Attribute | Value |
| --- | --- |
| Type | `integer` |
| Required | no |
| Default | `10000` |
| Constraints | minimum=1; maximum=10000 |

### yandexfree_Options.max_items

Key: `max_items`

Path: [`providers`](#configproviders).[`yandexfree`](#providersyandexfree_options).`max_items`

MaxItems limits items per one transitext batch.

| Attribute | Value |
| --- | --- |
| Type | `integer` |
| Required | no |
| Default | `10` |
| Constraints | minimum=1; maximum=100 |

### yandexfree_Options.max_text_chars

Key: `max_text_chars`

Path: [`providers`](#configproviders).[`yandexfree`](#providersyandexfree_options).`max_text_chars`

MaxTextChars limits one input text length.

| Attribute | Value |
| --- | --- |
| Type | `integer` |
| Required | no |
| Default | `10000` |
| Constraints | minimum=1; maximum=10000 |

### yandexfree_Options.transitext_HTTPRequestOptions

Key: `request`

Path: [`providers`](#configproviders).[`yandexfree`](#providersyandexfree_options).`request`

Request controls low-level HTTP header/cookie/user-agent shaping.

| Attribute | Value |
| --- | --- |
| Required | no |
| Reference | [`transitext_HTTPRequestOptions`](#transitext_httprequestoptions) (`#/$defs/transitext_HTTPRequestOptions`) |

### yandexfree_Options.timeout

Key: `timeout`

Path: [`providers`](#configproviders).[`yandexfree`](#providersyandexfree_options).`timeout`

Timeout is request timeout when HTTPClient is not provided.

| Attribute | Value |
| --- | --- |
| Type | `integer` |
| Required | no |
| Default | `20000000000` |
| Constraints | minimum=0 |

### yandexfree_Options.user_agent

Key: `user_agent`

Path: [`providers`](#configproviders).[`yandexfree`](#providersyandexfree_options).`user_agent`

UserAgent overrides default user-agent.

| Attribute | Value |
| --- | --- |
| Type | `string` |
| Required | no |
| Constraints | maxLength=512 |

## Example yaml document

```yaml
# yaml-language-server: $schema=https://raw.githubusercontent.com/woozymasta/transitext/HEAD/schema.json

# Providers contains per-provider transport and API options.
providers:
  # Azure configures official Azure Translator backend.
  azure:
    # BaseURL overrides Azure translate endpoint URL.
    base_url: https://api.cognitive.microsofttranslator.com/translate
    # BatchMaxChars limits request batch size by total chars.
    batch_max_chars: 50000
    # BatchMaxItems limits request batch size by item count.
    batch_max_items: 1000
    # Key is Azure Translator subscription key.
    key: <string>
    # Region is Azure Translator resource region.
    region: westeurope
    # Timeout is request timeout when HTTPClient is not provided.
    timeout: 20000000000
  # BingFree configures unofficial Bing web backend.
  bingfree:
    # HostURL overrides Bing host URL.
    host_url: https://www.bing.com
    # MaxChars limits total chars per one transitext batch.
    max_chars: 1000
    # MaxItems limits items per one transitext batch.
    max_items: 10
    # MaxTextChars limits one input text length.
    max_text_chars: 1000
    # Request controls low-level HTTP header/cookie/user-agent shaping.
    request:
      # Cookies sets default cookies applied when request cookie header is empty.
      cookies: {}
      # Headers sets default headers applied when missing in request.
      headers: {}
      # UserAgent sets default User-Agent when request does not set it explicitly.
      user_agent: <string>
    # Timeout is request timeout when HTTPClient is not provided.
    timeout: 20000000000
    # UserAgent overrides default request user agent.
    user_agent: <string>
  # DeepL configures official DeepL API backend.
  deepl:
    # Auth key for DeepL API.
    auth_key: <string>
    # BatchMaxChars limits request batch size by total chars.
    batch_max_chars: 50000
    # BatchMaxItems limits request batch size by item count.
    batch_max_items: 50
    # Formality controls formality mode where supported.
    formality: default
    # PreserveFormatting preserves source formatting when true.
    preserve_formatting: false
    # SourceLang sets default source language when request.SourceLang is empty.
    source_lang: EN
    # SplitSentences controls DeepL sentence splitting behavior.
    split_sentences: "0"
    # Timeout is request timeout when HTTPClient is not provided.
    timeout: 20000000000
    # URL overrides DeepL endpoint URL.
    url: https://api.deepl.com/v2/translate
    # UseFreeAPI selects free endpoint when URL is empty.
    use_free_api: false
  # DeepLFree configures unofficial DeepL web backend.
  deeplfree:
    # AcceptLanguage overrides Accept-Language header value.
    accept_language: en-US
    # DLSession optionally sends dl_session cookie for authenticated web mode.
    dl_session: <string>
    # MaxChars limits total chars per one transitext batch.
    max_chars: 5000
    # MaxItems limits items per one transitext batch.
    max_items: 20
    # MaxTextChars limits one input text length.
    max_text_chars: 5000
    # Request controls low-level HTTP header/cookie/user-agent shaping.
    request:
      # Cookies sets default cookies applied when request cookie header is empty.
      cookies: {}
      # Headers sets default headers applied when missing in request.
      headers: {}
      # UserAgent sets default User-Agent when request does not set it explicitly.
      user_agent: <string>
    # RequestAlternatives sets requestAlternatives for each text.
    request_alternatives: 3
    # SplitMode controls DeepL splitting mode (for example "newlines").
    split_mode: newlines
    # Timeout is request timeout when HTTPClient is not provided.
    timeout: 20000000000
    # URL overrides DeepL web endpoint URL.
    url: https://www2.deepl.com/jsonrpc
    # UserAgent overrides default request user agent.
    user_agent: <string>
  # Google configures official Google Translate API backend.
  google:
    # BaseURL overrides official Google endpoint URL.
    base_url: https://translation.googleapis.com/language/translate/v2
    # BatchMaxChars limits request batch size by total chars.
    batch_max_chars: 30000
    # BatchMaxItems limits request batch size by item count.
    batch_max_items: 0
    # Format controls source text format: "text" or "html".
    format: text
    # API key for Google Translate API.
    key: <string>
    # Model selects Google translation model when supported.
    model: <string>
    # Timeout is request timeout when HTTPClient is not provided.
    timeout: 20000000000
  # GoogleFree configures unofficial Google web backend.
  googlefree:
    # ClientValue overrides "client" query parameter.
    client_value: gtx
    # Concurrency limits parallel per-item HTTP calls.
    concurrency: 2
    # MaxChars limits total chars per one transitext batch.
    max_chars: 5000
    # MaxItems limits items per one transitext batch.
    max_items: 10
    # MaxTextChars limits one input text length.
    max_text_chars: 5000
    # Request controls low-level HTTP header/cookie/user-agent shaping.
    request:
      # Cookies sets default cookies applied when request cookie header is empty.
      cookies: {}
      # Headers sets default headers applied when missing in request.
      headers: {}
      # UserAgent sets default User-Agent when request does not set it explicitly.
      user_agent: <string>
    # ServiceHosts contains hosts for /translate_a/single endpoint.
    service_hosts:
      - translate.googleapis.com
    # Timeout is request timeout when HTTPClient is not provided.
    timeout: 20000000000
    # UserAgent overrides default request user agent.
    user_agent: <string>
  # Libre configures LibreTranslate backend.
  libre:
    # API key for LibreTranslate instance (optional for some deployments).
    api_key: <string>
    # BaseURL overrides LibreTranslate endpoint URL.
    base_url: https://libretranslate.com/translate
    # BatchMaxChars limits request batch size by total chars.
    batch_max_chars: 0
    # BatchMaxItems limits request batch size by item count.
    batch_max_items: 0
    # Format controls source text format: "text" or "html".
    format: text
    # Timeout is request timeout when HTTPClient is not provided.
    timeout: 20000000000
  # Microsoft configures unofficial Microsoft Edge backend.
  microsoft:
    # AuthURL overrides auth endpoint URL.
    auth_url: https://edge.microsoft.com/translate/auth
    # AuthenticationHeaders are applied directly in custom_headers mode.
    authentication_headers: {}
    # MaxChars limits total chars per one provider HTTP request.
    max_chars: 4000
    # MaxItems limits items per one provider HTTP request.
    max_items: 20
    # Mode selects auth mode: "edge_free" or "custom_headers".
    mode: edge_free
    # Request controls low-level HTTP header/cookie/user-agent shaping.
    request:
      # Cookies sets default cookies applied when request cookie header is empty.
      cookies: {}
      # Headers sets default headers applied when missing in request.
      headers: {}
      # UserAgent sets default User-Agent when request does not set it explicitly.
      user_agent: <string>
    # Timeout is request timeout when HTTPClient is not provided.
    timeout: 20000000000
    # TranslateOptions adds optional query params to translate endpoint.
    translate_options: {}
    # TranslateURL overrides translate endpoint URL.
    translate_url: https://api-edge.cognitive.microsofttranslator.com/translate
    # UserAgent overrides default request user agent.
    user_agent: <string>
  # OpenAI configures OpenAI-compatible chat translation backend.
  openai:
    # AuthToken is bearer token for API authentication.
    auth_token: <string>
    # BaseURL is OpenAI-compatible API base URL.
    base_url: https://api.openai.com/v1
    # BatchMaxChars limits request batch size by chars.
    batch_max_chars: 0
    # BatchMaxItems limits request batch size by items.
    batch_max_items: 0
    # InstructionPrefix is prepended to request instructions.
    instruction_prefix: <string>
    # InstructionSuffix is appended to request instructions.
    instruction_suffix: <string>
    # MaxTokens sets response token cap when supported.
    max_tokens: 0
    # Model is model identifier.
    model: gpt-4o-mini
    # StrictJSONArray requires strict JSON array response parsing.
    strict_json_array: false
    # SystemPrompt overrides default system prompt.
    system_prompt: <string>
    # Temperature sets sampling temperature.
    temperature: 0
    # Timeout is HTTP timeout when HTTPClient is not provided.
    timeout: 60000000000
    # TopP sets nucleus sampling parameter.
    top_p: 0
  # Yandex configures official Yandex Cloud Translate backend.
  yandex:
    # API key for Yandex Cloud Translate.
    api_key: <string>
    # BaseURL overrides Yandex translate endpoint URL.
    base_url: https://translate.api.cloud.yandex.net/translate/v2/translate
    # BatchMaxChars limits request batch size by total chars.
    batch_max_chars: 10000
    # BatchMaxItems limits request batch size by item count.
    batch_max_items: 0
    # FolderID is Yandex Cloud folder id for API key auth mode.
    folder_id: <string>
    # IAMToken can be used instead of APIKey.
    iam_token: <string>
    # Timeout is request timeout when HTTPClient is not provided.
    timeout: 20000000000
  # YandexFree configures unofficial Yandex web backend.
  yandexfree:
    # BaseURL overrides Yandex free API base URL.
    base_url: https://translate.yandex.net/api/v1/tr.json
    # MaxChars limits total chars per one transitext batch.
    max_chars: 10000
    # MaxItems limits items per one transitext batch.
    max_items: 10
    # MaxTextChars limits one input text length.
    max_text_chars: 10000
    # Request controls low-level HTTP header/cookie/user-agent shaping.
    request:
      # Cookies sets default cookies applied when request cookie header is empty.
      cookies: {}
      # Headers sets default headers applied when missing in request.
      headers: {}
      # UserAgent sets default User-Agent when request does not set it explicitly.
      user_agent: <string>
    # Timeout is request timeout when HTTPClient is not provided.
    timeout: 20000000000
    # UserAgent overrides default user-agent.
    user_agent: <string>
# Wrappers contains pipeline wrapper options applied around provider.
wrappers:
  # Cache configures in-memory response caching.
  cache:
    # IncludeHints adds `hints` fields into cache key.
    # Enable this when hints can change translation output.
    include_hints: false
    # IncludeMetadata adds request metadata into cache key.
    # Enable this only if metadata affects provider output.
    include_metadata: false
  # Context configures context marker injection mode.
  context:
    # Context is default translation context injected for each item.
    context: <string>
    # ContextByID overrides context for specific item IDs.
    context_by_id: {}
    # ContextToken is marker inserted before context payload.
    context_token: <string>
    # EndToken is marker appended after source text payload.
    end_token: <string>
    # TextToken is marker inserted before source text payload.
    text_token: <string>
  # LongText configures auto-splitting for oversized text items.
  long_text:
    # ErrorOnOverflow fails when item exceeds MaxTextChars.
    # By default long items are split and merged automatically.
    error_on_overflow: false
    # MaxTextChars overrides provider single-item limit.
    max_text_chars: 0
  # Partial configures per-item fallback behavior on batch failures.
  partial:
    # ContinueOnContext keeps processing after context cancellation/deadline.
    continue_on_context: false
    # ContinueOnTemporary keeps processing after temporary provider errors.
    continue_on_temporary: false
    # ItemRetries is extra retry count for one failed item.
    item_retries: 0
  # RateLimit configures minimum interval between outbound requests.
  rate_limit:
    # MinInterval is minimal delay between outbound requests.
    # Set `0` to disable throttling.
    min_interval: 0
  # Retry configures automatic retries for retryable errors.
  retry:
    # Attempts is total number of attempts for one request, including first try.
    attempts: 3
    # Backoff multiplies delay after each failed attempt.
    backoff: 2
    # Delay defines delay before first retry.
    delay: 300000000
    # MaxDelay defines maximum delay cap between retries.
    max_delay: 0
    # RetryPermanent allows retry for permanent provider errors.
    retry_permanent: false
```
<!-- Automatically generated file, do not modify! -->
