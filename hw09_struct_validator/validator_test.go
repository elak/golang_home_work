package hw09_struct_validator //nolint:golint,stylecheck

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

type UserRole string

// Test the function on different structures and other types.
type (
	User struct {
		ID     string `json:"id" validate:"len:36"`
		Name   string
		Age    int      `validate:"min:18|max:50"`
		Email  string   `validate:"regexp:^\\w+@\\w+\\.\\w+$"`
		Role   UserRole `validate:"in:admin,stuff"`
		Phones []string `validate:"len:11"`
		meta   json.RawMessage
	}

	Session struct {
		User User   `validate:"nested"`
		ID   string `json:"id" validate:"len:36"`
	}

	App struct {
		Version string `validate:"len:5"`
	}

	Token struct {
		Header    []byte
		Payload   []byte `validate:"len:11"`
		Signature []byte
	}

	Response struct {
		Code int    `validate:"in:200,404,500"`
		Body string `json:"omitempty"`
	}
)

func TestValidate(t *testing.T) {

	validUser := User{"4d425464-4e82-4150-8be8-01f807f69eb9", "Name", 22, "noreply@forget.me", "admin", []string{"322-223-322"}, []byte{}}
	chaosUser := User{"TODO:UUID", "Name", 12, "hide#my.mail", "🦉", []string{"322-223-322-51"}, []byte{}}

	badIDUser := validUser
	badIDUser.ID = chaosUser.ID

	tooYoungUser := validUser
	tooYoungUser.Age = 12

	tooOldUser := validUser
	tooOldUser.Age = 112

	badMailUser := validUser
	badMailUser.Email = chaosUser.Email

	badRoleUser := validUser
	badRoleUser.Role = chaosUser.Role

	badPhoneUser := validUser
	badPhoneUser.Phones = chaosUser.Phones

	testSession := Session{badIDUser, "663cef02-58a6-4eed-bec8-4364049910a1"}

	aprilThe4th := Response{404, "Day not found"}
	reDirection := Response{301, "Head some other way"}

	tests := []struct {
		in          interface{}
		expectedErr error
	}{
		{
			validUser,
			nil,
		},
		{
			badIDUser,
			ValidationErrors{ValidationError{"ID", ErrLengthNotMatched}},
		},
		{
			badPhoneUser,
			ValidationErrors{ValidationError{"Phones[0]", ErrLengthNotMatched}},
		},
		{
			tooYoungUser,
			ValidationErrors{ValidationError{"Age", ErrLessThenMin}},
		},
		{
			tooOldUser,
			ValidationErrors{ValidationError{"Age", ErrMoreThenMax}},
		},
		{
			badMailUser,
			ValidationErrors{ValidationError{"Email", ErrPatternNotMatched}},
		},
		{
			badRoleUser,
			ValidationErrors{ValidationError{"Role", ErrNotInList}},
		},
		{
			chaosUser,
			ValidationErrors{
				ValidationError{"ID", ErrLengthNotMatched},
				ValidationError{"Age", ErrLessThenMin},
				ValidationError{"Email", ErrPatternNotMatched},
				ValidationError{"Role", ErrNotInList},
				ValidationError{"Phones[0]", ErrLengthNotMatched},
			},
		},
		{
			testSession,
			ValidationErrors{ValidationError{"User.ID", ErrLengthNotMatched}},
		},
		{
			aprilThe4th,
			nil,
		},
		{
			reDirection,
			ValidationErrors{ValidationError{"Code", ErrNotInList}},
		},
		{
			"Trust me: I'm a struct",
			ErrNotAStruct,
		},
		{
			Token{[]byte{}, []byte{11}, []byte{}},
			ValidationErrors{ValidationError{"Payload[0]", ErrFieldTypeUnsupported}},
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			err := Validate(tt.in)
			require.Equal(t, tt.expectedErr, err)
		})
	}

}
