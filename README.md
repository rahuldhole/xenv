# Usage

Install (optional):
```bash
go install github.com/rahuldhole/xenv@latest
```

Run with a template (supports .xenv / .template / .example / plain .env):
```bash
xenv path/to/config.xenv
```

Flags:
```
-o, --output <file>   Write to specific output file (default: dot-prefixed name of template)
-d, --defaults        Use existing/ template defaults without interactive prompts
-r, --run-scripts     Run all inline scripts automatically (no confirmations)
-m, --merge           Merge with existing output (preserve undeclared keys, show conflicts)
-f, --force           Overwrite existing output file without prompting
-h, --help            Show detailed help and exit
```
Rules:
- Do not combine --merge (-m) and --force (-f); they are mutually exclusive.
- If neither -m nor -f is given and the output file exists, an interactive choice (overwrite / merge / cancel) is shown.
- You can combine -d (defaults) with -r (run all scripts) to produce a fully auto-generated file.
- With --defaults and no --run-scripts, script-capable fields keep existing/ template values.

Examples:
```bash
# Interactive (asks)
xenv config.xenv # You can use any file extension  (.xenv, .template, .example, etc.)

# Force overwrite using defaults only
xenv config.xenv -d -f

# Merge into existing file, run scripts non-interactively
xenv config.xenv -m -r

# Generate to custom file with defaults & scripts
xenv config.xenv -o .env.production -d -r

# Show help
xenv -h
```

Directive examples:
```env
# @text label="Database Host"
DB_HOST=localhost

# @password label="Database Password"
DB_PASSWORD=

# @select label="Database Type" options=mysql,postgres,sqlite
DB_TYPE=mysql

# @checkbox label="Enable SSL"
DB_SSL=false

# @text label="Show Greeting" script=`echo "Hello $DB_HOST"`
GREETING=
```

Supported directives (short list):
@text @number @float @date @file @url @password @select @boolean @list
@checkbox @color @datetime @email @image @month @radio @range @reset
@tel @time @week @readonly @hidden @skip @button (script only)

Scripts:
Add script="..." or script=`...` to any directive (e.g. @text or @button) to run a shell snippet. Use -r / --run-scripts to run all automatically (also honored in defaults mode).

Output:
Generates a dot-prefixed file derived from template name unless -o/--output is specified.
