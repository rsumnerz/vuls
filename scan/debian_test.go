/* Vuls - Vulnerability Scanner
Copyright (C) 2016  Future Architect, Inc. Japan.

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program.  If not, see <http://www.gnu.org/licenses/>.
*/

package scan

import (
	"reflect"
	"testing"

	"github.com/future-architect/vuls/config"
	"github.com/k0kubun/pp"
)

func TestParseScannedPackagesLineDebian(t *testing.T) {

	var packagetests = []struct {
		in      string
		name    string
		version string
	}{
		{"base-passwd	3.5.33", "base-passwd", "3.5.33"},
		{"bzip2	1.0.6-5", "bzip2", "1.0.6-5"},
		{"adduser	3.113+nmu3ubuntu3", "adduser", "3.113+nmu3ubuntu3"},
		{"bash	4.3-7ubuntu1.5", "bash", "4.3-7ubuntu1.5"},
		{"bsdutils	1:2.20.1-5.1ubuntu20.4", "bsdutils", "1:2.20.1-5.1ubuntu20.4"},
		{"ca-certificates	20141019ubuntu0.14.04.1", "ca-certificates", "20141019ubuntu0.14.04.1"},
		{"apt	1.0.1ubuntu2.8", "apt", "1.0.1ubuntu2.8"},
	}

	d := newDebian(config.ServerInfo{})
	for _, tt := range packagetests {
		n, v, _ := d.parseScannedPackagesLine(tt.in)
		if n != tt.name {
			t.Errorf("name: expected %s, actual %s", tt.name, n)
		}
		if v != tt.version {
			t.Errorf("version: expected %s, actual %s", tt.version, v)
		}
	}

}

