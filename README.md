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
- -o <file>            Specify output file
- --defaults           Use existing/ template defaults without prompts
- --run-scripts        Run all inline scripts automatically
- Combine --defaults --run-scripts to generate fully scripted defaults

Merge / overwrite:
If output exists you will be prompted: y (overwrite) / m (merge) / N (cancel).

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
Add script="..." or script=`...` to any directive (e.g. @text or @button) to run a shell snippet. Use --run-scripts to run them without confirmation.

Output:
Generates a dot-prefixed file derived from template name unless -o is specified.
