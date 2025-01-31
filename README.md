# 🚀 Docme-Ag: AI-Powered Commit & Versioning Agent  

**Version: v1.2.0**

Docme-Ag is an AI-powered CLI tool that automates commit message generation, versioning, and documentation updates in your software projects.

## 🎯 Features
- **AI-Generated Commit Messages** - Uses Google Gemini AI to generate commit messages based on code changes.
- **Automatic Versioning** - Determines the appropriate semantic version number.
- **README Auto-Updater** - Updates README.md based on new versions and commits.
- **GitHub Integration** - Works seamlessly with GitHub repositories.

---

## 📥 Installation  

### 🔹 **Install via Installer Script**  
You can install `docm` by running the following command:  
```sh
curl -fsSL https://raw.githubusercontent.com/Dert-Ops/Docme-Ag/main/installer.sh | bash
```

### 🔹 **Manual Installation**  
Alternatively, download the binary manually:  
```sh
LATEST_VERSION=$(curl -s https://api.github.com/repos/Dert-Ops/Docme-Ag/releases/latest | grep '"tag_name":' | sed -E 's/.*"([^"]+)".*/\1/')
wget -O docm "https://github.com/Dert-Ops/Docme-Ag/releases/download/$LATEST_VERSION/docm-linux-amd64"
chmod +x docm
mv docm /usr/local/bin/
```

After installation, make sure you set up your **GEMINI API KEY**:  
```sh
echo 'export GEMINI_API_KEY="your-api-key-here"' >> ~/.bashrc
source ~/.bashrc
```
For **Mac users** (zsh default shell):  
```sh
echo 'export GEMINI_API_KEY="your-api-key-here"' >> ~/.zshrc
source ~/.zshrc
```

---

## 🚀 Usage  

Run the following commands in your terminal:  

🔹 **Generate AI-powered commit messages:**  
```sh
docm cm
```

🔹 **Generate a new semantic version:**  
```sh
docm vs
```

🔹 **Update README with the latest changes:**  
```sh
docm update-readme
```

---

## 📝 Commit Message Standard (Conventional Commits)  
This project follows the **Conventional Commits** standard for clarity and consistency.  

```
<type>(<scope>): <description>
```

- **type**: `feat`, `fix`, `chore`, `docs`, `style`, `refactor`, `test`
- **scope**: Affected module or component
- **description**: A concise description of the change

Examples:
```sh
feat(api): add user authentication
fix(ui): resolve navbar layout issue
docs(readme): update installation instructions
```

---

## 🏷️ Versioning Standard (Semantic Versioning)  
Docme-Ag follows **Semantic Versioning (SemVer)**:  
```
MAJOR.MINOR.PATCH
```
- **MAJOR** - Breaking changes  
- **MINOR** - New features (backward-compatible)  
- **PATCH** - Bug fixes and minor improvements  

Examples:
```sh
1.0.0  # Initial stable release
1.1.0  # New feature added
1.1.1  # Minor bug fix
2.0.0  # Breaking change
```

---

## 🎯 Contributing  
Contributions are welcome! Open an issue or create a pull request.

## 📜 License  
This project is licensed under the MIT License.
