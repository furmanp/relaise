# Relaise


**AI-Powered Release Notes Generator**

Relaise is a CLI tool that automatically generates release notes by analyzing your Git commit history since the last annotated semantic version tag (e.g., `v1.2.3`). It uses AI to transform your commit messages into polished, human-readable release notes based on your preferences.

## Features

*   **Automatic Tag Detection:** Finds the latest annotated semantic version tag in your Git repository.
*   **Commit Analysis:** Extracts commit messages between the latest tag and the current HEAD.
*   **AI-Powered Generation:** Leverages Large Language Models (LLMs) via configurable providers (e.g., Mistral) to draft release notes.
*   **Highly Configurable:** Customize the output format, including:
   *   Tone/Mood (e.g., professional, casual, funny)
   *   Language
   *   Bullet point style
   *   Inclusion of structured sections (Features, Fixes, Other)
   *   Use of emojis
   *   Release type context (major, minor, patch) for version bumping hints.
*   **Simple Configuration:** Manage settings easily using the `relaise config` command, saving preferences to `~/.relaise/config.yaml`.
*   **Clipboard Integration:** Optionally copy the generated notes directly to your clipboard. 

## Installation

### Prerequisites

*   Go (version 1.23 or later recommended)
*   Git

### Steps

1.  **Clone the repository:**
    ```bash
    git clone https://github.com/furmanp/relaise.git
    cd relaise
    ```
2.  **Build the project:**
    ```bash
    go build -o relaise .
    ```
    (Optional) Place the built `relaise` binary in a directory included in your system's PATH for easier access, or use `go install`:
    ```bash
    go install .
    ```

## Configuration

Before first use, you need to configure Relaise, primarily to set your AI provider's API key.

1.  **Run the `config` command:**
    Use `relaise config` with flags to set your desired defaults. At a minimum, set your API key.
    ```bash
    relaise config --api-key YOUR_AI_PROVIDER_API_KEY
    ```

2.  **Set other preferences (optional):**
    You can configure other options simultaneously or update them later.
    ```bash
    # Example: Configure Mistral provider, model, language, and enable emojis
    relaise config \
      --provider mistral \
      --model mistral-small-latest \
      --language en \
      --emojis true \
      --api-key YOUR_MISTRAL_API_KEY
    ```

Configuration is saved to `~/.relaise/config.yaml`.

### Configuration Options

*   `--api-key`: (Required) Your API key for the chosen LLM provider.
*   `--provider`: AI provider (default: `mistral`). See [gollm providers](https://github.com/teilomillet/gollm?tab=readme-ov-file#supported-providers) for options.
*   `--model`: Specific model to use (default: `mistral-small-latest`).
*   `--mood`: Tone for the release notes (default: `professional`).
*   `--release-type`: Type of release (default: `minor`) - hints at version bumping.
*   `--bullet-style`: Bullet character for lists (default: `-`).
*   `--include-sections`: Group notes into Features/Fixes/Other (default: `false`).
*   `--language`: Output language code (default: `en`).
*   `--emojis`: Include relevant emojis (default: `false`).
*   `--copy`: Copy generated notes to clipboard (default: `false`).

## Usage

Navigate to your Git repository's root directory.

1.  **Generate release notes using saved configuration:**
    ```bash
    relaise
    ```

2.  **Generate notes overriding some configurations:**
    Flags provided during execution override the settings from the configuration file for that specific run.
    ```bash
    # Generate notes for a patch release with a funny mood and copy to clipboard
    relaise --release-type patch --mood funny --copy true
    ```

## How It Works

1.  **Find Last Tag:** Relaise looks for the most recent annotated tag matching a semantic version pattern (e.g., `v1.0.0`, `v2.3.4-rc1`).
2.  **Collect Commits:** It gathers all Git commit messages made *after* that tag up to the current `HEAD`.
3.  **Load Config:** It reads your settings from `~/.relaise/config.yaml`.
4.  **Apply Overrides:** Command-line flags temporarily override corresponding settings from the config file.
5.  **Build Prompt:** It constructs a detailed prompt for the AI, including the commit messages, the previous tag name, and all your formatting/style preferences.
6.  **Query AI:** The prompt is sent to the configured LLM provider and model.
7.  **Display Output:** The AI's response (the generated release notes) is printed to the console. If `--copy true` was used (either via config or flag), the output is also copied to the clipboard.

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.
