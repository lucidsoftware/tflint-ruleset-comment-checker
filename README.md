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
plugin "comment-checker" {

  enabled = true

  version = "0.1.1"
  source  = "github.com/lucidsoftware/tflint-ruleset-comment-checker"

  # Optionally omit this to use Keyless verification
  signing_key = <<-KEY
    -----BEGIN PGP PUBLIC KEY BLOCK-----

    mG8EaXACCRMFK4EEACIDAwQdSSKnORcu1YozK8MQMrLJ4LBN171J/Zf3G//FUxX8
    hvlh1CyPvcTgi1UYuj8wWCF19L2GazNv32MmPDk9ueGzfmTsp5ONHddg4Tiu6SZV
    zgfkfyrhJfq9h4A1FTYq0oC0JUx1Y2lkIFNvZnR3YXJlLCBJbmMuICh0ZmxpbnQg
    c2lnbmluZymIswQTEwkAOxYhBC/Y2UK9yTIieUACoT8p7Wh2KOJWBQJpcAIJAhsD
    BQsJCAcCAiICBhUKCQgLAgQWAgMBAh4HAheAAAoJED8p7Wh2KOJWUV4BgIuR7YuX
    ON5YSi9+XPbg6zxEihHRp9NWs76ipYHJdd5eXKRYeW69MrQgr7TY5TyDFwF/Q2Be
    1Xy0JwzT5zmg6vnwPc3+9I5oq7rWEEbKAP4PZd0pYLsH5MqghYrcE1FCXf7+
    =LNHa
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

It also recursively searches through any literal objects passed to the module for keys that match the name.

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
