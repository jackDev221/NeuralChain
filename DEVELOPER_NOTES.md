Developer Notes
===============

This document notes is refered from [Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments) and [Effective Go](https://golang.org/doc/effective_go.html)


<!-- markdown-toc start -->
**Table of Contents**

1. [Development Notes](#development-notes)
2. [Development environment](#development-environment)
3. [Repository structure](#repository-structure)
4. [Formatting and style](#formatting-and-style)
    - [Imports](#1-imports)
    - [Package Names](#2-package-names)
    - [Variable Names](#3-variable-names)
    - [Declaring Empty Slices](#4-declaring-empty-slices)
    - [Pass Values](#5-pass-values)
    - [Receiver Names](#6-receiver-names)
    - [Receiver Type](#7-receiver-type)
    - [Synchronous Functions](#8-synchronous-functions)
    - [Crypto Rand](#9-crypto-rand)
    - [Comments](#10-comments)
    - [Error Strings](#11-error-strings)
    - [Useful Test Failures](#12-useful-test-failures)
    - [Examples](#13-examples)
    - [Handle Errors](#14-handle-errors)
    - [In-Band Errors](#15-in-band-errors)
    - [Indent Error Flow](#16-indent-error-flow)
    - [Interfaces](#17-interfaces)
    - [Line Length](#18-line_length)
    - [Mixed Caps](#19-mixed-caps)
    - [Named Result Parameters](#20-named-result-parameters)
    - [Mocking](#21-mocking)

5. [Program design](#program-design)
    - [Use struct literal initialization](#1-use-struct-literal-initialization)
    - [Avoid nil checks via default no-op implementations](#2-avoid-nil-checks-via-default-no-op-implementations)
    - [Loggers are dependencies](#3-loggers-are-dependencies)

6. [Logging and instrumentation](#logging-and-instrumentation)
7. [Testings](#testings)
8. [Build and deploy](#build-and-deploy)
10. [References](#references)

<!-- markdown-toc end -->

Development environment
----------------------

```
Tip: Put $GOPATH/bin in your $PATH, so installed binaries are easily accessible.
```

Go has development environment conventions centered around the *GOPATH*. Put GOPATH/bin into your PATH to allow easily run binaries you get via *go get*, and makes the (preferred) go install mechanism of building code easier to work with.

For editors and IDEs, you should install below plugins to make your life easier.
- Vim: [vim-go](https://github.com/fatih/vim-go).
- Sublime Text: [GoSublime](https://github.com/DisposaBoy/GoSublime).
- Atom: [go-plus](https://atom.io/packages/go-plus).
- VS Code: [vscode-go](https://github.com/Microsoft/vscode-go).


Repository structure
------------------

There is no single best repo structure, however, we do have a good general model that fit for many projects, especially for projects with both binaries and libraries.

The basic idea is to have 2 top-level directories, pkg and cmd. Create directories for each of your libraries under pkg, and similar for binaries under cmd.

```
Tip: Put library code under a pkg/ subdirectory. Put binaries under a cmd/ subdirectory.
```

**Example**

```
github.com/example/foo/
circle.yml
Dockerfile
cmd/
foosrv/
main.go
foocli/
main.go
pkg/
fs/
fs.go
fs_test.go
mock.go
mock_test.go
merge/
merge.go
merge_test.go
api/
api.go
api_test.go
```

```
Tip: Always use fully-qualified import paths. Never use relative imports.
```

You should use `import "github.com/example/foo/pkg/fs"`. Relative imports can be messy, especially for shared projects where directory structures are likely to change.


Formatting and style
------------------

Please refer to [Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments).

I just highlighted some notes here:


**1. Imports**
--------

```
Tip: Use goimports as formatting tool. Standard library packages first, the other packages later.
```

Just use [`goimports`](https://godoc.org/golang.org/x/tools/cmd/goimports) as formatting tool and our standard (`gofmt` is deprecated soon)

Imports are organized in groups, with blank lines between them. The standard library packages are always in the first group.

Example:
```
package main

import (
"fmt"
"hash/adler32"
"os"

"appengine/foo"
"appengine/user"

"github.com/foo/bar"
"rsc.io/goversion/version"
)
```

**2. Package Names**
--------


```
Tip: Just use `golint` to check naming conventions.
```

Install [golint](https://github.com/golang/lint) using `go get -u github.com/golang/lint/golint`

The package name is the base name of its source directory; the package in `src/encoding/base64` is imported as `"encoding/base64"` but has name `base64`, not `encoding_base64` and not `encodingBase64`.

Avoid meaningless package names like `util`, `common`, `misc`, `api`, `types`, and `interfaces`

**3. Variable Names**
--------

```
Tip: Single letter for common variables, descriptive name for global variables.
```

Variable names in Go should be short rather than long. This is especially true for local variables with limited scope. Prefer `c` to `lineCount`. Prefer `i` to `sliceIndex`.

The basic rule: the further from its declaration that a name is used, the more descriptive the name must be. For a method receiver, one or two letters is sufficient. Common variables such as loop indices and readers can be a single letter (`i`, `r`). More unusual things and global variables need more descriptive names.

**4. Declaring Empty Slices**
--------

When declaring empty slices, prefer:
```
var t []string
```
over
```
t := []string{}
```

The former one is a nil slice value, while the latter is non-nil but zero-length.

**5. Pass Values**
--------

Consider to pass params by value when:
- Fixed size
- Memory insignificant (string/ numerical/ interfaces)
- Doesn't need to refer to its allocation, only its values

Otherwise use pointer.

**6. Receiver Names**
--------

```
Tip: Name should be short but need to be consistent.
```

The name of a method's receiver should be a reflection of its identity; often a one or two letter abbreviation of its type suffices (such as "c" or "cl" for "Client"). Don't use generic names such as "me", "this" or "self", identifiers typical of object-oriented languages that gives the method a special meaning. In Go, the receiver of a method is just another parameter and therefore, should be named accordingly. The name need not be as descriptive as that of a method argument, as its role is obvious and serves no documentary purpose. It can be very short as it will appear on almost every line of every method of the type; familiarity admits brevity. Be consistent, too: if you call the receiver "c" in one method, don't call it "cl" in another.

**7. Receiver Type**
--------

```
Tip: When in doubt, use a pointer receiver.
```

Pointer receiver:

- Shared objects (might be modified/ accessed in different go routines)
- Mutable objects or object with mutable fields.
- Large struct/ large slice.
- When in doubt, use pointer


Normal receiver:

- Chan, Map, Func.
- Anything that is small & non-mutable.

**8. Synchronous Functions**
--------

```
Tip: Prefer synchronous functin over asynchronous ones.
```

Prefer synchronous functions - functions which return their results directly or finish any callbacks or channel ops before returning - over asynchronous ones.

Synchronous functions keep goroutines localized within a call, making it easier to reason about their lifetimes and avoid leaks and data races. They're also easier to test: the caller can pass an input and check the output without the need for polling or synchronization.

If callers need more concurrency, they can add it easily by calling the function from a separate goroutine. But it is quite difficult - sometimes impossible - to remove unnecessary concurrency at the caller side.


**9. Crypto Rand**
--------

Use `crypto/rand` instead of `math/rand` to generate keys as the second one is unseeded, the generator is completely predictable.

**10. Comments**
--------

```
Tip: Comments should be in full sentences.
```

Go provides C-style /* */ block comments and C++-style // line comments.

Every package should have a package comment and must appear adjacent to the package clause, with no blank line. For multi-file packages, the package comment only needs to be present in one file, and any one will do. The package comment should introduce the package and provide information relevant to the package as a whole. It will appear first on the godoc page and should set up the detailed documentation that follows.

If the package is simple, the package comment can be brief using line comments.

```
/*
Package template implements data-driven templates for generating textual
output such as HTML.
....
*/
package template

// Package math provides basic constants and mathematical functions.
package math
```

For "package main" comments, other styles of comment are fine after the binary name (and it may be capitalized if it comes first), For example, for a `package main` in the directory `seedgen` you could write:

```
// Binary seedgen ...
package main
```
or
```
// Command seedgen ...
package main
```

All top-level, exported names should have doc comments, as should non-trivial unexported type or function declarations. Doc comments work best as complete sentences, which allow a wide variety of automated presentations. The first sentence should be a one-sentence summary that starts with the name being declared.

```
// Compile parses a regular expression and returns, if successful,
// a Regexp that can be used to match against text.
func Compile(str string) (*Regexp, error) { }
```

See [here](https://golang.org/doc/effective_go.html#commentary) for more information about commentary conventions.


**11. Error Strings**
--------

Error strings should not be capitalized (unless beginning with proper nouns or acronyms) or end with punctuation, since they are usually printed following other context. That is, use ```fmt.Errorf("something bad")``` not ```fmt.Errorf("Something bad")```, so that ```log.Printf("Reading %s: %v", filename, err)``` formats without a spurious capital letter mid-message. This does not apply to logging, which is implicitly line-oriented and not combined inside other messages.

**12. Useful Test Failures**

Tests should fail with helpful messages saying what was wrong, with what inputs, what was actually got, and what was expected. It may be tempting to write a bunch of assertFoo helpers, but be sure your helpers produce useful error messages. Assume that the person debugging your failing test is not you, and is not your team. A typical Go test fails like:

```
if got != tt.want {
t.Errorf("Foo(%q) = %d; want %d", tt.in, got, tt.want) // or Fatalf, if test can't test anything more past this point
}
```

Another common technique to disambiguate failing tests when using a test helper with different input is to wrap each caller with a different TestFoo function, so the test fails with that name:

```
func TestSingleValue(t *testing.T) { testHelper(t, []int{80}) }
func TestNoValues(t *testing.T)    { testHelper(t, []int{}) }
```


**13. Examples**
--------

When adding a new package, include examples of intended usage: a runnable Example, or a simple test demonstrating a complete call sequence.

Read more [here](https://blog.golang.org/examples)


**14. Handle Errors**
--------

Do not discard error using `_` variables. If a function returns an error, check it to make sure the function succeeded. Handle the error, return it, or, in truly exceptional situations, panic.


**15. In-Band Errors**
--------

Go's support for multiple return values provides a better solution. Instead of requiring clients to check for an in-band error value, a function should return an additional value to indicate whether its other return values are valid. This return value may be an error, or a boolean when no explanation is needed. It should be the final return value.

```
// Lookup returns the value for key or ok=false if there is no mapping for key.
func Lookup(key string) (value string, ok bool)
```

Encourages more robust and readable code
```
value, ok := Lookup(key)
if !ok  {
return fmt.Errorf("no value for %q", key)
}
return Parse(value)
```


**16. Indent Error Flow**
--------

Try to keep the normal code path at a minimal indentation, and indent the error handling, dealing with it first. This improves the readability of the code by permitting visually scanning the normal path quickly. For instance, don't write:

```
if err != nil {
// error handling
} else {
// normal code
}
```

Write this instead
```
if err != nil {
// error handling
return // or continue, etc.
}
// normal code
```
If the `if` statement has an initialization statement, such as:
```
if x, err := f(); err != nil {
// error handling
return
} else {
// use x
}
```
then this may require moving the short variable declaration to its own line:

```
x, err := f()
if err != nil {
// error handling
return
}
// use x
```


**17. Interfaces**
--------

As our last discussion, we agree with this option for defining intefaces.

```
package producer

type Thinger interface { Thing() bool }

type defaultThinger struct{ … }
func (t defaultThinger) Thing() bool { … }

func NewThinger() Thinger { return defaultThinger{ … } }
```


**18. Line Length**
--------

Avoid uncomfortably long lines but don't add line breaks to keep lines short when they are more readable long.
Long lines seem to go with long names, try to get rid of the long names instead.
Also try to avoid a long function even though we don't have any advices of how long a function should be.


**19. Mixed Caps**
--------

This convention is different conventions for other languages. For example an unexported constant is `maxLength` not `MaxLength` or `MAX_LENGTH`.


**20. Named Result Parameters**
--------

Good:
```
func (n *Node) Parent1() *Node
func (n *Node) Parent2() (*Node, error)
```

Bad:
```
func (n *Node) Parent1() (node *Node)
func (n *Node) Parent2() (node *Node, err error)
```

If a function returns two or three parameters of the same type, or if the meaning of a result isn't clear from context, adding names may be useful in some contexts.

Good:
```
// Location returns f's latitude and longitude.
// Negative values mean south and west, respectively.
func (f *Foo) Location() (lat, long float64, err error)
```

Bad:
```
func (f *Foo) Location() (float64, float64, error)
```

**21. Mocking**
--------

Using `gomock` for mocking

Program design
------------------


**1. Use struct literal initialization**
--------

```
Tip: Use struct literal initialization to avoid invalid intermediate state. Inline struct declarations where possible.
```

Take an example below:

```
// Don't do this.
cfg := fooConfig{}
cfg.Bar = bar
cfg.Period = 100 * time.Millisecond
cfg.Output = nil

foo, err := newFoo(*fooKey, cfg)
if err != nil {
log.Fatal(err)
}
defer foo.close()
```

It’s considerably nicer to leverage so-called struct initialization syntax to construct the object all at once, in a single statement.

```
// This is better.
cfg := fooConfig{
Bar:    bar,
Period: 100 * time.Millisecond,
Output: nil,
}

foo, err := newFoo(*fooKey, cfg)
if err != nil {
log.Fatal(err)
}
defer foo.close()
```

As we construct and immediately use `cgf` object, we can put it inside `newFoo` directly.

```
// This is even better.
foo, err := newFoo(*fooKey, fooConfig{
Bar:    bar,
Period: 100 * time.Millisecond,
Output: nil,
})
if err != nil {
log.Fatal(err)
}
defer foo.close()
```


**2. Avoid nil checks via default no-op implementations**
--------

It’s much safer, and nicer, to be able to use output without having to check it for existence.

Good:
```
func (f *foo) process() {
fmt.Fprintf(f.Output, "start\n")
// ...
}
```

Bad:
```
func (f *foo) process() {
if f.Output != nil {
fmt.Fprintf(f.Output, "start\n")
}
// ...
}
```


**3. Loggers are dependencies**
--------

```
Tip: Loggers are dependencies, just like references to other components, database handles, commandline flags, etc.
```

```
func (f *foo) process() {
fmt.Fprintf(f.Output, "start\n")
result := f.Bar.compute()
log.Printf("bar: %v", result) // Whoops!
// ...
}
```

`fmt.Printf` is self-contained and doesn’t affect or depend on global state; in functional terms, it has something like referential transparency. So it is not a dependency.

`log.Printf` acts on a package-global logger object, it’s just obscured behind the free function `Printf`. So it, too, is a dependency.

What do we do with dependencies? **We make them explicit**. Because the process method prints to a log as part of its work, either the method or the foo object itself needs to take a logger object as a dependency. For example, log.Printf should become f.Logger.Printf.



```
func (f *foo) process() {
fmt.Fprintf(f.Output, "start\n")
result := f.Bar.compute()
f.Logger.Printf("bar: %v", result) // Better.
// ...
}
```


Logging and instrumentation
------------------

- Log only actionable information, which will be read by a human or a machine
- Avoid fine-grained log levels — info and debug are probably enough
- Use structured logging — I’m biased, but I recommend go-kit/log
- Loggers are dependencies!

Where logging is expensive, instrumentation is cheap. You should be instrumenting every significant component of your codebase.

Let’s use loggers and metrics to pivot and address global state more directly. Here are some facts about Go:

- log.Print uses a fixed, global log.Logger
- http.Get uses a fixed, global http.Client
- http.Server, by default, uses a fixed, global log.Logger
- database/sql uses a fixed, global driver registry
- func init exists only to have side effects on package-global state

Example:
```
func foo() {
resp, err := http.Get("http://zombo.com")
// ...
}
```

`http.Get` calls on a global in package `http`. It has an implicit global dependency. Which we can eliminate pretty easily.

```
func foo(client *http.Client) {
resp, err := client.Get("http://zombo.com")
// ...
}
```

Just pass an `http.Client` as a parameter. But that is a concrete type, which means if we want to test this function we also need to provide a concrete `http.Client`, which likely forces us to do actual HTTP communication. Not great. We can do one better, by passing an interface which can Do (execute) HTTP requests.

```
type Doer interface {
Do(*http.Request) (*http.Response, error)
}

func foo(d Doer) {
req, _ := http.NewRequest("GET", "http://zombo.com", nil)
resp, err := d.Do(req)
// ...
}
```

`http.Client` satisfies our Doer interface automatically, but now we have the freedom to pass a mock Doer implementation in our test. And that’s great: a unit test for func foo is meant to test only the behavior of foo, it can safely assume that the `http.Client` is going to work as advertised.


Testings
----------------------

```
Tip: Tests only need to test the thing being tested.
``` 

As in the `http.Client` example just above, remember that unit tests should be written to test the thing being tested, and nothing more. If you’re testing a process function, there’s no reason to also test the HTTP transport the request came in on, or the path on disk the results get written to. Provide inputs and outputs as fake implementations of interface parameters, and focus on the business logic of the method or component exclusively.


```
Tip: Use many small interfaces to model dependencies.
```

In general, the thing that seems to work the best is to write Go in a generally functional style, where dependencies are explicitly enumerated, and provided as small, tightly-scoped interfaces whenever possible. Beyond being good software engineering discipline in itself, it feels like it automatically optimizes your code for easy testing.


Build and deploy
----------------------

```
Tip: Prefer go install to go build.
```

Prefer `go install` to `go build`. The install verb caches build artifacts from dependencies in `$GOPATH/pkg`, making builds faster. It also puts binaries in `$GOPATH/bin`, making them easier to find and invoke.


