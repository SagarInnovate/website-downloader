# Contributing to Website Downloader

Thank you for your interest in contributing to Website Downloader! 🎉

We love contributions from the community and want to make it as easy as possible for you to contribute.

## 🤔 How Can I Contribute?

### Reporting Bugs 🐛

Before creating a bug report, please check existing issues to avoid duplicates.

When creating a bug report, include:
- **Clear title** describing the issue
- **Detailed description** of the problem
- **Steps to reproduce** the behavior
- **Expected vs actual behavior**
- **Screenshots** if applicable
- **Environment details**:
  - OS version (e.g., Windows 11)
  - App version (e.g., 1.0.0)
  - Website URL (if relevant)

### Suggesting Features 💡

We welcome feature suggestions! When suggesting a feature:
- **Use a clear title** describing the feature
- **Explain the problem** it solves
- **Describe the solution** you'd like
- **Consider alternatives** you've thought about
- **Explain use cases** where this would be helpful

### Submitting Code 💻

We love pull requests! Here's how:

1. **Fork the repository**
2. **Create a branch** from `main`:
   ```bash
   git checkout -b feature/your-feature-name
   ```
3. **Make your changes**
4. **Test thoroughly**
5. **Commit your changes**:
   ```bash
   git commit -m "feat: Add awesome feature"
   ```
6. **Push to your fork**:
   ```bash
   git push origin feature/your-feature-name
   ```
7. **Open a Pull Request**

## 🛠️ Development Setup

### Prerequisites

- Go 1.18 or higher
- Node.js 16 or higher
- Wails CLI v2

### Setup Steps

1. **Clone the repository**
   ```bash
   git clone https://github.com/SagarInnovate/website-downloader.git
   cd website-downloader-app
   ```

2. **Install dependencies**
   ```bash
   # Install Wails CLI
   go install github.com/wailsapp/wails/v2/cmd/wails@latest
   
   # Install Go dependencies
   go mod download
   
   # Install frontend dependencies
   cd frontend
   npm install
   cd ..
   ```

3. **Run development server**
   ```bash
   wails dev
   ```

## 📝 Code Style

### Go Code
- Follow standard Go conventions
- Run `go fmt` before committing
- Use meaningful variable names
- Add comments for complex logic
- Keep functions focused and small

### JavaScript Code
- Use ES6+ features
- Use const/let, not var
- Add JSDoc comments for functions
- Keep functions pure when possible
- Follow existing code patterns

### CSS Code
- Follow BEM methodology where applicable
- Use CSS variables for theming
- Keep selectors specific but not overly nested
- Group related styles together

## 🧪 Testing

Before submitting a PR, test your changes with:

### Functional Testing
- [ ] Different website types (blog, docs, SPA)
- [ ] Both download modes (Static and Browser)
- [ ] Various configurations in Settings
- [ ] History functionality
- [ ] Stop/Cancel feature
- [ ] File saving and naming
- [ ] Native dialogs

### Edge Cases
- [ ] Invalid URLs
- [ ] Network errors
- [ ] Very large websites
- [ ] Slow loading sites
- [ ] Sites with authentication

### UI Testing
- [ ] All views (Home, History, Settings, About)
- [ ] Responsive behavior
- [ ] Scrolling smoothness
- [ ] Button interactions
- [ ] Form validation

## 📋 Pull Request Guidelines

### PR Title Format
Use conventional commits format:
- `feat: Add new feature`
- `fix: Fix bug in scraper`
- `docs: Update README`
- `style: Format code`
- `refactor: Refactor downloader`
- `test: Add tests`
- `chore: Update dependencies`

### PR Description
Include:
- **What** - What changes did you make?
- **Why** - Why did you make these changes?
- **How** - How did you implement it?
- **Testing** - How did you test it?
- **Screenshots** - If UI changes

### Example PR Description
```markdown
## What
Added support for custom user agents in browser mode

## Why
Some websites block the default Chrome user agent, causing downloads to fail

## How
- Added userAgent field to config
- Passed custom user agent to Chromedp
- Added UI input in Settings view

## Testing
- Tested with websites that require custom user agents
- Verified Settings persistence
- Checked both Static and Browser modes

## Screenshots
![Settings UI](screenshot.png)
```

## 🎨 Design Guidelines

When making UI changes:
- Match the existing design style
- Use the established color palette
- Maintain spacing consistency
- Ensure readability
- Test on different screen sizes
- Keep native OS feel

## 📚 Documentation

When adding features:
- Update README.md
- Update QUICK_START.md if user-facing
- Add code comments
- Update CHANGELOG.md
- Consider adding examples

## 🔄 Review Process

1. **Automated checks** run on your PR
2. **Maintainer review** within 48-72 hours
3. **Feedback** if changes needed
4. **Approval** when ready
5. **Merge** by maintainer

## 💬 Communication

- Be respectful and constructive
- Ask questions if unclear
- Respond to feedback promptly
- Help others in issues and discussions

## 🏆 Recognition

Contributors are recognized in:
- Release notes
- Contributors list
- Project documentation

## 📖 Resources

- [Wails Documentation](https://wails.io/)
- [Go Documentation](https://go.dev/doc/)
- [Chromedp Documentation](https://github.com/chromedp/chromedp)

## ❓ Questions?

- Open a [Discussion](https://github.com/SagarInnovate/website-downloader/discussions)
- Ask in an [Issue](https://github.com/SagarInnovate/website-downloader/issues)

## 📜 Code of Conduct

This project follows a simple code of conduct:
- Be respectful and inclusive
- Welcome newcomers
- Focus on constructive feedback
- Respect different viewpoints
- Report unacceptable behavior

## 🙏 Thank You!

Your contributions make this project better for everyone. We appreciate your time and effort! 

---

**Happy Coding! 🚀**
