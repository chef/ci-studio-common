studio:
	HAB_ORIGIN='chef' hab studio -k chef enter

test:
	docker run --rm --volume $(PWD):/ci-studio-common travisci/ci-garnet:packer-1503972846 bats /ci-studio-common/tests