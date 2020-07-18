package filerotater

// Option is function type used for setting console logging config attributes
type option func(*Rotater)

// WithFilename sets logging file and path
// It uses <processname>-nlogrotater.log in os.TempDir() if empty.
func WithFilename(filename string) option {
	return func(c *Rotater) {
		c.Filename = filename
	}
}

// WithMaxSize sets the maximum size in megabytes of the log file before it gets
// rotated. It defaults to 100 megabytes.
func WithMaxSize(size int) option {
	return func(c *Rotater) {
		c.MaxSize = size
	}
}

// WithMaxAge  sets the the maximum number of days to retain old log files based on the
// timestamp encoded in their filename. If not set, old files will not be deleted.
func WithMaxAge(age int) option {
	return func(c *Rotater) {
		c.MaxAge = age
	}
}

// WithMaxBackups sets the maximum number of old log files to retain.  The default
// is to retain all old log files (though MaxAge may still cause them to get
// deleted.)
func WithMaxBackups(cnt int) option {
	return func(c *Rotater) {
		c.MaxBackups = cnt
	}
}

// WithCompress determines if the rotated log files should be compressed
// using gzip. The default is not to perform compression.
func WithCompress() option {
	return func(c *Rotater) {
		c.Compress = true
	}
}

// WithUTC determines if the time used for formatting the timestamps in
// backup files is the computer's local time. The default is to use local
// time.
func WithUTC() option {
	return func(c *Rotater) {
		c.LocalTime = false
	}
}

// NewFileRotater returns a new instance of Rotater
func NewFileRotater(opts ...option) *Rotater {
	i := &Rotater{LocalTime: true}
	// Loop through each option and set
	for _, opt := range opts {
		opt(i)
	}
	return i
}
