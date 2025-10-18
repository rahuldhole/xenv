# xenv

## Overview
xenv is an interactive Ruby CLI tool for configuring environment variable files using a simple `.xenv` template. It guides users through input, password, select, and checkbox prompts, then writes the results to a `.env`-style file.

## Installation

Clone the repository and install dependencies:

```sh
git clone [repository-url]
cd xenv
bundle install
```

## Usage

### Run

To run the tool directly:

```sh
ruby xenv.rb your.env.xenv
```

This will interactively prompt you for values and write the output file (e.g., `.env`) in your current directory.

### Generate binary 

To package the tool into a standalone binary using `tebako`:

```sh
tebako press \
  --entry-point=xenv.rb \
  --root=. \
  --output=xenv \
  --Ruby=3.4.2
```

## Running Tests

To run tests (if you add them):

```sh
rspec
```

## Contributing

Contributions are welcome! Please open an issue or submit a pull request for improvements or bug fixes.

## License

This project is licensed under the [Your License] - see the LICENSE file for details.