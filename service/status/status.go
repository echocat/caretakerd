package status
import (
    "github.com/echocat/caretakerd/panics"
    "github.com/echocat/caretakerd/errors"
    "strconv"
    "strings"
    "encoding/json"
)

type Status int
const (
    Down = Status(0)
    Running = Status(1)
    Stopped = Status(2)
    Killed = Status(3)
)

var All []Status = []Status{
    Down,
    Running,
    Stopped,
    Killed,
}

func (this Status) String() string {
    switch this {
    case Down:
        return "down"
    case Running:
        return "running"
    case Stopped:
        return "stopped"
    case Killed:
        return "killed"
    }
    panic(panics.New("Illegal status: %d", this))
}

func (this *Status) Set(value string) error {
    if valueAsInt, err := strconv.Atoi(value); err == nil {
        for _, candidate := range All {
            if int(candidate) == valueAsInt {
                (*this) = candidate
                return nil
            }
        }
        return errors.New("Illegal status: " + value)
    } else {
        lowerValue := strings.ToLower(value)
        for _, candidate := range All {
            if strings.ToLower(candidate.String()) == lowerValue {
                (*this) = candidate
                return nil
            }
        }
        return errors.New("Illegal status: " + value)
    }
}

func (this Status) MarshalJSON() ([]byte, error) {
    return json.Marshal(this.String())
}

func (this *Status) UnmarshalJSON(b []byte) error {
    var value string
    if err := json.Unmarshal(b, &value); err != nil {
        return err
    }
    return this.Set(value)
}

func (this Status) Validate() {
    this.String()
}

func (this Status) IsGoDownRequest() bool {
    switch this {
    case Stopped: fallthrough
    case Killed:
        return true
    }
    return false
}
