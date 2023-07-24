package batch

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/benhoyt/goawk/interp"
	anon "github.com/octoberswimmer/batchforce/apex"
)

func (e *Execution) getApexContext() (map[string]any, error) {
	apex := e.Apex
	apexVars, err := anon.Vars(apex)
	if err != nil {
		return nil, err
	}
	lines := []string{"\n" + `Map<String, Object> b_f_c_t_x = new Map<String, Object>();`}
	for _, v := range apexVars {
		lines = append(lines, fmt.Sprintf(`b_f_c_t_x.put('%s', %s);`, v, v))
	}
	lines = append(lines, `System.debug(JSON.serialize(b_f_c_t_x));`)
	apex = apex + strings.Join(lines, "\n")

	session := e.session()
	debugLog, err := session.Partner.ExecuteAnonymous(apex)
	if err != nil {
		return nil, err
	}
	val, err := varFromDebugLog(debugLog)
	if err != nil {
		return nil, err
	}
	var n map[string]any
	err = json.Unmarshal(val, &n)
	if err != nil {
		return nil, err
	}
	return n, nil
}

func varFromDebugLog(log string) ([]byte, error) {
	input := strings.NewReader(log)
	output := new(bytes.Buffer)
	err := interp.Exec(`$2~/USER_DEBUG/ { var = $5 } END { print var }`, "|", input, output)
	if err != nil {
		return nil, err
	}
	return output.Bytes(), nil
}
