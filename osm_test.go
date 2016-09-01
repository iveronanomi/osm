package osm

import (
	"bytes"
	"encoding/xml"
	"reflect"
	"testing"
)

func TestOSMMarshal(t *testing.T) {
	c := loadChange(t, "testdata/changeset_38162206.osc")
	o1 := flattenOSM(c)
	o1.Bounds = &Bounds{1.1, 2.2, 3.3, 4.4}

	data, err := o1.Marshal()
	if err != nil {
		t.Fatalf("marshal error: %v", err)
	}

	o2, err := UnmarshalOSM(data)
	if err != nil {
		t.Fatalf("unmarshal error: %v", err)
	}

	if !reflect.DeepEqual(o1, o2) {
		t.Errorf("osm are not equal")
		t.Logf("%+v", o1)
		t.Logf("%+v", o2)
	}

	// second changeset
	c = loadChange(t, "testdata/changeset_38162210.osc")
	o1 = flattenOSM(c)

	data, err = o1.Marshal()
	if err != nil {
		t.Fatalf("marshal error: %v", err)
	}

	o2, err = UnmarshalOSM(data)
	if err != nil {
		t.Fatalf("unmarshal error: %v", err)
	}

	if !reflect.DeepEqual(o1, o2) {
		t.Errorf("osm are not equal")
		t.Logf("%+v", o1)
		t.Logf("%+v", o2)
	}
}

func TestOSMMarshalXML(t *testing.T) {
	o := &OSM{
		Nodes: Nodes{
			&Node{ID: 123},
		},
	}

	data, err := xml.Marshal(o)
	if err != nil {
		t.Fatalf("xml marshal error: %v", err)
	}

	expected := `<osm version="0.6" generator="go.osm"><node id="123" lat="0" lon="0" user="" uid="0" visible="false" version="0" changeset="0" timestamp="0001-01-01T00:00:00Z"></node></osm>`

	if !bytes.Equal(data, []byte(expected)) {
		t.Errorf("incorrect marshal, got: %s", string(data))
	}
}

func flattenOSM(c *Change) *OSM {
	o := c.Create
	if o == nil {
		o = &OSM{}
	}

	if c.Modify != nil {
		o.Nodes = append(o.Nodes, c.Modify.Nodes...)
		o.Ways = append(o.Ways, c.Modify.Ways...)
		o.Relations = append(o.Relations, c.Modify.Relations...)
	}

	if c.Delete != nil {
		o.Nodes = append(o.Nodes, c.Delete.Nodes...)
		o.Ways = append(o.Ways, c.Delete.Ways...)
		o.Relations = append(o.Relations, c.Delete.Relations...)
	}

	return o
}

func cleanXMLNameFromOSM(o *OSM) {
	for _, n := range o.Nodes {
		n.XMLName = xml.Name{}
	}

	for _, w := range o.Ways {
		w.XMLName = xml.Name{}
	}

	for _, r := range o.Relations {
		r.XMLName = xml.Name{}
	}
}
