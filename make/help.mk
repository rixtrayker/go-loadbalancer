# Help system for all makefiles

.PHONY: help help-all

# Main help target that includes all other help targets
help: help-all

# Show all available help sections
help-all:
	@echo "Go Load Balancer - Available targets:"
	@echo ""
	@echo "Use 'make help-<section>' for detailed help on each section:"
	@echo "  help-env        - Environment setup and dependencies"
	@echo "  help-build      - Build targets"
	@echo "  help-test       - Testing targets"
	@echo "  help-quality    - Code quality targets"
	@echo "  help-docker     - Docker targets"
	@echo "  help-clean      - Cleanup targets"
	@echo "  help-k6         - Load testing targets"
	@echo ""
	@echo "For general help:"
	@echo "  help            - Show this help message"
	@echo "  help-all        - Show all help sections" 