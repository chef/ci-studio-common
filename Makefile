studio:
	HAB_ORIGIN='chef' hab studio -k chef enter

bats:
	docker run --rm --volume $(PWD):/ci-studio-common --tty buildkite/plugin-tester bats /ci-studio-common/tests
