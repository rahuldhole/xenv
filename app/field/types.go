package field

type Type int

const (
	NoField Type = iota
	TextField
	NumberField
	FloatField
	DateField
	DateTimeField
	FileField
	URLField
	PasswordField
	SelectField
	BooleanField
	ListField
	HiddenField
	SkipField
	ButtonField
	CheckboxField
	ColorField
	EmailField
	ImageField
	MonthField
	RadioField
	RangeField
	ResetField
	TelField
	TimeField
	WeekField
	ReadonlyField
)

type Info struct {
	Type     Type
	Label    string
	Options  []string
	Note     string
	Pattern  string // regex pattern for validation
	Required bool
	Readonly bool
	Default  string
	Script   string // shell script to run
}