func TestgetCveIDParsingChangelog(t *testing.T) {

	var tests = []struct {
		in       []string
		expected []DetectedCveID
	}{
		{
			// verubuntu1
			[]string{
				"systemd",
				"228-4ubuntu1",
				`systemd (229-2) unstable; urgency=medium
systemd (229-1) unstable; urgency=medium
systemd (228-6) unstable; urgency=medium
CVE-2015-2325: heap buffer overflow in compile_branch(). (Closes: #781795)
CVE-2015-2326: heap buffer overflow in pcre_compile2(). (Closes: #783285)
CVE-2015-3210: heap buffer overflow in pcre_compile2() /
systemd (228-5) unstable; urgency=medium
systemd (228-4) unstable; urgency=medium
systemd (228-3) unstable; urgency=medium
systemd (228-2) unstable; urgency=medium
systemd (228-1) unstable; urgency=medium
systemd (227-3) unstable; urgency=medium
systemd (227-2) unstable; urgency=medium
systemd (227-1) unstable; urgency=medium`,
			},
			[]DetectedCveID{
				{"CVE-2015-2325", models.ChangelogExactMatch},
				{"CVE-2015-2326", models.ChangelogExactMatch},
				{"CVE-2015-3210", models.ChangelogExactMatch},
			},
		},

		{
			// ver
			[]string{
				"libpcre3",
				"2:8.38-1ubuntu1",
				`pcre3 (2:8.38-2) unstable; urgency=low
pcre3 (2:8.38-1) unstable; urgency=low
pcre3 (2:8.35-8) unstable; urgency=low
pcre3 (2:8.35-7.4) unstable; urgency=medium
pcre3 (2:8.35-7.3) unstable; urgency=medium
pcre3 (2:8.35-7.2) unstable; urgency=low
CVE-2015-2325: heap buffer overflow in compile_branch(). (Closes: #781795)
CVE-2015-2326: heap buffer overflow in pcre_compile2(). (Closes: #783285)
CVE-2015-3210: heap buffer overflow in pcre_compile2() /
pcre3 (2:8.35-7.1) unstable; urgency=medium
pcre3 (2:8.35-7) unstable; urgency=medium`,
			},
			[]DetectedCveID{
				{"CVE-2015-2325", models.ChangelogExactMatch},
				{"CVE-2015-2326", models.ChangelogExactMatch},
				{"CVE-2015-3210", models.ChangelogExactMatch},
			},
		},

		{
			// ver-ubuntu3
			[]string{
				"sysvinit",
				"2.88dsf-59.2ubuntu3",
				`sysvinit (2.88dsf-59.3ubuntu1) xenial; urgency=low
sysvinit (2.88dsf-59.3) unstable; urgency=medium
CVE-2015-2325: heap buffer overflow in compile_branch(). (Closes: #781795)
CVE-2015-2326: heap buffer overflow in pcre_compile2(). (Closes: #783285)
CVE-2015-3210: heap buffer overflow in pcre_compile2() /
sysvinit (2.88dsf-59.2ubuntu3) xenial; urgency=medium
sysvinit (2.88dsf-59.2ubuntu2) wily; urgency=medium
sysvinit (2.88dsf-59.2ubuntu1) wily; urgency=medium
CVE-2015-2321: heap buffer overflow in pcre_compile2(). (Closes: #783285)
sysvinit (2.88dsf-59.2) unstable; urgency=medium
sysvinit (2.88dsf-59.1ubuntu3) wily; urgency=medium
CVE-2015-2322: heap buffer overflow in pcre_compile2(). (Closes: #783285)
sysvinit (2.88dsf-59.1ubuntu2) wily; urgency=medium
sysvinit (2.88dsf-59.1ubuntu1) wily; urgency=medium
sysvinit (2.88dsf-59.1) unstable; urgency=medium
CVE-2015-2326: heap buffer overflow in pcre_compile2(). (Closes: #783285)
sysvinit (2.88dsf-59) unstable; urgency=medium
sysvinit (2.88dsf-58) unstable; urgency=low
sysvinit (2.88dsf-57) unstable; urgency=low`,
			},
			[]DetectedCveID{
				{"CVE-2015-2325", models.ChangelogExactMatch},
				{"CVE-2015-2326", models.ChangelogExactMatch},
				{"CVE-2015-3210", models.ChangelogExactMatch},
			},
		},
		{
			// 1:ver-ubuntu3
			[]string{
				"bsdutils",
				"1:2.27.1-1ubuntu3",
				` util-linux (2.27.1-3ubuntu1) xenial; urgency=medium
util-linux (2.27.1-3) unstable; urgency=medium
CVE-2015-2325: heap buffer overflow in compile_branch(). (Closes: #781795)
CVE-2015-2326: heap buffer overflow in pcre_compile2(). (Closes: #783285)
CVE-2015-3210: heap buffer overflow in pcre_compile2() /
util-linux (2.27.1-2) unstable; urgency=medium
util-linux (2.27.1-1ubuntu4) xenial; urgency=medium
util-linux (2.27.1-1ubuntu3) xenial; urgency=medium
util-linux (2.27.1-1ubuntu2) xenial; urgency=medium
util-linux (2.27.1-1ubuntu1) xenial; urgency=medium
util-linux (2.27.1-1) unstable; urgency=medium
util-linux (2.27-3ubuntu1) xenial; urgency=medium
util-linux (2.27-3) unstable; urgency=medium
util-linux (2.27-2) unstable; urgency=medium
util-linux (2.27-1) unstable; urgency=medium
util-linux (2.27~rc2-2) experimental; urgency=medium
util-linux (2.27~rc2-1) experimental; urgency=medium
util-linux (2.27~rc1-1) experimental; urgency=medium
util-linux (2.26.2-9) unstable; urgency=medium
util-linux (2.26.2-8) experimental; urgency=medium
util-linux (2.26.2-7) experimental; urgency=medium
util-linux (2.26.2-6ubuntu3) wily; urgency=medium
CVE-2015-2329: heap buffer overflow in compile_branch(). (Closes: #781795)
util-linux (2.26.2-6ubuntu2) wily; urgency=medium
util-linux (2.26.2-6ubuntu1) wily; urgency=medium
util-linux (2.26.2-6) unstable; urgency=medium`,
			},
			[]DetectedCveID{
				{"CVE-2015-2325", models.ChangelogExactMatch},
				{"CVE-2015-2326", models.ChangelogExactMatch},
				{"CVE-2015-3210", models.ChangelogExactMatch},
				{"CVE-2016-1000000", models.ChangelogExactMatch},
			},
		},
		{
			// https://github.com/future-architect/vuls/pull/350
			[]string{
				"CVE-2015-2325",
				"CVE-2015-2326",
				"CVE-2015-3210",
			},
		},
	}

	d := newDebian(config.ServerInfo{})
	for _, tt := range tests {
		actual, _ := d.getCveIDParsingChangelog(tt.in[2], tt.in[0], tt.in[1])
		if len(actual) != len(tt.expected) {
			t.Errorf("Len of return array are'nt same. expected %#v, actual %#v", tt.expected, actual)
			continue
		}
		for i := range tt.expected {
			if !reflect.DeepEqual(tt.expected[i], actual[i]) {
				t.Errorf("expected %v, actual %v", tt.expected[i], actual[i])
			}
		}
	}

	for _, tt := range tests {
		_, err := d.getCveIDParsingChangelog(tt.in[2], tt.in[0], "version number do'nt match case")
		if err != nil {
			t.Errorf("Returning error is unexpected")
		}
	}
}

