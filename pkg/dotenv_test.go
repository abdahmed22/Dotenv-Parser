package dotenv

import (
	"os"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

type LoadFromStringTestCase struct {
	desc          string
	input         string
	expectedError error
	expectedMap   map[string]string
}

type LoadFromFileTestCase struct {
	desc          string
	path          string
	expectedError error
	expectedMap   map[string]string
}

type LoadFromFilesTestCase struct {
	desc          string
	paths         []string
	expectedError error
	expectedMap   map[string]string
}

type GetEnvTestCase struct {
	desc          string
	input         string
	expectedError error
	expectedMap   map[string]string
}

type SetEnvTestCase struct {
	desc          string
	path          string
	expectedError error
}

type GetTestCase struct {
	desc          string
	input         string
	key           string
	value         string
	expectedError error
}

func TestENV_LoadFromString(t *testing.T) {
	parser := EnvContent{}
	emptyMap := make(map[string]string)
	testCases := []LoadFromStringTestCase{
		{
			desc:          "Empty string as input",
			input:         "",
			expectedError: errFileIsEmpty,
			expectedMap:   emptyMap,
		},
		{
			desc:          "Only one comment line start with # as input",
			input:         "#This is a comment",
			expectedError: errFileIsEmpty,
			expectedMap:   emptyMap,
		},
		{
			desc: "Only comments as input",
			input: "#This is a comment 1\n " +
				"#This is a comment 2\n" +
				"#This is a comment 3\n" +
				"#This is a comment 4",
			expectedError: errFileIsEmpty,
			expectedMap:   emptyMap,
		},
		{
			desc: "Only comments as input",
			input: "# This is a comment\n " +
				"\n" +
				"key value",
			expectedError: errWrongFormat,
			expectedMap:   emptyMap,
		},
		{
			desc:          "Only one key = value line as input with spaces",
			input:         " key     =     value",
			expectedError: nil,
			expectedMap: map[string]string{
				"key": "value",
			},
		},
		{
			desc:          "Only one key = value line as input with spaces",
			input:         " key :    value ",
			expectedError: nil,
			expectedMap: map[string]string{
				"key": "value",
			},
		},
		{
			desc:          "Only one key = value line as input with no spaces",
			input:         "key=value",
			expectedError: nil,
			expectedMap: map[string]string{
				"key": "value",
			},
		},
		{
			desc:          "Only one key = value line as input with  no spaces",
			input:         "key:value",
			expectedError: nil,
			expectedMap: map[string]string{
				"key": "value",
			},
		},
		{
			desc: "Only keys = values as input with spaces",
			input: "key1 =     value1\n " +
				"key2    = value2\n" +
				"  key3 = value3\n" +
				"key4 =    value4   ",
			expectedError: nil,
			expectedMap: map[string]string{
				"key1": "value1",
				"key2": "value2",
				"key3": "value3",
				"key4": "value4",
			},
		},
		{
			desc: "Only keys : values as input with spaces",
			input: "key1 :     value1\n " +
				"key2    : value2\n" +
				"  key3 : value3\n" +
				"key4 :    value4   ",
			expectedError: nil,
			expectedMap: map[string]string{
				"key1": "value1",
				"key2": "value2",
				"key3": "value3",
				"key4": "value4",
			},
		},
		{
			desc: "Only keys = values as input with no spaces",
			input: "key1=value1\n " +
				"key2=value2\n" +
				"key3=value3\n" +
				"key4=value4",
			expectedError: nil,
			expectedMap: map[string]string{
				"key1": "value1",
				"key2": "value2",
				"key3": "value3",
				"key4": "value4",
			},
		},
		{
			desc: "Only keys : values as input with no spaces",
			input: "key1:value1\n " +
				"key2:value2\n" +
				"key3:value3\n" +
				"key4:value4   ",
			expectedError: nil,
			expectedMap: map[string]string{
				"key1": "value1",
				"key2": "value2",
				"key3": "value3",
				"key4": "value4",
			},
		},
		{
			desc: "keys : values & keys = values as input with spaces",
			input: "key1 =     value1\n " +
				"key2    : value2\n" +
				"  key3 : value3\n" +
				"key4 =    value4   ",
			expectedError: nil,
			expectedMap: map[string]string{
				"key1": "value1",
				"key2": "value2",
				"key3": "value3",
				"key4": "value4",
			},
		},
		{
			desc: "keys : values & keys = values as input with no spaces",
			input: "key1:value1\n " +
				"key2:value2\n" +
				"key3=value3\n" +
				"key4=value4",
			expectedError: nil,
			expectedMap: map[string]string{
				"key1": "value1",
				"key2": "value2",
				"key3": "value3",
				"key4": "value4",
			},
		},
		{
			desc: "Normal test case 1",
			input: "\n\n\n\n" +
				"#comment 1\n" +
				"key1:value1\n " +
				"key2:value2\n" +
				"\n\n\n\n" +
				"#comment 2\n" +
				"key3=value3\n" +
				"key4=value4",
			expectedError: nil,
			expectedMap: map[string]string{
				"key1": "value1",
				"key2": "value2",
				"key3": "value3",
				"key4": "value4",
			},
		},
		{
			desc: "Normal test case 2",
			input: "\n" +
				"#comment 1\n" +
				"#comment 2\n" +
				"#comment 3\n" +
				"key1:value1\n " +
				"key2:value2\n" +
				"\n" +
				"#Comment 4\n" +
				"key3=value3\n" +
				"key4=value4\n" +
				"key5=value5\n" +
				"key6=value6",
			expectedError: nil,
			expectedMap: map[string]string{
				"key1": "value1",
				"key2": "value2",
				"key3": "value3",
				"key4": "value4",
				"key5": "value5",
				"key6": "value6",
			},
		},
		{
			desc: "Normal test case 3",
			input: "\n\n\n\n" +
				"#comment 1\n" +
				"key1:value1\n " +
				"key2:value2\n" +
				"\n\n\n\n" +
				"#Comment 2\n" +
				"key3=value3\n" +
				"key4:value4\n" +
				"\n\n\n\n" +
				"#Comment 3\n" +
				"key5=value5\n" +
				"key6:value6\n" +
				"\n\n\n\n" +
				"#Comment 4\n" +
				"key7=value7\n" +
				"key8=value8\n",
			expectedError: nil,
			expectedMap: map[string]string{
				"key1": "value1",
				"key2": "value2",
				"key3": "value3",
				"key4": "value4",
				"key5": "value5",
				"key6": "value6",
				"key7": "value7",
				"key8": "value8",
			},
		},
		{
			desc: "Normal test case 4",
			input: "# comment 1\n" +
				"# comment 2\n" +
				"key1:value1\n " +
				"key2:value2\n" +
				"key3=value3\n" +
				"key4=value4",
			expectedError: nil,
			expectedMap: map[string]string{
				"key1": "value1",
				"key2": "value2",
				"key3": "value3",
				"key4": "value4",
			},
		},
		{
			desc: "Normal test case 5",
			input: "\n\n\n\n" +
				"#comment 1\n" +
				"key1:value1\n " +
				"key2:value2\n" +
				"\n\n\n\n" +
				"#Comment 2\n" +
				"key3=value3\n" +
				"key4=value4\n" +
				"#comment 3\n" +
				"#comment 4\n" +
				"\n\n\n\n",
			expectedError: nil,
			expectedMap: map[string]string{
				"key1": "value1",
				"key2": "value2",
				"key3": "value3",
				"key4": "value4",
			},
		},
	}

	for _, test := range testCases {
		t.Run(test.desc, func(t *testing.T) {
			resultedMap, resultedError := parser.LoadFromString(test.input)

			assert.Equal(t, test.expectedError, resultedError)
			if !reflect.DeepEqual(test.expectedMap, resultedMap) {
				t.Fail()
			}

		})
	}

}

func TestENV_LoadFromFile(t *testing.T) {
	parser := EnvContent{}
	emptyMap := make(map[string]string)
	testCases := []LoadFromFileTestCase{
		{
			desc:          "Wrong path as an input",
			path:          "no path",
			expectedError: errReadingFile,
			expectedMap:   emptyMap,
		},
		{
			desc:          "Empty file as input",
			path:          "testdata/test_00.txt",
			expectedError: errFileIsEmpty,
			expectedMap:   emptyMap,
		},
		{
			desc:          "Only one comment line start with # as input",
			path:          "testdata/test_01.txt",
			expectedError: errFileIsEmpty,
			expectedMap:   emptyMap,
		},
		{
			desc:          "Only comments as input",
			path:          "testdata/test_02.txt",
			expectedError: errFileIsEmpty,
			expectedMap:   emptyMap,
		},
		{
			desc:          "Only comments as input",
			path:          "testdata/test_03.txt",
			expectedError: errWrongFormat,
			expectedMap:   emptyMap,
		},
		{
			desc:          "Only comments as input",
			path:          "testdata/test_04.txt",
			expectedError: errWrongFormat,
			expectedMap:   emptyMap,
		},
		{
			desc:          "Only one key = value line as input with spaces",
			path:          "testdata/test_05.txt",
			expectedError: nil,
			expectedMap: map[string]string{
				"key": "value",
			},
		},
		{
			desc:          "Only one key = value line as input with spaces",
			path:          "testdata/test_06.txt",
			expectedError: nil,
			expectedMap: map[string]string{
				"key": "value",
			},
		},
		{
			desc:          "Only one key = value line as input with no spaces",
			path:          "testdata/test_07.txt",
			expectedError: nil,
			expectedMap: map[string]string{
				"key": "value",
			},
		},
		{
			desc:          "Only one key = value line as input with  no spaces",
			path:          "testdata/test_08.txt",
			expectedError: nil,
			expectedMap: map[string]string{
				"key": "value",
			},
		},
		{
			desc:          "Only keys = values as input with spaces",
			path:          "testdata/test_09.txt",
			expectedError: nil,
			expectedMap: map[string]string{
				"key1": "value1",
				"key2": "value2",
				"key3": "value3",
				"key4": "value4",
			},
		},
		{
			desc:          "Only keys : values as input with spaces",
			path:          "testdata/test_10.txt",
			expectedError: nil,
			expectedMap: map[string]string{
				"key1": "value1",
				"key2": "value2",
				"key3": "value3",
				"key4": "value4",
			},
		},
		{
			desc:          "Only keys = values as input with no spaces",
			path:          "testdata/test_11.txt",
			expectedError: nil,
			expectedMap: map[string]string{
				"key1": "value1",
				"key2": "value2",
				"key3": "value3",
				"key4": "value4",
			},
		},
		{
			desc:          "Only keys : values as input with no spaces",
			path:          "testdata/test_12.txt",
			expectedError: nil,
			expectedMap: map[string]string{
				"key1": "value1",
				"key2": "value2",
				"key3": "value3",
				"key4": "value4",
			},
		},
		{
			desc:          "keys : values & keys = values as input with spaces",
			path:          "testdata/test_13.txt",
			expectedError: nil,
			expectedMap: map[string]string{
				"key1": "value1",
				"key2": "value2",
				"key3": "value3",
				"key4": "value4",
			},
		},
		{
			desc:          "keys : values & keys = values as input with no spaces",
			path:          "testdata/test_14.txt",
			expectedError: nil,
			expectedMap: map[string]string{
				"key1": "value1",
				"key2": "value2",
				"key3": "value3",
				"key4": "value4",
			},
		},
		{
			desc:          "Normal test case 1",
			path:          "testdata/test_15.txt",
			expectedError: nil,
			expectedMap: map[string]string{
				"key1": "value1",
				"key2": "value2",
				"key3": "value3",
				"key4": "value4",
			},
		},
		{
			desc:          "Normal test case 2",
			path:          "testdata/test_16.txt",
			expectedError: nil,
			expectedMap: map[string]string{
				"key1": "value1",
				"key2": "value2",
				"key3": "value3",
				"key4": "value4",
				"key5": "value5",
				"key6": "value6",
			},
		},
		{
			desc:          "Normal test case 3",
			path:          "testdata/test_17.txt",
			expectedError: nil,
			expectedMap: map[string]string{
				"key1": "value1",
				"key2": "value2",
				"key3": "value3",
				"key4": "value4",
				"key5": "value5",
				"key6": "value6",
				"key7": "value7",
				"key8": "value8",
			},
		},
		{
			desc:          "Normal test case 4",
			path:          "testdata/test_18.txt",
			expectedError: nil,
			expectedMap: map[string]string{
				"key1": "value1",
				"key2": "value2",
				"key3": "value3",
				"key4": "value4",
			},
		},
		{
			desc:          "Normal test case 5",
			path:          "testdata/test_19.txt",
			expectedError: nil,
			expectedMap: map[string]string{
				"key1": "value1",
				"key2": "value2",
				"key3": "value3",
				"key4": "value4",
			},
		},
	}

	for _, test := range testCases {
		t.Run(test.desc, func(t *testing.T) {
			resultedMap, resultedError := parser.LoadFromFile(test.path)

			assert.Equal(t, test.expectedError, resultedError)
			if !reflect.DeepEqual(test.expectedMap, resultedMap) {
				t.Fail()
			}

		})
	}
}

func TestENV_GetEnv(t *testing.T) {
	parser := EnvContent{}
	emptyMap := make(map[string]string)
	testCases := []GetEnvTestCase{
		{
			desc:          "Empty Map",
			input:         "",
			expectedError: errEmptyMap,
			expectedMap:   emptyMap,
		},
		{
			desc:          "Only one key = value",
			input:         "key=value",
			expectedError: nil,
			expectedMap: map[string]string{
				"key": "value",
			},
		},
		{
			desc: "Normal test case 1",
			input: "key1:value1\n " +
				"key2:value2\n" +
				"key3:value3\n" +
				"key4:value4",
			expectedError: nil,
			expectedMap: map[string]string{
				"key1": "value1",
				"key2": "value2",
				"key3": "value3",
				"key4": "value4",
			},
		},
		{
			desc: "Normal test case 2",
			input: "key1:value1\n " +
				"key2:value2\n" +
				"key3:value3\n" +
				"key4:value4\n" +
				"key5:value5\n" +
				"key6:value6",
			expectedError: nil,
			expectedMap: map[string]string{
				"key1": "value1",
				"key2": "value2",
				"key3": "value3",
				"key4": "value4",
				"key5": "value5",
				"key6": "value6",
			},
		},
		{
			desc: "Normal test case 3",
			input: "key1:value1\n " +
				"key2:value2\n" +
				"key3:value3\n" +
				"key4:value4\n" +
				"key5:value5\n" +
				"key6:value6\n" +
				"key7:value7\n" +
				"key8:value8",
			expectedError: nil,
			expectedMap: map[string]string{
				"key1": "value1",
				"key2": "value2",
				"key3": "value3",
				"key4": "value4",
				"key5": "value5",
				"key6": "value6",
				"key7": "value7",
				"key8": "value8",
			},
		},
		{
			desc: "Normal test case 4",
			input: "key1:value1\n " +
				"key2:value2\n" +
				"key3:value3",
			expectedError: nil,
			expectedMap: map[string]string{
				"key1": "value1",
				"key2": "value2",
				"key3": "value3",
			},
		},
		{
			desc: "Normal test case 5",
			input: "key1:value1\n " +
				"key2:value2\n" +
				"key3:value3\n" +
				"key4:value4\n" +
				"key5:value5\n" +
				"key6:value6\n" +
				"key7:value7\n" +
				"key8:value8\n" +
				"key9:value9\n" +
				"key10:value10",
			expectedError: nil,
			expectedMap: map[string]string{
				"key1":  "value1",
				"key2":  "value2",
				"key3":  "value3",
				"key4":  "value4",
				"key5":  "value5",
				"key6":  "value6",
				"key7":  "value7",
				"key8":  "value8",
				"key9":  "value9",
				"key10": "value10",
			},
		},
	}

	for _, test := range testCases {
		t.Run(test.desc, func(t *testing.T) {

			_, _ = parser.LoadFromString(test.input)
			resultedMap, resultedError := parser.GetEnv()

			assert.Equal(t, test.expectedError, resultedError)
			if !reflect.DeepEqual(test.expectedMap, resultedMap) {
				t.Fail()
			}

		})
	}
}

func TestENV_SetEnv(t *testing.T) {
	parser := EnvContent{}
	testCases := []SetEnvTestCase{
		{
			desc:          "Empty file as input",
			path:          "testdata/test_00.txt",
			expectedError: errFileIsEmpty,
		},
	}
	for _, test := range testCases {
		t.Run(test.desc, func(t *testing.T) {
			envMap, _ := parser.LoadFromFile(test.path)
			actualError := parser.SetEnv()

			assert.Equal(t, test.expectedError, actualError)

			for key, expectedValue := range envMap {
				actualValue := os.Getenv(key)

				assert.Equal(t, expectedValue, actualValue)
			}

		})
	}
}

func TestINI_Get(t *testing.T) {
	parser := EnvContent{}
	testCases := []GetTestCase{
		{
			desc:          "Empty map as input",
			input:         "",
			key:           "key1",
			value:         "",
			expectedError: errMissingValue,
		},
		{
			desc:          "Normal case 1",
			input:         "key1:value1",
			key:           "key1",
			value:         "value1",
			expectedError: nil,
		},
		{
			desc: "Normal case 2",
			input: "key1:value1\n " +
				"key2:value2\n" +
				"key3:value3\n" +
				"key4:value4\n" +
				"key5:value5\n" +
				"key6:value6\n" +
				"key7:value7\n" +
				"key8:value8\n" +
				"key9:value9\n" +
				"key10:value10",
			key:           "key2",
			value:         "value2",
			expectedError: nil,
		},
		{
			desc: "Normal case 3",
			input: "key1:value1\n " +
				"key2:value2\n" +
				"key3:value3\n" +
				"key4:value4",
			key:           "key5",
			value:         "",
			expectedError: errMissingValue,
		},
		{
			desc: "Normal case 4",
			input: "key1:value1\n " +
				"key2:value2\n" +
				"key3:value3\n" +
				"key4:value4\n" +
				"key5:value5\n" +
				"key6:value6\n" +
				"key7:value7",
			key:           "key1",
			value:         "value1",
			expectedError: nil,
		},
		{
			desc: "Normal case 5",
			input: "key1:value1\n " +
				"key2:value2\n" +
				"key3:value3",
			key:           "key4",
			value:         "",
			expectedError: errMissingValue,
		},
	}

	for _, test := range testCases {
		t.Run(test.desc, func(t *testing.T) {
			_, _ = parser.LoadFromString(test.input)
			resultedValue, resultedError := parser.Get(test.key)

			assert.Equal(t, test.expectedError, resultedError)
			if resultedValue != test.value {
				t.Fail()
			}

		})
	}
}
