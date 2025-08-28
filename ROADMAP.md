## **Features**

### **Configuration & Setup**

* Interactive setup wizard (`salah-cli setup`).

### **User Experience Enhancements**

* Command for config documentation (`salah-cli config-docs`).

---

## **Project Quality**

### **Code Structure**

* Refactor into proper packages (`cmd`, `internal`, `pkg`) for testability & scalability.
* Abstract dependencies (time, filesystem, environment variables) to enable mocking in tests.

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

