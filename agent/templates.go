package agent

// defines tengo templates that are useful

var (
	// CatPasswd returns the contents of /etc/passwd
	CatPasswd = `
	os := import("os")
	fmt := import("fmt")
	data := os.read_file("/etc/passwd")
	fmt.println(string(data))
	`
	// CatShadow attempts to cat the contents of /etc/shadow
	CatShadow = `
	os := import("os")
	fmt := import("fmt")
	data := os.read_file("/etc/shadow")
	if is_error(data) {
		fmt.println("failed to cat shadow")
	} else {
		fmt.println(string(data))
	}
	`
)
