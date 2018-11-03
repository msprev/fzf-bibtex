package cache

import (
	"testing"
)

func TestCacheName(t *testing.T) {
	tables := []struct {
		in  string
		out string
	}{
		{"/Users/msprevak/Dropbox/msprevak/-library-/texmf/bibtex/bib/mds-bib/refs.bib", "F5KXGZLSOMXW243QOJSXMYLLF5CHE33QMJXXQL3NONYHEZLWMFVS6LLMNFRHEYLSPEWS65DFPBWWML3CNFRHIZLYF5RGSYRPNVSHGLLCNFRC64TFMZZS4YTJMI======"},
	}

	for _, table := range tables {
		out := cacheName(table.in)
		if out != table.out {
			t.Errorf("cacheName '%s' failed, got: '%s', want: '%s'.", table.in, out, table.out)
		}
	}
}

func TestIslocked(t *testing.T) {
	testname := "cachetest"
	if islocked(".", testname) {
		t.Errorf("step 1")
	}
	lock(".", testname)
	if !islocked(".", testname) {
		t.Errorf("step 2")
	}
	unlock(".", testname)
	if islocked(".", testname) {
		t.Errorf("step 3")
	}
}

func FormatBool(b bool) string {
	if b {
		return "true"
	}
	return "false"
}
