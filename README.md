# xenv

A simple interactive tool for configuring environment variables from template files.

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

Create a template file with `.xenv` extension:

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
```

This will generate `.config` with your configured values.

## Template Format

xenv uses special comment directives to specify how each value should be prompted:

- `@input`: Simple text input
  - Options: `label="Display Label"`
  
- `@password`: Password input (text will be hidden)
  - Options: `label="Display Label"`
  
- `@select`: Selection from predefined options
  - Options: `label="Display Label" options=option1,option2,option3`
  
- `@checkbox`: Yes/No boolean input
  - Options: `label="Display Label"`

Each directive should be placed in a comment line directly above the variable definition.

## Examples

See the `examples/` directory for sample configuration templates.

## License

MIT