func TestGetUpdatablePackNames(t *testing.T) {

	var tests = []struct {
		in       string
		expected []string
	}{
		{ // Ubuntu 12.04
			`Reading package lists... Done
Building dependency tree
Reading state information... Done
The following packages will be upgraded:
  apt ca-certificates cpio dpkg e2fslibs e2fsprogs gnupg gpgv libc-bin libc6 libcomerr2 libpcre3
  libpng12-0 libss2 libssl1.0.0 libudev0 multiarch-support openssl tzdata udev upstart
21 upgraded, 0 newly installed, 0 to remove and 0 not upgraded.`,
			[]string{
				"apt",
				"ca-certificates",
				"cpio",
				"dpkg",
				"e2fslibs",
				"e2fsprogs",
				"gnupg",
				"gpgv",
				"libc-bin",
				"libc6",
				"libcomerr2",
				"libpcre3",
				"libpng12-0",
				"libss2",
				"libssl1.0.0",
				"libudev0",
				"multiarch-support",
				"openssl",
				"tzdata",
				"udev",
				"upstart",
			},
		},
		{ // Ubuntu 14.04
			`Reading package lists... Done
Building dependency tree
Reading state information... Done
Calculating upgrade... Done
The following packages will be upgraded:
  apt apt-utils base-files bsdutils coreutils cpio dh-python dpkg e2fslibs
  e2fsprogs gcc-4.8-base gcc-4.9-base gnupg gpgv ifupdown initscripts iproute2
  isc-dhcp-client isc-dhcp-common libapt-inst1.5 libapt-pkg4.12 libblkid1
  libc-bin libc6 libcgmanager0 libcomerr2 libdrm2 libexpat1 libffi6 libgcc1
  libgcrypt11 libgnutls-openssl27 libgnutls26 libmount1 libpcre3 libpng12-0
  libpython3.4-minimal libpython3.4-stdlib libsqlite3-0 libss2 libssl1.0.0
  libstdc++6 libtasn1-6 libudev1 libuuid1 login mount multiarch-support
  ntpdate passwd python3.4 python3.4-minimal rsyslog sudo sysv-rc
  sysvinit-utils tzdata udev util-linux
59 upgraded, 0 newly installed, 0 to remove and 0 not upgraded.
`,
			[]string{
				"apt",
				"apt-utils",
				"base-files",
				"bsdutils",
				"coreutils",
				"cpio",
				"dh-python",
				"dpkg",
				"e2fslibs",
				"e2fsprogs",
				"gcc-4.8-base",
				"gcc-4.9-base",
				"gnupg",
				"gpgv",
				"ifupdown",
				"initscripts",
				"iproute2",
				"isc-dhcp-client",
				"isc-dhcp-common",
				"libapt-inst1.5",
				"libapt-pkg4.12",
				"libblkid1",
				"libc-bin",
				"libc6",
				"libcgmanager0",
				"libcomerr2",
				"libdrm2",
				"libexpat1",
				"libffi6",
				"libgcc1",
				"libgcrypt11",
				"libgnutls-openssl27",
				"libgnutls26",
				"libmount1",
				"libpcre3",
				"libpng12-0",
				"libpython3.4-minimal",
				"libpython3.4-stdlib",
				"libsqlite3-0",
				"libss2",
				"libssl1.0.0",
				"libstdc++6",
				"libtasn1-6",
				"libudev1",
				"libuuid1",
				"login",
				"mount",
				"multiarch-support",
				"ntpdate",
				"passwd",
				"python3.4",
				"python3.4-minimal",
				"rsyslog",
				"sudo",
				"sysv-rc",
				"sysvinit-utils",
				"tzdata",
				"udev",
				"util-linux",
			},
		},
		{
			//Ubuntu12.04
			`Reading package lists... Done
Building dependency tree
Reading state information... Done
0 upgraded, 0 newly installed, 0 to remove and 0 not upgraded.`,
			[]string{},
		},
		{
			//Ubuntu14.04
			`Reading package lists... Done
Building dependency tree
Reading state information... Done
Calculating upgrade... Done
0 upgraded, 0 newly installed, 0 to remove and 0 not upgraded.`,
			[]string{},
		},
	}

	d := newDebian(config.ServerInfo{})
	for _, tt := range tests {
		actual, err := d.parseAptGetUpgrade(tt.in)
		if err != nil {
			t.Errorf("Returning error is unexpected")
		}
		if len(tt.expected) != len(actual) {
			t.Errorf("Result length is not as same as expected. expected: %d, actual: %d", len(tt.expected), len(actual))
			pp.Println(tt.expected)
			pp.Println(actual)
			return
		}
		for i := range tt.expected {
			if tt.expected[i] != actual[i] {
				t.Errorf("[%d] expected %s, actual %s", i, tt.expected[i], actual[i])
			}
		}
	}
}

