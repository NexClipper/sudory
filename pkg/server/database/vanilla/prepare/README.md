# Usage Prepare Query

no case sensitive

## Query

`WHERE condition`

| operator | keyword
|---|---
| AND | and
| OR | or
| NOT | not
| = | eq, equal
| &gt; | gt
| &lt; | lt
| &ge; | ge, gte
| &le; | le, lte
| LIKE | like
| IS NULL | isnull
| IN |in
| BETWEEN | between

### Query Example

- WHERE foo = 123

    ```json
    {"EQ" : {"foo" : 123}}
    ```

    ```lisp
    (EQ foo 123)
    ```

- WHERE foo = "bar"

    ```json
    {"EQ" : {"foo" : "bar"}}
    ```

    ```lisp
    (EQ foo "bar")
    or
    (EQ foo bar)
    ```

- WHERE foobar = "foo bar"

    ```json
    {"EQ" : {"foo" : "foo bar"}}
    ```

    ```lisp
    (EQ foobar "foo bar")
    ```

- WHERE (foo = 123)

    ```json
    {"AND" : {"EQ" : {"foo" : 123}}}
    ```

    ```lisp
    (AND (EQ foo 123))
    ```

- WHERE NOT foo = 123

    ```json
    {"NOT" : {"EQ" : {"foo" : 123}}}
    ```

    ```lisp
    (NOT (EQ foo 123))
    ```

- WHERE (foo = 123 AND bar = 456)

    ```json
    {"AND" : [{"EQ" : {"foo" : 123}}, {"EQ" : {"bar" : 456}}]}
    ```

    ```lisp
    (AND (EQ foo 123) (EQ bar 456))
    ```

- WHERE (foo = 123 OR bar = 456)

    ```json
    {"OR" : [{"EQ" : {"foo" : 123}}, {"EQ" : {"bar" : 456}}]}
    ```

    ```lisp
    (OR (EQ foo 123) (EQ bar 456))
    ```

- WHERE (foo = 123 AND bar = 456 AND (baz = 789 OR foobar = 123456))

    ```json
    {"AND" : [{"EQ" : {"foo" : 123}}, {"EQ" : {"bar" : 456}}, {"OR" : [{"EQ" : {"baz" : 789}}, {"EQ" : {"foobar" : 123456}}]}]}
    ```

    ```lisp
    (AND (EQ foo 123) (EQ bar 456) (OR (EQ baz 789) (EQ foobar 123456)))
    ```

- WHERE foo &gt; 123

    ```json
    {"GT" : {"foo" : 123}}
    ```

    ```lisp
    (GT foo 123)
    ```

- WHERE foo &lt; 123

    ```json
    {"LT" : {"foo" : 123}}
    ```

    ```lisp
    (LT foo 123)
    ```

- WHERE foo &ge; 123

    ```json
    {"GE" : {"foo" : 123}}
    ```

    ```lisp
    (GE foo 123) 
    ```

- WHERE foo &le; 123

    ```json
    {"LE" : {"foo" : 123}}
    ```

    ```lisp
    (LE foo 123) 
    ```

- WHERE foo LIKE 'string%'

    ```json
    {"LIKE" : {"foo" : "string%"}}
    ```

    ```lisp
    (LIKE foo "string%") 
    or 
    (LIKE foo string%) 
    ```

- WHERE foo IS NULL

    ```json
    {"ISNULL" : "foo"}
    ```

    ```lisp
    (ISNULL foo) 
    ```

- WHERE foo IN (123)

    ```json
    {"IN" : {"foo" : 123}}

    {"IN" : {"foo" : [123]}}
    ```

    ```lisp
    (IN foo 123) 
    or
    (IN foo `(123)) 
    ```

- WHERE foo IN (123, 456, 789)

    ```json
    {"IN" : {"foo" : [123, 456, 789]}}
    ```

    ```lisp
    (IN foo 123 456 789) 
    or
    (IN foo `(123 456 789)) 
    ```

- WHERE foo BETWEEN 123 AND 456

    ```json
    {"BETWEEN" : {"foo" : [123, 456]}}
    ```

    ```lisp
    (BETWEEN foo 123 456) 
    or
    (BETWEEN foo `(123 456)) 
    ```

## Pagination

`LIMIT [OFFSET, ] ROW_COUNT`

| operator | keyword
|---|---
| ROW_COUNT | LIMIT
| PAGE | PAGE
| OFFSET | (PAGE - 1) * LIMIT

### Pagination Example

- LIMIT 0, 255

    ```json
    {}
    ```

- LIMIT 10

    ```json
    {"limit":10}
    ```

- LIMIT 0, 10

    ```json
    {"limit" : 10, "page" : 1}
    ```

## Order

`ORDER BY column_1, column_2 [ASC|DESC]`

| operator | keyword
|---|---
| ASC | ASC
| DESC | DESC

### Order Example

- ORDER BY column_name_1, column_name_2

    ```json
    ["column_name_1", "column_name_2"]
    ```

- ORDER BY column_name_1, column_name_2 DESC

    ```json
    ["column_name_1", "column_name_2", "DESC"]

    {"DESC" : ["column_name_1", "column_name_2"]}
    ```

- ORDER BY column_name_1, column_name_2 ASC

    ```json
    ["column_name_1", "column_name_2", "ASC"]

    {"ASC" : ["column_name_1", "column_name_2"]}
    ```

- ORDER BY column_name_1 ASC, column_name_2 DESC

    ```json
    ["column_name_1", "ASC", "column_name_2", "DESC"]

    [{"ASC" : "column_name_1"}, {"DESC" : "column_name_2"}]
    ```

- ORDER BY column_name_1, column_name_2 ASC, column_name_3 DESC

    ```json
    ["column_name_1", "column_name_2", "ASC", "column_name_3", "DESC"]

    [{"ASC" : ["column_name_1", "column_name_2"]}, {"DESC" : "column_name_2"}]
    ```
