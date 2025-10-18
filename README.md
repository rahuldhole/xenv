# xenv

## Overview
xenv is an interactive Ruby CLI tool for configuring environment variable files using a simple `.form` template. It guides users through input, password, select, and checkbox prompts, then writes the results to a `.env`-style file.

## Installation

Clone the repository and install dependencies:

```bash
git clone [repository-url]
cd xenv
bundle install
```

## Usage

### With Ruby

To run the tool directly:

```bash
ruby xenv.rb your.env.form
```

### With Docker

Build the Docker image:

```bash
docker build -t xenv .
```

Run the tool (replace `your.env.form` with your form file):

```bash
docker run -it --rm -v "$PWD":/app xenv your.env.form
```

This will interactively prompt you for values and write the output file (e.g., `.env`) in your current directory.

## Running Tests

To run tests (if you add them):

```bash
rspec
```

## Contributing

Contributions are welcome! Please open an issue or submit a pull request for improvements or bug fixes.

## License

This project is licensed under the [Your License] - see the LICENSE file for details.