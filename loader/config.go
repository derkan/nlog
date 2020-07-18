package loader

import (
	"time"

	"github.com/derkan/nlog"
)

var LevelCodes = map[string]nlog.Level{
	nlog.FatalStr: nlog.FATAL,
	"FATAL":       nlog.FATAL,
	"fatal":       nlog.FATAL,
	nlog.ErrorStr: nlog.ERROR,
	"ERROR":       nlog.ERROR,
	"error":       nlog.ERROR,
	nlog.WarnStr:  nlog.WARNING,
	"WARNING":     nlog.WARNING,
	"warning":     nlog.WARNING,
	nlog.InfoStr:  nlog.INFO,
	"INFO":        nlog.INFO,
	"info":        nlog.INFO,
	nlog.DebugStr: nlog.DEBUG,
	"DEBUG":       nlog.DEBUG,
	"debug":       nlog.DEBUG,
}

var FormatterTypes = []string{"console", "json"}
var LeveledTypes = []string{"normal", "parallel"}
var WriterTypes = []string{"stdout", "stderr", "syslog", "filerotator"}

type FileRotatorConfig struct {
	// Filename is the file to write logs to.  Backup log files will be retained
	// in the same directory.  It uses <processname>-nlogrotater.log in
	// os.TempDir() if empty.
	Filename string `json:"filename" yaml:"filename"`

	// MaxSize is the maximum size in megabytes of the log file before it gets
	// rotated. It defaults to 100 megabytes.
	MaxSize int `json:"max_size" yaml:"max_size"`

	// MaxAge is the maximum number of days to retain old log files based on the
	// timestamp encoded in their filename.  Note that a day is defined as 24
	// hours and may not exactly correspond to calendar days due to daylight
	// savings, leap seconds, etc. The default is not to remove old log files
	// based on age.
	MaxAge int `json:"max_age" yaml:"max_age"`

	// MaxBackups is the maximum number of old log files to retain.  The default
	// is to retain all old log files (though MaxAge may still cause them to get
	// deleted.)
	MaxBackups int `json:"max_backups" yaml:"max_backups"`

	// UTC determines if the time used for formatting the timestamps in
	// backup files is the computer's local time.  The default is to use UTC
	// time.
	UTC bool `json:"utc" yaml:"utc"`

	// Compress determines if the rotated log files should be compressed
	// using gzip. The default is not to perform compression.
	Compress bool `json:"compress" yaml:"compress"`
}

// Writer is for final writers
type Writer struct {
	// Type is type of writer. Available options are syslog, stdout, stderr, filerotater
	TypeStr string `json:"type" yaml:"type"`
	Type    string
	FileRotatorConfig
	// QueueLen defines queue length for buffered writers.
	// Only valid for Parallel writers which uses buffered channels
	QueueLen int `json:"queue_len" yaml:"queue_len"`
	// Level is logging level
	LevelStr string `json:"level" yaml:"level"`
	Level    nlog.Level
}

// Formatter holds common formatter configs
type Formatter struct {
	// Type is type of formatter. Can be json or console
	TypeStr string `json:"type" yaml:"type"`
	Type    string
	// Level is logging level
	LevelStr string `json:"level" yaml:"level"`
	Level    nlog.Level
	// Colored determines whether to print in color only valid for console
	Colored bool `json:"colored" yaml:"colored"`
	// NoPrintLevel if set level string will not be written
	NoPrintLevel bool `json:"no_print_level" yaml:"print_level"`
	// Date sets whether to print date or not in layout 2006-01-02
	Date bool `json:"date" yaml:"date"`
	// Whether to print as string time, resolution can be set using TimeResolution
	Time bool `json:"time" yaml:"time"`
	// TimeUTC Prints time in UTC if defined
	TimeUTC bool `json:"time_utc" yaml:"time_utc"`
	// Whether to print in unix time (as int), resolution can be set using TimeResolution
	UnixTime bool `json:"unix_time" yaml:"unix_time"`
	// TimeResolution is used to getting time with specific resolution
	TimeResolutionStr string `json:"time_resolution" yaml:"time_resolution"`
	TimeResolution    time.Duration
	// FileLoc is for printing caller file location or not
	FileLoc bool `json:"file_loc" yaml:"file_loc"`
	// FileLocStrip is for stripping given path prefix from file location
	FileLocStrip string `json:"file_loc_strip" yaml:"file_loc_strip"`
	// FileLocCallerDepth is for setting depth of caller to find file loc
	FileLocCallerDepth int `json:"file_loc_caller_depth" yaml:"file_loc_caller_depth"`
	// LeveledTypeStr defines how to call write method of writers. Can be one of normal or parallel
	LeveledTypeStr string `json:"leveled" yaml:"leveled"`
	// LeveledTypeStr defines how to call write method of writers. Can be one of normal or parallel
	LeveledType string
	// Writers holds leveled writer config
	Writers []Writer `json:"writers" yaml:"writers"`
}

// Loader loads logging configs
type Loader struct {
	// Prefix adds prefix for each line log
	Prefix string `json:"prefix" yaml:"prefix"`
	// MinLevel is minimal level of logging
	MinLevelStr string `json:"min_level" yaml:"min_level"`
	MinLevel    nlog.Level
	// Formatters holds formatters
	Formatters []Formatter `json:"console_formatters" yaml:"console_formatters"`
}

func AsLevel(levelStr string, defaultValue nlog.Level) nlog.Level {
	if lvl, ok := LevelCodes[levelStr]; ok {
		return lvl
	}
	return defaultValue
}

