package state

import (
	"fmt"

	"gopkg.in/yaml.v2"

	"github.com/helmfile/helmfile/pkg/tmpl"
)

func (r ReleaseSpec) ExecuteTemplateExpressions(renderer *tmpl.FileRenderer) (*ReleaseSpec, error) {
	var result *ReleaseSpec
	var err error

	result, err = r.Clone()
	if err != nil {
		return nil, fmt.Errorf("failed executing template expressions in release %q: %v", r.Name, err)
	}

	{
		ts := result.Name
		result.Name, err = renderer.RenderTemplateContentToString([]byte(ts))
		if err != nil {
			return nil, fmt.Errorf("failed executing template expressions in release %q.name = %q: %v", r.Name, ts, err)
		}
	}

	{
		ts := result.Chart
		result.Chart, err = renderer.RenderTemplateContentToString([]byte(ts))
		if err != nil {
			return nil, fmt.Errorf("failed executing template expressions in release %q.chart = %q: %v", r.Name, ts, err)
		}
	}

	{
		ts := result.Namespace
		result.Namespace, err = renderer.RenderTemplateContentToString([]byte(ts))
		if err != nil {
			return nil, fmt.Errorf("failed executing template expressions in release %q.namespace = %q: %v", r.Name, ts, err)
		}
	}

	{
		ts := result.Version
		result.Version, err = renderer.RenderTemplateContentToString([]byte(ts))
		if err != nil {
			return nil, fmt.Errorf("failed executing template expressions in release %q.version = %q: %v", r.Name, ts, err)
		}
	}

	if result.WaitTemplate != nil {
		ts := *result.WaitTemplate
		resultTmpl, err := renderer.RenderTemplateContentToString([]byte(ts))
		if err != nil {
			return nil, fmt.Errorf("failed executing template expressions in release %q.version = %q: %v", r.Name, ts, err)
		}
		result.WaitTemplate = &resultTmpl
	}

	if result.InstalledTemplate != nil {
		ts := *result.InstalledTemplate
		resultTmpl, err := renderer.RenderTemplateContentToString([]byte(ts))
		if err != nil {
			return nil, fmt.Errorf("failed executing template expressions in release %q.version = %q: %v", r.Name, ts, err)
		}
		result.InstalledTemplate = &resultTmpl
	}

	if result.TillerlessTemplate != nil {
		ts := *result.TillerlessTemplate
		resultTmpl, err := renderer.RenderTemplateContentToString([]byte(ts))
		if err != nil {
			return nil, fmt.Errorf("failed executing template expressions in release %q.version = %q: %v", r.Name, ts, err)
		}
		result.TillerlessTemplate = &resultTmpl
	}

	if result.VerifyTemplate != nil {
		ts := *result.VerifyTemplate
		resultTmpl, err := renderer.RenderTemplateContentToString([]byte(ts))
		if err != nil {
			return nil, fmt.Errorf("failed executing template expressions in release %q.version = %q: %v", r.Name, ts, err)
		}
		result.VerifyTemplate = &resultTmpl
	}

	for key, val := range result.Labels {
		ts := val
		s, err := renderer.RenderTemplateContentToBuffer([]byte(ts))
		if err != nil {
			return nil, fmt.Errorf("failed executing template expressions in release %q.labels[%s] = %q: %v", r.Name, key, ts, err)
		}
		result.Labels[key] = s.String()
	}

	if len(result.ValuesTemplate) > 0 {
		for i, t := range result.ValuesTemplate {
			switch ts := t.(type) {
			case map[interface{}]interface{}:
				serialized, err := yaml.Marshal(ts)
				if err != nil {
					return nil, fmt.Errorf("failed executing template expressions in release %q.values[%d] = \"%v\": %v", r.Name, i, ts, err)
				}

				s, err := renderer.RenderTemplateContentToBuffer(serialized)
				if err != nil {
					return nil, fmt.Errorf("failed executing template expressions in release %q.values[%d] = \"%v\": %v", r.Name, i, serialized, err)
				}

				var deserialized map[interface{}]interface{}

				if err := yaml.Unmarshal(s.Bytes(), &deserialized); err != nil {
					return nil, fmt.Errorf("failed executing template expressions in release %q.values[%d] = \"%v\": %v", r.Name, i, ts, err)
				}

				result.ValuesTemplate[i] = deserialized
			}
		}

		var newvals []interface{}
		newvals = append(newvals, result.ValuesTemplate...)
		result.Values = append(newvals, result.Values...)
	}

	for i, t := range result.Values {
		switch ts := t.(type) {
		case string:
			s, err := renderer.RenderTemplateContentToBuffer([]byte(ts))
			if err != nil {
				return nil, fmt.Errorf("failed executing template expressions in release %q.values[%d] = %q: %v", r.Name, i, ts, err)
			}
			result.Values[i] = s.String()
		}
	}

	for i, t := range result.Secrets {
		switch ts := t.(type) {
		case string:
			s, err := renderer.RenderTemplateContentToBuffer([]byte(ts))
			if err != nil {
				return nil, fmt.Errorf("failed executing template expressions in release %q.secrets[%d] = %q: %v", r.Name, i, ts, err)
			}
			result.Secrets[i] = s.String()
		}
	}

	if len(result.SetValuesTemplate) > 0 {
		for i, val := range result.SetValuesTemplate {
			{
				// name
				ts := val.Name
				s, err := renderer.RenderTemplateContentToBuffer([]byte(ts))
				if err != nil {
					return nil, fmt.Errorf("failed executing template expressions in release %q.set[%d].name = %q: %v", r.Name, i, ts, err)
				}
				result.SetValuesTemplate[i].Name = s.String()
			}
			{
				// value
				ts := val.Value
				s, err := renderer.RenderTemplateContentToBuffer([]byte(ts))
				if err != nil {
					return nil, fmt.Errorf("failed executing template expressions in release %q.set[%d].value = %q: %v", r.Name, i, ts, err)
				}
				result.SetValuesTemplate[i].Value = s.String()
			}
			{
				// file
				ts := val.File
				s, err := renderer.RenderTemplateContentToBuffer([]byte(ts))
				if err != nil {
					return nil, fmt.Errorf("failed executing template expressions in release %q.set[%d].file = %q: %v", r.Name, i, ts, err)
				}
				result.SetValuesTemplate[i].File = s.String()
			}
			for j, ts := range val.Values {
				// values
				s, err := renderer.RenderTemplateContentToBuffer([]byte(ts))
				if err != nil {
					return nil, fmt.Errorf("failed executing template expressions in release %q.set[%d].values[%d] = %q: %v", r.Name, i, j, ts, err)
				}
				result.SetValuesTemplate[i].Values[j] = s.String()
			}
		}

		var newvals []SetValue
		newvals = append(newvals, result.SetValuesTemplate...)
		result.SetValues = append(newvals, result.SetValues...)
	}

	for i, n := range result.Needs {
		s, err := renderer.RenderTemplateContentToBuffer([]byte(n))
		if err != nil {
			return nil, fmt.Errorf("failed executing template expressions in release %q.needs[%d] = %q: %v", r.Name, i, n, err)
		}
		result.Needs[i] = s.String()
	}

	return result, nil
}

func (r ReleaseSpec) Clone() (*ReleaseSpec, error) {
	serialized, err := yaml.Marshal(r)
	if err != nil {
		return nil, fmt.Errorf("failed cloning release %q: %v", r.Name, err)
	}

	var deserialized ReleaseSpec

	if err := yaml.Unmarshal(serialized, &deserialized); err != nil {
		return nil, fmt.Errorf("failed cloning release %q: %v", r.Name, err)
	}

	return &deserialized, nil
}

func (r ReleaseSpec) Desired() bool {
	return r.Installed == nil || *r.Installed
}
