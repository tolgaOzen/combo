# Combo CLI

Combo is a command-line tool designed to simplify Git workflows by leveraging AI for generating commit messages, branch names, changelogs, and linking commits to issue trackers. It integrates seamlessly into your development process to save time and ensure consistency, especially with conventional commit standards.

---

## Features

- **Commit Message Generation**: Create meaningful and consistent commit messages based on conventional commit standards and your staged changes with AI assistance.
- **Branch Name Suggestions**: Generate branch names aligned with your changes or ticket references.
- **Issue Tracker Integration**: Link commits to issue trackers by specifying an issue ID.
- **Changelog Creation**: Automatically generate changelogs based on your Git history.

---

## Installation

### Using Homebrew

1. Add the Combo tap:
   ```bash
   brew tap tolgaOzen/combo
   ```

2. Install Combo:
   ```bash
   brew install combo
   ```

3. Verify installation:
   ```bash
   combo --version
   ```

---

## Usage

### Commands

1. **Commit Message**:
   Generate and confirm a commit message:
   ```bash
   combo commit
   ```

   Example:
   ```
   Here’s your commit message:

   ➤ feat(config): add support for managing API keys and locales

   Would you like to use this message? (Y/n):
   ```

    - `feat`: Indicates a new feature.
    - `config`: Scope of the change (optional but recommended).
    - Description: Clear and concise explanation of the change.

2. **Branch Name**:
   Generate a branch name based on your changes or tickets:
   ```bash
   combo branch
   ```

   Example:
   ```
   feat/config-api-key-support
   ```

3. **Issue Tracker**:
   Link a commit to an issue tracker by providing the issue ID:
   ```bash
   combo issue 123
   ```

   Example:
   ```
   Linked commit to issue #123.
   ```

4. **Changelog**:
   Generate a changelog from your Git history:
   ```bash
   combo changelog
   ```

---

## Configuration

Combo stores its configuration in `~/.combo/config`. The configuration file is created automatically if it doesn’t exist.

### Example Configuration

```plaintext
# ~/.combo/config
openai_api_key=sk-xxxxxxxxxx
prompt_locale=en-US
prompt_max_length=72
```

### Managing Configuration

- **Set a key**:
  ```bash
  combo config set <key> <value>
  ```
  Example:
  ```bash
  combo config set openai_api_key sk-xxxxxxxxxx
  ```

- **Get a key**:
  ```bash
  combo config get <key>
  ```
  Example:
  ```bash
  combo config get prompt_locale
  ```

---

## Example Workflow

1. **Stage Changes**:
   ```bash
   git add .
   ```

2. **Generate and Confirm a Commit Message**:
   ```bash
   combo commit
   ```

   Example Interaction:
   ```
   Here’s your commit message:

   ➤ fix(auth): resolve token expiration issue for multi-tenant support

   Would you like to use this message? (Y/n):
   ```

3. **Generate a Branch Name**:
   ```bash
   combo branch
   ```

   Example:
   ```
   fix/auth-token-expiration
   ```

4. **Link a Commit to an Issue**:
   ```bash
   combo issue 456
   ```

   Example:
   ```
   Linked commit to issue #456.
   ```

5. **Generate a Changelog**:
   ```bash
   combo changelog
   ```

   Example:
   ```
   ## [1.2.0] - 2024-01-01

   ### Features
   - feat(config): add support for managing API keys and locales

   ### Bug Fixes
   - fix(auth): resolve token expiration issue for multi-tenant support
   ```

---

## Requirements

- **Git**: Must be installed and configured.
- **OpenAI API Key**: Required for AI-powered features.

---

## Contributing

Contributions are welcome! Submit issues or pull requests on the GitHub repository.

---

## License

Combo is licensed under the MIT License. See the `LICENSE` file for more details.
```