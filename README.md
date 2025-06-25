# ticktick-tui - A terminal-based interface for TickTick

This is a terminal-based interface for TickTick, a popular task management application. The goal of this project is to provide a simple and efficient way to manage your tasks directly from the terminal instead of using a web browser or app.

Windows, macOS, and Linux are supported.

## Features

- View and manage tasks in a terminal.
- Beautiful TUI interface for better user experience.
- Using official TickTick API for task management.

Since the official TickTick API does not cover all features of original application, this project can only provide basic functions. Sub-tasks, attachments, countdown, habits and many other features are not supported yet. Eisenhower Matrix and TimeLine view may be supported unofficially.

## Installation

Download from releases page or build from source.

Building from source requires Go version 1.24.4 or higher.

```bash
go build -o your-program ./main.go
```

## Usage

**It is highly recommended to use TUI interface for better user experience.**
**Commands are also available, but created for scripting purposes.**

## Roadmap

- [ ] TUI interface (Create and delete tasks)
- [ ] README.md
- [ ] I18n
- [ ] Code refactoring and reuse
- [ ] Improve user experience
- [ ] Prefetching tasks at background
- [ ] Eisenhower Matrix view
- [ ] TimeLine view
- [ ] Next `n` days view
- [ ] AI completion integration
- [ ] Search tasks

## Contributing

Contributions are welcome! If you have suggestions for improvements or new features, please open an issue or submit a pull request.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.