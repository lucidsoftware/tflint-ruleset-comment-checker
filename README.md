# TFLint Ruleset For Requiring Comments
[![Build Status](https://github.com/lucidsoftware/tflint-ruleset-comment-checker/actions/workflows/build.yml/badge.svg?branch=main)](https://github.com/lucidsoftware/tflint-ruleset-comment-checker/actions)

This is a [tflint](https://github.com/terraform-linters/tflint) plugin that allows you to specify attributes on modules that are required to have comments on them
(for example to explain why a non-default value is being used).

In the future thsi may be expanded to support requiring comments in other places as well.

## Requirements

- TFLint v0.46+
- Go v1.25.3

## Installation

You can install the plugin with `tflint --init`. Declare a config in `.tflint.hcl` as follows:

```hcl
plugin "template" {

  enabled = true

  version = "0.1.0"
  source  = "github.com/lucidsoftware/tflint-ruleset-comment-checker"

  signing_key = <<-KEY
    -----BEGIN PGP PUBLIC KEY BLOCK-----

    mDMEaWtJCxYJKwYBBAHaRw8BAQdAJZ4V4dCqaxSlpC+BEz6xe7I6dqezSPVh/dgS
    T/4UY3u0JUx1Y2lkIFNvZnR3YXJlLCBJbmMuICh0ZmxpbnQgc2lnbmluZymIlgQT
    FgoAPgIbAwULCQgHAgIiAgYVCgkICwIEFgIDAQIeBwIXgBYhBOxebUuPZXpkg4Vq
    +MBvK4kkvJxJBQJpa0nFAhkBAAoJEMBvK4kkvJxJ58MA+gM3Z5LLsk5FA/1UvNpA
    5a+g+roGFd7G0x1zL23vFQEwAP9QaUn96ez9XHdvVaq9q0RAeft+STQV91YCwv1V
    O/r8Aw==
    =sc5Q
    -----END PGP PUBLIC KEY BLOCK-----
  KEY
}
```

## Rules

|Name|Description|Severity|Enabled|Link|
| --- | --- | --- | --- | --- |
|module_attribute_comments|Checks if specified module call attributes have comments immediately preceding them|ERROR|-||

### module_attribute_comments

This rule checks that specified attributes in Terraform module calls have comments immediately preceding them. This is useful for ensuring that important configuration decisions are documented.

**Configuration:**

The attribute names are configured at the plugin level:

```hcl
rule "module_attribute_comments" {
  enabled = true

  attribute {
    name = "instance_type"
    message = "Explain why default instance_type was overriden."
  }


  attribute {
    name = "count"
    # message is optional, but recommended
  }
}
```

**Example of valid code:**

```hcl
module "example" {
  source = "./modules/example"

  # Specifying t2.micro for cost optimization
  instance_type = "t2.micro"

  # Running 3 instances for high availability
  count = 3
}
```

**Example of invalid code:**

```hcl
module "example" {
  source = "./modules/example"
  instance_type = "t2.micro"
  count = 3
}
```

## Building the plugin

Clone the repository locally and run the following command:

```
$ make
```

You can easily install the built plugin with the following:

```
$ make install
```

You can run the built plugin like the following:

```
$ cat << EOS > .tflint.hcl
plugin "comment-checker" {
  enabled = true
}
EOS
$ tflint
```
