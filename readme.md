# Ask.md
Go script to run [ask.md](https://x.com/yacineMTB/status/1789368952956555430) anywhere. Too many browser specific extensions to do this, when it should be a simple go script to watch your directory.

Create an `ask.md` in your project and start talking to it.

## Requirements
- OpenAI API Key
- Go >= 1.23.2

## Usage
1) Clone this project
2) Set env var `OPENAI_API_KEY` to your openai key
3) Build `go build -o ask`
4) Run `./ask watch` to watch your working directory `ask.md`


Put the built binary in your bin or path for better results.
