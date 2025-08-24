## **Features**

### **Configuration & Setup**

* Interactive setup wizard (`salah-cli setup`).
* Config validation (`salah-cli validate-config`).
* Config overrides via CLI flags (`--latitude`, `--longitude`, etc.).
* Support multiple saved locations (`--location=home`).

### **Prayer Times**

* Weekly/monthly timetable commands (`salah-cli week`, `salah-cli month`).
* Countdown to next prayer (`salah-cli next --countdown`).
* Highlight ongoing/current prayer in output.
* Export to `.ics` calendar files.

### **Location Handling**

* Auto-detect location via IP/geolocation API.

### **Notifications & Integrations**

* Desktop notifications (`notify-send`, native Mac/Windows notifications).
* Cronjob integration to notify automatically.
* Local REST API mode (`salah-cli serve`) for programmatic access.

### **User Experience Enhancements**

* Custom output formats (JSON, table, plain text).
* Colorized/pretty CLI output.
* Command for config documentation (`salah-cli config-docs`).
* Verbose/debug output to show internal calculations.

---

## **Project Quality**

### **Code Structure**

* Refactor into proper packages (`cmd`, `internal`, `pkg`) for testability & scalability.
* Abstract dependencies (time, filesystem, environment variables) to enable mocking in tests.
* Consistent error handling (centralized error logger or structured errors).

### **Testing & Reliability**

* Unit tests for all core logic (currently adding).
* Integration tests with sample configs & expected outputs.
* CI pipeline (GitHub Actions) for automated testing on each push.

### **Release & Distribution**

* GoReleaser setup for cross-platform binaries.
* Automated Homebrew tap & Linux package builds.
* Install script (curl/bash) for user convenience.

### **Documentation**

* Developer-focused README (already started, expand with config usage examples).
* CONTRIBUTING.md for external contributors.
* Example configs for common use cases (various methods/madhabs).

