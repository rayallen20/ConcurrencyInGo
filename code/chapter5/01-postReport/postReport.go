func PostReport(id string) error {
	result, err := lowlevel.DoWork()
	if err != nil {
		if _, ok := err.(lowlevel.Error); ok {
			err = WrapErr(err, "cannot post report with id %q", id)
		}
		// ...
	}
}
