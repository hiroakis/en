# en

The "en" is an environment variables management tool for Circle CI. "en" means "Circle" in Japanese(円). Kanji character(円) is the same as Japansese currency "Yen".

# Installation

```
go get github.com/hiroakis/en
```

# How to use

You have to get Circle CI API token before you use en. You can set it to `CIRCLE_TOKEN` environment variable or `-token` argument option.

* show variables
```
en -export
```

Output will be json. If you would like to export variables, you can redirect stdout to file. Like this `en > en.json`.
Note: Exported variables will be masked by Circle CI as follows. Edit it before you apply.

```
{
    "name": "AWS_SECRET_ACCESS_KEY",
    "value": "xxxxS73Z"
}
```

* dry-run
```
en -apply -dry-run -file path/to/en.json
```

* apply
```
en -apply -file path/to/en.json
```

* print version
```
en -version
```

* print usage
```
en -help
```

# Configuration

`en.example.json` is an example.

# License

MIT