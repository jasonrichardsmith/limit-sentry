package tags

import (
	"io/ioutil"
	"log"
	"strings"
	"testing"

	"k8s.io/api/admission/v1beta1"
	"k8s.io/apimachinery/pkg/runtime"
)

var (
	podnotag     []byte
	podlatesttag []byte
	podpass      []byte
)

func init() {
	var err error
	podpass, err = ioutil.ReadFile("podtest.json.pass")
	if err != nil {
		log.Fatal(err)
	}
	podnotag, err = ioutil.ReadFile("podtest.json.notag")
	if err != nil {
		log.Fatal(err)
	}
	podlatesttag, err = ioutil.ReadFile("podtest.json.latesttag")
	if err != nil {
		log.Fatal(err)
	}
}

func TestType(t *testing.T) {
	ts := TagsSentry{}
	if is.Type() != "Pod" {
		t.Fatal("Failed type test")
	}
}

func TestAdmit(t *testing.T) {
	c := Config{}
	is, err := c.LoadSentry()
	if err != nil {
		log.Fatal(err)
	}
	ar := v1beta1.AdmissionReview{
		Request: &v1beta1.AdmissionRequest{
			Object: runtime.RawExtension{
				Raw: podpass,
			},
		},
	}
	resp := is.Admit(ar)
	if !resp.Allowed {
		t.Fatal("expected passing review")
	}
	ar.Request.Object.Raw = podnotag
	resp = is.Admit(ar)
	if resp.Allowed {
		t.Fatal("Expected no tag to fail")
	}
	ar.Request.Object.Raw = podlatesttag
	resp = is.Admit(ar)
	if resp.Allowed {
		t.Fatal("Expected latest tag to fail")
	}
	ar.Request.Object.Raw = podpass[0:5]
	resp = is.Admit(ar)
	if !strings.Contains(resp.Result.Message, "json parse error") {
		t.Fatal("Expecting json parse error")
	}
}