func TestParseAptCachePolicy(t *testing.T) {

	var tests = []struct {
		stdout   string
		name     string
		expected packCandidateVer
	}{
		{
			// Ubuntu 16.04
			`openssl:
  Installed: 1.0.2f-2ubuntu1
  Candidate: 1.0.2g-1ubuntu2
  Version table:
     1.0.2g-1ubuntu2 500
        500 http://archive.ubuntu.com/ubuntu xenial/main amd64 Packages
 *** 1.0.2f-2ubuntu1 100
        100 /var/lib/dpkg/status`,
			"openssl",
			packCandidateVer{
				Name:      "openssl",
				Installed: "1.0.2f-2ubuntu1",
				Candidate: "1.0.2g-1ubuntu2",
			},
		},
		{
			// Ubuntu 14.04
			`openssl:
  Installed: 1.0.1f-1ubuntu2.16
  Candidate: 1.0.1f-1ubuntu2.17
  Version table:
     1.0.1f-1ubuntu2.17 0
        500 http://archive.ubuntu.com/ubuntu/ trusty-updates/main amd64 Packages
        500 http://archive.ubuntu.com/ubuntu/ trusty-security/main amd64 Packages
 *** 1.0.1f-1ubuntu2.16 0
        100 /var/lib/dpkg/status
     1.0.1f-1ubuntu2 0
        500 http://archive.ubuntu.com/ubuntu/ trusty/main amd64 Packages`,
			"openssl",
			packCandidateVer{
				Name:      "openssl",
				Installed: "1.0.1f-1ubuntu2.16",
				Candidate: "1.0.1f-1ubuntu2.17",
			},
		},
		{
			// Ubuntu 12.04
			`openssl:
  Installed: 1.0.1-4ubuntu5.33
  Candidate: 1.0.1-4ubuntu5.34
  Version table:
     1.0.1-4ubuntu5.34 0
        500 http://archive.ubuntu.com/ubuntu/ precise-updates/main amd64 Packages
        500 http://archive.ubuntu.com/ubuntu/ precise-security/main amd64 Packages
 *** 1.0.1-4ubuntu5.33 0
        100 /var/lib/dpkg/status
     1.0.1-4ubuntu3 0
        500 http://archive.ubuntu.com/ubuntu/ precise/main amd64 Packages`,
			"openssl",
			packCandidateVer{
				Name:      "openssl",
				Installed: "1.0.1-4ubuntu5.33",
				Candidate: "1.0.1-4ubuntu5.34",
			},
		},
	}

	d := newDebian(config.ServerInfo{})
	for _, tt := range tests {
		actual, err := d.parseAptCachePolicy(tt.stdout, tt.name)
		if err != nil {
			t.Errorf("Error has occurred: %s, actual: %#v", err, actual)
		}
		if !reflect.DeepEqual(tt.expected, actual) {
			e := pp.Sprintf("%v", tt.expected)
			a := pp.Sprintf("%v", actual)
			t.Errorf("expected %s, actual %s", e, a)
		}
	}
}