func AsDuration(durationStr string, defautValue time.Duration) time.Duration {
	switch durationStr {
	case "ns":
		return time.Nanosecond
	case "mcs":
		return time.Microsecond
	case "mls":
		return time.Millisecond
	case "s":
		return time.Second
	case "h":
		return time.Hour
	case "m":
		return time.Minute
	default:
		return defautValue
	}
}

// StrInSlice returns true if str is in list
func StrInSlice(str string, list []string) bool {
	for _, b := range list {
		if b == str {
			return true
		}
	}
	return false
}

// CleanType checks if given type is in lst and returns it if exists.
// Else defaultVal is returned
func CleanType(lst []string, typ string, defautlVal string) string {
	if StrInSlice(typ, lst) {
		return typ
	}
	return defautlVal
}

// FromContent loads yaml content from content.
// baseKey is root key for logs if multiple root keys are existing
func FromContent(content string, baseKey string) (*Loader, error) {
	config, err := Config(content)
	if err != nil {
		return nil, err
	}
	return fromCfg(config, baseKey), nil
}

// FromFile loads yaml content from filename.
// baseKey is root key for logs if multiple root keys are existing
func FromFile(filename string, baseKey string) (*Loader, error) {
	config, err := ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return fromCfg(config, baseKey), nil
}
func fromCfg(config *File, baseKey string) *Loader {
	if baseKey != "" {
		baseKey += "."
	}
	l := new(Loader)
	l.Prefix, _ = config.Get(baseKey + "prefix")
	l.MinLevelStr, _ = config.Get(baseKey + "min_level")
	l.MinLevel = AsLevel(l.MinLevelStr, nlog.INFO)
	fmtCnt, _ := config.Count(baseKey + "formatters")
	if fmtCnt > 0 {
		l.Formatters = make([]Formatter, fmtCnt)
		for i := 0; i < fmtCnt; i++ {
			l.Formatters[i].TypeStr, _ = config.Get(baseKey+"formatters[%d].type", i)
			l.Formatters[i].Type = CleanType(FormatterTypes, l.Formatters[i].TypeStr, "stdout")
			l.Formatters[i].LevelStr, _ = config.Get(baseKey+"formatters[%d].level", i)
			l.Formatters[i].Level = AsLevel(l.Formatters[i].LevelStr, l.MinLevel)
			l.Formatters[i].NoPrintLevel, _ = config.GetBool(baseKey+"formatters[%d].no_print_level", i)
			l.Formatters[i].Date, _ = config.GetBool(baseKey+"formatters[%d].date", i)
			l.Formatters[i].Time, _ = config.GetBool(baseKey+"formatters[%d].time", i)
			l.Formatters[i].TimeUTC, _ = config.GetBool(baseKey+"formatters[%d].time_utc", i)
			l.Formatters[i].UnixTime, _ = config.GetBool(baseKey+"formatters[%d].unix_time", i)
			l.Formatters[i].TimeResolutionStr, _ = config.Get(baseKey+"formatters[%d].time_resolution", i)
			l.Formatters[i].TimeResolution = AsDuration(l.Formatters[i].TimeResolutionStr, time.Second)
			l.Formatters[i].FileLoc, _ = config.GetBool(baseKey+"formatters[%d].file_loc", i)
			l.Formatters[i].FileLocStrip, _ = config.Get(baseKey+"formatters[%d].file_loc_strip", i)
			l.Formatters[i].FileLocCallerDepth, _ = config.GetInt(baseKey+"formatters[%d].file_loc_caller_depth", i)
			l.Formatters[i].Colored, _ = config.GetBool(baseKey+"formatters[%d].colored", i)
			l.Formatters[i].LeveledTypeStr, _ = config.Get(baseKey+"formatters[%d].leveled", i)
			l.Formatters[i].LeveledType = CleanType(LeveledTypes, l.Formatters[i].LeveledTypeStr, "normal")
			writerCnt, _ := config.Count(baseKey+"formatters[%d].writers", i)
			if writerCnt > 0 {
				l.Formatters[i].Writers = make([]Writer, writerCnt)
				for j := 0; j < writerCnt; j++ {
					l.Formatters[i].Writers[j].TypeStr, _ = config.Get(baseKey+"formatters[%d].writers[%d].type", i, j)
					l.Formatters[i].Writers[j].Type = CleanType(WriterTypes, l.Formatters[i].Writers[j].TypeStr, "stdout")
					l.Formatters[i].Writers[j].LevelStr, _ = config.Get(baseKey+"formatters[%d].writers[%d].level", i, j)
					l.Formatters[i].Writers[j].Level = AsLevel(l.Formatters[i].Writers[j].LevelStr, l.Formatters[i].Level)
					l.Formatters[i].Writers[j].Filename, _ = config.Get(baseKey+"formatters[%d].writers[%d].filename", i, j)
					l.Formatters[i].Writers[j].MaxSize, _ = config.GetInt(baseKey+"formatters[%d].writers[%d].max_size", i, j)
					l.Formatters[i].Writers[j].MaxAge, _ = config.GetInt(baseKey+"formatters[%d].writers[%d].max_age", i, j)
					l.Formatters[i].Writers[j].MaxBackups, _ = config.GetInt(baseKey+"formatters[%d].writers[%d].max_backups", i, j)
					l.Formatters[i].Writers[j].UTC, _ = config.GetBool(baseKey+"formatters[%d].writers[%d].utc", i, j)
					l.Formatters[i].Writers[j].Compress, _ = config.GetBool(baseKey+"formatters[%d].writers[%d].compress", i, j)
					l.Formatters[i].Writers[j].QueueLen, _ = config.GetInt(baseKey+"formatters[%d].writers[%d].queue_len", i, j)
				}
			}
		}
	}
	return l
}
