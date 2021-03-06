# Changelog

## [v0.2.0](https://github.com/itchyny/gojq/compare/v0.1.0..v0.2.0) (2019-05-06)

* implement binding variable syntax (`... as $var`)
* implement `try` `catch` syntax
* implement alternative operator (`//`)
* implement various functions (`in`, `to_entries`, `startswith`, `endswith`,
  `ltrimstr`, `rtrimstr`, `combinations`, `ascii_downcase`, `ascii_upcase`,
  `tojson`, `fromjson`)
* support query for object indexing
* support object construction with variables
* support indexing against strings

## [v0.1.0](https://github.com/itchyny/gojq/compare/v0.0.1..v0.1.0) (2019-05-02)

* implement binary operators (`+`, `-`, `*`, `/`, `%`, `==`, `!=`, `>`, `<`,
  `>=`, `<=`, `and`, `or`)
* implement unary operators (`+`, `-`)
* implement booleans (`false`, `true`), `null`, number and string constant
  values
* implement `empty` value
* implement conditional syntax (`if` `then` `elif` `else` `end`)
* implement various functions (`length`, `utf8bytelength`, `not`, `keys`,
  `has`, `map`, `select`, `recurse`, `while`, `until`, `range`, `tonumber`,
  `type`, `arrays`, `objects`, `iterables`, `booleans`, `numbers`, `strings`,
  `nulls`, `values`, `scalars`, `reverse`, `explode`, `implode`, `join`)
* support function declaration
* support iterators in object keys
* support object construction shortcut
* support query in array indices
* support negative number indexing against arrays
* support json file name arguments

## [v0.0.1](https://github.com/itchyny/gojq/compare/0fa3241..v0.0.1) (2019-04-14)

* initial implementation
