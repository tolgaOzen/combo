# 🚀 Combo CLI
  
**AI-Powered Git Workflow Automation**
  
[![Version](https://img.shields.io/badge/version-v0.2.3-blue.svg)](https://github.com/tolgaOzen/combo/releases)
[![License](https://img.shields.io/badge/license-MIT-green.svg)](LICENSE)
[![Go Report Card](https://goreportcard.com/badge/github.com/tolgaOzen/combo)](https://goreportcard.com/report/github.com/tolgaOzen/combo)
  
*Simplify your Git workflows with AI-generated commit messages, branch names, and more.*


---

## ✨ Features

🤖 **Smart Commit Messages** - Generate meaningful commit messages following conventional commit standards  
🌿 **Intelligent Branch Names** - Create descriptive branch names based on your changes  
🔗 **Issue Tracker Integration** - Link commits to issues automatically  
📚 **Changelog Generation** - Auto-generate changelogs from Git history  
🌍 **Multi-language Support** - Support for 15+ languages and locales  
⚙️ **Configurable** - Customize prompts, formats, and AI behavior

## 📦 Installation

### 🍺 Homebrew (Recommended)

```bash
# Add the Combo tap
brew tap tolgaOzen/tap

# Install Combo
brew install combo

# Verify installation
combo version
```

### 📥 Direct Download

Download the latest release from [GitHub Releases](https://github.com/tolgaOzen/combo/releases).

### 🔧 Build from Source

```bash
git clone https://github.com/tolgaOzen/combo.git
cd combo
make build
```

## 🚀 Quick Start

1. **Set up your OpenAI API key:**
   ```bash
   combo config set openai_api_key sk-your-api-key-here
   ```

2. **Stage some changes:**
   ```bash
   git add .
   ```

3. **Generate and commit:**
   ```bash
   combo commit
   ```

## 💻 Usage

### 📋 Available Commands

| Command | Description | Example |
|---------|-------------|---------|
| `combo commit` | Generate AI-powered commit messages | `combo commit` |
| `combo branch` | Create intelligent branch names | `combo branch` |
| `combo config` | Manage configuration settings | `combo config set key value` |
| `combo version` | Show version information | `combo version` |

### 🎯 Command Details

#### 💬 Commit Messages

Generate conventional commit messages based on your staged changes:

```bash
combo commit
```

**Interactive Preview:**
```
Generating your commit message...

Here's your commit message:

➤ feat(auth): add OAuth2 integration with Google provider

Would you like to use this message? (Y/n):
```

**Commit Types Supported:**
- `feat` - New features
- `fix` - Bug fixes  
- `docs` - Documentation changes
- `style` - Code style changes
- `refactor` - Code refactoring
- `test` - Adding tests
- `chore` - Maintenance tasks

#### 🌿 Branch Names

Create descriptive branch names from your changes:

```bash
combo branch
```

**Example Output:**
```
feat/oauth2-google-integration
fix/memory-leak-user-service
docs/api-authentication-guide
```

## ⚙️ Configuration

Combo stores configuration in `~/.combo/config`. The file is created automatically with defaults.

### 🔧 Configuration Options

| Setting | Description | Default | Example |
|---------|-------------|---------|---------|
| `openai_api_key` | Your OpenAI API key | *Required* | `sk-xxx...` |
| `prompt_locale` | Language for prompts | `en-US` | `en-US`, `fr-FR`, `es-ES` |
| `prompt_max_length` | Max commit message length | `72` | `50`, `72`, `100` |

### 🛠️ Managing Configuration

```bash
# Set configuration values
combo config set openai_api_key sk-your-key-here
combo config set prompt_locale en-US
combo config set prompt_max_length 72

# Get configuration values  
combo config get openai_api_key
combo config get prompt_locale
```

### 🌍 Supported Languages

| Language | Code | Language | Code |
|----------|------|----------|------|
| English (US) | `en-US` | Korean | `ko-KR` |
| English (UK) | `en-GB` | Japanese | `ja-JP` |
| French | `fr-FR` | Chinese (Simplified) | `zh-CN` |
| Spanish | `es-ES` | Chinese (Traditional) | `zh-TW` |
| German | `de-DE` | Portuguese (Brazil) | `pt-BR` |
| Italian | `it-IT` | Russian | `ru-RU` |
| Arabic | `ar-SA` | Hindi | `hi-IN` |

## 🔄 Example Workflow

```bash
# 1. Make your changes
echo "console.log('Hello World');" > app.js

# 2. Stage changes
git add .

# 3. Generate commit message
combo commit
# Output: feat: add hello world console output

# 4. Create a new branch for next feature
combo branch
# Output: feat/user-authentication

# 5. Switch to the new branch
git checkout feat/user-authentication
```

## 📋 Requirements

- **Git**: Version 2.0 or higher
- **OpenAI API Key**: Required for AI-powered features ([Get yours here](https://platform.openai.com/api-keys))
- **Internet Connection**: For API calls to OpenAI

## 🐛 Troubleshooting

### Common Issues

**Error: "missing or empty 'openai_api_key' in configuration"**
```bash
# Solution: Set your OpenAI API key
combo config set openai_api_key sk-your-key-here
```

**Error: "no staged changes found"**
```bash
# Solution: Stage your changes first
git add .
# Then run combo commit
```

**Error: "git command failed"**
```bash
# Solution: Ensure you're in a git repository
git init
# Or check if git is installed
git --version
```

### 🔍 Debug Mode

Enable verbose logging for troubleshooting:
```bash
export COMBO_DEBUG=1
combo commit
```

## 🤝 Contributing

We welcome contributions! Here's how to get started:

1. **Fork the repository**
2. **Create a feature branch:**
   ```bash
   git checkout -b feature/amazing-feature
   ```
3. **Make your changes and commit:**
   ```bash
   combo commit  # Use combo to generate your commit message!
   ```
4. **Push to your branch:**
   ```bash
   git push origin feature/amazing-feature
   ```
5. **Open a Pull Request**

### 🏗️ Development Setup

```bash
# Clone the repository
git clone https://github.com/tolgaOzen/combo.git
cd combo

# Install dependencies
go mod tidy

# Run tests
make test

# Build the project
make build

# Run linters
make lint-all
```

### 📝 Code Style

- Follow Go conventions and `gofmt`
- Add tests for new functionality
- Update documentation for any API changes
- Use conventional commits for your contributions

## 📊 Performance

- **Fast**: Commit message generation typically takes 1-3 seconds
- **Lightweight**: Binary size ~10MB
- **Efficient**: Minimal API calls with smart caching
- **Offline-ready**: Configuration and Git operations work offline

## 🛡️ Security

- API keys are stored securely in `~/.combo/config`
- Configuration files are created with restrictive permissions (0750)
- No sensitive data is sent to external services except OpenAI API
- All file operations are validated and sanitized

## 📈 Roadmap

- [ ] Support for additional AI providers (Claude, Gemini)
- [ ] GitHub/GitLab integration for automatic PR descriptions
- [ ] Team-based configuration sharing
- [ ] Plugin system for custom commit types
- [ ] Advanced diff analysis and suggestions
- [ ] Integration with popular IDEs

## 🆘 Support

- **Issues**: [GitHub Issues](https://github.com/tolgaOzen/combo/issues)
- **Discussions**: [GitHub Discussions](https://github.com/tolgaOzen/combo/discussions)
- **Documentation**: [Wiki](https://github.com/tolgaOzen/combo/wiki)

## 📄 License

Combo is licensed under the MIT License. See the [LICENSE](LICENSE) file for more details.

---

<div align="center">
  <p>Made with ❤️ by <a href="https://github.com/tolgaOzen">Tolga Ozen</a></p>
  <p>⭐ Star this project if you find it helpful!</p>
</div>
