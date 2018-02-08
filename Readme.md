
  `terraform-docs(1)` &sdot; a quick utility to generate docs from terraform modules.

This repository is a fork of `github.com/segmentio/terraform-docs` with some changes.

![Terraform docs](docs.gif)

## Installation

  - `go get github.com/Gufran/terraform-docs`
  - [Binaries](https://github.com/Gufran/terraform-docs/releases)

## Usage

```bash

  Usage:
    terraform-docs [--no-required] [--json] <path>...
    terraform-docs -h | --help

  Examples:

    # View inputs and outputs in markdown
    $ terraform-docs ./my-module

    # Generate JSON docs for module
    $ terraform-docs --json ./my-module

    # Generate markdown docs for module
    $ terraform-docs md ./my-module

    # Generate markdown but don't print "Required" column
    $ terraform-docs --no-required md ./my-module

  Options:
    -h, --help     show help information

```

## Syntax

Every comment group with `@doc` annotation on first line will be included in documentation, the
attribute can be written in following ways:

 - `@doc()`: Comment group is the introduction block at the top of the page
 - `@doc(resource)`: Comment group is the introduction block at the top on *Resources* section
 - `@doc(data-provider)`: Comment group is the introduction block at the top of *Data Providers* section
 - `@doc(iam-policy)`: Comment group is the introduction block at the top of *IAM Policies* section
 - `@doc(input)`: Comment group is the introduction block at the top of *Inputs* section
 - `@doc(output)`: Comment group is the introduction block at the top of *Outputs* section

All resources, data providers, iam policies, inputs and outputs can further have their own
description.

For inputs and outputs, the value of `description` attribute is used. If `description` attribute is
not available or has empty value then the associated comment (see below) is used as description.

For resources, data providers, and iam policies, the associated comment is used.

An associated comment is a comment group immediately followed by a configuration block without any
line break in between. For example:

``` hcl
// This comment will be used as description because there is no
// new line between comment and configuration
resource "aws_instance" "app" {
  // configuration...
}

// This comment will be ignore because it is followed by a line
// break and then configuration and does not have a valid @doc
// attribute

data "aws_ami" "app" {
  // configuration...
}
```

[ex]: ./_example/main.tf


## License

MIT License

Copyright (c) 2017 Mohammad Gufran

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
