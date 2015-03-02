// Copyright 2011 Google Inc. All Rights Reserved.
// This file is available under the Apache license.

package main

import (
	"os"
	"reflect"
	"sort"
	"testing"
	"time"
)

func FakeSocketWrite(f formatter, m *Metric) []string {
	var ret []string
	lc := make(chan *LabelSet)
	go m.EmitLabelSets(lc)
	for l := range lc {
		ret = append(ret, f(m, l))
	}
	sort.Strings(ret)
	return ret
}

func TestMetricToCollectd(t *testing.T) {
	ts, terr := time.Parse("2006/01/02 15:04:05", "2012/07/24 10:14:00")
	if terr != nil {
		t.Errorf("time parse error: %s", terr)
	}
	hostname, herr := os.Hostname()
	if herr != nil {
		t.Errorf("hostname error: %s", herr)
	}

	scalar_metric := NewMetric("foo", "prog", Counter)
	d, _ := scalar_metric.GetDatum()
	d.Set(37, ts)
	r := FakeSocketWrite(MetricToCollectd, scalar_metric)
	expected := []string{"PUTVAL \"" + hostname + "/mtail-prog/counter-foo\" interval=60 1343124840:37\n"}
	if !reflect.DeepEqual(expected, r) {
		t.Errorf("String didn't match:\n\texpected: %v\n\treceived: %v", expected, r)
	}

	dimensioned_metric := NewMetric("bar", "prog", Gauge, "label")
	d, _ = dimensioned_metric.GetDatum("quux")
	d.Set(37, ts)
	d, _ = dimensioned_metric.GetDatum("snuh")
	d.Set(37, ts)
	r = FakeSocketWrite(MetricToCollectd, dimensioned_metric)
	expected = []string{
		"PUTVAL \"" + hostname + "/mtail-prog/gauge-bar-label-quux\" interval=60 1343124840:37\n",
		"PUTVAL \"" + hostname + "/mtail-prog/gauge-bar-label-snuh\" interval=60 1343124840:37\n"}
	if !reflect.DeepEqual(expected, r) {
		t.Errorf("String didn't match:\n\texpected: %v\n\treceived: %v", expected, r)
	}
}

func TestMetricToGraphite(t *testing.T) {
	ts, terr := time.Parse("2006/01/02 15:04:05", "2012/07/24 10:14:00")
	if terr != nil {
		t.Errorf("time parse error: %s", terr)
	}

	scalar_metric := NewMetric("foo", "prog", Counter)
	d, _ := scalar_metric.GetDatum()
	d.Set(37, ts)
	r := FakeSocketWrite(MetricToGraphite, scalar_metric)
	expected := []string{"prog.foo 37 1343124840\n"}
	if !reflect.DeepEqual(expected, r) {
		t.Errorf("String didn't match:\n\texpected: %v\n\treceived: %v", expected, r)
	}

	dimensioned_metric := NewMetric("bar", "prog", Gauge, "l")
	d, _ = dimensioned_metric.GetDatum("quux")
	d.Set(37, ts)
	d, _ = dimensioned_metric.GetDatum("snuh")
	d.Set(37, ts)
	r = FakeSocketWrite(MetricToGraphite, dimensioned_metric)
	expected = []string{
		"prog.bar.l.quux 37 1343124840\n",
		"prog.bar.l.snuh 37 1343124840\n"}
	if !reflect.DeepEqual(expected, r) {
		t.Errorf("String didn't match:\n\texpected: %v\n\treceived: %v", expected, r)
	}
}

func TestMetricToStatsd(t *testing.T) {
	ts, terr := time.Parse("2006/01/02 15:04:05", "2012/07/24 10:14:00")
	if terr != nil {
		t.Errorf("time parse error: %s", terr)
	}

	scalar_metric := NewMetric("foo", "prog", Counter)
	d, _ := scalar_metric.GetDatum()
	d.Set(37, ts)
	r := FakeSocketWrite(MetricToStatsd, scalar_metric)
	expected := []string{"prog.foo:37|c"}
	if !reflect.DeepEqual(expected, r) {
		t.Errorf("String didn't match:\n\texpected: %v\n\treceived: %v", expected, r)
	}

	dimensioned_metric := NewMetric("bar", "prog", Gauge, "l")
	d, _ = dimensioned_metric.GetDatum("quux")
	d.Set(37, ts)
	d, _ = dimensioned_metric.GetDatum("snuh")
	d.Set(42, ts)
	r = FakeSocketWrite(MetricToStatsd, dimensioned_metric)
	expected = []string{
		"prog.bar.l.quux:37|c",
		"prog.bar.l.snuh:42|c"}
	if !reflect.DeepEqual(expected, r) {
		t.Errorf("String didn't match:\n\texpected: %v\n\treceived: %v", expected, r)
	}
}
