# xenv

A simple interactive tool for configuring environment variables from template files with built-in validation.

## Installation

```bash
go install github.com/rahuldhole/xenv@latest
```

Or clone and build from source:

```bash
git clone https://github.com/rahuldhole/xenv.git
cd xenv
go build -o xenv
```

## Usage

Create a template file with `.xenv`, `.template`, or `.example` extension, or use a plain `.env` file:

```
# Database Configuration
# @input label="Database Host"
DB_HOST=localhost

# @password label="Database Password"
DB_PASS=default_password

# @select label="Database Type" options=mysql,postgres,sqlite
DB_TYPE=mysql

# @checkbox label="Enable SSL"
DB_SSL=false
```

Run xenv with your template file:

```bash
xenv env.config.xenv
# or you can re-execute existing .env file
xenv .env
```

This will generate `.config` with your configured values.

## Template Format

xenv uses special comment directives to specify how each value should be prompted:

- `@text`: Simple text input
  - Options: `label="Display Label"` `note="Help text"` `pattern="regex"`
  - Validation: Any text, custom regex pattern optional
  
- `@number`: Integer number input
  - Options: `label="Display Label"` `note="Help text"`
  - Validation: Must be a valid integer (e.g., 123, -456)
  
- `@float`: Floating point number input
  - Options: `label="Display Label"` `note="Help text"`
  - Validation: Must be a valid float (e.g., 3.14, -0.5, 42)
  
- `@date`: Date input
  - Options: `label="Display Label"` `note="Help text"` `pattern="regex"`
  - Validation: Default YYYY-MM-DD format, custom regex pattern optional
  
- `@file`: File path input
  - Options: `label="Display Label"` `note="Help text"`
  - Validation: Must look like a valid file path
  
- `@url`: URL input
  - Options: `label="Display Label"` `note="Help text"`
  - Validation: Must be a valid http:// or https:// URL
  
- `@password`: Password input (text will be hidden)
  - Options: `label="Display Label"` `note="Help text"`
  - Validation: Any text
  
- `@select`: Selection from predefined options
  - Options: `label="Display Label"` `options=option1,option2,option3` `note="Help text"`
  - Validation: Must select from provided options

- `@boolean`: Boolean input (true/false)
  - Options: `label="Display Label"` `note="Help text"`
  - Validation: Must be true/false, yes/no, or y/n

- `@list`: Comma-separated list input
  - Options: `label="Display Label"` `note="Help text"`
  - Validation: Must be comma-separated values

- `@skip`: Skip all variables below this directive (preserve existing values)
  - Options: `note="Reason for skipping"`

Each directive should be placed in a comment line directly above the variable definition.

**Custom Validation:**
- Use `pattern="regex"` for `@text` and `@date` fields to enforce custom validation rules
- Example: `pattern="^\d{3}-\d{2}-\d{4}$"` for SSN format

## Examples

See the `examples/` directory for sample configuration templates.

## Building Platform-Independent Binaries

You can build binaries for different platforms using Go's cross-compilation:

```sh
# Linux (amd64)
GOOS=linux GOARCH=amd64 go build -o xenv-linux-amd64

# macOS (amd64)
GOOS=darwin GOARCH=amd64 go build -o xenv-darwin-amd64

# macOS (Apple Silicon)
GOOS=darwin GOARCH=arm64 go build -o xenv-darwin-arm64

# Windows (amd64)
GOOS=windows GOARCH=amd64 go build -o xenv-windows-amd64.exe
```

You can set `GOOS` and `GOARCH` to any valid Go target. See the [Go documentation](https://golang.org/doc/install/source#environment) for more options.

## License

MIT

## Special: Inline Shell Scripts

You can use `script="..."` on any field (commonly with `@button`) to run a shell script.  
The script will be shown before execution and you can choose to run or skip it.

Example:
