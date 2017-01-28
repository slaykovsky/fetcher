// Copyright 2017 Alexey Slaykovsky
//
// Redistribution and use in source and binary forms, with or without
// modification, are permitted provided that the following conditions are met:
//
// 1. Redistributions of source code must retain the above copyright notice,
// this list of conditions and the following disclaimer.
//
// 2. Redistributions in binary form must reproduce the above copyright notice,
// this list of conditions and the following disclaimer in the documentation
// and/or other materials provided with the distribution.
//
// 3. Neither the name of the copyright holder nor the names of its
// contributors may be used to endorse or promote products derived from this
// software without specific prior written permission.
//
// THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS"
// AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE
// IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE
// ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE
// LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR
// CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF
// SUBSTITUTE GOODS OR SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS
// INTERRUPTION) HOWEVER CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN
// CONTRACT, STRICT LIABILITY, OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE)
// ARISING IN ANY WAY OUT OF THE USE OF THIS SOFTWARE, EVEN IF ADVISED OF
// THE POSSIBILITY OF SUCH DAMAGE.

package fetcher

import (
	"io/ioutil"
	"os"
	"testing"
)

const (
	Location     string = "http://mirror.yandex.ru"
	ScratchDir   string = "/tmp"
	FileName     string = "centos/7/os/x86_64/isolinux/initrd.img"
	FileSize     int64  = 43372552
	TempFileName string = "fetcher_test"
	ExpectedUrl  string = "http://mirror.yandex.ru/centos/7/os/x86_64/isolinux/initrd.img"
)

func TestMakeURL(t *testing.T) {
	fetcher := &HTTPFetcher{
		Location:   Location,
		ScratchDir: ScratchDir,
	}

	url, err := MakeURL(FileName, fetcher)
	if err != nil {
		t.Fatal(err)
	}
	if url != ExpectedUrl {
		t.Fatal("Expected url: ", ExpectedUrl,
			" doesn't match actual url: ", url)
	}
}

func TestWriteFile(t *testing.T) {
	fetcher := &HTTPFetcher{
		Location:   Location,
		ScratchDir: ScratchDir,
	}

	url, err := MakeURL(FileName, fetcher)
	if err != nil {
		t.Fatal(err)
	}

	tempFile, err := ioutil.TempFile(fetcher.ScratchDir, TempFileName)
	if err != nil {
		t.Fatal(err)
	}
	defer tempFile.Close()
	defer os.Remove(tempFile.Name())

	fileSize, err := fetcher.WriteFile(url, tempFile)
	if err != nil {
		t.Fatal(err)
	}
	if fileSize != FileSize {
		t.Fatal("Expected file size: ", FileSize,
			" doesn't match actual file size: ", fileSize)
	}
}
